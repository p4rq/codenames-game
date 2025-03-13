package chat

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "codenames-game/internal/domain/chat"
    "codenames-game/internal/infrastructure/persistence"
)

func TestSendMessage(t *testing.T) {
    repo := persistence.NewChatRepository()
    service := NewChatService(repo)

    message := chat.Message{Content: "Hello, World!", Sender: "User1"}
    err := service.SendMessage(message)

    assert.NoError(t, err)

    messages, err := service.GetMessages()
    assert.NoError(t, err)
    assert.Contains(t, messages, message)
}

func TestGetMessages(t *testing.T) {
    repo := persistence.NewChatRepository()
    service := NewChatService(repo)

    messages, err := service.GetMessages()

    assert.NoError(t, err)
    assert.NotNil(t, messages)
}