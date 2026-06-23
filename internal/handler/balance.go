package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/model"
	"github.com/Den8319/go-musthave-diploma-tpl/pkg/luhn"
<<<<<<< HEAD

=======
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
)

// GetBalanceHandler обрабатывает GET /api/user/balance
func (s *Server) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		log.Warn().Msg("Отсутствует userID в контексте")
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	balance, err := s.BalanceService.GetBalance(r.Context(), userID)
	if err != nil {
		log.Err(err).Int64("user_id", userID).Msg("Ошибка получения баланса")
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
}
// WithdrawHandler обрабатывает POST /api/user/balance/withdraw
func (s *Server) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	var req model.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Err(err).Msg("Ошибка при декодировании запроса списания")
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

<<<<<<< HEAD
	// Проверка валидности номера заказа по алгоритму Луна
	if !luhn.Valid(req.Order) {
		log.Warn().Str("order", req.Order).Msg("Неверный номер заказа (не прошёл проверку Luhn)")
		http.Error(w, "Неверный номер заказа", http.StatusUnprocessableEntity)
		return // ← это было упущено!
=======
	// Проверка номера заказа по алгоритму Луна
	if !luhn.Valid(req.Order) {
		http.Error(w, "Неверный номер заказа", http.StatusUnprocessableEntity)
		return
>>>>>>> 2b1a5516c0a0d0c48d0447c35020101ddd378b53
	}

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		log.Warn().Msg("Отсутствует userID в контексте")
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	err := s.BalanceService.Withdraw(r.Context(), userID, req.Order, req.Sum)
	if err != nil {
		log.Err(err).Int64("user_id", userID).Int("sum", req.Sum).Msg("Ошибка списания средств")
		switch {
		case errors.Is(err, model.ErrorInsufficientFunds):
			http.Error(w, "Недостаточно средств на счету", http.StatusPaymentRequired)
		default:
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"order": req.Order,
		"sum":   req.Sum,
	})
}

func (s *Server) GetWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
		}

	userID, ok := r.Context().Value("user_id").(int64)
		if !ok {
			log.Warn().Msg("Отсутствует userID в контексте")
			http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
			return
		}

	withdraws, err := s.BalanceService.GetWithdrawals(r.Context(), userID)
		if err != nil {
			log.Err(err).Int64("user_id", userID).Msg("Ошибка получения списаний")
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
			}
	
			if len(withdraws) == 0 {
				http.Error(w, "Списаний не найдено", http.StatusNoContent)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(withdraws)
}