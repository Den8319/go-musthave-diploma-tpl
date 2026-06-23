package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/auth"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/model"
)

// UserService предоставляет бизнес-логику для работы с пользователями
type UserService struct {
	repo model.UserRepository
}

// NewUserService создаёт новый экземпляр UserService
func NewUserService(repo model.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Register регистрирует нового пользователя
func (s *UserService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	// Проверка обязательных полей
	if req.Login == "" || req.Password == "" {
		return nil, errors.New("login и password обязательны")
	}

	

	// Хеширование пароля
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка при хешировании пароля: %w", err)
	}

	// Создание пользователя с хешем пароля
	user := &model.User{
		Login:        req.Login,
		PasswordHash: hash,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Генерация JWT токена для ответа
	token, err := auth.NewToken(req.Login, req.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка при генерации токена: %w", err)
	}

	return &model.AuthResponse{
		Login: req.Login,
		Token: token,
	}, nil
}

// Login аутентифицирует пользователя
func (s *UserService) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	// Проверка обязательных полей
	if req.Login == "" || req.Password == "" {
		return nil, errors.New("login и password обязательны")
	}

	// Получение пользователя
	user, err := s.repo.GetByLogin(ctx, req.Login)
	if err != nil {
		return nil, model.ErrorNotFound
	}

	// Проверка пароля по хешу
	passwordValid, err := auth.CheckPasswordHash(req.Password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке пароля: %w", err)
	}
	if !passwordValid {
		return nil, model.ErrorNotFound
	}

	// Генерация JWT токена для ответа
	token, err := auth.NewToken(req.Login, req.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка при генерации токена: %w", err)
	}

	return &model.AuthResponse{
		Login: req.Login,
		Token: token,
	}, nil
}

// GetUserByLogin получает пользователя по логину
func (s *UserService) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	return s.repo.GetByLogin(ctx, login)
}
