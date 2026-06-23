package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// UploadOrderHandler обрабатывает POST /api/user/orders
func (s *Server) UploadOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Err(err).Msg("Ошибка при чтении тела запроса загрузки заказа")
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	orderNumber := strings.TrimSpace(string(body))

	if orderNumber == "" {
		log.Warn().Msg("Пустой номер заказа")
		http.Error(w, "Неверный формат номера заказа", http.StatusUnprocessableEntity)
		return
	}

	if !isNumeric(orderNumber) {
		log.Warn().Str("order", orderNumber).Msg("Некорректный формат номера заказа (содержатся нецифровые символы)")
		http.Error(w, "Неверный формат номера заказа", http.StatusUnprocessableEntity)
		return
	}

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		log.Warn().Msg("Отсутствует userID в контексте")
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	status, err := s.OrderService.UploadOrder(r.Context(), userID, orderNumber)
	if err != nil {
		log.Err(err).Str("order", orderNumber).Msg("Ошибка загрузки заказа")
		if err.Error() == "некорректный номер заказа" {
			http.Error(w, "Неверный формат номера заказа", http.StatusUnprocessableEntity)
		} else {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(status)
	if status == http.StatusAccepted {
		json.NewEncoder(w).Encode(map[string]int{"status": status})
	}
}

// isNumeric проверяет, что строка содержит только цифры
func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

// GetOrdersHandler обрабатывает GET /api/user/orders
func (s *Server) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
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

	orders, err := s.OrderService.GetOrdersByUserID(r.Context(), userID)
	if err != nil {
		log.Err(err).Int64("user_id", userID).Msg("Ошибка получения заказов")
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
