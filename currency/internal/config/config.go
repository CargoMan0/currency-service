package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Service     ServiceConfig     `mapstructure:"service"`
	CurrencyAPI CurrencyAPIConfig `mapstructure:"api"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Worker      WorkerConfig      `mapstructure:"worker"`
}

type ServiceConfig struct {
	ServerPort string `mapstructure:"server_port"`
}

type CurrencyAPIConfig struct {
	BaseURL        string `mapstructure:"base_url"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
}

type DatabaseConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Name           string `mapstructure:"name"`
	SSLMode        string `mapstructure:"ssl_mode"`
	MigrationsPath string `mapstructure:"migrations_path"`
}

func (cfg DatabaseConfig) ToDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode,
	)
}

type WorkerConfig struct {
	Schedule      string `mapstructure:"schedule"`
	TimoutSeconds int    `mapstructure:"timeout_seconds"`
	CurrencyPair  struct {
		BaseCurrency   string `mapstructure:"base_currency"`
		TargetCurrency string `mapstructure:"target_currency"`
	} `mapstructure:"currency_pair"`
}

func Load() (Config, error) {
	path := getPath()

	var cfg Config
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	return cfg, nil
}

func getPath() string {
	configPath := flag.String("config", "./config", "path to the config file")
	flag.Parse()

	return *configPath
}
