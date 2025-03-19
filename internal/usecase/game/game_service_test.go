package game

import (
	"codenames-game/internal/domain/game"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockRepository struct {
	games map[string]*game.GameState
	words []string
}

func (m *MockRepository) Create(g *game.GameState) error {
	m.games[g.ID] = g
	return nil
}

func (m *MockRepository) FindByID(id string) (*game.GameState, error) {
	if g, ok := m.games[id]; ok {
		return g, nil
	}
	return nil, errors.New("game not found")
}

func (m *MockRepository) FindAll() ([]*game.GameState, error) {
	var games []*game.GameState
	for _, g := range m.games {
		games = append(games, g)
	}
	return games, nil
}

func (m *MockRepository) Update(g *game.GameState) error {
	m.games[g.ID] = g
	return nil
}

func (m *MockRepository) Delete(id string) error {
	delete(m.games, id)
	return nil
}

func (m *MockRepository) GetWords() ([]string, error) {
	return m.words, nil
}

func (m *MockRepository) AddWord(word string) error {
	m.words = append(m.words, word)
	return nil
}

func (m *MockRepository) AddWords(words []string) error {
	m.words = append(m.words, words...)
	return nil
}

func (m *MockRepository) DeleteWord(word string) error {
	for i, w := range m.words {
		if w == word {
			m.words = append(m.words[:i], m.words[i+1:]...)
			break
		}
	}
	return nil
}

func TestCreateGame(t *testing.T) {
	repo := &MockRepository{games: make(map[string]*game.GameState)}
	service := NewServiceWithRepo(repo)

	req := game.CreateGameRequest{
		CreatorID: "creator1",
		Username:  "player1",
	}

	gameState, err := service.CreateGame(req)
	assert.NoError(t, err)
	assert.NotNil(t, gameState)
	assert.Equal(t, "creator1", gameState.Players[0].ID)
	assert.Equal(t, "player1", gameState.Players[0].Username)
}

func TestJoinGame(t *testing.T) {
	repo := &MockRepository{games: make(map[string]*game.GameState)}
	service := NewServiceWithRepo(repo)

	req := game.CreateGameRequest{
		CreatorID: "creator1",
		Username:  "player1",
	}

	gameState, err := service.CreateGame(req)
	assert.NoError(t, err)

	joinReq := game.JoinGameRequest{
		GameID:   gameState.ID,
		PlayerID: "player2",
		Username: "player2",
		Team:     game.RedTeam,
	}

	updatedGameState, err := service.JoinGame(joinReq)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(updatedGameState.Players))
	assert.Equal(t, "player2", updatedGameState.Players[1].ID)
	assert.Equal(t, game.RedTeam, updatedGameState.Players[1].Team)
}
