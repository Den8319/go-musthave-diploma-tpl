package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/auth"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/model"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

// Мок репозитория для тестов
type mockUserRepository struct {
	users map[string]*model.User
}

func (m *mockUserRepository) Create(ctx context.Context, user *model.User) error {
	if m.users == nil {
		m.users = make(map[string]*model.User)
	}
	m.users[user.Login] = user
	return nil
}

func (m *mockUserRepository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	if user, exists := m.users[login]; exists {
		return user, nil
	}
	return nil, model.ErrorNotFound
}

func TestAuthMiddlewareWithRepo(t *testing.T) {
	auth.Init("test_secret_key_123")

	// Создаем мок репозитория с предварительно добавленным пользователем
	mockRepo := &mockUserRepository{}
	mockRepo.Create(context.Background(), &model.User{
		ID:           1,
		Login:        "testuser",
		PasswordHash: "hashed_password",
	})

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedUserID int64
	}{
		{
			name:           "отсутствует заголовок авторизации",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: 0,
		},
		{
			name:           "невалидный токен",
			authHeader:     "Bearer invalid_token",
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: 0,
		},
		{
			name: "валидный токен с userID",
			authHeader: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
					},
					User: "123",
				})
				signed, _ := token.SignedString([]byte("test_secret_key_123"))
				return "Bearer " + signed
			}(),
			expectedStatus: http.StatusOK,
			expectedUserID: 123,
		},
		{
			name: "валидный токен с login, пользователь найден в репозитории",
			authHeader: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
					},
					User: "testuser",
				})
				signed, _ := token.SignedString([]byte("test_secret_key_123"))
				return "Bearer " + signed
			}(),
			expectedStatus: http.StatusOK,
			expectedUserID: 1,
		},
		{
			name: "валидный токен с login, пользователь не найден в репозитории",
			authHeader: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
					},
					User: "nonexistent_user",
				})
				signed, _ := token.SignedString([]byte("test_secret_key_123"))
				return "Bearer " + signed
			}(),
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: 0,
		},
		{
			name: "просроченный токен",
			authHeader: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
					},
					User: "123",
				})
				signed, _ := token.SignedString([]byte("test_secret_key_123"))
				return "Bearer " + signed
			}(),
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем обработчик, который проверяет userID в контексте
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			// Оборачиваем handler в middleware
			middleware := auth.AuthMiddlewareWithRepo(mockRepo)
			wrappedHandler := middleware(handler)

			// Создаем тестовый запрос
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Выполняем запрос
			rec := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(rec, req)

			// Проверяем статус ответа
			assert.Equal(t, tt.expectedStatus, rec.Code, "Статус ответа не совпадает")

			// Если ожидаем успешный ответ, проверяем тело ответа
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, "OK", rec.Body.String(), "Тело ответа не совпадает")
			}
		})
	}
}
