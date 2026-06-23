package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
<<<<<<< HEAD
		SELECT  current_balance, withdrawn_balance
=======
		SELECT current_balance, withdrawn_balance
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
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
<<<<<<< HEAD
		
=======
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
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
<<<<<<< HEAD
		
=======
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
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

	err := r.db.QueryRowContext(ctx, query,
<<<<<<< HEAD
		 withdrawal.UserID,
		 withdrawal.Order,
		 float64(withdrawal.Sum),
		 withdrawal.ProcessedAt,
		 ).Scan(&withdrawal.ID)
=======
		withdrawal.UserID,
		withdrawal.Order,
		float64(withdrawal.Sum),
		withdrawal.ProcessedAt,
	).Scan(&withdrawal.ID)
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
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

<<<<<<< HEAD
=======
// Withdraw выполняет атомарное списание в рамках одной транзакции:
// 1. Проверяет баланс (с блокировкой SELECT FOR UPDATE)
// 2. Списание + обновление withdrawn_balance
// 3. Создание записи о списании
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
func (r *BalanceDb) Withdraw(ctx context.Context, userID int64, order string, sum int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
<<<<<<< HEAD
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

=======
	defer tx.Rollback() // безопасно — no-op после Commit

	// 1. Проверка баланса с блокировкой строки
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
	var current, withdrawn float64
	err = tx.QueryRowContext(ctx, `
		SELECT current_balance, withdrawn_balance
		FROM balances
		WHERE user_id = $1
		FOR UPDATE`, userID,
	).Scan(&current, &withdrawn)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		current = 0
		withdrawn = 0
	}

	if current < float64(sum) {
		return model.ErrorInsufficientFunds
	}

<<<<<<< HEAD
=======
	// 2. Обновляем баланс
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
	current -= float64(sum)
	withdrawn += float64(sum)

	result, err := tx.ExecContext(ctx, `
		UPDATE balances
		SET current_balance = $1, withdrawn_balance = $2
		WHERE user_id = $3`, current, withdrawn, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO balances (user_id, current_balance, withdrawn_balance)
			VALUES ($1, $2, $3)`, userID, current, withdrawn)
		if err != nil {
			return err
		}
	}

<<<<<<< HEAD
	_, err = tx.ExecContext(ctx, `
		INSERT INTO withdrawals (user_id, order_number, amount, processed_at)
		VALUES ($1, $2, $3, $4)`, userID, order, float64(sum), time.Now())
=======
	// 3. Создаём запись о списании
	_, err = tx.ExecContext(ctx, `
		INSERT INTO withdrawals (user_id, order_number, amount, processed_at)
		VALUES ($1, $2, $3, $4)`,
		userID, order, float64(sum), time.Now().Format(time.RFC3339))
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
	if err != nil {
		return err
	}

<<<<<<< HEAD
	err = tx.Commit()
	return err
}


		



=======
	return tx.Commit()
}
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
