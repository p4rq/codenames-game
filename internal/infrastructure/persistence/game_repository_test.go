package persistence

import (
	"codenames-game/internal/domain/game"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGame(t *testing.T) {
	repo := NewGameRepository()

	gameState := &game.GameState{
		ID: "game1",
	}

	err := repo.Create(gameState)
	assert.NoError(t, err)

	retrievedGame, err := repo.FindByID("game1")
	assert.NoError(t, err)
	assert.Equal(t, "game1", retrievedGame.ID)
}
