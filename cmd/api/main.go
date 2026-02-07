package main

import (
	"github.com/albin6/api/config"
	"log"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Initialize Server (DI Wiring)
	server := NewServer(cfg)

	// 3. Run Server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}