package model

import (
	"context"
	"errors"
)

type Order struct {
	ID          int64   `json:"-"`
	UserID      int64   `json:"-"`
	OrderNumber string  `json:"number"`
	Status      string  `json:"status"`
	Accrual     float64 `json:"accrual"`
	UploadedAt  string  `json:"uploaded_at"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByOrderNumber(ctx context.Context, orderNumber string) (*Order, error)
	GetByUserID(ctx context.Context, userID int64) ([]*Order, error)
}

var ErrorOrderExists = errors.New("заказ уже загружен")
