package service

import (
	"context"
	"errors"
	"time"

<<<<<<< HEAD
	
    "github.com/Den8319/go-musthave-diploma-tpl/internal/accrual"
=======
	"github.com/rs/zerolog/log"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/accrual"
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
	"github.com/Den8319/go-musthave-diploma-tpl/internal/model"
	"github.com/Den8319/go-musthave-diploma-tpl/pkg/luhn"

	"github.com/rs/zerolog/log"
)


type OrderService struct {
	OrderRepo model.OrderRepository
	UserRepo  model.UserRepository
	accrualClient *accrual.Client
}

<<<<<<< HEAD
 
func NewOrderService(orderRepo model.OrderRepository, userRepo model.UserRepository, accrualURL string) *OrderService {
	return &OrderService{
		OrderRepo: orderRepo, 
		UserRepo: userRepo,
=======

func NewOrderService(orderRepo model.OrderRepository, userRepo model.UserRepository, accrualURL string) *OrderService {
	return &OrderService{
		OrderRepo: orderRepo,
		UserRepo:  userRepo,
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
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
	if err != nil && !errors.Is(err, model.ErrorNotFound) {
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
<<<<<<< HEAD
//фоновый процесс опроса статусов заказов
=======

// StartAccrualPolling запускает фоновый процесс опроса статусов заказов
// Проверяет заказы со статусом NEW, отправляет их в accrual-сервис и обновляет статусы
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
func (s *OrderService) StartAccrualPolling(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
<<<<<<< HEAD
			case <-ctx.Done():
			log.Info().Msg("Опрос статусов заказов завершен")
			return
			case <-ticker.C:
				s.pollAccrualStatuses(ctx)
			}
=======
		case <-ctx.Done():
			log.Info().Msg("Фоновый процесс опроса начислений остановлен")
			return
		case <-ticker.C:
			s.pollAccrualStatuses(ctx)
		}
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
	}
}

func (s *OrderService) pollAccrualStatuses(ctx context.Context) {
<<<<<<< HEAD

	orders, err := s.OrderRepo.GetByStatus(ctx, "NEW", "PROGRESSING")
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
				//204 заказ еще не зарегистрирован
				continue
			}
	
	switch accrualResult.Status {
		case "REGISTERED":

			err := s.OrderRepo.UpdateStatus(ctx, order.OrderNumber, "REGISTERED",0)
			if err != nil {
				log.Err(err).Str("order", order.OrderNumber).Msg("Ошибка обновления статуса заказа")
				}
		case "PROCESSING":
			err := s.OrderRepo.UpdateStatus(ctx, order.OrderNumber, "PROCESSING",0)
			if err != nil {
				log.Err(err).Str("order", order.OrderNumber).Msg("Ошибка обновления статуса заказа")
				}
			log.Info().Str("order", order.OrderNumber).Float64("accrual", accrualResult.Accural).Msg("Начислено баллов")	
		case "PROCESSED":
			err := s.OrderRepo.UpdateStatus(ctx, order.OrderNumber, "PROCESSED",0)
			if err != nil {
				log.Err(err).Str("order", order.OrderNumber).Msg("Ошибка обновления статуса заказа")
				}
		case "INVALID":
			err := s.OrderRepo.UpdateStatus(ctx, order.OrderNumber, "INVALID",0)
			if err != nil {
				log.Err(err).Str("order", order.OrderNumber).Msg("Ошибка обновления статуса заказа")
				}
		default:
			log.Err(err).Str("status", accrualResult.Status).Str("order", order.OrderNumber).Msg("Неизвестный статус заказа")
			}
		}
	}
=======
	// Получаем заказы в статусе NEW (ещё не обработаны)
	orders, err := s.OrderRepo.GetByStatus(ctx, "NEW", "PROCESSING")
	if err != nil {
		log.Err(err).Msg("Ошибка при получении заказов для опроса начислений")
		return
	}

	if len(orders) == 0 {
		return
	}

	log.Info().Int("count", len(orders)).Msg("Опрос статусов заказов в системе начислений")

	for _, order := range orders {
		// Запрашиваем статус у внешнего сервиса
		accrualResult, err := s.accrualClient.GetOrderInfo(ctx, order.OrderNumber)
		if err != nil {
			if errors.Is(err, accrual.ErrTooManyRequests) {
				// Превышен лимит запросов — ждём следующего цикла
				log.Warn().Str("order", order.OrderNumber).Msg("Превышен лимит запросов к сервису начислений, пропускаем")
				continue
			}
			// Другие ошибки — оставляем заказ для повторной попытки
			log.Err(err).Str("order", order.OrderNumber).Msg("Ошибка при запросе начисления")
			continue
		}

		if accrualResult == nil {
			// 204 — заказ ещё не зарегистрирован в системе расчёта
			// Пропускаем, статус остаётся NEW
			continue
		}

		// Обновляем статус заказа на основе ответа от сервиса начислений
		switch accrualResult.Status {
		case "REGISTERED":
			// Заказ зарегистрирован в системе, но расчёт ещё не выполнен
			// Меняем статус на PROCESSING — начисление в процессе
			err := s.OrderRepo.UpdateStatus(ctx, order.OrderNumber, "PROCESSING", 0)
			if err != nil {
				log.Err(err).Str("order", order.OrderNumber).Msg("Ошибка при обновлении статуса заказа")
			}
		case "PROCESSING":
			// Расчёт в процессе
			err := s.OrderRepo.UpdateStatus(ctx, order.OrderNumber, "PROCESSING", 0)
			if err != nil {
				log.Err(err).Str("order", order.OrderNumber).Msg("Ошибка при обновлении статуса заказа")
			}
		case "PROCESSED":
			// Расчёт завершён — начисляем баллы
			err := s.OrderRepo.UpdateStatus(ctx, order.OrderNumber, "PROCESSED", accrualResult.Accrual)
			if err != nil {
				log.Err(err).Str("order", order.OrderNumber).Msg("Ошибка при обновлении статуса заказа")
			}
			log.Info().Str("order", order.OrderNumber).Float64("accrual", accrualResult.Accrual).Msg("Начисление баллов за заказ")
		case "INVALID":
			// Заказ не принят к расчёту — финальный статус
			err := s.OrderRepo.UpdateStatus(ctx, order.OrderNumber, "INVALID", 0)
			if err != nil {
				log.Err(err).Str("order", order.OrderNumber).Msg("Ошибка при обновлении статуса заказа")
			}
		default:
			log.Warn().Str("status", accrualResult.Status).Str("order", order.OrderNumber).Msg("Неизвестный статус от сервиса начислений")
		}
	}
}
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
