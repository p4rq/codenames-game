package chat

// MessageRequest represents a request to send a chat message
type MessageRequest struct {
	Content  string `json:"content"`
	SenderID string `json:"sender_id"`
	Username string `json:"username"`
	ChatID   string `json:"chat_id,omitempty"` // Optional, for game-specific chat
}
