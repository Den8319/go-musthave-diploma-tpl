package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/accrual"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/model"
	"github.com/stretchr/testify/mock"
)

// Моки для репозиториев и клиента

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*model.Order, error) {
	args := m.Called(ctx, orderNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByUserID(ctx context.Context, userID int64) ([]*model.Order, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByStatus(ctx context.Context, statuses ...string) ([]*model.Order, error) {
	args := m.Called(ctx, statuses)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
}

func (m *MockOrderRepository) UpdateStatus(ctx context.Context, orderNumber string, status string, accrual float64) error {
	args := m.Called(ctx, orderNumber, status, accrual)
	return args.Error(0)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	args := m.Called(ctx, login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

type MockBalanceRepository struct {
	mock.Mock
}

func (m *MockBalanceRepository) GetBalance(ctx context.Context, userID int64) (*model.Balance, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Balance), args.Error(1)
}

func (m *MockBalanceRepository) UpdateBalance(ctx context.Context, userID int64, amount float64) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

func (m *MockBalanceRepository) CreateWithdrawal(ctx context.Context, withdrawal *model.Withdrawal) error {
	args := m.Called(ctx, withdrawal)
	return args.Error(0)
}

func (m *MockBalanceRepository) GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]*model.Withdrawal, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Withdrawal), args.Error(1)
}

func (m *MockBalanceRepository) Withdraw(ctx context.Context, userID int64, order string, sum float64) error {
	args := m.Called(ctx, userID, order, sum)
	return args.Error(0)
}

type MockAccrualClient struct {
	mock.Mock
}

func (m *MockAccrualClient) GetOrderInfo(ctx context.Context, orderNumber string) (*accrual.AccuralOrder, error) {
	args := m.Called(ctx, orderNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*accrual.AccuralOrder), args.Error(1)
}

func TestOrderService_UploadOrder(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		userID         int64
		orderNumber    string
		setupMocks     func(*MockOrderRepository, *MockUserRepository, *MockBalanceRepository)
		expectedStatus int
		expectedError  error
	}{
		{
			name:        "Пустой номер заказа",
			userID:      1,
			orderNumber: "",
			setupMocks:  func(mockOrder *MockOrderRepository, mockUser *MockUserRepository, mockBalance *MockBalanceRepository) {},
			expectedStatus: 0,
			expectedError:  errors.New("номер заказа не указан"),
		},
		{
			name:        "Некорректный номер заказа (не прошел проверку Луна)",
			userID:      1,
			orderNumber: "123456",
			setupMocks:  func(mockOrder *MockOrderRepository, mockUser *MockUserRepository, mockBalance *MockBalanceRepository) {},
			expectedStatus: 0,
			expectedError:  errors.New("некорректный номер заказа"),
		},
		{
			name:        "Заказ уже загружен этим пользователем",
			userID:      1,
			orderNumber: "542875634233",
			setupMocks: func(mockOrder *MockOrderRepository, mockUser *MockUserRepository, mockBalance *MockBalanceRepository) {
				mockOrder.On("GetByUserID", ctx, int64(1)).Return([]*model.Order{
					{ID: 1, UserID: 1, OrderNumber: "542875634233", Status: "NEW"},
				}, nil)
			},
			expectedStatus: 200,
			expectedError:  nil,
		},
		{
			name:        "Номер заказа уже был загружен другим пользователем",
			userID:      1,
			orderNumber: "542875634233",
			setupMocks: func(mockOrder *MockOrderRepository, mockUser *MockUserRepository, mockBalance *MockBalanceRepository) {
				mockOrder.On("GetByUserID", ctx, int64(1)).Return([]*model.Order{}, nil)
				mockOrder.On("GetByOrderNumber", ctx, "542875634233").Return(&model.Order{
					ID: 1, UserID: 2, OrderNumber: "542875634233", Status: "NEW",
				}, nil)
			},
			expectedStatus: 409,
			expectedError:  nil,
		},
		{
			name:        "Успешное создание нового заказа",
			userID:      1,
			orderNumber: "542875634233",
			setupMocks: func(mockOrder *MockOrderRepository, mockUser *MockUserRepository, mockBalance *MockBalanceRepository) {
				mockOrder.On("GetByUserID", ctx, int64(1)).Return([]*model.Order{}, nil)
				mockOrder.On("GetByOrderNumber", ctx, "542875634233").Return(nil, model.ErrorNotFound)
				mockOrder.On("Create", ctx, mock.AnythingOfType("*model.Order")).Return(nil)
			},
			expectedStatus: 202,
			expectedError:  nil,
		},
		{
			name:        "Заказ уже загружен (дублирование)",
			userID:      1,
			orderNumber: "542875634233",
			setupMocks: func(mockOrder *MockOrderRepository, mockUser *MockUserRepository, mockBalance *MockBalanceRepository) {
				mockOrder.On("GetByUserID", ctx, int64(1)).Return([]*model.Order{}, nil)
				mockOrder.On("GetByOrderNumber", ctx, "542875634233").Return(nil, model.ErrorNotFound)
				mockOrder.On("Create", ctx, mock.AnythingOfType("*model.Order")).Return(model.ErrorOrderExists)
			},
			expectedStatus: 409,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderRepo := &MockOrderRepository{}
			mockUserRepo := &MockUserRepository{}
			mockBalanceRepo := &MockBalanceRepository{}

			tt.setupMocks(mockOrderRepo, mockUserRepo, mockBalanceRepo)

			s := &OrderService{
				OrderRepo:   mockOrderRepo,
				UserRepo:    mockUserRepo,
				BalanceRepo: mockBalanceRepo,
			}

			status, err := s.UploadOrder(ctx, tt.userID, tt.orderNumber)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Ожидалась ошибка: %v, получено: nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Ожидалась ошибка: %v, получено: %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Ожидался nil, получена ошибка: %v", err)
				}
			}

			if status != tt.expectedStatus {
				t.Errorf("Ожидался статус: %d, получено: %d", tt.expectedStatus, status)
			}

			mockOrderRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
			mockBalanceRepo.AssertExpectations(t)
		})
	}
}

func TestOrderService_GetOrdersByUserID(t *testing.T) {
	ctx := context.Background()

	testOrders := []*model.Order{
		{ID: 1, UserID: 1, OrderNumber: "542875634233", Status: "NEW"},
		{ID: 2, UserID: 1, OrderNumber: "542875634234", Status: "PROCESSING"},
	}

	tests := []struct {
		name          string
		userID        int64
		setupMocks    func(*MockOrderRepository)
		expectedOrders []*model.Order
		expectedError error
	}{
		{
			name:     "Успешное получение заказов",
			userID:   1,
			setupMocks: func(mockOrder *MockOrderRepository) {
				mockOrder.On("GetByUserID", ctx, int64(1)).Return(testOrders, nil)
			},
			expectedOrders: testOrders,
			expectedError:  nil,
		},
		{
			name:     "Ошибка при получении заказов",
			userID:   1,
			setupMocks: func(mockOrder *MockOrderRepository) {
				mockOrder.On("GetByUserID", ctx, int64(1)).Return(nil, errors.New("database error"))
			},
			expectedOrders: nil,
			expectedError:  errors.New("database error"),
		},
		{
			name:     "Нет заказов у пользователя",
			userID:   1,
			setupMocks: func(mockOrder *MockOrderRepository) {
				mockOrder.On("GetByUserID", ctx, int64(1)).Return([]*model.Order{}, nil)
			},
			expectedOrders: []*model.Order{},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderRepo := &MockOrderRepository{}
			mockUserRepo := &MockUserRepository{}
			mockBalanceRepo := &MockBalanceRepository{}

			tt.setupMocks(mockOrderRepo)

			s := &OrderService{
				OrderRepo:   mockOrderRepo,
				UserRepo:    mockUserRepo,
				BalanceRepo: mockBalanceRepo,
			}

			orders, err := s.GetOrdersByUserID(ctx, tt.userID)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Ожидалась ошибка: %v, получено: nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Ожидалась ошибка: %v, получено: %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Ожидался nil, получена ошибка: %v", err)
				}
			}

			if len(orders) != len(tt.expectedOrders) {
				t.Errorf("Ожидалось %d заказов, получено: %d", len(tt.expectedOrders), len(orders))
			}

			mockOrderRepo.AssertExpectations(t)
		})
	}
}

func TestOrderService_pollAccrualStatuses(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name               string
		setupMocks         func(*MockOrderRepository, *MockAccrualClient, *MockBalanceRepository)
		expectLogMessages  bool
	}{
		{
			name: "Нет заказов для опроса",
			setupMocks: func(mockOrder *MockOrderRepository, mockAccrual *MockAccrualClient, mockBalance *MockBalanceRepository) {
				mockOrder.On("GetByStatus", ctx, mock.Anything).Return([]*model.Order{}, nil)
			},
			expectLogMessages: false,
		},
		{
			name: "Заказ с статусом REGISTERED",
			setupMocks: func(mockOrder *MockOrderRepository, mockAccrual *MockAccrualClient, mockBalance *MockBalanceRepository) {
				mockOrder.On("GetByStatus", ctx, mock.Anything).Return([]*model.Order{
					{ID: 1, UserID: 1, OrderNumber: "542875634233", Status: "NEW"},
				}, nil)
				mockAccrual.On("GetOrderInfo", ctx, "542875634233").Return(&accrual.AccuralOrder{
					Order:   "542875634233",
					Status:  "REGISTERED",
					Accural: 0,
				}, nil)
				mockOrder.On("UpdateStatus", ctx, "542875634233", "REGISTERED", float64(0)).Return(nil)
			},
			expectLogMessages: true,
		},
		{
			name: "Заказ с статусом PROCESSED и начислением баллов",
			setupMocks: func(mockOrder *MockOrderRepository, mockAccrual *MockAccrualClient, mockBalance *MockBalanceRepository) {
				mockOrder.On("GetByStatus", ctx, mock.Anything).Return([]*model.Order{
					{ID: 1, UserID: 1, OrderNumber: "542875634233", Status: "PROCESSING"},
				}, nil)
				mockAccrual.On("GetOrderInfo", ctx, "542875634233").Return(&accrual.AccuralOrder{
					Order:   "542875634233",
					Status:  "PROCESSED",
					Accural: 123.45,
				}, nil)
				mockOrder.On("UpdateStatus", ctx, "542875634233", "PROCESSED", float64(123.45)).Return(nil)
				mockBalance.On("UpdateBalance", ctx, int64(1), float64(123.45)).Return(nil)
			},
			expectLogMessages: true,
		},
		{
			name: "Заказ с статусом INVALID",
			setupMocks: func(mockOrder *MockOrderRepository, mockAccrual *MockAccrualClient, mockBalance *MockBalanceRepository) {
				mockOrder.On("GetByStatus", ctx, mock.Anything).Return([]*model.Order{
					{ID: 1, UserID: 1, OrderNumber: "542875634233", Status: "NEW"},
				}, nil)
				mockAccrual.On("GetOrderInfo", ctx, "542875634233").Return(&accrual.AccuralOrder{
					Order:   "542875634233",
					Status:  "INVALID",
					Accural: 0,
				}, nil)
				mockOrder.On("UpdateStatus", ctx, "542875634233", "INVALID", float64(0)).Return(nil)
			},
			expectLogMessages: true,
		},
		{
			name: "Превышен лимит запросов (ErrTooManyRequests)",
			setupMocks: func(mockOrder *MockOrderRepository, mockAccrual *MockAccrualClient, mockBalance *MockBalanceRepository) {
				mockOrder.On("GetByStatus", ctx, mock.Anything).Return([]*model.Order{
					{ID: 1, UserID: 1, OrderNumber: "542875634233", Status: "NEW"},
				}, nil)
				mockAccrual.On("GetOrderInfo", ctx, "542875634233").Return(nil, accrual.ErrTooManyRequests)
			},
			expectLogMessages: true,
		},
		{
			name: "Заказ не зарегистрирован (204 No Content)",
			setupMocks: func(mockOrder *MockOrderRepository, mockAccrual *MockAccrualClient, mockBalance *MockBalanceRepository) {
				mockOrder.On("GetByStatus", ctx, mock.Anything).Return([]*model.Order{
					{ID: 1, UserID: 1, OrderNumber: "542875634233", Status: "NEW"},
				}, nil)
				mockAccrual.On("GetOrderInfo", ctx, "542875634233").Return(nil, nil)
			},
			expectLogMessages: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderRepo := &MockOrderRepository{}
			mockAccrualClient := &MockAccrualClient{}
			mockBalanceRepo := &MockBalanceRepository{}

			tt.setupMocks(mockOrderRepo, mockAccrualClient, mockBalanceRepo)

			s := &OrderService{
				OrderRepo:   mockOrderRepo,
				UserRepo:    &MockUserRepository{},
				BalanceRepo: mockBalanceRepo,
				accrualClient: mockAccrualClient,
			}

			s.pollAccrualStatuses(ctx)

			mockOrderRepo.AssertExpectations(t)
			mockAccrualClient.AssertExpectations(t)
			mockBalanceRepo.AssertExpectations(t)
		})
	}
}

func TestOrderService_StartAccrualPolling(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	mockOrderRepo := &MockOrderRepository{}
	mockAccrualClient := &MockAccrualClient{}
	mockBalanceRepo := &MockBalanceRepository{}

	mockOrderRepo.On("GetByStatus", ctx, mock.Anything).Return([]*model.Order{}, nil)

	s := &OrderService{
		OrderRepo:   mockOrderRepo,
		UserRepo:    &MockUserRepository{},
		BalanceRepo: mockBalanceRepo,
		accrualClient: mockAccrualClient,
	}

	// Запускаем опрос в отдельной горутине
	go s.StartAccrualPolling(ctx, 100*time.Millisecond)

	// Ждем немного
	time.Sleep(300 * time.Millisecond)

	// Останавливаем опрос
	cancel()

	// Даем время на завершение
	time.Sleep(200 * time.Millisecond)

	mockOrderRepo.AssertExpectations(t)
}
