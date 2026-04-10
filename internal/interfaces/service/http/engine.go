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
		cfg: cfg,
	}
}

func (e *HttpEngine) Start() error {
	port := fmt.Sprintf(":%d", e.cfg.AppPort)

	fmt.Println("server started listening on port : ", port)
	err := http.ListenAndServe(port, e.mux)
	if err != nil {
		return err
	}

	return nil
}
