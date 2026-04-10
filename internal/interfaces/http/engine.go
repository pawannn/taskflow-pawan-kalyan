package engine

import (
	"fmt"
	"net/http"
)

type HttpEngine struct {
	mux *http.ServeMux
}

func NewHttpEngine() *HttpEngine {
	mux := http.NewServeMux()

	return &HttpEngine{
		mux: mux,
	}
}

func (e *HttpEngine) Start(port uint16) error {
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), e.mux)
	if err != nil {
		return err
	}

	return nil
}
