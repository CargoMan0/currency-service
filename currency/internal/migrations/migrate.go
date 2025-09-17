package migrations

import (
	"errors"
	"fmt"
	"github.com/BernsteinMondy/currency-service/currency/internal/config"
	"github.com/golang-migrate/migrate"
)

func RunPgMigrations(cfg config.DatabaseConfig) error {
	dsn := cfg.ToDSN()
	if dsn == "" {
		return errors.New("DSN must not be empty")
	}
	if cfg.MigrationsPath == "" {
		return errors.New("path must not be empty")
	}

	m, err := migrate.New(cfg.MigrationsPath, dsn)

	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run up migrations: %w", err)
	}

	return nil
}
