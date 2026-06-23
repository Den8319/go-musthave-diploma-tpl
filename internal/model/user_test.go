package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserStruct(t *testing.T) {
	// Положительный сценарий: создание корректной структуры User
	user := &User{
		ID:           1,
		Login:        "testuser",
		PasswordHash: "hashed_password",
		AuthToken:    "auth_token",
	}

	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "testuser", user.Login)
	assert.Equal(t, "hashed_password", user.PasswordHash)
	assert.Equal(t, "auth_token", user.AuthToken)
}

func TestRegisterRequestStruct(t *testing.T) {
	// Положительный сценарий: создание корректной структуры RegisterRequest
	req := &RegisterRequest{
		Login:    "testuser",
		Password: "password123",
	}

	assert.Equal(t, "testuser", req.Login)
	assert.Equal(t, "password123", req.Password)
}

func TestLoginRequestStruct(t *testing.T) {
	// Положительный сценарий: создание корректной структуры LoginRequest
	req := &LoginRequest{
		Login:    "testuser",
		Password: "password123",
	}

	assert.Equal(t, "testuser", req.Login)
	assert.Equal(t, "password123", req.Password)
}

func TestAuthResponseStruct(t *testing.T) {
	// Положительный сценарий: создание корректной структуры AuthResponse
	resp := &AuthResponse{
		Login: "testuser",
		Token: "auth_token",
	}

	assert.Equal(t, "testuser", resp.Login)
	assert.Equal(t, "auth_token", resp.Token)
}

func TestRegisterRequest_EmptyLogin(t *testing.T) {
	// Негативный сценарий: пустой логин
	req := &RegisterRequest{
		Login:    "",
		Password: "password123",
	}

	assert.Equal(t, "", req.Login)
}

func TestRegisterRequest_EmptyPassword(t *testing.T) {
	// Негативный сценарий: пустой пароль
	req := &RegisterRequest{
		Login:    "testuser",
		Password: "",
	}

	assert.Equal(t, "", req.Password)
}

func TestLoginRequest_EmptyLogin(t *testing.T) {
	// Негативный сценарий: пустой логин
	req := &LoginRequest{
		Login:    "",
		Password: "password123",
	}

	assert.Equal(t, "", req.Login)
}

func TestLoginRequest_EmptyPassword(t *testing.T) {
	// Негативный сценарий: пустой пароль
	req := &LoginRequest{
		Login:    "testuser",
		Password: "",
	}

	assert.Equal(t, "", req.Password)
}
