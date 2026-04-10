package main

import (
	"fmt"
	"log"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastrcture/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(cfg)
}
