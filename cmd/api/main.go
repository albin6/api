package main

import (
	"github.com/albin6/api/config"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	server := NewServer(cfg)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}