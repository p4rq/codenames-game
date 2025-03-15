package chat

// MessageRequest represents a request to send a chat message
type MessageRequest struct {
	Content  string `json:"content"`
	SenderID string `json:"sender_id"`
	Username string `json:"username"`
	ChatID   string `json:"chat_id"` // Game ID
	Team     string `json:"team"`    // Team: "red", "blue", or empty for global chat
}
