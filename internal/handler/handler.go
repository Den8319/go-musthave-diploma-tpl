package handler

import (
	"github.com/Den8319/go-musthave-diploma-tpl/internal/repository"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/service"
)


type Server struct {
	UserService  *service.UserService
	OrderService *service.OrderService
}

func New(db *repository.DB) (*Server, error) {
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	userService := service.NewUserService(userRepo)
	orderService := service.NewOrderService(orderRepo, userRepo)

	server := &Server{
		UserService:  userService,
		OrderService: orderService,
	}

	return server, nil
}
