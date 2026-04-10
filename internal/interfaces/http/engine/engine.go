package engine

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/config"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/logger"
)

type HttpEngine struct {
	mux *http.ServeMux
	cfg *config.Config
	Log *logger.Logger
}

func NewHttpEngine(cfg *config.Config) *HttpEngine {
	mux := http.NewServeMux()

	return &HttpEngine{
		mux: mux,
		cfg: cfg,
	}
}

func (e *HttpEngine) Start() error {
	port := fmt.Sprintf(":%d", e.cfg.AppPort)

	log.Printf("server started listening on port : %s \n", port)
	err := http.ListenAndServe(port, e.mux)
	if err != nil {
		return err
	}

	return nil
}
