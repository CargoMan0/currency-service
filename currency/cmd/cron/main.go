package main

import (
	"fmt"
	"github.com/BernsteinMondy/currency-service/currency/internal/config"
	"github.com/BernsteinMondy/currency-service/currency/internal/service"
	"github.com/BernsteinMondy/currency-service/currency/internal/worker"
	"github.com/BernsteinMondy/currency-service/pkg/database"
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
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		return fmt.Errorf("create new connection: %w", err)
	}

	return nil
}

func main() {

	repo, err := repository.NewCurrency(db)
	if err != nil {
		log.Fatalf("error creating repository: %v", err)
	}

	// Logger инициировать как можно раньше.
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("init logger: %w", err)
	}

	client, err := currency.New(cfg.API, logger)
	if err != nil {
		log.Fatalf("error creating currency client: %v", err)
	}

	svc := service.NewCurrency(repo, client, logger)

	c := cron.New()

	currencyWorker := worker.NewCurrency(cfg.Worker, svc, c, logger)

	if err != nil {
		log.Fatalf("error adding cron job: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := currencyWorker.StartFetchingCurrencyRates(); err != nil {
		log.Fatalf("error start fetching currency rates: %v", err)
	}

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	currencyWorker.Stop()
}
