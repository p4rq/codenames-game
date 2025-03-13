package persistence

import (
    "codenames-game/internal/domain/chat"
    "sync"
)

type InMemoryChatRepository struct {
    messages []chat.Message
    mu       sync.Mutex
}

func NewInMemoryChatRepository() *InMemoryChatRepository {
    return &InMemoryChatRepository{
        messages: []chat.Message{},
    }
}

func (r *InMemoryChatRepository) SaveMessage(message chat.Message) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.messages = append(r.messages, message)
    return nil
}

func (r *InMemoryChatRepository) GetMessages() ([]chat.Message, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    return r.messages, nil
}