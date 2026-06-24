package accrual

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Тестируем создание клиента
func TestNewClient(t *testing.T) {
	baseURL := "http://localhost:8080"
	client := NewClient(baseURL)

	if client.baseURL != baseURL {
		t.Errorf("Ожидаемый baseURL: %s, получили: %s", baseURL, client.baseURL)
	}

	if client.httpClient == nil {
		t.Error("httpClient не должен быть nil")
	}

	if client.httpClient.Timeout != 10*time.Second {
		t.Errorf("Ожидаемый timeout: 10s, получили: %v", client.httpClient.Timeout)
	}
}

// Тестируем GetOrderInfo с успешным ответом
func TestGetOrderInfo_Success(t *testing.T) {
	// Создаем тестовый HTTP-сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что запрос пришел на правильный путь
		if r.Method != http.MethodGet {
			t.Errorf("Ожидаемый метод: GET, получили: %s", r.Method)
		}

		// Возвращаем успешный ответ
		order := AccuralOrder{
			Order:   "order123",
			Status:  "REGISTERED",
			Accural: 100.5,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(order)
	}))
	defer ts.Close()

	client := NewClient(ts.URL)

	ctx := context.Background()
	order, err := client.GetOrderInfo(ctx, "order123")

	if err != nil {
		t.Errorf("Ожидаемый err: nil, получили: %v", err)
	}

	if order == nil {
		t.Error("Ожидаемый order: не nil")
	}

	if order.Order != "order123" {
		t.Errorf("Ожидаемый Order: order123, получили: %s", order.Order)
	}

	if order.Status != "REGISTERED" {
		t.Errorf("Ожидаемый Status: REGISTERED, получили: %s", order.Status)
	}

	if order.Accural != 100.5 {
		t.Errorf("Ожидаемый Accural: 100.5, получили: %f", order.Accural)
	}
}

// Тестируем GetOrderInfo с ответом 204 (заказ не зарегистрирован)
func TestGetOrderInfo_NoContent(t *testing.T) {
	// Создаем тестовый HTTP-сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client := NewClient(ts.URL)

	ctx := context.Background()
	order, err := client.GetOrderInfo(ctx, "order123")

	if err != nil {
		t.Errorf("Ожидаемый err: nil, получили: %v", err)
	}

	if order != nil {
		t.Error("Ожидаемый order: nil, получили: не nil")
	}
}

// Тестируем GetOrderInfo с ответом 429 (превышено количество запросов)
func TestGetOrderInfo_TooManyRequests(t *testing.T) {
	// Создаем тестовый HTTP-сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer ts.Close()

	client := NewClient(ts.URL)

	ctx := context.Background()
	order, err := client.GetOrderInfo(ctx, "order123")

	if err == nil {
		t.Error("Ожидаемый err: не nil, получили: nil")
	}

	if order != nil {
		t.Error("Ожидаемый order: nil, получили: не nil")
	}
}

// Тестируем GetOrderInfo с неожиданным статусом ответа
func TestGetOrderInfo_UnexpectedStatus(t *testing.T) {
	// Создаем тестовый HTTP-сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := NewClient(ts.URL)

	ctx := context.Background()
	order, err := client.GetOrderInfo(ctx, "order123")

	if err == nil {
		t.Error("Ожидаемый err: не nil, получили: nil")
	}

	if order != nil {
		t.Error("Ожидаемый order: nil, получили: не nil")
	}
}

// Тестируем GetOrderInfo с ошибкой контекста
func TestGetOrderInfo_ContextCancelled(t *testing.T) {
	// Создаем тестовый HTTP-сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ничего не делаем, просто ждем
	}))
	defer ts.Close()

	client := NewClient(ts.URL)

	// Создаем отменяемый контекст
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Отменяем контекст

	order, err := client.GetOrderInfo(ctx, "order123")

	if err == nil {
		t.Error("Ожидаемый err: не nil, получили: nil")
	}

	if order != nil {
		t.Error("Ожидаемый order: nil, получили: не nil")
	}
}
