package migrations

import (
	"embed"
	"errors"
	"fmt"
	"github.com/BernsteinMondy/currency-service/currency/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var fs embed.FS

func RunPgMigrations(cfg config.DatabaseConfig) error {
	dsn := cfg.ToDSN()
	if dsn == "" {
		return errors.New("DSN must not be empty")
	}
	if cfg.MigrationsPath == "" {
		return errors.New("path must not be empty")
	}

	d, err := iofs.New(fs, cfg.MigrationsPath)
	if err != nil {
		return fmt.Errorf("iofs.New(): %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("could not create migrations instance: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run up migrations: %w", err)
	}

	return nil
}
