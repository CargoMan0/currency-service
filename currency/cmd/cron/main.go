package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/CargoMan0/currency-service/currency/internal/clients/currency"
	"github.com/CargoMan0/currency-service/currency/internal/config"
	"github.com/CargoMan0/currency-service/currency/internal/repository"
	"github.com/CargoMan0/currency-service/currency/internal/service"
	"github.com/CargoMan0/currency-service/currency/internal/worker"
	"github.com/CargoMan0/currency-service/pkg/database"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	err := run()
	if err != nil {
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
		log.Fatalf("init logger: %w", err)
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// New database connection
	db, err := openDatabaseConnection(cfg.Database)
	if err != nil {
		return fmt.Errorf("create new connection: %w", err)
	}

	// New repository
	repo := repository.New(db)

	// New currency client
	client, err := currency.NewClient(cfg.CurrencyAPI, logger)
	if err != nil {
		log.Fatalf("error creating currency client: %v", err)
	}

	// New service
	svc := service.NewCurrency(repo, client, logger)

	// New cron job and start worker
	c := cron.New()

	currencyWorker := worker.NewCurrency(cfg.Worker, svc, c, logger)
	err = currencyWorker.StartFetchingCurrencyRates()
	if err != nil {
		log.Fatalf("error start fetching currency rates: %v", err)
	}

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	currencyWorker.Stop()
	return nil
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
