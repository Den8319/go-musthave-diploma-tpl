package handler

import (
	"net/http"
	
	"github.com/Den8319/go-musthave-diploma-tpl/internal/auth"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/repository"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/service"
)

type Server struct {
	*http.ServeMux
	UserService    *service.UserService
	OrderService   *service.OrderService
	BalanceService *service.BalanceService
}

func New(db *repository.DB) (*Server, error) {
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	balanceRepo := repository.NewBalanceRepository(db)

	userService := service.NewUserService(userRepo)
	orderService := service.NewOrderService(orderRepo, userRepo)
	balanceService := service.NewBalanceService(balanceRepo, userRepo)

	mux := &Server{
		ServeMux:       http.NewServeMux(),
		UserService:    userService,
		OrderService:   orderService,
		BalanceService: balanceService,
	}

	mux.registerHandlers()

	return mux, nil
}

func (s *Server) registerHandlers() {
	s.HandleFunc("/api/user/register", s.RegisterHandler)
	s.HandleFunc("/api/user/login", s.LoginHandler)

	authMiddleware := auth.AuthMiddlewareWithRepo(repository.NewUserRepository(nil))

	s.HandleFunc("/api/user/orders", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authMiddleware(http.HandlerFunc(s.UploadOrderHandler)).ServeHTTP(w, r)
	}))
	s.HandleFunc("/api/user/orders", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authMiddleware(http.HandlerFunc(s.GetOrdersHandler)).ServeHTTP(w, r)
	}))

	s.HandleFunc("/api/user/balance", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authMiddleware(http.HandlerFunc(s.GetBalanceHandler)).ServeHTTP(w, r)
	}))

	s.HandleFunc("/api/user/balance/withdraw", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authMiddleware(http.HandlerFunc(s.WithdrawHandler)).ServeHTTP(w, r)
	}))

	s.HandleFunc("/api/user/withdrawals", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authMiddleware(http.HandlerFunc(s.GetWithdrawalsHandler)).ServeHTTP(w, r)
	}))
}

 
 
 
 

 

 

 

 