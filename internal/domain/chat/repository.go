package chat

// Repository defines the interface for chat data storage
type Repository interface {
	SaveMessage(message *Message) error
	GetMessages(chatID string) ([]*Message, error)
	GetAllMessages() ([]*Message, error)
}
