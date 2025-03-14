package chat

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ServiceImpl implements the chat Service interface
type ServiceImpl struct {
	messages []Message
	mutex    sync.RWMutex
}

// NewService creates a new chat service
func NewService() Service {
	return &ServiceImpl{
		messages: make([]Message, 0),
	}
}

// SendMessage adds a new message to the chat
func (s *ServiceImpl) SendMessage(req MessageRequest) (*Message, error) {
	if req.Content == "" {
		return nil, errors.New("message content cannot be empty")
	}

	msg := Message{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Username:  req.Username,
		Content:   req.Content,
		Timestamp: time.Now(),
	}

	s.mutex.Lock()
	s.messages = append(s.messages, msg)
	s.mutex.Unlock()

	return &msg, nil
}

// GetMessages retrieves chat messages
func (s *ServiceImpl) GetMessages(limit int, before time.Time) ([]Message, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []Message
	var zeroTime time.Time
	useTimeFilter := before != zeroTime

	// Start from the most recent messages
	for i := len(s.messages) - 1; i >= 0 && len(result) < limit; i-- {
		msg := s.messages[i]
		if !useTimeFilter || msg.Timestamp.Before(before) {
			// Insert at the beginning to maintain chronological order
			result = append([]Message{msg}, result...)
		}
	}

	return result, nil
}
