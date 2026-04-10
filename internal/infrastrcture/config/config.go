package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort   uint16 `mapstructure:"APP_PORT"`
	DBUrl     string `mapstructure:"DB_URL"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Read env file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Unable to load `.env` file: %s", err.Error())
	}

	// Automatically read env vars
	viper.AutomaticEnv()

	// sometimes viper doesn't map env vars automatically without binding.
	viper.BindEnv("APP_PORT")
	viper.BindEnv("DB_URL")
	viper.BindEnv("JWT_SECRET")

	cfg := new(Config)
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("Unable to decode config: %s", err.Error())
	}

	if cfg.AppPort == 0 {
		cfg.AppPort = 1337
	}

	if cfg.DBUrl == "" {
		return nil, fmt.Errorf("`DB_URL` is required in `.env`")
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("`JWT_SECRET` is required in `.env`")
	}

	return cfg, nil
}
