package main

import (
	"log"
	"stoke/internal/cfg"
	"stoke/internal/web"
)

func main() {
	log.Println("Starting Stoke Server...")

	config := cfg.Server {
		Address : "",
		Port: 8080,
	}

	server := web.Server {
		Config: config,
	}

	server.Init()
	if err := server.Run(); err != nil {
		log.Printf("An error stopped the server: %v", err)
	}
	
	log.Println("Stoke Server Terminated.")
}
