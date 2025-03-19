package repository

import (
	"errors"
	"strings"
	"sync"
	"time"

	"codenames-game/internal/domain/game"
)

// InMemoryRepository implements Repository with in-memory storage
type InMemoryRepository struct {
	games       map[string]*game.GameState
	words       []string
	activeWords map[string]bool // Track which words are active
	mutex       sync.RWMutex
}

// NewInMemoryRepository creates a new repository with in-memory storage
func NewInMemoryRepository() *InMemoryRepository {
	// Default words
	defaultWords := []string{
		"AFRICA", "AGENT", "AIR", "ALIEN", "ALPS", "AMAZON", "AMBULANCE", "AMERICA", "ANGEL",
		"ANTARCTICA", "APPLE", "ARM", "ATLANTIS", "AUSTRALIA", "AZTEC", "BACK", "BALL", "BAND",
		"BANK", "BAR", "BARK", "BAT", "BATTERY", "BEACH", "BEAR", "BEAT", "BED", "BEIJING",
		"BELL", "BELT", "BERLIN", "BERMUDA", "BERRY", "BILL", "BLOCK", "BOARD", "BOLT", "BOMB",
		"BOND", "BOOM", "BOOT", "BOTTLE", "BOW", "BOX", "BRIDGE", "BRUSH", "BUCK", "BUFFALO",
		"BUG", "BUGLE", "BUTTON", "CALF", "CANADA", "CAP", "CAPITAL", "CAR", "CARD", "CARROT",
		"CASINO", "CAST", "CAT", "CELL", "CENTAUR", "CENTER", "CHAIR", "CHANGE", "CHARGE", "CHECK",
		// Add more words as needed
	}

	// Create active words map
	activeWords := make(map[string]bool)
	for _, word := range defaultWords {
		activeWords[word] = true
	}

	return &InMemoryRepository{
		games:       make(map[string]*game.GameState),
		words:       defaultWords,
		activeWords: activeWords,
		mutex:       sync.RWMutex{},
	}
}

// Save stores a game in memory
func (r *InMemoryRepository) Save(gameState *game.GameState) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.games[gameState.ID] = gameState
	return nil
}

// FindByID retrieves a game from memory by ID
func (r *InMemoryRepository) FindByID(id string) (*game.GameState, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	gameState, exists := r.games[id]
	if !exists {
		return nil, errors.New("game not found")
	}
	return gameState, nil
}

// FindAll retrieves all games from memory
func (r *InMemoryRepository) FindAll() ([]*game.GameState, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	games := make([]*game.GameState, 0, len(r.games))
	for _, gameState := range r.games {
		games = append(games, gameState)
	}
	return games, nil
}

// Update modifies a game in memory
func (r *InMemoryRepository) Update(gameState *game.GameState) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.games[gameState.ID]; !exists {
		return errors.New("game not found")
	}

	gameState.UpdatedAt = time.Now()
	r.games[gameState.ID] = gameState
	return nil
}

// Delete removes a game from memory
func (r *InMemoryRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.games[id]; !exists {
		return errors.New("game not found")
	}

	delete(r.games, id)
	return nil
}

// GetWords retrieves all active words from memory
func (r *InMemoryRepository) GetWords() ([]string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Filter only active words
	activeWords := make([]string, 0, len(r.activeWords))
	for word, active := range r.activeWords {
		if active {
			activeWords = append(activeWords, word)
		}
	}

	return activeWords, nil
}

// AddWord adds a word to memory
func (r *InMemoryRepository) AddWord(word string) error {
	word = strings.TrimSpace(strings.ToUpper(word))
	if word == "" {
		return errors.New("word cannot be empty")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// If word doesn't exist in the list, add it
	if _, exists := r.activeWords[word]; !exists {
		r.words = append(r.words, word)
	}

	// Mark the word as active
	r.activeWords[word] = true

	return nil
}

// AddWords adds multiple words to memory
func (r *InMemoryRepository) AddWords(words []string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, word := range words {
		word = strings.TrimSpace(strings.ToUpper(word))
		if word == "" {
			continue
		}

		// If word doesn't exist in the list, add it
		if _, exists := r.activeWords[word]; !exists {
			r.words = append(r.words, word)
		}

		// Mark the word as active
		r.activeWords[word] = true
	}

	return nil
}

// DeleteWord deactivates a word in memory
func (r *InMemoryRepository) DeleteWord(word string) error {
	word = strings.TrimSpace(strings.ToUpper(word))

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Mark the word as inactive
	r.activeWords[word] = false

	return nil
}
