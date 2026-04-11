package engine

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/config"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/logger"
)

type HttpEngine struct {
	router *chi.Mux
	cfg    *config.Config
	Log    *logger.Logger
}

func NewHttpEngine(cfg *config.Config, logger *logger.Logger) *HttpEngine {

	return &HttpEngine{
		router: chi.NewRouter(),
		cfg:    cfg,
		Log:    logger,
	}
}

func (e *HttpEngine) Start() error {
	port := fmt.Sprintf(":%d", e.cfg.AppPort)

	server := &http.Server{
		Addr:    port,
		Handler: e.router,
	}

	go func() {
		fmt.Println()
		log.Printf("server started on port :%s\n", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("server failed to start:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	log.Println("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("graceful shutdown failed", err)
		return err
	}

	log.Println("server exited properly")
	return nil
}
