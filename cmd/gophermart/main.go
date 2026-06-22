package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	

	"github.com/rs/zerolog/log"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/auth"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/compress"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/config"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/handler"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/logger"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/repository"
)

func main() {
	
	auth.Init("SECRET_KEY")

	
	cfg := config.New()

	
	logger.InitLogger(cfg.LogLevel)

	var db *repository.DB
	var err error

	if cfg.DatabaseDSN != "" {
		db, err = repository.New(cfg.DatabaseDSN)
		if err != nil {
			log.Fatal().Err(err).Msg("Ошибка подключения к базе данных")
		}
		defer db.Close()
	} else {
		log.Info().Msg("DATABASE_DSN не указан, работа с БД невозможна")
		return
	}

	if err := db.Migrate("migrations"); err != nil {
		log.Fatal().Err(err).Msg("Ошибка применения миграций")
	}

	server, err := handler.New(db)
	if err != nil {
		log.Fatal().Err(err).Msg("incorrect server")
		return
	}

	
	userHandler := &handler.User{Server: server}
	ordersHandler := &handler.Orders{Server: server}

	
	userRepo := repository.NewUserRepository(db)

	
	router := chi.NewRouter()

	
	router.Use(logger.WithLogging)
	router.Use(compress.WithCompression)

	
	router.Post("/api/user/register", userHandler.RegisterHandler)
	router.Post("/api/user/login", userHandler.LoginHandler)

	
	authMiddleware := auth.AuthMiddlewareWithRepo(userRepo)

	
	router.With(authMiddleware).Post("/api/user/orders", ordersHandler.UploadOrderHandler)
	router.With(authMiddleware).Get("/api/user/orders", ordersHandler.GetOrdersHandler)

	
	addr := cfg.ServerAddress
	log.Info().Str("address", addr).Msg("Запуск сервера")
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal().Err(err).Msg("Ошибка при запуске сервера")
	}
}
