package chat

import "codenames-game/internal/domain/chat"

type ChatService struct {
    repository chat.ChatRepository
}

func NewChatService(repository chat.ChatRepository) *ChatService {
    return &ChatService{repository: repository}
}

func (s *ChatService) SendMessage(message chat.Message) error {
    return s.repository.SaveMessage(message)
}

func (s *ChatService) GetMessages() ([]chat.Message, error) {
    return s.repository.FetchMessages()
}