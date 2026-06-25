package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/Den8319/go-musthave-diploma-tpl/internal/auth"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/compress"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/config"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/handler"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/logger"
	"github.com/Den8319/go-musthave-diploma-tpl/internal/repository"
)

func main() {

	cfg := config.New()
	auth.Init(cfg.SecretKey)


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

	server, err := handler.New(db,cfg.AccrualSystemAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("incorrect server")
		return
	}

	
	wrappedHandler := logger.WithLogging(compress.WithCompression(server))
	 
	httpServer := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: wrappedHandler,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Info().Str("address", cfg.ServerAddress).Msg("Запуск сервера")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Ошибка запуска сервера")
		}
	}()
	
	sig := <-quit
	log.Info().Str("signal", sig.String()).Msg("Получен сигнал для завершения работы сервера")

	server.Shutdown()

	shutdownCtx, shutdowncancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdowncancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Ошибка при завершении работы сервера")
		}
	log.Info().Msg("Сервер успешно остановлен")
	
}
