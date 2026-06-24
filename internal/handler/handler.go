package handler

import (
	"context"
	"net/http"
	"time"
	"errors"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/auth"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/repository"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/service"
	"github.com/rs/zerolog/log"
)

type Server struct {
	*http.ServeMux
	UserService    *service.UserService
	OrderService   *service.OrderService
	BalanceService *service.BalanceService
	cancelPoolling context.CancelFunc
}

func New(db *repository.DB, accrualAddress string) (*Server, error) {

    if db == nil {
        return nil, errors.New("database is not define")
    }
    
    
    if accrualAddress == "" {
        return nil, errors.New("accrual system URL is empty")
    }

	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	balanceRepo := repository.NewBalanceRepository(db)

	userService := service.NewUserService(userRepo)
	orderService := service.NewOrderService(orderRepo, userRepo,balanceRepo, accrualAddress)
	balanceService := service.NewBalanceService(balanceRepo, userRepo)

	mux := &Server{
		ServeMux:       http.NewServeMux(),
		UserService:    userService,
		OrderService:   orderService,
		BalanceService: balanceService,
	}

	mux.registerHandlers()

	ctx, cancel := context.WithCancel(context.Background())
	mux.cancelPoolling = cancel
	go func() {
		orderService.StartAccrualPolling (ctx, 10*time.Second) // Периодичность опроса
	}()

	log.Info().Msg("Фоновый процесс опроса начисленний запущен")

	return mux, nil
}

// Останавливает фоновый polling
func (s *Server) Shutdown() {
	if s.cancelPoolling != nil {
		s.cancelPoolling()
	}
}


func (s *Server) registerHandlers() {
	s.HandleFunc("/api/user/register", s.RegisterHandler)
	s.HandleFunc("/api/user/login", s.LoginHandler)

	// Оборачиваем защищённые хендлеры в middleware аутентификации
	// Передаём UserService, который уже имеет доступ к БД через репозиторий
	authHandler := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Получаем токен из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
				return
			}

			tokenString := authHeader
			if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
				tokenString = tokenString[7:]
			}

			// Проверяем и парсим JWT токен
			userID, login, err := auth.GetUser(r.Context(), tokenString)
			if err != nil {
				http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
				return
			}

			// Если userID равен 0, пытаемся получить по логину
			if userID == 0 && login != "" {
				user, err := s.UserService.GetUserByLogin(r.Context(), login)
				if err != nil {
					http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
					return
				}
				userID = user.ID
			}

			// Сохраняем userID в контекст
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}

	// Единый хендлер для /api/user/orders — диспатчим по HTTP-методу
	s.HandleFunc("/api/user/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			authHandler(s.UploadOrderHandler).ServeHTTP(w, r)
		case http.MethodGet:
			authHandler(s.GetOrdersHandler).ServeHTTP(w, r)
		default:
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		}
	})

	s.HandleFunc("/api/user/balance", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}
		authHandler(s.GetBalanceHandler).ServeHTTP(w, r)
	})

	s.HandleFunc("/api/user/balance/withdraw", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}
		authHandler(s.WithdrawHandler).ServeHTTP(w, r)
	})

	s.HandleFunc("/api/user/withdrawals", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}
		authHandler(s.GetWithdrawalsHandler).ServeHTTP(w, r)
	})
}

 
 
 
 

 

 

 

 
 
 
 
 

 

 

 

 
