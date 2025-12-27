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
	"github.com/adexcell/go-tutorial/internal/handler"
	"github.com/adexcell/go-tutorial/internal/repository/postgres"
	"github.com/adexcell/go-tutorial/internal/repository/cache"
	"github.com/adexcell/go-tutorial/internal/service"
	"github.com/adexcell/go-tutorial/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type App struct {
	cfg     *config.Config
	logger  zerolog.Logger
	storage *pgxpool.Pool
	cache   *redis.Client
}

func New(cfg *config.Config, logger zerolog.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage, err := postgres.New(ctx, a.cfg.Postgres)
	if err != nil {
		return fmt.Errorf("не удалось запустить базу данных: %v", err)
	}
	a.storage = storage

	redisCache, err := cache.New(ctx, a.cfg.Redis)
	if err != nil {
		return fmt.Errorf("не удалось запустить redis: %v", err)
	}
	a.cache = redisCache

	userRepo := postgres.NewUserRepository(a.storage)
	userCache := cache.NewUserCache(a.cache)

	manager, err := auth.NewManager(a.cfg.Auth.JWTSecret)
	if err != nil {
		return fmt.Errorf("не удалось запустить manager: %w", err)
	}

	userService := service.NewUserService(
		userRepo,
		manager,
		userCache,
		a.cfg.Auth.TokenTTL,
		a.cfg.Redis.TTL,
	)
	userHandler := handler.NewUserHandler(userService)

	router := gin.New()
	// middleware
	router.Use(handler.Logger(a.logger))
	router.Use(gin.Recovery())

	// static
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")

	// routs
	router.POST("/auth/register", userHandler.Register)
	router.POST("/auth/login", userHandler.Login)

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
	a.logger.Info().Msg("Server down")

	a.cache.Close()
	a.logger.Info().Msg("Cache down")

	a.storage.Close()
	a.logger.Info().Msg("Storage down")

	return nil
}
