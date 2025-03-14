package chat

import (
	"testing"

	"codenames-game/internal/domain/chat"
	"codenames-game/internal/infrastructure/persistence"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	repo := persistence.NewChatRepository()
	service := NewChatService(repo)

	req := chat.MessageRequest{
		Content:  "Hello, World!",
		SenderID: "user-1",
		Username: "User1",
	}

	err := service.SendMessage(req)
	assert.NoError(t, err)

	messages, err := service.GetAllMessages()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(messages))
	assert.Equal(t, "Hello, World!", messages[0].Content)
	assert.Equal(t, "user-1", messages[0].SenderID)
	assert.Equal(t, "User1", messages[0].Username)
}

func TestGetMessages(t *testing.T) {
	repo := persistence.NewChatRepository()
	service := NewChatService(repo)

	// Add a test message with a specific chat ID
	req1 := chat.MessageRequest{
		Content:  "Game message",
		SenderID: "user-1",
		Username: "User1",
		ChatID:   "game-123",
	}

	req2 := chat.MessageRequest{
		Content:  "General message",
		SenderID: "user-2",
		Username: "User2",
	}

	service.SendMessage(req1)
	service.SendMessage(req2)

	// Test getting messages for a specific chat
	gameMessages, err := service.GetMessages("game-123")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(gameMessages))
	assert.Equal(t, "Game message", gameMessages[0].Content)

	// Test getting all messages
	allMessages, err := service.GetAllMessages()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(allMessages))
}
