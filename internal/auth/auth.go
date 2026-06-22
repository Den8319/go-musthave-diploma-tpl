package auth

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/rs/zerolog/log"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/model"
	
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	User string
}

var secretKey string

// Init инициализирует секретный ключ для JWT из строки секрета конфигурации
func Init(secret string) {
	if secret == "" {
		log.Fatal().Msg("SECRET_KEY должен быть заполнен")
	}
	
	// Используем секрет напрямую
	secretKey = secret
	log.Info().Msg("Секретный ключ успешно инициализирован")
}

// GenerateToken генерирует JWT токен
func GenerateToken(user, password string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(model.Expire)),
		},
		User: user,
	})

	auth, err := token.SignedString([]byte(password))
	if err != nil {
		return "", err
	}

	return auth, nil
}

// NewToken аутентифицирует пользователя и генерирует токен
func NewToken(login, password string) (string, error) {
	// Генерируем токен, используя секретный ключ
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(model.Expire)),
		},
		User: login,
	})

	auth, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("ошибка при подписи токена: %w", err)
	}

	return auth, nil
}

// GetUser извлекает пользователя из токена
func GetUser(ctx context.Context, tokenString string) (int64, string, error) {
	// Парсим и валидируем токен
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неподдерживаемый алгоритм подписи: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return 0, "", model.ErrorNotFound
	}

	// Извлекаем claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, "", model.ErrorNotFound
	}

	// Возвращаем userID (в формате int64) и login
	var userID int64
	_, err = fmt.Sscanf(claims.User, "%d", &userID)
	if err != nil {
		// Если не удалось распарсить как число, используем login напрямую
		userID = 0
	}

	return userID, claims.User, nil
}

// HashPassword хеширует пароль с помощью bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("ошибка при хешировании пароля: %w", err)
	}
	return string(hash), nil
}

// CheckPasswordHash проверяет, соответствует ли пароль хешу
func CheckPasswordHash(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, fmt.Errorf("ошибка при проверке пароля: %w", err)
	}
	return true, nil
}
