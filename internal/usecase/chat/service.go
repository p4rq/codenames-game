package chat

import (
	"codenames-game/internal/domain/chat"
)

// Service defines the interface for chat functionality
type Service interface {
	// SendMessage sends a chat message
	SendMessage(req chat.MessageRequest) error

	// GetMessages retrieves chat messages for a specific chat
	GetMessages(chatID string) ([]*chat.Message, error)

	// GetAllMessages retrieves all chat messages
	GetAllMessages() ([]*chat.Message, error)
}
