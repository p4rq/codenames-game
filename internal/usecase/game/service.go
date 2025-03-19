package game

import (
	"codenames-game/internal/domain/game"
)

// Service defines the game service interface
type Service interface {
	CreateGame(req game.CreateGameRequest) (*game.GameState, error)
	GetGame(gameID string) (*game.GameState, error)
	JoinGame(req game.JoinGameRequest) (*game.GameState, error)
	RevealCard(req game.RevealCardRequest) (*game.GameState, error)
	SetSpymaster(gameID string, playerID string) (*game.GameState, error)
	EndTurn(gameID string, playerID string) (*game.GameState, error)
	ChangeTeam(gameID string, playerID string, team game.Team) (*game.GameState, error)

	// Add these methods for word management
	GetAllWords() ([]string, error)
	AddNewWord(word string) error
	DeleteExistingWord(word string) error
}
