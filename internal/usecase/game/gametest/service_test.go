package gametest

import (
	"testing"

	"codenames-game/internal/domain/game"
	gameservice "codenames-game/internal/usecase/game"

	"github.com/stretchr/testify/assert"
)

// mockGameRepository implements a mock repository for testing
type mockGameRepository struct {
	games map[string]*game.GameState
}

// Let's adapt the mock repository to work with our current game service implementation
func (m *mockGameRepository) Create(g *game.GameState) error {
	m.games[g.ID] = g
	return nil
}

func (m *mockGameRepository) FindByID(id string) (*game.GameState, error) {
	g, exists := m.games[id]
	if !exists {
		return nil, nil // return nil for not found, real implementation would return error
	}
	return g, nil
}

func (m *mockGameRepository) Update(g *game.GameState) error {
	m.games[g.ID] = g
	return nil
}

func (m *mockGameRepository) Delete(id string) error {
	delete(m.games, id)
	return nil
}

// TestCreateGame tests game creation
func TestCreateGame(t *testing.T) {
	repo := &mockGameRepository{games: make(map[string]*game.GameState)}
	service := gameservice.NewServiceWithRepo(repo)

	req := game.CreateGameRequest{
		CreatorID: "user-1",
		Username:  "Player 1",
	}

	gameState, err := service.CreateGame(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, gameState.ID)

	// Verify the game was stored in repository
	createdGame, err := repo.FindByID(gameState.ID)
	assert.NoError(t, err)
	assert.Equal(t, gameState.ID, createdGame.ID)
}

// TestGetGame tests retrieving a game
func TestGetGame(t *testing.T) {
	repo := &mockGameRepository{games: make(map[string]*game.GameState)}
	service := gameservice.NewServiceWithRepo(repo)

	// First create a game
	req := game.CreateGameRequest{
		CreatorID: "user-1",
		Username:  "Player 1",
	}

	createdGame, _ := service.CreateGame(req)
	gameID := createdGame.ID

	// Now get the game
	fetchedGame, err := service.GetGame(gameID)

	assert.NoError(t, err)
	assert.Equal(t, gameID, fetchedGame.ID)
}
