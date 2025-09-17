package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/BernsteinMondy/currency-service/gateway/internal/clients/auth"
	"github.com/BernsteinMondy/currency-service/gateway/internal/config"
	"github.com/BernsteinMondy/currency-service/gateway/internal/handler"
	"github.com/BernsteinMondy/currency-service/gateway/internal/middleware"
	"github.com/BernsteinMondy/currency-service/gateway/internal/repository"
	"github.com/BernsteinMondy/currency-service/gateway/internal/service"
	"github.com/BernsteinMondy/currency-service/pkg/grpc_client"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}

func run() (err error) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	logger.Info("Server initializing with config", zap.Any("config", cfg))

	router := gin.New()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	authClient, err := auth.NewClient(cfg.Auth)
	if err != nil {
		return fmt.Errorf("create new auth client: %w", err)
	}

	resp, err := authClient.Ping()
	if err != nil {
		return fmt.Errorf("check if auth service is alive: %w", err)
	}

	if resp != "pong" {
		return fmt.Errorf("auth client answered with invalid response: %w", err)
	}

	authMiddleware := middleware.NewAuthorization(authClient, logger, shouldSkipAuthMiddleware)
	router.Use(authMiddleware.Authorize())

	currencyClient, conn, err := grpc_client.NewCurrencyGRPCClient(cfg.GRPC.CurrencyServiceURL)
	if err != nil {
		return fmt.Errorf("create new grpc client %w", err)
	}
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close connection: %w", closeErr))
			logger.Warn("Cannot close GRPC Client for auth service", zap.Error(err))
		}
	}()

	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo, authClient)
	currencyService := service.NewCurrencyService(currencyClient)

	err = prepareTestUser(ctx, userRepo, cfg.TestUserCredentials)
	if err != nil {
		return fmt.Errorf("prepare test user: %w", err)
	}

	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}

	handler.RegisterRoutes(authService, currencyService, router, logger)

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	log.Println("Server shut down")
	return nil
}

func shouldSkipAuthMiddleware(c *gin.Context) bool {
	if strings.HasSuffix(c.Request.URL.Path, "/login") ||
		strings.HasSuffix(c.Request.URL.Path, "/register") {
		return true
	}

	return false
}

func prepareTestUser(ctx context.Context, repo *repository.UserRepository, cfg config.TestUserCredentials) error {
	user := repository.User{
		Login:    cfg.Login,
		Password: cfg.Password,
	}

	err := repo.SaveUser(ctx, user)
	if err != nil {
		return fmt.Errorf("repository: save user: %w", err)
	}

	return nil
}
