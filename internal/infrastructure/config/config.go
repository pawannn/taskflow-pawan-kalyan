// Package config handles application configuration loading and defaults.
package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config represents application configuration loaded from environment variables.
type Config struct {
	Env                 string `mapstructure:"ENV"`
	AppName             string `mapstructure:"APP_NAME"`
	AppPort             uint16 `mapstructure:"APP_PORT"`
	DBUrl               string `mapstructure:"DB_URL"`
	JWTSecret           string `mapstructure:"JWT_SECRET"`
	JWTExpiry           int    `mapstructure:"JWT_EXPIRY"`
	BCryptCost          int    `mapstructure:"BCRYPT_COST"`
	RateLimitIntervalMS int    `mapstructure:"RATE_LIMIT_INTERVAL_MS"`
	RateLimitBurst      int    `mapstructure:"RATE_LIMIT_BURST"`
}

// Load reads configuration from environment and .env file, applies defaults, and validates required fields.
func Load() (*Config, error) {
	viper.SetConfigType("env")

	// Defaults
	viper.SetDefault("ENV", "PROD")
	viper.SetDefault("APP_NAME", "taskflow")
	viper.SetDefault("APP_PORT", 1337)

	viper.SetDefault("JWT_EXPIRY", 24)
	viper.SetDefault("BCRYPT_COST", 12)

	// General rate limit defaults
	viper.SetDefault("RATE_LIMIT_INTERVAL_MS", 100) // ~10 RPS
	viper.SetDefault("RATE_LIMIT_BURST", 20)

	// Auth rate limit defaults (strict)
	viper.SetDefault("LOGIN_RATE_LIMIT_INTERVAL_MS", 500) // ~2 RPS
	viper.SetDefault("LOGIN_RATE_LIMIT_BURST", 5)

	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	cfg := new(Config)
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	if cfg.JWTExpiry < 24 {
		cfg.JWTExpiry = 24
	}

	if cfg.BCryptCost < 12 {
		cfg.BCryptCost = 12
	}

	if cfg.RateLimitIntervalMS <= 0 {
		cfg.RateLimitIntervalMS = 100
	}
	if cfg.RateLimitBurst <= 0 {
		cfg.RateLimitBurst = 20
	}

	if cfg.DBUrl == "" {
		return nil, fmt.Errorf("DB_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	log.Println("Config loaded successfully")
	return cfg, nil
}
