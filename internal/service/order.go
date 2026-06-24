package service

import (
	"context"
	"errors"
	"time"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/accrual"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/model"
	"github.com/Den8319/go-musthave-diploma-tpl/pkg/luhn"

	"github.com/rs/zerolog/log"
)

// AccrualClient - интерфейс для клиента начисления баллов
type AccrualClient interface {
	GetOrderInfo(ctx context.Context, orderNumber string) (*accrual.AccuralOrder, error)
}

type OrderService struct {
	OrderRepo     model.OrderRepository
	UserRepo      model.UserRepository
	BalanceRepo   model.BalanceRepository
	accrualClient AccrualClient
}

func NewOrderService(orderRepo model.OrderRepository, userRepo model.UserRepository, BalanceRepo model.BalanceRepository, accrualURL string) *OrderService {
	return &OrderService{
		OrderRepo:     orderRepo,
		UserRepo:      userRepo,
		BalanceRepo:   BalanceRepo,
		accrualClient: accrual.NewClient(accrualURL),
	}
}

func (s *OrderService) UploadOrder(ctx context.Context, userID int64, orderNumber string) (int, error) {
	// Валидация номера заказа
	if orderNumber == "" {
		return 0, errors.New("номер заказа не указан")
	}

	// Проверка по алгоритму Луна
	if !luhn.Valid(orderNumber) {
		return 0, errors.New("некорректный номер заказа")
	}

	// Проверка, загружен ли этот заказ уже этим пользователем
	userOrders, err := s.OrderRepo.GetByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	for _, order := range userOrders {
		if order.OrderNumber == orderNumber {
			return 200, nil // заказ уже загружен этим пользователем
		}
	}

	// Проверка, загружен ли заказ другим пользователем
	existingOrder, err := s.OrderRepo.GetByOrderNumber(ctx, orderNumber)
	if err != nil && !errors.Is(err, model.ErrorOrderNotFound) {
		return 0, err
	}

	if existingOrder != nil {
		return 409, nil // номер заказа уже был загружен другим пользователем
	}

	// Создание нового заказа
	newOrder := &model.Order{
		UserID:      userID,
		OrderNumber: orderNumber,
	}

	if err := s.OrderRepo.Create(ctx, newOrder); err != nil {
		if errors.Is(err, model.ErrorOrderExists) {
			return 409, nil
		}
		return 0, err
	}

	return 202, nil // новый номер заказа принят в обработку
}

func (s *OrderService) GetOrdersByUserID(ctx context.Context, userID int64) ([]*model.Order, error) {
	return s.OrderRepo.GetByUserID(ctx, userID)
}

// фоновый процесс опроса статусов заказов
func (s *OrderService) StartAccrualPolling(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Опрос статусов заказов завершен")
			return
		case <-ticker.C:
			s.pollAccrualStatuses(ctx)
		}
	}
}

func (s *OrderService) pollAccrualStatuses(ctx context.Context) {
	orders, err := s.OrderRepo.GetByStatus(ctx, "NEW", "PROCESSING")
	if err != nil {
		log.Err(err).Msg("Ошибка получения заказов для опроса")
		return
	}

	if len(orders) == 0 {
		log.Info().Msg("Нет заказов для опроса")
		return
	}

	log.Info().Int("count", len(orders)).Msg("Начинаем опрос статусов заказов")

	for _, order := range orders {
		accrualResult, err := s.accrualClient.GetOrderInfo(ctx, order.OrderNumber)
		if err != nil {
			if errors.Is(err, accrual.ErrTooManyRequests) {
				log.Warn().Str("order", order.OrderNumber).Msg("Превышен лимит запросов к сервису начисления")
				continue
			}
			log.Err(err).Str("order", order.OrderNumber).Msg("Ошибка получения статуса заказа")
			continue
		}

		if accrualResult == nil {
			// 204 заказ еще не зарегистрирован
			continue
		}

		switch accrualResult.Status {
		case "REGISTERED":
			err := s.OrderRepo.UpdateStatus(ctx, order.OrderNumber, "REGISTERED", accrualResult.Accural)
			if err != nil {
				log.Err(err).Str("order", order.OrderNumber).Msg("Ошибка обновления статуса заказа")
				continue
			}
		case "PROCESSING":
			err := s.OrderRepo.UpdateStatus(ctx, order.OrderNumber, "PROCESSING", accrualResult.Accural)
			if err != nil {
				log.Err(err).Str("order", order.OrderNumber).Msg("Ошибка обновления статуса заказа")
				continue
			}
			log.Info().Str("order", order.OrderNumber).Float64("accrual", accrualResult.Accural).Msg("Начислено баллов")
		case "PROCESSED":
			err := s.OrderRepo.UpdateStatus(ctx, order.OrderNumber, "PROCESSED", accrualResult.Accural)
			if err != nil {
				log.Err(err).Str("order", order.OrderNumber).Msg("Ошибка обновления статуса заказа")
				continue
			}
			if err := s.BalanceRepo.UpdateBalance(ctx, order.UserID, accrualResult.Accural); err != nil {
				log.Err(err).Str("order", order.OrderNumber).Float64("accrual", accrualResult.Accural).Msg("Ошибка начисления баллов")
				continue
			}
		case "INVALID":
			err := s.OrderRepo.UpdateStatus(ctx, order.OrderNumber, "INVALID", accrualResult.Accural)
			if err != nil {
				log.Err(err).Str("order", order.OrderNumber).Msg("Ошибка обновления статуса заказа")
				continue
			}
		default:
			log.Err(err).Str("status", accrualResult.Status).Str("order", order.OrderNumber).Msg("Неизвестный статус заказа")
			continue
		}
	}
}
