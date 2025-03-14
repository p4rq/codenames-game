package chat

import "time"

type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	SenderID  string    `json:"sender_id"`
	Username  string    `json:"username"`
	ChatID    string    `json:"chat_id,omitempty"` // For game-specific chats
	Timestamp time.Time `json:"timestamp"`
}

type Chat struct {
	ID       string    `json:"id"`
	Messages []Message `json:"messages"`
}
