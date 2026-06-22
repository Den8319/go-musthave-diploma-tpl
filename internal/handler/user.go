package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/model"
)

// User предоставляет хендлеры для работы с пользователями
type User struct {
	*Server
}

// RegisterHandler обрабатывает POST /api/user/register
func (h *User) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Err(err).Msg("Ошибка при декодировании запроса регистрации")
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	authResp, err := h.UserService.Register(r.Context(), &req)
	if err != nil {
		log.Err(err).Str("login", req.Login).Msg("Ошибка регистрации пользователя")
		switch {
		case err.Error() == "пользователь с логином уже существует":
			http.Error(w, "Логин уже занят", http.StatusConflict)
		default:
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	// Устанавливаем заголовок Auth вместо cookies
	w.Header().Set(model.HeaderAuth, authResp.Token)
	w.Header().Set(model.ContentTypeJSON, "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authResp)
}

// LoginHandler обрабатывает POST /api/user/login
func (h *User) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Err(err).Msg("Ошибка при чтении тела запроса")
		http.Error(w, "Ошибка при чтении запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Декодируем JSON
	var req model.LoginRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Err(err).Msg("Ошибка при декодировании запроса входа")
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверяем пользователя через service
	authResp, err := h.UserService.Login(r.Context(), &req)
	if err != nil {
		if errors.Is(err, model.ErrorNotFound) {
			log.Err(err).Str("login", req.Login).Msg("HandlerPostLogin error GetUser: пользователь не найден")
			http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
			return
		}
		log.Err(err).Str("login", req.Login).Msg("HandlerPostLogin error GetUser: внутренняя ошибка")
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	log.Info().Str("login", authResp.Login).Msg("Пользователь успешно аутентифицирован")

	// Устанавливаем заголовок Auth вместо cookies
	w.Header().Set(model.HeaderAuth, authResp.Token)
	w.Header().Set(model.ContentTypeJSON, "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authResp)
}
