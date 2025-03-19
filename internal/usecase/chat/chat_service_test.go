package chat

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"codenames-game/internal/domain/chat"

	"github.com/stretchr/testify/assert"
)

// MockChatRepository implements the chat.Repository interface for testing
type MockChatRepository struct {
	messages map[string][]*chat.Message
}

func NewMockChatRepository() *MockChatRepository {
	return &MockChatRepository{
		messages: make(map[string][]*chat.Message),
	}
}

func (m *MockChatRepository) SaveMessage(message *chat.Message) error {
	chatID := message.ChatID
	m.messages[chatID] = append(m.messages[chatID], message)
	return nil
}

// Add GetAllMessages method to satisfy the Repository interface
func (m *MockChatRepository) GetAllMessages() ([]*chat.Message, error) {
	var allMessages []*chat.Message
	for _, msgs := range m.messages {
		allMessages = append(allMessages, msgs...)
	}
	return allMessages, nil
}

// GetMessages method with the correct signature to match Repository interface
func (m *MockChatRepository) GetMessages(chatID string) ([]*chat.Message, error) {
	if chatID == "" {
		// Return all messages if no chatID specified
		return m.GetAllMessages()
	} else {
		// Return messages for a specific chat
		return m.messages[chatID], nil
	}
}

// Updated GetMessagesByTeam method with correct signature to match Repository interface
func (m *MockChatRepository) GetMessagesByTeam(chatID string, team string) ([]*chat.Message, error) {
	var teamMessages []*chat.Message

	if chatID == "" {
		// Filter by team across all chats
		for _, msgs := range m.messages {
			for _, msg := range msgs {
				if msg.Team == team {
					teamMessages = append(teamMessages, msg)
				}
			}
		}
	} else {
		// Filter by team in a specific chat
		messages, exists := m.messages[chatID]
		if !exists {
			return []*chat.Message{}, nil
		}

		for _, msg := range messages {
			if msg.Team == team {
				teamMessages = append(teamMessages, msg)
			}
		}
	}

	return teamMessages, nil
}

func TestSendMessage(t *testing.T) {
	repo := NewMockChatRepository()
	service := NewChatService(repo)

	// Create a message request
	msgRequest := chat.MessageRequest{
		Content:  "Hello, world!",
		SenderID: "player1",
		Username: "player1",
		ChatID:   "game1",
		Team:     "Red",
	}

	// Send the message
	err := service.SendMessage(msgRequest)
	assert.NoError(t, err)

	// Get all messages - passing both required parameters
	messages, err := service.GetMessages("", "")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(messages))
	assert.Equal(t, "Hello, world!", messages[0].Content)

	// Get messages for a specific game/chat
	chatMessages, err := service.GetMessages("game1", "")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(chatMessages))
	assert.Equal(t, "Hello, world!", chatMessages[0].Content)
}

func TestGetMessages(t *testing.T) {
	repo := NewMockChatRepository()
	service := NewChatService(repo)

	// Create message requests for different chats
	msgRequest1 := chat.MessageRequest{
		Content:  "Hello from game 1!",
		SenderID: "player1",
		Username: "player1",
		ChatID:   "game1",
		Team:     "Red",
	}

	msgRequest2 := chat.MessageRequest{
		Content:  "Hello from game 2!",
		SenderID: "player2",
		Username: "player2",
		ChatID:   "game2",
		Team:     "Blue",
	}

	// Send the messages
	err := service.SendMessage(msgRequest1)
	assert.NoError(t, err)
	err = service.SendMessage(msgRequest2)
	assert.NoError(t, err)

	// Get all messages
	allMessages, err := service.GetMessages("", "")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(allMessages))

	// Get messages for specific chats
	chat1Messages, err := service.GetMessages("game1", "")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(chat1Messages))
	assert.Equal(t, "Hello from game 1!", chat1Messages[0].Content)

	chat2Messages, err := service.GetMessages("game2", "")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(chat2Messages))
	assert.Equal(t, "Hello from game 2!", chat2Messages[0].Content)

	// Test team filtering
	redTeamMessages, err := service.GetMessages("", "Red")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(redTeamMessages))
	assert.Equal(t, "Hello from game 1!", redTeamMessages[0].Content)

	// Test specific team in specific chat
	redTeamGame1Messages, err := service.GetMessages("game1", "Red")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(redTeamGame1Messages))
	assert.Equal(t, "Hello from game 1!", redTeamGame1Messages[0].Content)
}

// Add these test functions after your existing tests

// MockErrorRepository simulates repository with errors
type MockErrorRepository struct {
	shouldFailSave    bool
	shouldFailGet     bool
	shouldFailGetTeam bool
	shouldFailGetAll  bool
}

func NewMockErrorRepository() *MockErrorRepository {
	return &MockErrorRepository{}
}

func (m *MockErrorRepository) SaveMessage(message *chat.Message) error {
	if m.shouldFailSave {
		return fmt.Errorf("simulated database error")
	}
	return nil
}

func (m *MockErrorRepository) GetMessages(chatID string) ([]*chat.Message, error) {
	if m.shouldFailGet {
		return nil, fmt.Errorf("simulated database error")
	}
	return []*chat.Message{}, nil
}

func (m *MockErrorRepository) GetMessagesByTeam(chatID string, team string) ([]*chat.Message, error) {
	if m.shouldFailGetTeam {
		return nil, fmt.Errorf("simulated database error")
	}
	return []*chat.Message{}, nil
}

func (m *MockErrorRepository) GetAllMessages() ([]*chat.Message, error) {
	if m.shouldFailGetAll {
		return nil, fmt.Errorf("simulated database error")
	}
	return []*chat.Message{}, nil
}

// Test error handling for SaveMessage
func TestSendMessageError(t *testing.T) {
	repo := NewMockErrorRepository()
	repo.shouldFailSave = true
	service := NewChatService(repo)

	msgRequest := chat.MessageRequest{
		Content:  "Hello, world!",
		SenderID: "player1",
		Username: "player1",
		ChatID:   "game1",
		Team:     "Red",
	}

	err := service.SendMessage(msgRequest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "simulated database error")
}

// Test error handling for GetMessages
func TestGetMessagesError(t *testing.T) {
	repo := NewMockErrorRepository()
	service := NewChatService(repo)

	// Test error when fetching all messages
	repo.shouldFailGetAll = true
	_, err := service.GetMessages("", "") // Используем _ для игнорирования первого возвращаемого значения

	// Проверяем, что возникла ошибка, перед тем как продолжить
	if assert.Error(t, err, "Expected an error when GetAllMessages fails") {
		assert.Contains(t, err.Error(), "simulated database error")
	}

	// Test error when fetching messages by chat
	repo.shouldFailGetAll = false
	repo.shouldFailGet = true
	_, err = service.GetMessages("game1", "") // Используем _ для игнорирования первого возвращаемого значения

	// Проверяем, что возникла ошибка, перед тем как продолжить
	if assert.Error(t, err, "Expected an error when GetMessages fails") {
		assert.Contains(t, err.Error(), "simulated database error")
	}

	// Test error when fetching messages by team
	repo.shouldFailGet = false
	repo.shouldFailGetTeam = true
	_, err = service.GetMessages("", "Red") // Используем _ для игнорирования первого возвращаемого значения

	// Проверяем, что возникла ошибка, перед тем как продолжить
	if assert.Error(t, err, "Expected an error when GetMessagesByTeam fails") {
		assert.Contains(t, err.Error(), "simulated database error")
	}
}

// Test invalid parameters
func TestSendMessageInvalidParams(t *testing.T) {
	repo := NewMockChatRepository()
	service := NewChatService(repo)

	testCases := []struct {
		name    string
		request chat.MessageRequest
		errMsg  string
	}{
		{
			name: "Empty Content",
			request: chat.MessageRequest{
				Content:  "",
				SenderID: "player1",
				Username: "player1",
				ChatID:   "game1",
				Team:     "Red",
			},
			errMsg: "message content cannot be empty",
		},
		{
			name: "Empty SenderID",
			request: chat.MessageRequest{
				Content:  "Hello",
				SenderID: "",
				Username: "player1",
				ChatID:   "game1",
				Team:     "Red",
			},
			errMsg: "sender ID cannot be empty",
		},
		{
			name: "Empty Username",
			request: chat.MessageRequest{
				Content:  "Hello",
				SenderID: "player1",
				Username: "",
				ChatID:   "game1",
				Team:     "Red",
			},
			errMsg: "username cannot be empty",
		},
		{
			name: "Empty ChatID",
			request: chat.MessageRequest{
				Content:  "Hello",
				SenderID: "player1",
				Username: "player1",
				ChatID:   "",
				Team:     "Red",
			},
			errMsg: "chat ID cannot be empty",
		},
		{
			name: "Invalid Team",
			request: chat.MessageRequest{
				Content:  "Hello",
				SenderID: "player1",
				Username: "player1",
				ChatID:   "game1",
				Team:     "InvalidTeam", // Assuming only Red, Blue are valid
			},
			errMsg: "invalid team",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := service.SendMessage(tc.request)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", tc.name)
			} else if !strings.Contains(err.Error(), tc.errMsg) {
				t.Errorf("Expected error containing '%s', got '%s'", tc.errMsg, err.Error())
			}
		})
	}
}

// Test performance with large number of messages
func TestGetMessagesPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	repo := NewMockChatRepository()
	service := NewChatService(repo)

	// Create a large number of messages
	const messageCount = 10000
	const chatCount = 10
	teams := []string{"Red", "Blue"}

	// Populate repository with test data
	for i := 0; i < messageCount; i++ {
		chatID := fmt.Sprintf("game%d", i%chatCount)
		team := teams[i%2]

		msgRequest := chat.MessageRequest{
			Content:  fmt.Sprintf("Message %d", i),
			SenderID: fmt.Sprintf("player%d", i%100),
			Username: fmt.Sprintf("player%d", i%100),
			ChatID:   chatID,
			Team:     team,
		}

		err := service.SendMessage(msgRequest)
		assert.NoError(t, err)
	}

	// Measure performance for different scenarios
	testCases := []struct {
		name   string
		chatID string
		team   string
	}{
		{"Get All Messages", "", ""},
		{"Get Single Chat", "game1", ""},
		{"Get Single Team", "", "Red"},
		{"Get Team in Chat", "game1", "Red"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()

			messages, err := service.GetMessages(tc.chatID, tc.team)
			assert.NoError(t, err)

			duration := time.Since(start)

			// Log results
			t.Logf("%s: Retrieved %d messages in %v", tc.name, len(messages), duration)

			// Simple performance assertion - should complete in reasonable time
			assert.Less(t, duration, 1*time.Second, "Operation took too long")
		})
	}
}

// Test message validation helper (assuming it exists in your service)
func TestValidateMessage(t *testing.T) {
	testCases := []struct {
		name    string
		message chat.MessageRequest
		valid   bool
	}{
		{
			name: "Valid Message",
			message: chat.MessageRequest{
				Content:  "Hello",
				SenderID: "player1",
				Username: "player1",
				ChatID:   "game1",
				Team:     "Red",
			},
			valid: true,
		},
		{
			name: "Message Too Long",
			message: chat.MessageRequest{
				Content:  strings.Repeat("a", 1001), // Assuming 1000 char limit
				SenderID: "player1",
				Username: "player1",
				ChatID:   "game1",
				Team:     "Red",
			},
			valid: false,
		},
		{
			name: "Message With HTML",
			message: chat.MessageRequest{
				Content:  "<script>alert('xss')</script>",
				SenderID: "player1",
				Username: "player1",
				ChatID:   "game1",
				Team:     "Red",
			},
			valid: false, // Assuming HTML is not allowed
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMockChatRepository()
			service := NewChatService(repo)

			// Вместо прямого вызова метода validateMessage, используем SendMessage
			// SendMessage должен вызывать validateMessage внутри
			err := service.SendMessage(tc.message)

			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
