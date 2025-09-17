package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Auth   AuthConfig   `mapstructure:"auth"`
	GRPC   GRPCConfig   `mapstructure:"grpc"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type AuthConfig struct {
	BaseURL       string `mapstructure:"base_url"`
	ClientTimeout int    `mapstructure:"client_timeout"`
}

type GRPCConfig struct {
	CurrencyServiceURL string `mapstructure:"currency_service_url"`
}

func Load(path string) (Config, error) {
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
