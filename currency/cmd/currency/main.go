package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	currencyapi "github.com/BernsteinMondy/currency-service/currency/internal/clients/currency"
	"github.com/BernsteinMondy/currency-service/currency/internal/config"
	"github.com/BernsteinMondy/currency-service/currency/internal/handler"
	"github.com/BernsteinMondy/currency-service/currency/internal/metrics"
	"github.com/BernsteinMondy/currency-service/currency/internal/middleware"
	"github.com/BernsteinMondy/currency-service/currency/internal/repository"
	"github.com/BernsteinMondy/currency-service/currency/internal/service"
	"github.com/BernsteinMondy/currency-service/pkg/currency"
	"github.com/BernsteinMondy/currency-service/pkg/database"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("run() returned error: %v", err)
	}
}

func run() (err error) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	// Init logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("error initializing logger: %w", err)
	}

	// Load config
	logger.Info("Loading config...")
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}
	logger.Info("Config loaded")

	// New database connection
	logger.Info("Opening database connection...")
	db, err := openDatabaseConnection(cfg.Database)
	if err != nil {
		return fmt.Errorf("error creating db connection: %w", err)
	}
	logger.Info("Database connection opened")

	// Close database connection before exiting run()
	defer func() {
		logger.Info("Closing database connection")
		closeErr := db.Close()
		if closeErr != nil {
			logger.Warn("Error closing database connection", zap.Error(closeErr))
		}

		logger.Info("Database connection closed")
	}()

	// New repository
	repo := repository.New(db)

	// Currency API client
	currencyClient, err := currencyapi.NewClient(cfg.CurrencyAPI, logger)
	if err != nil {
		return fmt.Errorf("error creating currency client: %w", err)
	}

	// Currency Service
	currencyService := service.NewCurrency(repo, currencyClient, logger)

	// Currency Server
	currencyServer := handler.NewCurrencyServer(
		currencyService,
		logger,
	)

	// Init metrics
	mtrcs, err := metrics.Init()
	if err != nil {
		return fmt.Errorf("error initializing metrics: %w", err)
	}

	// Init middleware
	mw := middleware.Init(mtrcs, logger)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		logger.Info("Launching gRPC server...")
		err = launchGRPCServer(ctx, cfg.Service, currencyServer, mw)
		if err != nil {
			logger.Error("error starting gRPC server", zap.Error(err))
		}

		logger.Info("gRPC server shut down gracefully")
	}(ctx)

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		logger.Info("Launching Prometheus Metrics server...")
		err = launchPrometheusServer(ctx)
		if err != nil {
			logger.Error("error starting Prometheus metrics server", zap.Error(err))
		}

		logger.Info("Prometheus Metrics shut down gracefully")
	}(ctx)

	<-ctx.Done()
	wg.Wait()
	return nil
}

func launchPrometheusServer(ctx context.Context) error {
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Prometheus metrics server running on :8081")

	serveErr := make(chan error, 1)

	srv := http.Server{
		Addr:    ":8081",
		Handler: nil,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			serveErr <- fmt.Errorf("error serving Prometheus metrics server: %w", err)
			close(serveErr)
		}
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
		_ = srv.Shutdown(ctx)
		return nil
	}
}

func launchGRPCServer(ctx context.Context, cfg config.ServiceConfig, srv currency.CurrencyServiceServer, mw *middleware.Middleware) error {
	listener, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			mw.LoggingMiddleware,
			mw.RequestDurationMiddleware,
		),
	)

	currency.RegisterCurrencyServiceServer(s, srv)

	log.Printf("gRPC server is listening on :%s", cfg.ServerPort)

	serveErr := make(chan error, 1)

	go func() {
		err = s.Serve(listener)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			serveErr <- fmt.Errorf("failed to serve: %w", err)
			close(serveErr)
		}
	}()

	select {
	case err = <-serveErr:
		return err
	case <-ctx.Done():
		s.GracefulStop()
		return nil
	}
}

func openDatabaseConnection(cfg config.DatabaseConfig) (*sql.DB, error) {
	c := &database.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		DBName:   cfg.Name,
		SSLMode:  cfg.SSLMode,
	}

	db, err := database.NewConnection(c)
	if err != nil {
		return nil, err
	}

	return db, nil
}
