package repository

import (
	"context"
	"database/sql"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/model"
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
	if err != nil {
		// Проверка на дублирование уникального ключа
		if err.Error() == `pq: duplicate key value violates unique constraint "orders_order_number_key"` {
			return model.ErrorOrderExists
		}
		return err
	}

	order.Status = "NEW"
	order.Accrual = 0

	return nil
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
