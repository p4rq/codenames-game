package persistence

import (
	"codenames-game/internal/domain/chat"
	"sync"
)

// ChatRepository is an in-memory implementation of chat.Repository
type ChatRepository struct {
	messages []*chat.Message
	mutex    sync.RWMutex
}

// NewChatRepository creates a new chat repository
func NewChatRepository() *ChatRepository {
	return &ChatRepository{
		messages: make([]*chat.Message, 0),
	}
}

// SaveMessage saves a chat message to the repository
func (r *ChatRepository) SaveMessage(message *chat.Message) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.messages = append(r.messages, message)
	return nil
}

// GetMessages retrieves all messages for a specific game
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

// GetMessagesByTeam retrieves messages for a specific game and team
func (r *ChatRepository) GetMessagesByTeam(chatID string, team string) ([]*chat.Message, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*chat.Message

	for _, msg := range r.messages {
		if msg.ChatID == chatID && msg.Team == team {
			result = append(result, msg)
		}
	}

	return result, nil
}

// GetAllMessages retrieves all chat messages
func (r *ChatRepository) GetAllMessages() ([]*chat.Message, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]*chat.Message, len(r.messages))
	copy(result, r.messages)

	return result, nil
}
