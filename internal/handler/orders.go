package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// Orders предоставляет хендлеры для работы с заказами
type Orders struct {
	*Server
}

// UploadOrderHandler обрабатывает POST /api/user/orders
func (h *Orders) UploadOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	// Читаем тело запроса как plain text
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Err(err).Msg("Ошибка при чтении тела запроса загрузки заказа")
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Удаляем возможные пробельные символы в начале и конце
	orderNumber := strings.TrimSpace(string(body))

	// Валидация: номер заказа должен быть непустой строкой из цифр
	if orderNumber == "" {
		log.Warn().Msg("Пустой номер заказа")
		http.Error(w, "Неверный формат номера заказа", http.StatusUnprocessableEntity)
		return
	}

	// Проверка, что номер состоит только из цифр
	if !isNumeric(orderNumber) {
		log.Warn().Str("order", orderNumber).Msg("Некорректный формат номера заказа (содержатся нецифровые символы)")
		http.Error(w, "Неверный формат номера заказа", http.StatusUnprocessableEntity)
		return
	}

	// Получаем userID из контекста (после аутентификации через middleware)
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		log.Warn().Msg("Отсутствует userID в контексте")
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	status, err := h.OrderService.UploadOrder(r.Context(), userID, orderNumber)
	if err != nil {
		log.Err(err).Str("order", orderNumber).Msg("Ошибка загрузки заказа")
		// Проверяем, связана ли ошибка с неверным форматом номера
		if err.Error() == "некорректный номер заказа" {
			http.Error(w, "Неверный формат номера заказа", http.StatusUnprocessableEntity)
		} else {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(status)
	// Отправляем JSON-ответ только для 202 (новый заказ принят в обработку)
	if status == http.StatusAccepted {
		json.NewEncoder(w).Encode(map[string]int{"status": status})
	}
	// Для кода 200 (заказ уже загружен) ответ не требуется по спецификации
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
func (h *Orders) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	// Получаем userID из контекста (после аутентификации через middleware)
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		log.Warn().Msg("Отсутствует userID в контексте")
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	orders, err := h.OrderService.GetOrdersByUserID(r.Context(), userID)
	if err != nil {
		log.Err(err).Int64("user_id", userID).Msg("Ошибка получения заказов")
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
