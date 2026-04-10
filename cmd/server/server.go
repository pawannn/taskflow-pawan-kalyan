package main

import (
	"log"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/config"
	database "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/db"
	userRepository "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/repository/user_repository"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
	authhandler "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/handler/auth"
	authservice "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/auth"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := database.NewPostgresDB(cfg.DBUrl)

	engine := engine.NewHttpEngine(cfg)

	userRepo := userRepository.NewUserRepository(db)

	authService := authservice.NewAuthService(cfg, userRepo)

	authHandler := authhandler.NewAuthHandler(engine, authService)

	authHandler.AddUserRoutes()

	engine.Start()
}
