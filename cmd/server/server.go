package main

import (
	"log"

	auth "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/auth/jwt"
	config "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/config"
	database "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/database"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/database/userRepository"
	logger "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/logger"
	engine "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
	authHandler "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/handler/auth"
	authservice "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/auth"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	logger := logger.New()

	engine := engine.NewHttpEngine(cfg, logger)

	// Init DB
	db, err := database.NewPostgresDB(cfg.DBUrl)
	if err != nil {
		log.Fatal(err.Error())
	}

	tokenService := auth.NewTokenService(cfg.AppName, cfg.JWTSecret, cfg.JWTExpiry)

	// init user repository
	userService := userRepository.NewUserRepository(db)

	// init auth service
	authService := authservice.NewAuthService(cfg.BCryptCost, userService, tokenService)

	// init auth handler
	authHandler := authHandler.NewAuthHandler(engine, authService)

	// Add Auth routes
	authHandler.AddAuthRoutes()

	engine.Start()
}
