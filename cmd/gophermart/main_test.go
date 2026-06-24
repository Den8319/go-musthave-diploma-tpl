package main

import (
	"testing"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/handler"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/config"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name       string
		dsn        string
		accrualURL string
		wantErr    bool
		wantMsg    string
	}{
		{
			name:       "корректная инициализация все параметры",
			dsn:        "postgres://user:pass@localhost:5432/dbname?sslmode=disable",
			accrualURL: "http://localhost:8080",
			wantErr:    false,
		},
		{
			name:       "БД пусто",
			dsn:        "",
			accrualURL: "http://localhost:8080",
			wantErr:    true,
			wantMsg:    "database is not define",
		},
		{
			name:       "accrual URL пусто",
			dsn:        "postgres://user:pass@localhost:5432/dbname?sslmode=disable",
			accrualURL: "",
			wantErr:    true,
			wantMsg:    "accrual system URL is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var db *repository.DB
			var err error

			if tt.dsn != "" {
				db, err = repository.New(tt.dsn)
				if err != nil {
					t.Logf("Ошибка при создании подключения к БД: %v", err)
					if !tt.wantErr {
						t.Errorf("Ожидался nil, но получена ошибка: %v", err)
					}
					return
				}
				defer db.Close()
			}

			cfg := &config.Config{
				AccrualSystemAddress: tt.accrualURL,
			}

			h, err := handler.New(db, cfg.AccrualSystemAddress)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantMsg != "" {
					assert.Contains(t, err.Error(), tt.wantMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, h)
			}
		})
	}
}