package chat

import (
	"time"
)

// Message represents a chat message
type Message struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// MessageRequest represents the request to send a new message
type MessageRequest struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

// Service defines the interface for chat functionality
type Service interface {
	// SendMessage adds a new message to the chat
	SendMessage(req MessageRequest) (*Message, error)

	// GetMessages retrieves chat messages, optionally filtered by parameters
	GetMessages(limit int, before time.Time) ([]Message, error)
}
