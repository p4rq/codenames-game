package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	// Remove the API import and use the interfaces instead
	"codenames-game/internal/domain/game"
	"codenames-game/internal/interfaces/websocket" // Use the interface package
)

// Repository defines the storage operations for games
type Repository interface {
	// Game operations
	Create(game *game.GameState) error // Changed from Save to Create
	FindByID(id string) (*game.GameState, error)
	FindAll() ([]*game.GameState, error) // Added missing FindAll method
	Update(game *game.GameState) error
	Delete(id string) error

	// Word operations
	GetWords() ([]string, error)
	AddWord(word string) error
	AddWords(words []string) error
	DeleteWord(word string) error
}

// ServiceImpl implements the game Service interface
type ServiceImpl struct {
	games     map[string]*game.GameState
	mutex     sync.RWMutex
	wordList  []string
	repo      Repository                  // Optional repository for persistent storage
	wsHandler websocket.UpdateBroadcaster // Use the interface instead of concrete type
}

// NewService creates a new game service with in-memory storage
func NewService() Service {
	return newService(nil, nil)
}

// NewServiceWithRepo creates a new game service with the provided repository
func NewServiceWithRepo(repo Repository) Service {
	return newService(repo, nil)
}

// NewServiceWithWebSocket creates a new game service with WebSocket support
func NewServiceWithWebSocket(repo Repository, wsHandler websocket.UpdateBroadcaster) Service {
	return newService(repo, wsHandler)
}

// Private helper to initialize a service
func newService(repo Repository, wsHandler websocket.UpdateBroadcaster) *ServiceImpl {
	var wordList []string
	var err error

	// Try to load words from repository if available
	if repo != nil {
		words, err := repo.GetWords()
		if err == nil && len(words) >= 25 {
			wordList = words
			fmt.Printf("Loaded %d words from repository\n", len(wordList))
		} else if err != nil {
			fmt.Printf("Error loading words from repository: %v\n", err)
		} else if len(words) < 25 {
			fmt.Printf("Not enough words in repository (%d). Need at least 25.\n", len(words))
		}
	}

	// If no words loaded from repository, use default list
	if len(wordList) < 25 {
		// Fallback to default word list
		wordList = []string{
			"AFRICA", "AGENT", "AIR", "ALIEN", "ALPS", "AMAZON", "AMBULANCE", "AMERICA", "ANGEL",
			"ANTARCTICA", "APPLE", "ARM", "ATLANTIS", "AUSTRALIA", "AZTEC", "BACK", "BALL", "BAND",
			"BANK", "BAR", "BARK", "BAT", "BATTERY", "BEACH", "BEAR", "BEAT", "BED", "BEIJING",
			"BELL", "BELT", "BERLIN", "BERMUDA", "BERRY", "BILL", "BLOCK", "BOARD", "BOLT", "BOMB",
			"BOND", "BOOM", "BOOT", "BOTTLE", "BOW", "BOX", "BRIDGE", "BRUSH", "BUCK", "BUFFALO",
			"BUG", "BUGLE", "BUTTON", "CALF", "CANADA", "CAP", "CAPITAL", "CAR", "CARD", "CARROT",
			"CASINO", "CAST", "CAT", "CELL", "CENTAUR", "CENTER", "CHAIR", "CHANGE", "CHARGE", "CHECK",
			// Add more words as needed
		}
		fmt.Printf("Using default word list with %d words\n", len(wordList))

		// Save the default words to the repository if available
		if repo != nil {
			err = repo.AddWords(wordList)
			if err != nil {
				fmt.Printf("Failed to save default words to repository: %v\n", err)
			} else {
				fmt.Println("Default words saved to repository")
			}
		}
	}

	return &ServiceImpl{
		games:     make(map[string]*game.GameState),
		wordList:  wordList,
		repo:      repo,
		mutex:     sync.RWMutex{},
		wsHandler: wsHandler,
	}
}

// broadcastGameUpdate sends game state updates to all connected clients
func (s *ServiceImpl) broadcastGameUpdate(gameState *game.GameState) {
	if s.wsHandler != nil {
		gameData, err := json.Marshal(gameState)
		if err != nil {
			fmt.Printf("Error marshaling game state: %v\n", err)
			return
		}
		s.wsHandler.BroadcastGameUpdate(gameState.ID, gameData)
	}
}

// CreateGame creates a new game
func (s *ServiceImpl) CreateGame(req game.CreateGameRequest) (*game.GameState, error) {
	if req.CreatorID == "" || req.Username == "" {
		return nil, errors.New("creator ID and username are required")
	}

	// Create a new game with 25 random cards
	cards := s.generateCards()

	// Determine which team goes first (randomly)
	var firstTeam game.Team
	if rand.Intn(2) == 0 {
		firstTeam = game.RedTeam
	} else {
		firstTeam = game.BlueTeam
	}

	// Count cards per team
	redCards := 0
	blueCards := 0
	for _, card := range cards {
		if card.Type == game.RedCard { // Fixed comparison
			redCards++
		} else if card.Type == game.BlueCard { // Fixed comparison
			blueCards++
		}
	}

	// Generate a unique game ID
	gameID := generateGameID()

	// Create the new game state
	newGame := &game.GameState{
		ID:            gameID,
		Cards:         cards,
		Players:       make([]game.Player, 0),
		CurrentTurn:   firstTeam,
		RedCardsLeft:  redCards,
		BlueCardsLeft: blueCards,
		WinningTeam:   nil,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Add the creator as the first player - CHANGED TO SPECTATOR
	creator := game.Player{
		ID:          req.CreatorID,
		Username:    req.Username,
		Team:        game.Spectator, // Start as spectator instead of first team
		IsSpymaster: false,
	}
	newGame.Players = append(newGame.Players, creator)

	// Store in memory
	s.mutex.Lock() // Add mutex lock before modifying shared state
	s.games[gameID] = newGame
	s.mutex.Unlock()

	// If repository is available, store the game there too
	if s.repo != nil {
		err := s.repo.Create(newGame)
		if err != nil {
			return nil, err
		}
	}

	// Broadcast the new game
	s.broadcastGameUpdate(newGame)

	fmt.Printf("Created new game with ID: %s\n", newGame.ID) // Add debug output

	return newGame, nil
}

// Helper function to generate a unique game ID
func generateGameID() string {
	// A simple implementation - for production, use a more robust method
	return "game-" + strconv.FormatInt(time.Now().UnixNano(), 36)
}

// GetGame retrieves a game by ID
func (s *ServiceImpl) GetGame(gameID string) (*game.GameState, error) {
	// If we have a repository, try to fetch from there first
	if s.repo != nil {
		if gameState, err := s.repo.FindByID(gameID); err != nil {
			return nil, err
		} else if gameState != nil {
			return gameState, nil
		}
	}

	// Otherwise try in-memory cache
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	gameState, exists := s.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}

	return gameState, nil
}

// JoinGame adds a player to a game
func (s *ServiceImpl) JoinGame(req game.JoinGameRequest) (*game.GameState, error) {
	if req.GameID == "" || req.PlayerID == "" || req.Username == "" {
		return nil, errors.New("game ID, player ID and username are required")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	gameState, exists := s.games[req.GameID]
	if !exists {
		return nil, errors.New("game not found")
	}

	// Check if player is already in the game
	for i, player := range gameState.Players {
		if player.ID == req.PlayerID {
			// If player is already in game but wants to change their name or team, update it
			if req.Username != "" && req.Username != player.Username {
				gameState.Players[i].Username = req.Username
			}

			// If team is specified and different from current, update it
			if req.Team != "" && req.Team != player.Team {
				gameState.Players[i].Team = req.Team
			}

			gameState.UpdatedAt = time.Now()
			s.broadcastGameUpdate(gameState)
			return gameState, nil
		}
	}

	// Use spectator team if no team specified
	team := req.Team
	if team == "" {
		team = game.Spectator
	}

	// Add the new player
	player := game.Player{
		ID:          req.PlayerID,
		Username:    req.Username,
		Team:        team,
		IsSpymaster: false,
	}
	gameState.Players = append(gameState.Players, player)
	gameState.UpdatedAt = time.Now()

	// Broadcast the update
	s.broadcastGameUpdate(gameState)

	return gameState, nil
}

// RevealCard reveals a card
func (s *ServiceImpl) RevealCard(req game.RevealCardRequest) (*game.GameState, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	gameState, exists := s.games[req.GameID]
	if !exists {
		return nil, errors.New("game not found")
	}

	// Check if the game is already over
	if gameState.WinningTeam != nil {
		return nil, errors.New("game is already over")
	}

	// Find the player
	var player *game.Player
	for i := range gameState.Players {
		if gameState.Players[i].ID == req.PlayerID {
			player = &gameState.Players[i]
			break
		}
	}

	if player == nil {
		return nil, errors.New("player not found in this game")
	}

	// Spectators can't reveal cards
	if player.Team == game.Spectator {
		return nil, errors.New("spectators cannot reveal cards")
	}

	// Spymasters can't reveal cards
	if player.IsSpymaster {
		return nil, errors.New("spymasters cannot reveal cards")
	}

	// Check if it's the player's team's turn
	if player.Team != gameState.CurrentTurn {
		return nil, errors.New("it's not your team's turn")
	}

	// Find and reveal the card
	var cardRevealed *game.Card
	for i := range gameState.Cards {
		if gameState.Cards[i].ID == req.CardID {
			cardRevealed = &gameState.Cards[i]
			break
		}
	}

	if cardRevealed == nil {
		return nil, errors.New("card not found")
	}

	if cardRevealed.Revealed {
		return nil, errors.New("card is already revealed")
	}

	// Reveal the card
	cardRevealed.Revealed = true
	gameState.UpdatedAt = time.Now()

	// Handle the consequences of revealing this card
	switch cardRevealed.Type {
	case game.RedCard:
		gameState.RedCardsLeft--
		if gameState.RedCardsLeft == 0 {
			redTeam := game.RedTeam
			gameState.WinningTeam = &redTeam
		}
		if gameState.CurrentTurn != game.RedTeam {
			gameState.CurrentTurn = game.RedTeam // Switch turns
		}
	case game.BlueCard:
		gameState.BlueCardsLeft--
		if gameState.BlueCardsLeft == 0 {
			blueTeam := game.BlueTeam
			gameState.WinningTeam = &blueTeam
		}
		if gameState.CurrentTurn != game.BlueTeam {
			gameState.CurrentTurn = game.BlueTeam // Switch turns
		}
	case game.AssassinCard:
		// Game over - the team that revealed the assassin loses
		var winningTeam game.Team
		if gameState.CurrentTurn == game.RedTeam {
			winningTeam = game.BlueTeam
		} else {
			winningTeam = game.RedTeam
		}
		gameState.WinningTeam = &winningTeam
	default: // NeutralCard
		// Switch turns
		if gameState.CurrentTurn == game.RedTeam {
			gameState.CurrentTurn = game.BlueTeam
		} else {
			gameState.CurrentTurn = game.RedTeam
		}
	}

	// Broadcast the update
	s.broadcastGameUpdate(gameState)

	return gameState, nil
}

// SetSpymaster sets a player as a spymaster
func (s *ServiceImpl) SetSpymaster(gameID string, playerID string) (*game.GameState, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	gameState, exists := s.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}

	// Find the player
	var player *game.Player
	for i := range gameState.Players {
		if gameState.Players[i].ID == playerID {
			player = &gameState.Players[i]
			break
		}
	}

	if player == nil {
		return nil, errors.New("player not found in this game")
	}

	// Spectators can't be spymasters
	if player.Team == game.Spectator {
		return nil, errors.New("spectators cannot be spymasters")
	}

	// Check if there's already a spymaster for this team
	for _, p := range gameState.Players {
		if p.Team == player.Team && p.IsSpymaster && p.ID != playerID {
			return nil, fmt.Errorf("team %s already has a spymaster", player.Team)
		}
	}

	player.IsSpymaster = true
	gameState.UpdatedAt = time.Now()

	// Broadcast the update
	s.broadcastGameUpdate(gameState)

	return gameState, nil
}

// EndTurn ends the current team's turn
func (s *ServiceImpl) EndTurn(gameID string, playerID string) (*game.GameState, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	gameState, exists := s.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}

	// Check if the game is already over
	if gameState.WinningTeam != nil {
		return nil, errors.New("game is already over")
	}

	// Find the player
	var player *game.Player
	for i := range gameState.Players {
		if gameState.Players[i].ID == playerID {
			player = &gameState.Players[i]
			break
		}
	}

	if player == nil {
		return nil, errors.New("player not found in this game")
	}

	// Spectators can't end turns
	if player.Team == game.Spectator {
		return nil, errors.New("spectators cannot end turns")
	}

	// Check if it's the player's team's turn
	if player.Team != gameState.CurrentTurn {
		return nil, errors.New("it's not your team's turn")
	}

	// Switch turns
	if gameState.CurrentTurn == game.RedTeam {
		gameState.CurrentTurn = game.BlueTeam
	} else {
		gameState.CurrentTurn = game.RedTeam
	}

	gameState.UpdatedAt = time.Now()

	// Broadcast the update
	s.broadcastGameUpdate(gameState)

	return gameState, nil
}

// ChangeTeam changes a player's team
func (s *ServiceImpl) ChangeTeam(gameID string, playerID string, team game.Team) (*game.GameState, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Validate team
	if team != game.RedTeam && team != game.BlueTeam && team != game.Spectator {
		return nil, fmt.Errorf("invalid team: %s", team)
	}

	gameState, exists := s.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}

	playerIndex := -1
	for i, p := range gameState.Players {
		if p.ID == playerID {
			playerIndex = i
			break
		}
	}

	if playerIndex == -1 {
		return nil, errors.New("player not found in this game")
	}

	// Don't allow spymasters to change teams unless they're becoming spectators
	if gameState.Players[playerIndex].IsSpymaster && team != game.Spectator {
		return nil, errors.New("spymasters cannot change teams (must become spectator first)")
	}

	// Update the player's team
	gameState.Players[playerIndex].Team = team

	// If changing to spectator, remove spymaster status
	if team == game.Spectator {
		gameState.Players[playerIndex].IsSpymaster = false
	}

	gameState.UpdatedAt = time.Now()

	// Update repository if available
	if s.repo != nil {
		err := s.repo.Update(gameState)
		if err != nil {
			return nil, err
		}
	}

	// Broadcast the update
	s.broadcastGameUpdate(gameState)

	return gameState, nil
}

// Helper function to generate 25 random cards for a new game
func (s *ServiceImpl) generateCards() []game.Card {
	// Shuffle the word list
	rand.Seed(time.Now().UnixNano())
	shuffled := make([]string, len(s.wordList))
	copy(shuffled, s.wordList)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// Pick 25 words
	words := shuffled[:25]

	// Determine the first team (gets 9 cards instead of 8)
	var firstTeamColor, secondTeamColor game.CardType
	if rand.Intn(2) == 0 {
		firstTeamColor = game.RedCard
		secondTeamColor = game.BlueCard
	} else {
		firstTeamColor = game.BlueCard
		secondTeamColor = game.RedCard
	}

	// Assign card types: 9 for first team, 8 for second team, 7 neutral, 1 assassin
	cards := make([]game.Card, 25)
	for i := 0; i < 25; i++ {
		cardType := game.NeutralCard
		if i < 9 {
			cardType = firstTeamColor
		} else if i < 17 {
			cardType = secondTeamColor
		} else if i < 24 {
			cardType = game.NeutralCard
		} else {
			cardType = game.AssassinCard
		}

		cards[i] = game.Card{
			ID:       uuid.New().String(),
			Word:     words[i],
			Type:     cardType,
			Revealed: false,
		}
	}

	// Shuffle the cards to randomize the distribution
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	return cards
}

// Add these methods to your ServiceImpl

// GetAllWords returns all available words
func (s *ServiceImpl) GetAllWords() ([]string, error) {
	if s.repo != nil {
		return s.repo.GetWords()
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.wordList, nil
}

// AddNewWord adds a new word
func (s *ServiceImpl) AddNewWord(word string) error {
	if s.repo != nil {
		err := s.repo.AddWord(word)
		if err != nil {
			return err
		}

		// Update the in-memory word list
		words, err := s.repo.GetWords()
		if err == nil {
			s.mutex.Lock()
			s.wordList = words
			s.mutex.Unlock()
		}

		return nil
	}

	// If no repository, just update the in-memory list
	word = strings.TrimSpace(strings.ToUpper(word))
	if word == "" {
		return errors.New("word cannot be empty")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if word already exists
	for _, w := range s.wordList {
		if w == word {
			return nil // Word already exists
		}
	}

	s.wordList = append(s.wordList, word)
	return nil
}

// DeleteExistingWord removes a word
func (s *ServiceImpl) DeleteExistingWord(word string) error {
	if s.repo != nil {
		err := s.repo.DeleteWord(word)
		if err != nil {
			return err
		}

		// Update the in-memory word list
		words, err := s.repo.GetWords()
		if err == nil {
			s.mutex.Lock()
			s.wordList = words
			s.mutex.Unlock()
		}

		return nil
	}

	// If no repository, just update the in-memory list
	word = strings.TrimSpace(strings.ToUpper(word))

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Filter out the word
	var newList []string
	for _, w := range s.wordList {
		if w != word {
			newList = append(newList, w)
		}
	}

	s.wordList = newList
	return nil
}
