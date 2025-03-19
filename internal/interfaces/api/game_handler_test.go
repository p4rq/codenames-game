package api

import (
	"bytes"
	"codenames-game/internal/domain/game"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockGameRepository implements the game.Repository interface for testing
type MockRepository struct {
	games map[string]*game.Game
	words []string
}

// Implement Repository methods...
func (m *MockRepository) Save(game *game.Game) error {
	m.games[game.ID] = game
	return nil
}

func (m *MockRepository) FindByID(id string) (*game.Game, error) {
	game, exists := m.games[id]
	if !exists {
		return nil, nil
	}
	return game, nil
}

func (m *MockRepository) Update(game *game.Game) error {
	m.games[game.ID] = game
	return nil
}

func (m *MockRepository) AddWord(word string) error {
	m.words = append(m.words, word)
	return nil
}

func (m *MockRepository) GetRandomWords(count int) ([]string, error) {
	if count > len(m.words) {
		count = len(m.words)
	}
	return m.words[:count], nil
}

// MockGameService implements the gameUsecase.Service interface for testing
type MockGameService struct {
	repo *MockRepository
}

// Implement the Service interface methods - we only need to implement
// the ones that will be called in our tests
func (s *MockGameService) CreateGame(req game.CreateGameRequest) (*game.GameState, error) {
	// Simple implementation for testing
	gameState := &game.GameState{
		ID: "test-game-id",
		// Add other required fields
	}
	return gameState, nil
}

// Implement other methods from the Service interface with minimal functionality
func (s *MockGameService) GetGame(gameID string) (*game.GameState, error) {
	return &game.GameState{ID: gameID}, nil
}

func (s *MockGameService) JoinGame(req game.JoinGameRequest) (*game.GameState, error) {
	return &game.GameState{ID: req.GameID}, nil
}

func (s *MockGameService) RevealCard(req game.RevealCardRequest) (*game.GameState, error) {
	return &game.GameState{ID: req.GameID}, nil
}

func (s *MockGameService) SetSpymaster(gameID string, playerID string) (*game.GameState, error) {
	return &game.GameState{ID: gameID}, nil
}

func (s *MockGameService) EndTurn(gameID string, playerID string) (*game.GameState, error) {
	return &game.GameState{ID: gameID}, nil
}

func (s *MockGameService) ChangeTeam(gameID string, playerID string, team game.Team) (*game.GameState, error) {
	return &game.GameState{ID: gameID}, nil
}

func (s *MockGameService) GetAllWords() ([]string, error) {
	return s.repo.words, nil
}

func (s *MockGameService) AddNewWord(word string) error {
	return s.repo.AddWord(word)
}

func (s *MockGameService) DeleteExistingWord(word string) error {
	return nil
}

func TestStartGame(t *testing.T) {
	repo := &MockRepository{
		games: make(map[string]*game.Game),
		words: []string{"apple", "banana", "cherry", "dog", "elephant"},
	}

	// Use a mock service instead of the real one
	service := &MockGameService{repo: repo}

	handler := NewGameHandler(service)

	reqBody := `{"creatorID":"creator1","username":"player1"}`
	req, err := http.NewRequest("POST", "/game/start", bytes.NewBufferString(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.StartGame(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Print the response for debugging
	t.Logf("Response body: %s", rr.Body.String())

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check for the correct field name based on your GameState struct's JSON tags
	assert.NotNil(t, response["id"], "Game ID should be present in response")

	// Alternatively, you can unmarshal directly to a GameState to validate the structure
	var gameState game.GameState
	err = json.Unmarshal(rr.Body.Bytes(), &gameState)
	assert.NoError(t, err, "Response should unmarshal to GameState")
	assert.Equal(t, "test-game-id", gameState.ID)
}
