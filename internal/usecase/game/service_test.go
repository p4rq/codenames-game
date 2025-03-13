package game

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "codenames-game/internal/domain/game"
    "codenames-game/internal/usecase/game"
)

type mockGameRepository struct {
    games map[string]*game.Game
}

func (m *mockGameRepository) Create(game *game.Game) error {
    m.games[game.ID] = game
    return nil
}

func (m *mockGameRepository) FindByID(id string) (*game.Game, error) {
    return m.games[id], nil
}

func (m *mockGameRepository) Update(game *game.Game) error {
    m.games[game.ID] = game
    return nil
}

func (m *mockGameRepository) Delete(id string) error {
    delete(m.games, id)
    return nil
}

func TestStartGame(t *testing.T) {
    repo := &mockGameRepository{games: make(map[string]*game.Game)}
    service := game.NewService(repo)

    gameID, err := service.StartGame()
    assert.NoError(t, err)
    assert.NotEmpty(t, gameID)

    createdGame, err := repo.FindByID(gameID)
    assert.NoError(t, err)
    assert.Equal(t, gameID, createdGame.ID)
}

func TestGetGameState(t *testing.T) {
    repo := &mockGameRepository{games: make(map[string]*game.Game)}
    service := game.NewService(repo)

    gameID, _ := service.StartGame()
    state, err := service.GetGameState(gameID)
    
    assert.NoError(t, err)
    assert.Equal(t, gameID, state.ID)
}