package chat

import (
	"codenames-game/internal/domain/chat"
	"time"

	"github.com/google/uuid"
)

// ServiceImpl implements the chat Service interface
type ServiceImpl struct {
	repo chat.Repository
}

// NewService creates a new chat service with an in-memory repository
func NewService() Service {
	// This would be replaced with a real repository in production
	return &ServiceImpl{
		repo: newInMemoryRepository(),
	}
}

// NewChatService creates a new chat service with the provided repository
func NewChatService(repo chat.Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

// SendMessage sends a chat message
func (s *ServiceImpl) SendMessage(req chat.MessageRequest) error {
	message := &chat.Message{
		ID:        uuid.New().String(),
		Content:   req.Content,
		SenderID:  req.SenderID,
		Username:  req.Username,
		ChatID:    req.ChatID,
		Timestamp: time.Now(),
	}

	return s.repo.SaveMessage(message)
}

// GetMessages retrieves chat messages for a specific chat
func (s *ServiceImpl) GetMessages(chatID string) ([]*chat.Message, error) {
	return s.repo.GetMessages(chatID)
}

// GetAllMessages retrieves all chat messages
func (s *ServiceImpl) GetAllMessages() ([]*chat.Message, error) {
	return s.repo.GetAllMessages()
}

// inMemoryRepository is a simple in-memory implementation of Repository for testing
type inMemoryRepository struct {
	messages []*chat.Message
}

func newInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{
		messages: make([]*chat.Message, 0),
	}
}

func (r *inMemoryRepository) SaveMessage(message *chat.Message) error {
	r.messages = append(r.messages, message)
	return nil
}

func (r *inMemoryRepository) GetMessages(chatID string) ([]*chat.Message, error) {
	var result []*chat.Message

	for _, msg := range r.messages {
		if msg.ChatID == chatID {
			result = append(result, msg)
		}
	}

	return result, nil
}

func (r *inMemoryRepository) GetAllMessages() ([]*chat.Message, error) {
	result := make([]*chat.Message, len(r.messages))
	copy(result, r.messages)
	return result, nil
}
