package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"codenames-game/configs"
	"codenames-game/internal/infrastructure/persistence"
	"codenames-game/internal/interfaces/api"
	chatService "codenames-game/internal/usecase/chat"
	gameService "codenames-game/internal/usecase/game"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// spaHandler implements the http.Handler interface for serving a Single Page Application
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP serves the React app, routing all non-API requests to index.html
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the absolute path to prevent directory traversal
	path := filepath.Join(h.staticPath, r.URL.Path)

	// Check if file exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// File doesn't exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	}

	if err != nil {
		// Some other error occurred
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// File exists, serve it
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	// Load environment variables
	godotenv.Load()

	// Load configuration
	config := configs.LoadConfig()

	// Initialize repositories
	gameRepo := persistence.NewGameRepository()

	// First, create the WebSocket handler
	wsHandler := api.NewWebSocketHandler()

	// Initialize chat service
	chatSvc := chatService.NewChatService(nil)
	chatHandler := api.NewChatHandler(chatSvc)

	// Create game service with WebSocket handler
	gameSvc := gameService.NewServiceWithWebSocket(gameRepo, wsHandler)

	// IMPORTANT: Create game handler just once with the WebSocket-enabled service
	gameHandler := api.NewGameHandler(gameSvc)

	// Setup router
	router := mux.NewRouter()

	// API routes need to be registered BEFORE the SPA handler
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Game routes
	apiRouter.HandleFunc("/game/start", gameHandler.StartGame).Methods("POST")
	apiRouter.HandleFunc("/game/join", gameHandler.JoinGame).Methods("POST")
	apiRouter.HandleFunc("/game/state", gameHandler.GetGameState).Methods("GET")
	apiRouter.HandleFunc("/game/reveal", gameHandler.RevealCard).Methods("POST")
	apiRouter.HandleFunc("/game/set-spymaster", gameHandler.SetSpymaster).Methods("POST")
	apiRouter.HandleFunc("/game/end-turn", gameHandler.EndTurn).Methods("POST")
	// Add the missing change-team route
	apiRouter.HandleFunc("/game/change-team", gameHandler.ChangeTeam).Methods("POST")

	// Chat routes
	apiRouter.HandleFunc("/chat/send", chatHandler.SendMessage).Methods("POST")
	apiRouter.HandleFunc("/chat/messages", chatHandler.GetMessages).Methods("GET")

	// Register WebSocket routes directly on the main router (not under /api)
	wsHandler.RegisterRoutes(router)

	// CORS middleware with more permissive settings for development
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept", "Origin", "X-Requested-With"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300,   // Maximum value not ignored by any of major browsers
		Debug:            false, // Set to false to reduce log noise
	})

	// SPA handler should come AFTER the API routes
	spa := spaHandler{staticPath: "frontend/build", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	// Add middleware to log all incoming requests
	loggedRouter := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		router.ServeHTTP(w, r)
	})

	// Apply CORS middleware to both API routes and WebSocket routes
	// Use a handler chain: CORS -> Logging -> Router
	corsHandler := corsMiddleware.Handler(loggedRouter)

	// Setup HTTP server with proper timeouts from config
	server := &http.Server{
		Addr:         ":" + config.Server.Port, // Use Server.Port from the nested structure
		Handler:      corsHandler,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server
	fmt.Printf("Starting Codenames Game server on :%s\n", config.Server.Port)
	log.Fatal(server.ListenAndServe())
}
