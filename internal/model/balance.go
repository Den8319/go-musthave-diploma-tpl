package model

import (
	"context"
	"errors"
)

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawRequest struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

type WithdrawResponse struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

type Withdrawal struct {
	Order       string  `json:"order"`
	Sum         int     `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
	UserID      int64   `json:"-"`
	ID          int64   `json:"-"`
}

type BalanceRepository interface {
	GetBalance(ctx context.Context, userID int64) (*Balance, error)
	UpdateBalance(ctx context.Context, userID int64, amount float64) error
	CreateWithdrawal(ctx context.Context, withdrawal *Withdrawal) error
	GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]*Withdrawal, error)
<<<<<<< HEAD
	// Выполняет проверку баланса + списание + создание записи в одной транзакции
=======
	// Withdraw выполняет атомарное списание средств: проверка баланса + списание + создание записи
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
	Withdraw(ctx context.Context, userID int64, order string, sum int) error
}

var ErrorInsufficientFunds = errors.New("недостаточно средств на счету")
