package main

import (
	"log"

	"gitlab.com/g6834/team28/auth/config"
	"gitlab.com/g6834/team28/auth/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig("./config/config.yml")
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
