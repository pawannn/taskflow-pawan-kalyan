package main

import (
	"fmt"
	"log"

	auth "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/auth/jwt"
	config "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/config"
	database "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/database"
	projectRepository "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/database/projectRepository"
	taskRepository "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/database/taskRepository"
	userRepository "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/database/userRepository"
	logger "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/logger"
	engine "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
	authHandler "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/handler/auth"
	projectHandler "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/handler/project"
	taskHandler "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/handler/task"
	middlewares "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/middlewares"
	authservice "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/auth"
	projectService "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/project"
	taskService "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/task"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("%+v\n\n", cfg)

	logger := logger.New(cfg.Env)

	// Init DB
	db, err := database.NewPostgresDB(cfg.DBUrl)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Init HTTP Engine
	engine := engine.NewHttpEngine(cfg, logger)
	rateLimiter := middlewares.NewRateLimiter(cfg.RateLimitIntervalMS, cfg.RateLimitBurst)
	engine.Use(rateLimiter.Limit)

	tokenService := auth.NewTokenService(cfg.AppName, cfg.JWTSecret, cfg.JWTExpiry)
	middlwareHandler := middlewares.NewMiddlewareHandler(engine, *tokenService)

	// init repositories
	userRepository := userRepository.NewUserRepository(db)
	projectRepository := projectRepository.NewProjectRepository(db)
	taskRepository := taskRepository.NewTaskRepository(db)

	// init services
	authService := authservice.NewAuthService(cfg.BCryptCost, userRepository, tokenService)
	projectService := projectService.NewProjectService(projectRepository, taskRepository)
	taskService := taskService.NewTaskService(taskRepository, projectRepository, userRepository)

	// init handlers
	authHandler := authHandler.NewAuthHandler(engine, authService)
	projectHandler := projectHandler.NewProjectHandler(engine, projectService, middlwareHandler)
	taskHandler := taskHandler.NewTaskHandler(engine, middlwareHandler, taskService)

	// Add routes
	authHandler.AddRoutes()
	projectHandler.AddRoutes()
	taskHandler.AddRoutes()

	engine.Start()
}
