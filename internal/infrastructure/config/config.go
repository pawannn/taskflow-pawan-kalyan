// Package config handles application configuration loading and defaults.
package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config represents application configuration loaded from environment variables.
type Config struct {
	Env        string `mapstructure:"ENV"`
	AppName    string `mapstructure:"APP_NAME"`
	AppPort    uint16 `mapstructure:"APP_PORT"`
	DBUrl      string `mapstructure:"DB_URL"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
	JWTExpiry  int    `mapstructure:"JWT_EXPIRY"`
	BCryptCost int    `mapstructure:"BCRYPT_COST"`
}

// Load reads configuration from environment and .env file, applies defaults, and validates required fields.
func Load() (*Config, error) {
	viper.SetConfigType("env")

	viper.SetDefault("ENV", "PROD")
	viper.SetDefault("APP_NAME", "")
	viper.SetDefault("APP_PORT", 0)
	viper.SetDefault("DB_URL", "")
	viper.SetDefault("JWT_SECRET", "")
	viper.SetDefault("JWT_EXPIRY", 0)
	viper.SetDefault("BCRYPT_COST", 0)

	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	cfg := new(Config)
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	// Defaults
	if cfg.AppPort == 0 {
		cfg.AppPort = 1337
	}
	if cfg.AppName == "" {
		cfg.AppName = "taskflow"
	}
	if cfg.Env == "" {
		cfg.Env = "PROD"
	}

	if cfg.JWTExpiry < 24 {
		cfg.JWTExpiry = 24
	}

	if cfg.BCryptCost < 12 {
		cfg.BCryptCost = 12
	}

	// Required fields
	if cfg.DBUrl == "" {
		return nil, fmt.Errorf("DB_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	log.Println("Config loaded successfully")
	return cfg, nil
}
