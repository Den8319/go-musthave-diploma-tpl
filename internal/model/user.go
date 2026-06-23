package model

import (
	"context"
	"errors"
)

type User struct {
	ID           int64  `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"-"`
	AuthToken    string `json:"-"`
}

 
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByLogin(ctx context.Context, login string) (*User, error)
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Login string `json:"login"`
	Token string `json:"token"`
}

var ErrorNotFound = errors.New("user not found")

var ErrorUserExists = errors.New("пользователь с логином уже существует")
