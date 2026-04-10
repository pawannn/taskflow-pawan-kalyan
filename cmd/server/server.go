package main

import (
	"log"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastrcture/config"
	engine "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/service/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	engine := engine.NewHttpEngine(cfg)

	engine.Start()

}
