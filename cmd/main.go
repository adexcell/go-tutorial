package main

import (
	"context"
	"log"

	"github.com/adexcell/go-tutorial/cmd/app"
	"github.com/adexcell/go-tutorial/internal/config"
	"github.com/adexcell/go-tutorial/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("не удалось загрузить конфиг: %v", err)
	}

	l := logger.New(cfg.Logger.Level, cfg.Logger.JSONFormat)

	myapp := app.New(cfg, l)
	myapp.Run(ctx)
}
