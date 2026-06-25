package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/model"

	"github.com/jackc/pgx/v5/pgconn"
)
 
type OrderDb struct {
	db *DB
}

 
func NewOrderRepository(db *DB) *OrderDb {
	return &OrderDb{db: db}
}

 
func (r *OrderDb) Create(ctx context.Context, order *model.Order) error {
	query := `
		INSERT INTO orders (user_id, order_number, status)
		VALUES ($1, $2, 'NEW')
		RETURNING id, uploaded_at
	`

	err := r.db.QueryRowContext(ctx, query, order.UserID, order.OrderNumber).Scan(&order.ID, &order.UploadedAt)
	var pgErr *pgconn.PgError
	 
	if errors.As(err, &pgErr) && pgErr.Code == pgErrUniqueViolation {
		return model.ErrorUserExists
	}
	order.Status = "NEW"
	order.Accrual = 0

	return err // Если err == nil, вернется nil. Если другая ошибка — вернется она.
}
 
func (r *OrderDb) GetByOrderNumber(ctx context.Context, orderNumber string) (*model.Order, error) {
	query := `
		SELECT id, user_id, order_number, status, accrual, uploaded_at
		FROM orders
		WHERE order_number = $1
	`

	order := &model.Order{}
	err := r.db.QueryRowContext(ctx, query, orderNumber).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNumber,
		&order.Status,
		&order.Accrual,
		&order.UploadedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrorNotFound
		}
		return nil, err
	}

	return order, nil
}

 
func (r *OrderDb) GetByUserID(ctx context.Context, userID int64) ([]*model.Order, error) {
	query := `
		SELECT id, user_id, order_number, status, accrual, uploaded_at
		FROM orders
		WHERE user_id = $1
		ORDER BY uploaded_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		order := &model.Order{}
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderNumber,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderDb) GetByStatus(ctx context.Context, statuses ...string) ([]*model.Order, error) {
	
	if len(statuses) == 0 {
		return nil, nil
	}

	query := `SELECT id, user_id, order_number, status, accrual, uploaded_at
		FROM orders
		WHERE status = ANY($1)
		ORDER BY uploaded_at ASC`

	rows, err := r.db.QueryContext(ctx, query, statuses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		order := &model.Order{}
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderNumber,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderDb) UpdateStatus(ctx context.Context, orderNumber string, status string, accrual float64) error {
	query := `
		UPDATE orders
		SET status = $1, accrual = $2
		WHERE order_number = $3
	`
	_, err := r.db.ExecContext(ctx, query, status, accrual, orderNumber)
	return err
}



			

