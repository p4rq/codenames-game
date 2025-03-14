package game

import (
	"codenames-game/internal/domain/game"
)

// Service defines the interface for game functionality
type Service interface {
	// CreateGame creates a new game
	CreateGame(req game.CreateGameRequest) (*game.GameState, error)

	// GetGame retrieves a game by ID
	GetGame(gameID string) (*game.GameState, error)

	// JoinGame adds a player to a game
	JoinGame(req game.JoinGameRequest) (*game.GameState, error)

	// RevealCard reveals a card
	RevealCard(req game.RevealCardRequest) (*game.GameState, error)

	// SetSpymaster sets a player as a spymaster
	SetSpymaster(gameID string, playerID string) (*game.GameState, error)

	// EndTurn ends the current team's turn
	EndTurn(gameID string, playerID string) (*game.GameState, error)
}
