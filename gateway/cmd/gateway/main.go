package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/BernsteinMondy/currency-service/gateway/internal/clients/auth"
	"github.com/BernsteinMondy/currency-service/gateway/internal/clients/currency"
	"github.com/BernsteinMondy/currency-service/gateway/internal/config"
	"github.com/BernsteinMondy/currency-service/gateway/internal/handler"
	"github.com/BernsteinMondy/currency-service/gateway/internal/middleware"
	errors2 "github.com/BernsteinMondy/currency-service/gateway/internal/repository"
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
		log.Fatalf("run() returned error: %v", err)
	}
}

func run() (err error) {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// Init logger
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	logger.Info("Loading config...")
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	logger.Info("Config loaded")

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	logger.Info("Server initializing with config", zap.Any("config", cfg))

	// Init router
	router := gin.New()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	// New auth client
	authClient, err := auth.NewClient(cfg.Auth)
	if err != nil {
		return fmt.Errorf("create new auth client: %w", err)
	}
	defer func() {
		logger.Info("Closing idle connections with AuthClient...")
		authClient.CloseIdleConnections()
		logger.Info("Idle connections with AuthClient closed")
	}()

	// Check if auth client is alive
	resp, err := authClient.Ping()
	if err != nil {
		return fmt.Errorf("check if auth service is alive: %w", err)
	}

	if resp != "pong" {
		return fmt.Errorf("auth client answered with invalid response: %w", err)
	}

	// New auth middleware
	authMiddleware := middleware.NewAuthorization(authClient, logger, shouldSkipAuthMiddleware)
	router.Use(authMiddleware.Authorize())

	// New grpc client for Currency.
	currencyGRPCClient, conn, err := grpc_client.NewCurrencyGRPCClient(cfg.GRPC.CurrencyServiceURL)
	if err != nil {
		return fmt.Errorf("create new grpc client %w", err)
	}
	defer func() {
		logger.Info("Closing gRPC connection with Currency...")
		closeErr := conn.Close()
		if closeErr != nil {
			logger.Warn("Failed to close gRPC connection with Currency", zap.Error(closeErr))
		} else {
			logger.Info("gRPC connection with Currency closed")
		}
	}()

	currencyClient := currency.NewClient(currencyGRPCClient)

	// Repository
	userRepo := errors2.NewUserRepository()

	// Service
	authService := service.NewAuthService(userRepo, authClient)
	currencyService := service.NewCurrencyService(currencyClient)

	// Prepare test user
	err = prepareTestUser(ctx, userRepo, cfg.TestUserCredentials)
	if err != nil {
		return fmt.Errorf("prepare test user: %w", err)
	}

	// HTTP server
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

func prepareTestUser(ctx context.Context, repo *errors2.UserRepository, cfg config.TestUserCredentials) error {
	user := service.User{
		Login:    cfg.Login,
		Password: cfg.Password,
	}

	err := repo.SaveUser(ctx, user)
	if err != nil {
		return fmt.Errorf("repository: save user: %w", err)
	}

	return nil
}
