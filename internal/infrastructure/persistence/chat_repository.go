package persistence

import (
	"codenames-game/internal/domain/chat"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ChatRepository implements chat.Repository interface with in-memory storage
type ChatRepository struct {
	messages []*chat.Message
	mutex    sync.RWMutex
}

// NewChatRepository creates a new chat repository instance
func NewChatRepository() *ChatRepository {
	return &ChatRepository{
		messages: make([]*chat.Message, 0),
	}
}

// SaveMessage stores a message
func (r *ChatRepository) SaveMessage(message *chat.Message) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Ensure the message has required fields
	if message.ID == "" {
		message.ID = uuid.New().String()
	}

	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}

	r.messages = append(r.messages, message)
	return nil
}

// GetMessages retrieves messages for a specific chat
func (r *ChatRepository) GetMessages(chatID string) ([]*chat.Message, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*chat.Message

	for _, msg := range r.messages {
		if msg.ChatID == chatID {
			result = append(result, msg)
		}
	}

	return result, nil
}

// GetAllMessages retrieves all messages
func (r *ChatRepository) GetAllMessages() ([]*chat.Message, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Create a copy to prevent race conditions
	result := make([]*chat.Message, len(r.messages))
	copy(result, r.messages)

	return result, nil
}
