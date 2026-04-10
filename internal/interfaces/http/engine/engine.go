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

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/config"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/logger"
)

type HttpEngine struct {
	mux *http.ServeMux
	cfg *config.Config
	Log *logger.Logger
}

func NewHttpEngine(cfg *config.Config, logger *logger.Logger) *HttpEngine {
	mux := http.NewServeMux()

	return &HttpEngine{
		mux: mux,
		cfg: cfg,
		Log: logger,
	}
}

func (e *HttpEngine) Start() error {
	port := fmt.Sprintf(":%d", e.cfg.AppPort)

	server := &http.Server{
		Addr:    port,
		Handler: e.mux,
	}

	go func() {
		fmt.Println()
		log.Printf("server started on port :%s\n", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("server failed to start:", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

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
