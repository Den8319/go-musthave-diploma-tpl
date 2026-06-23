package repository

import (
	"context"
	"database/sql"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/model"
)

type BalanceDb struct {
	db *DB
}

func NewBalanceRepository(db *DB) *BalanceDb {
	return &BalanceDb{db: db}
}

func (r *BalanceDb) GetBalance(ctx context.Context, userID int64) (*model.Balance, error) {
	query := `
		SELECT user_id, current_balance, withdrawn_balance
		FROM balances
		WHERE user_id = $1
	`

	balance := &model.Balance{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&balance.Current,
		&balance.Withdrawn,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если баланса нет, создаем новый с нулевым балансом
			return &model.Balance{Current: 0, Withdrawn: 0}, nil
		}
		return nil, err
	}

	return balance, nil
}

func (r *BalanceDb) UpdateBalance(ctx context.Context, userID int64, amount float64) error {
	query := `
		UPDATE balances
		SET current_balance = current_balance + $1
		WHERE user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, amount, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Если баланса нет, создаем новый
		query = `
			INSERT INTO balances (user_id, current_balance, withdrawn_balance)
			VALUES ($1, $2, 0)
		`
		_, err = r.db.ExecContext(ctx, query, userID, amount)
		return err
	}

	return nil
}

func (r *BalanceDb) CreateWithdrawal(ctx context.Context, withdrawal *model.Withdrawal) error {
	query := `
		INSERT INTO withdrawals (user_id, order_number, amount, processed_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, withdrawal.UserID, withdrawal.Order, float64(withdrawal.Sum), withdrawal.ProcessedAt).Scan(&withdrawal.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *BalanceDb) GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]*model.Withdrawal, error) {
	query := `
		SELECT id, user_id, order_number, amount, processed_at
		FROM withdrawals
		WHERE user_id = $1
		ORDER BY processed_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []*model.Withdrawal
	for rows.Next() {
		withdrawal := &model.Withdrawal{}
		var amount sql.NullFloat64
		err := rows.Scan(&withdrawal.ID, &withdrawal.UserID, &withdrawal.Order, &amount, &withdrawal.ProcessedAt)
		if err != nil {
			return nil, err
		}
		if amount.Valid {
			withdrawal.Sum = int(amount.Float64)
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return withdrawals, nil
}
