package main

import (
	"log"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/config"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	engine := engine.NewHttpEngine(cfg)

	engine.Start()
}
