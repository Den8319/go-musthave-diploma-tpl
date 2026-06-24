package main

import (
	"testing"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/handler"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/config"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name       string
		db         *repository.DB
		accrualURL string
		wantErr    bool
		wantMsg    string
	}{
		{
			name:       "корректная инициализация все параметры",
			db:         nil,
			accrualURL: "http://localhost:8080",
			wantErr:    true,
			wantMsg:    "database is not define",
		},
		{
			name:       "БД пусто",
			db:         nil,
			accrualURL: "http://localhost:8080",
			wantErr:    true,
			wantMsg:    "database is not define",
		},
		{
			name:       "accrual URL пусто",
			db:         &repository.DB{},
			accrualURL: "",
			wantErr:    true,
			wantMsg:    "accrual system URL is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				AccrualSystemAddress: tt.accrualURL,
			}

			h, err := handler.New(tt.db, cfg.AccrualSystemAddress)

			if tt.wantErr {
				if assert.Error(t, err) && tt.wantMsg != "" {
					assert.Contains(t, err.Error(), tt.wantMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, h)
			}
		})
	}
}

func TestNewHandlerWithMockDB(t *testing.T) {
	// Создаем mock-подключение к PostgreSQL
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании mock-подключения: %v", err)
	}
	defer db.Close()

	// Создаем repository.DB из sql.DB
	repoDB := &repository.DB{DB: db}

	cfg := &config.Config{
		AccrualSystemAddress: "http://localhost:8080",
	}

	// handler.New не будет использовать реальную БД, так как мы передаем mock
	h, err := handler.New(repoDB, cfg.AccrualSystemAddress)

	assert.NoError(t, err)
	assert.NotNil(t, h)
}
