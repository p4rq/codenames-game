package main

import (
	"codenames-game/configs"
	api "codenames-game/internal/interfaces"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// Load configuration
	config := configs.LoadConfig()

	// Set up logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Codenames Game Server...")

	// Initialize router and handlers
	router := http.NewServeMux()

	// Register static file server for web assets
	fs := http.FileServer(http.Dir("./web/static"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register API handlers
	api.RegisterHandlers(router)

	// Configure and start HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
	}

	log.Printf("Server listening on %s:%s", config.Server.Host, config.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
