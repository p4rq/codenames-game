package chat

// Repository defines the interface for chat storage
type Repository interface {
	// SaveMessage saves a chat message
	SaveMessage(message *Message) error

	// GetMessages retrieves all messages for a specific game
	GetMessages(chatID string) ([]*Message, error)

	// GetMessagesByTeam retrieves messages for a specific game and team
	GetMessagesByTeam(chatID string, team string) ([]*Message, error)

	// GetAllMessages retrieves all messages
	GetAllMessages() ([]*Message, error)
}
