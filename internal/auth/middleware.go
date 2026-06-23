package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/model"
)

func AuthMiddlewareWithRepo(userRepo model.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем токен из заголовка Authorization
			authHeader := r.Header.Get(model.HeaderAuth)
			if authHeader == "" {
				log.Warn().Msg("Отсутствует заголовок авторизации")
				http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			userID, login, err := GetUser(r.Context(), tokenString)
			if err != nil {
				log.Warn().Err(err).Msg("Неверный токен авторизации")
				http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
				return
			}

			if userID == 0 && login != "" {
				user, err := userRepo.GetByLogin(r.Context(), login)
				if err != nil {
					log.Warn().Err(err).Str("login", login).Msg("Пользователь не найден")
					http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
					return
				}
				userID = user.ID
			}

			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
