package game

import (
	"time"
)

// Team represents a team in the game
type Team string

const (
	RedTeam  Team = "red"
	BlueTeam Team = "blue"
)

// CardType represents the type of a card
type CardType string

const (
	RedCard      CardType = "red"
	BlueCard     CardType = "blue"
	NeutralCard  CardType = "neutral"
	AssassinCard CardType = "assassin"
)

// Card represents a word card in the game
type Card struct {
	ID       string   `json:"id"`
	Word     string   `json:"word"`
	Type     CardType `json:"type,omitempty"`
	Revealed bool     `json:"revealed"`
}

// Player represents a player in the game
type Player struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Team        Team   `json:"team"`
	IsSpymaster bool   `json:"is_spymaster"`
}

// GameState represents the current state of a game
type GameState struct {
	ID            string    `json:"id"`
	Cards         []Card    `json:"cards"`
	Players       []Player  `json:"players"`
	CurrentTurn   Team      `json:"current_turn"`
	WinningTeam   *Team     `json:"winning_team,omitempty"`
	RedCardsLeft  int       `json:"red_cards_left"`
	BlueCardsLeft int       `json:"blue_cards_left"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateGameRequest represents the request to create a new game
type CreateGameRequest struct {
	CreatorID string `json:"creator_id"`
	Username  string `json:"username"`
}

// JoinGameRequest represents the request to join a game
type JoinGameRequest struct {
	GameID   string `json:"game_id"`
	PlayerID string `json:"player_id"`
	Username string `json:"username"`
	Team     Team   `json:"team"`
}

// RevealCardRequest represents the request to reveal a card
type RevealCardRequest struct {
	GameID   string `json:"game_id"`
	CardID   string `json:"card_id"`
	PlayerID string `json:"player_id"`
}
