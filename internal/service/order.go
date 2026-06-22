package service

import (
	"context"
	"errors"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/model"
	"github.com/Den8319/go-musthave-diploma-tpl/pkg/luhn"
)

 
type OrderService struct {
	OrderRepo model.OrderRepository
	UserRepo  model.UserRepository
}

 
func NewOrderService(orderRepo model.OrderRepository, userRepo model.UserRepository) *OrderService {
	return &OrderService{OrderRepo: orderRepo, UserRepo: userRepo}
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
