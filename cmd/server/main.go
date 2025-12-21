package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adexcell/go-tutorial/internal/config"
	"github.com/adexcell/go-tutorial/internal/handler"
	"github.com/adexcell/go-tutorial/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.NewLogger()
	cfg := config.Load()

	r := gin.Default()
	r.GET("/health", handler.Health)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Info().Msgf("Server starting on port %s", cfg.Port)
		srv.ListenAndServe()
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msgf("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}
	log.Info().Msg("Server shutdown complete")

}
