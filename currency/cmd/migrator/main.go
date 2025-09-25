package main

import (
	"flag"
	"fmt"
	"github.com/BernsteinMondy/currency-service/currency/internal/config"
	"github.com/BernsteinMondy/currency-service/currency/internal/migrations"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
	"log"
)

type appConfig struct {
	Database config.DatabaseConfig `mapstructure:"database"`
}

func main() {
	err := run()
	if err != nil {
		log.Fatalf("run() returned error: %v", err)
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	err = migrations.RunPgMigrations(cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func loadConfig() (appConfig, error) {
	path := getPath()

	var cfg appConfig
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return cfg, fmt.Errorf("error reading config file: %w", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	return cfg, nil
}

func getPath() string {
	configPath := flag.String("config", "./config", "path to the config file")
	flag.Parse()

	return *configPath
}
