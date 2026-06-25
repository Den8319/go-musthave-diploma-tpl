package repository

import (
	"database/sql"
	"path/filepath"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const pgErrUniqueViolation = "23505"


type DB struct {
	*sql.DB
}

 
func New(dsn string) (*DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	log.Info().Msg("Подключение к базе данных успешно установлено")
	return &DB{DB: db}, nil
}

 
func (d *DB) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

 
func (d *DB) Migrate(migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	absMigrationsDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		return err
	}

	if err := goose.Up(d.DB, absMigrationsDir); err != nil {
		return err
	}
	return nil
}
