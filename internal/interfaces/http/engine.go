package engine

import (
	"fmt"
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastrcture/config"
)

type HttpEngine struct {
	mux *http.ServeMux
	cfg *config.Config
}

func NewHttpEngine(cfg *config.Config) *HttpEngine {
	mux := http.NewServeMux()

	return &HttpEngine{
		mux: mux,
	}
}

func (e *HttpEngine) Start() error {
	port := e.cfg.AppPort
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), e.mux)
	if err != nil {
		return err
	}

	return nil
}
