package main

import (
	"fmt"
	"github.com/BernsteinMondy/currency-service/currency/internal/config"
	"github.com/BernsteinMondy/currency-service/currency/internal/migrations"
	"log"
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
		return fmt.Errorf("failed to load config: %w", err)
	}

	err = migrations.RunPgMigrations(cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
