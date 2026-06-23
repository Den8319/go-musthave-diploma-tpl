package service

import (
	"context"
	"fmt"
	"time"

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
	// Получаем текущий баланс
	balance, err := s.balanceRepo.GetBalance(ctx, userID)
	if err != nil {
		return fmt.Errorf("ошибка при получении баланса: %w", err)
	}

	// Проверяем, достаточно ли средств
	if balance.Current < float64(sum) {
		return model.ErrorInsufficientFunds
	}

	// Обновляем баланс
	err = s.balanceRepo.UpdateBalance(ctx, userID, float64(-sum))
	if err != nil {
		return fmt.Errorf("ошибка при обновлении баланса: %w", err)
	}

	// Создаем запись о списании
	withdrawal := &model.Withdrawal{
		UserID:      userID,
		Order:       order,
		Sum:         sum,
		ProcessedAt: time.Now().Format(time.RFC3339),
	}

	err = s.balanceRepo.CreateWithdrawal(ctx, withdrawal)
	if err != nil {
		return fmt.Errorf("ошибка при создании записи о списании: %w", err)
	}

	return nil
}

func (s *BalanceService) GetWithdrawals(ctx context.Context, userID int64) ([]*model.Withdrawal, error) {
	withdrawals, err := s.balanceRepo.GetWithdrawalsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списаний: %w", err)
	}

	return withdrawals, nil
}
