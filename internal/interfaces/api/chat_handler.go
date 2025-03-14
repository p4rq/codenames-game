package api

import (
	"encoding/json"
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

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var req chat.MessageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" || req.SenderID == "" || req.Username == "" {
		http.Error(w, "Content, sender ID, and username are required", http.StatusBadRequest)
		return
	}

	err = h.chatService.SendMessage(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chat_id")

	var messages []*chat.Message
	var err error

	if chatID != "" {
		messages, err = h.chatService.GetMessages(chatID)
	} else {
		messages, err = h.chatService.GetAllMessages()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *ChatHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/chat/send", h.SendMessage).Methods("POST")
	r.HandleFunc("/api/chat/messages", h.GetMessages).Methods("GET")
}
