package repository

import (
	"context"
	 "errors"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/model"

	
	"github.com/jackc/pgx/v5/pgconn"
	
)

 
type UserDb struct {
	db *DB
}

 
func NewUserRepository(db *DB) *UserDb {
	return &UserDb{db: db}
}

 func (r *UserDb) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (login, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, user.Login, user.PasswordHash).Scan(&user.ID)
	var pgErr *pgconn.PgError
	 
	if errors.As(err, &pgErr) && pgErr.Code == pgErrUniqueViolation {
		return model.ErrorUserExists
	}

	return err // Если err == nil, вернется nil. Если другая ошибка — вернется она.

}

 
func (r *UserDb) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	query := `
		SELECT id, login, password_hash
		FROM users
		WHERE login = $1
	`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, login).Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		return nil, model.ErrorNotFound
	}

	return user, nil
}
