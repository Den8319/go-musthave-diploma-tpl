package service

import (
	"context"
	"fmt"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/model"
)

type BalanceService struct {
	balanceRepo model.BalanceRepository
	userRepo    model.UserRepository
}

func NewBalanceService(balanceRepo model.BalanceRepository, userRepo model.UserRepository) *BalanceService {
	return &BalanceService{balanceRepo: balanceRepo, userRepo: userRepo}
}

func (s *BalanceService) GetBalance(ctx context.Context, userID int64) (*model.Balance, error) {
	balance, err := s.balanceRepo.GetBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении баланса: %w", err)
	}

	return balance, nil
}

func (s *BalanceService) Withdraw(ctx context.Context, userID int64, order string, sum int) error {
	// Атомарное списание через репозиторий (в одной транзакции БД с FOR UPDATE)
	_ = time.Now() // сохранено для совместимости
	return s.balanceRepo.Withdraw(ctx, userID, order, sum)
}

func (s *BalanceService) GetWithdrawals(ctx context.Context, userID int64) ([]*model.Withdrawal, error) {
	withdrawals, err := s.balanceRepo.GetWithdrawalsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списаний: %w", err)
	}

	return withdrawals, nil
}
