package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"codenames-game/configs"
	"codenames-game/internal/infrastructure/persistence"
	"codenames-game/internal/interfaces/api"
	chatService "codenames-game/internal/usecase/chat"
	gameService "codenames-game/internal/usecase/game"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Load configuration
	config := configs.LoadConfig()

	// Initialize repositories
	gameRepo := persistence.NewGameRepository()
	chatRepo := persistence.NewChatRepository()

	// Initialize services
	gameSvc := gameService.NewServiceWithRepo(gameRepo)
	chatSvc := chatService.NewChatService(chatRepo)

	// Initialize handlers
	gameHandler := api.NewGameHandler(gameSvc)
	chatHandler := api.NewChatHandler(chatSvc)

	// Setup router
	router := mux.NewRouter()

	// Add a welcome page
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
        <html>
            <head>
                <title>Codenames Game API</title>
                <style>
                    body { font-family: Arial, sans-serif; line-height: 1.6; margin: 0; padding: 20px; max-width: 800px; margin: 0 auto; }
                    h1 { color: #333; }
                    h2 { color: #555; }
                    code { background: #f4f4f4; padding: 5px; border-radius: 3px; }
                    pre { background: #f4f4f4; padding: 10px; border-radius: 5px; overflow-x: auto; }
                    .endpoint { margin-bottom: 20px; border-bottom: 1px solid #eee; padding-bottom: 10px; }
                </style>
            </head>
            <body>
                <h1>Codenames Game API</h1>
                <p>Welcome to the Codenames Game API! Below are the available endpoints:</p>
                
                <h2>Game Endpoints</h2>
                
                <div class="endpoint">
                    <h3>Create Game</h3>
                    <code>POST /api/game/start</code>
                    <pre>
{
  "creator_id": "user-123",
  "username": "Player1"
}
                    </pre>
                </div>
                
                <div class="endpoint">
                    <h3>Join Game</h3>
                    <code>POST /api/game/join</code>
                    <pre>
{
  "game_id": "game-123",
  "player_id": "user-456",
  "username": "Player2",
  "team": "red" // or "blue"
}
                    </pre>
                </div>
                
                <div class="endpoint">
                    <h3>Get Game State</h3>
                    <code>GET /api/game/state?id=game-123</code>
                </div>
                
                <div class="endpoint">
                    <h3>Set Spymaster</h3>
                    <code>POST /api/game/set-spymaster?game_id=game-123&player_id=user-123</code>
                </div>
                
                <div class="endpoint">
                    <h3>Reveal Card</h3>
                    <code>POST /api/game/reveal</code>
                    <pre>
{
  "game_id": "game-123",
  "card_id": "card-456",
  "player_id": "user-123"
}
                    </pre>
                </div>
                
                <div class="endpoint">
                    <h3>End Turn</h3>
                    <code>POST /api/game/end-turn?game_id=game-123&player_id=user-123</code>
                </div>
                
                <h2>Chat Endpoints</h2>
                
                <div class="endpoint">
                    <h3>Send Message</h3>
                    <code>POST /api/chat/send</code>
                    <pre>
{
  "content": "Hello everyone",
  "sender_id": "user-123",
  "username": "Player1",
  "chat_id": "game-123"
}
                    </pre>
                </div>
                
                <div class="endpoint">
                    <h3>Get Messages</h3>
                    <code>GET /api/chat/messages?chat_id=game-123</code>
                </div>
                
                <p>Use these endpoints to interact with the Codenames Game API.</p>
                <p>For a simple test client, open the test-client.html file in your browser.</p>
            </body>
        </html>
        `)
	})

	// Register API routes
	gameHandler.RegisterRoutes(router)
	chatHandler.RegisterRoutes(router)

	// Add middleware for CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins for testing
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Setup HTTP server
	serverAddr := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      corsMiddleware.Handler(router),
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
	}

	// Start the server
	log.Printf("Starting Codenames Game server on %s", serverAddr)
	log.Printf("Visit http://%s to see the API documentation", serverAddr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
