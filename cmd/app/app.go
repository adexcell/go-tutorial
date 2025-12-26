package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/adexcell/go-tutorial/internal/config"
	"github.com/adexcell/go-tutorial/internal/repository/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type App struct {
	cfg     *config.Config
	logger  zerolog.Logger
	storage *pgxpool.Pool
}

func New(cfg *config.Config, logger zerolog.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	storage, err := postgres.New(ctx, a.cfg.Postgres)
	if err != nil {
		return fmt.Errorf("не удалось запустить базу данных: %v", err)
	}

	a.storage = storage

	router := gin.Default()
	srv := &http.Server{
		Addr:           a.cfg.HTTPServer.Addr,
		Handler:        router,
		ReadTimeout:    a.cfg.HTTPServer.ReadTimeout,
		WriteTimeout:   a.cfg.HTTPServer.WriteTimeout,
		IdleTimeout:    a.cfg.HTTPServer.IdleTimeout,
		MaxHeaderBytes: a.cfg.HTTPServer.MaxHeaderBytes,
	}

	go func() {
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			a.logger.Info().Msg("Server stopped gracefully")
		} else if err != nil {
			a.logger.Err(err).Msg("непредвиденная ошибка")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	ctxStop, cancel := context.WithTimeout(context.Background(), a.cfg.HTTPServer.ShutdownTimeout)
	defer cancel()
	srv.Shutdown(ctxStop)
	a.storage.Close()
	a.logger.Info().Msg("Server down")

	return nil
}
