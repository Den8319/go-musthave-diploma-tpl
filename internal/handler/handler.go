package handler

import (
	"context"
	"net/http"
	"time"

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
 
	s.HandleFunc("POST /api/user/register", s.RegisterHandler)
	s.HandleFunc("POST /api/user/login", s.LoginHandler)

	authMiddleware := auth.AuthMiddlewareWithRepo(s.UserService.Repository())


	s.Handle("POST /api/user/orders", authMiddleware(http.HandlerFunc(s.UploadOrderHandler)))
	s.Handle("GET /api/user/orders", authMiddleware(http.HandlerFunc(s.GetOrdersHandler)))

	s.Handle("GET /api/user/balance", authMiddleware(http.HandlerFunc(s.GetBalanceHandler)))
	s.Handle("POST /api/user/balance/withdraw", authMiddleware(http.HandlerFunc(s.WithdrawHandler)))
	s.Handle("GET /api/user/withdrawals", authMiddleware(http.HandlerFunc(s.GetWithdrawalsHandler)))
}