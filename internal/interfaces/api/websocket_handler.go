package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	gorillaWs "github.com/gorilla/websocket" // Alias for Gorilla's WebSocket package

	customWs "codenames-game/internal/infrastructure/websocket" // Alias for your custom WebSocket package
	wsinterfaces "codenames-game/internal/interfaces/websocket" // Import the interfaces
)

var upgrader = gorillaWs.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// WebSocketHandler handles WebSocket connections and implements UpdateBroadcaster
type WebSocketHandler struct {
	hub *customWs.Hub
}

// Verify WebSocketHandler implements the UpdateBroadcaster interface
var _ wsinterfaces.UpdateBroadcaster = (*WebSocketHandler)(nil)

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler() *WebSocketHandler {
	hub := customWs.NewHub()
	go hub.Run()

	return &WebSocketHandler{
		hub: hub,
	}
}

// RegisterRoutes registers the WebSocket routes
func (h *WebSocketHandler) RegisterRoutes(r *mux.Router) {
	// Make sure this matches what the frontend expects
	r.HandleFunc("/ws/game/{gameID}", h.ServeWS)
	log.Println("WebSocket routes registered at /ws/game/{gameID}")
}

// ServeWS handles WebSocket requests from clients
func (h *WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]

	if gameID == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}

	// Generate a client ID
	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		http.Error(w, "Client ID is required", http.StatusBadRequest)
		return
	}

	log.Printf("WebSocket connection request for client %s in game %s", clientID, gameID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	wsConn := customWs.NewConnection(conn)
	client := customWs.NewClient(clientID, wsConn, h.hub, gameID)

	// Register client with the hub
	h.hub.RegisterClient(client, gameID)

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()

	log.Printf("WebSocket client %s connected for game %s", clientID, gameID)
}

// BroadcastGameUpdate sends a game update to all clients in a game
func (h *WebSocketHandler) BroadcastGameUpdate(gameID string, data []byte) {
	log.Printf("Broadcasting update for game %s", gameID)
	h.hub.Broadcast(gameID, data)
}
