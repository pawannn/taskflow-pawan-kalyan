package main

import (
	"context"
	"log"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastrcture/config"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastrcture/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := db.NewPostgresDB(cfg.DBUrl)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal(err.Error())
	}
}
