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


type AccuralOrder struct {
	Order			string		`json:"order"`
	Status			string		`json:"status"`
	Accural			float64			`json:"accrual"`
}

type Client struct {
	baseURL string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
			    baseURL: baseURL,
			    httpClient: &http.Client{
						Timeout: 10 * time.Second,
				},
		}
	}

var ErrTooManyRequests = fmt.Errorf("too many requests")


func (s *Client) GetOrderInfo(ctx context.Context, orderNumber string) (*AccuralOrder, error) {
	url := fmt.Sprintf("%s/api/user/orders/%s", s.baseURL, orderNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса к сервису начислений: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса к сервису начислений: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var order AccuralOrder
		if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
			return nil, fmt.Errorf("ошибка при декодировании ответа: %w", err)
		}
		return &order, nil

	case http.StatusNoContent: // 204 — заказ не зарегистрирован
		return nil, nil

	case http.StatusTooManyRequests: // 429 — превышено количество подключений
		retryAfter := resp.Header.Get("Retry-After")
		log.Warn().Str("order", orderNumber).Str("retry_after", retryAfter).Msg("Превышено количество запросов к сервису начислений")
		return nil, fmt.Errorf("превышено количество запросов к сервису начислений")

	default:
		body, _ := io.ReadAll(resp.Body)
		log.Warn().Int("status_code", resp.StatusCode).Str("body", string(body)).Msg("Неожиданный статус ответа от сервиса начислений")
		return nil, fmt.Errorf("неожиданный статус ответа от сервиса начислений: %d", resp.StatusCode)
	}
}

 


 
