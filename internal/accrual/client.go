package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// AccrualOrder представляет ответ от внешнего сервиса начислений
type AccrualOrder struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

// Client — HTTP-клиент для взаимодействия с внешним сервисом расчёта баллов
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient создаёт новый клиент для работы с сервисом начислений
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetOrderInfo запрашивает информацию о начислении баллов по номеру заказа
func (c *Client) GetOrderInfo(ctx context.Context, orderNumber string) (*AccrualOrder, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, orderNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса к сервису начислений: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса к сервису начислений: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var order AccrualOrder
		if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
			return nil, fmt.Errorf("ошибка при декодировании ответа: %w", err)
		}
		return &order, nil

	case http.StatusNoContent:
		// 204 — заказ не зарегистрирован в системе расчёта
		return nil, nil

	case http.StatusTooManyRequests:
		// 429 — превышено количество запросов
		// Извлекаем Retry-After
		retryAfter := resp.Header.Get("Retry-After")
		log.Warn().Str("order", orderNumber).Str("retry-after", retryAfter).Msg("Превышен лимит запросов к сервису начислений")
		return nil, fmt.Errorf("превышен лимит запросов к сервису начислений: %w", ErrTooManyRequests)

	default:
		body, _ := io.ReadAll(resp.Body)
		log.Warn().Int("status", resp.StatusCode).Str("body", string(body)).Msg("Неожиданный ответ от сервиса начислений")
		return nil, fmt.Errorf("неожиданный ответ от сервиса начислений: статус %d", resp.StatusCode)
	}
}

var ErrTooManyRequests = fmt.Errorf("too many requests")