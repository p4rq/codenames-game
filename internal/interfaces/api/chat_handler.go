package api

import (
    "net/http"
    "github.com/gorilla/mux"
    "codenames-game/internal/usecase/chat"
)

type ChatHandler struct {
    chatService chat.Service
}

func NewChatHandler(chatService chat.Service) *ChatHandler {
    return &ChatHandler{chatService: chatService}
}

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
    // Implementation for sending a message
}

func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
    // Implementation for retrieving messages
}

func RegisterChatRoutes(router *mux.Router, chatService chat.Service) {
    handler := NewChatHandler(chatService)
    router.HandleFunc("/api/chat/send", handler.SendMessage).Methods("POST")
    router.HandleFunc("/api/chat/messages", handler.GetMessages).Methods("GET")
}