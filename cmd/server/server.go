package main

import (
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

	logger := logger.New()

	engine := engine.NewHttpEngine(cfg, logger)

	// Init DB
	db, err := database.NewPostgresDB(cfg.DBUrl)
	if err != nil {
		log.Fatal(err.Error())
	}

	tokenService := auth.NewTokenService(cfg.AppName, cfg.JWTSecret, cfg.JWTExpiry)
	middlewares := middlewares.NewMiddlewareHadler(engine, *tokenService)

	// init repositories
	userRepository := userRepository.NewUserRepository(db)
	projectRepository := projectRepository.NewProjectRepository(db)
	taskRepository := taskRepository.NewTaskRepository(db)

	// init services
	authService := authservice.NewAuthService(cfg.BCryptCost, userRepository, tokenService)
	projectService := projectService.NewProjectRepository(projectRepository, taskRepository)
	taskService := taskService.NewTaskService(taskRepository, projectRepository, userRepository)

	// init handlers
	authHandler := authHandler.NewAuthHandler(engine, authService)
	projectHandler := projectHandler.NewProjectHandler(engine, projectService, middlewares)
	taskHandler := taskHandler.NewTaskHandler(engine, middlewares, taskService)

	// Add routes
	authHandler.AddRoutes()
	projectHandler.AddRoutes()
	taskHandler.AddRoutes()

	engine.Start()
}
