package chat

type ChatRepository interface {
    SaveMessage(message *ChatMessage) error
    GetMessages(chatID string) ([]*ChatMessage, error)
}