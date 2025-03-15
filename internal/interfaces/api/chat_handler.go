package api

import (
	"encoding/json"
	"log"
	"net/http"

	"codenames-game/internal/domain/chat"
	chatservice "codenames-game/internal/usecase/chat"

	"github.com/gorilla/mux"
)

type ChatHandler struct {
	chatService chatservice.Service
}

func NewChatHandler(cs chatservice.Service) *ChatHandler {
	return &ChatHandler{chatService: cs}
}

// GetMessages handler for retrieving chat messages
func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	// Extract game ID from query parameters
	gameId := r.URL.Query().Get("game_id")
	team := r.URL.Query().Get("team") // Get team from query parameters

	log.Printf("GetMessages called with gameId=%s, team=%s", gameId, team)

	messages, err := h.chatService.GetMessages(gameId, team)
	if err != nil {
		log.Printf("Error fetching messages: %v", err)
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Printf("Error encoding messages: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// SendMessage handler for sending new chat messages
func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var req chat.MessageRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding message request: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Get game ID from query parameters
	req.ChatID = r.URL.Query().Get("game_id")

	log.Printf("Sending message: gameId=%s, team=%s, sender=%s", req.ChatID, req.Team, req.Username)

	if err := h.chatService.SendMessage(req); err != nil {
		log.Printf("Error sending message: %v", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// RegisterRoutes registers the chat API routes with the router
func (h *ChatHandler) RegisterRoutes(router *mux.Router) {
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Game-specific messages using path variables
	apiRouter.HandleFunc("/games/{gameId}/messages", h.GetGameMessages).Methods("GET")
	apiRouter.HandleFunc("/games/{gameId}/messages", h.SendGameMessage).Methods("POST")

	// Legacy routes for backward compatibility
	apiRouter.HandleFunc("/chat/messages", h.GetMessages).Methods("GET")
	apiRouter.HandleFunc("/chat/send", h.SendMessage).Methods("POST")
}

// GetGameMessages handler for retrieving game-specific chat messages
func (h *ChatHandler) GetGameMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameId := vars["gameId"]
	team := r.URL.Query().Get("team") // Get team from query parameters

	log.Printf("GetGameMessages called with gameId=%s, team=%s", gameId, team)

	messages, err := h.chatService.GetMessages(gameId, team)
	if err != nil {
		log.Printf("Error fetching game messages: %v", err)
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		return
	}

	// Return an empty array instead of null if no messages
	if messages == nil {
		messages = []*chat.Message{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Printf("Error encoding game messages: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// SendGameMessage handler for sending game-specific chat messages
func (h *ChatHandler) SendGameMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameId := vars["gameId"]

	var req chat.MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding game message request: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	req.ChatID = gameId // Set the game ID from URL path

	log.Printf("Sending game message: gameId=%s, team=%s, sender=%s", req.ChatID, req.Team, req.Username)

	if err := h.chatService.SendMessage(req); err != nil {
		log.Printf("Error sending game message: %v", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	// Return the created message
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
