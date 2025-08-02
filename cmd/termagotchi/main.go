package main

import (
	"log"

	"github.com/ezeoleaf/termagotchi/internal/app"
	"github.com/ezeoleaf/termagotchi/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	a := app.NewApp(cfg)
	a.Run()
}
