package persistence

import (
	"errors"
	"strings"
	"sync"
	"time"

	"codenames-game/internal/domain/game"
)

// GameRepository implements the repository interface for games
type GameRepository struct {
	games       map[string]*game.GameState
	words       []string
	activeWords map[string]bool
	mutex       sync.RWMutex
}

// NewGameRepository creates a new in-memory game repository
func NewGameRepository() *GameRepository {
	// Default word list
	defaultWords := []string{
		"AFRICA", "AGENT", "AIR", "ALIEN", "ALPS", "AMAZON", "AMBULANCE", "AMERICA", "ANGEL",
		"ANTARCTICA", "APPLE", "ARM", "ATLANTIS", "AUSTRALIA", "AZTEC", "BACK", "BALL", "BAND",
		"BANK", "BAR", "BARK", "BAT", "BATTERY", "BEACH", "BEAR", "BEAT", "BED", "BEIJING",
		"BELL", "BELT", "BERLIN", "BERMUDA", "BERRY", "BILL", "BLOCK", "BOARD", "BOLT", "BOMB",
		"BOND", "BOOM", "BOOT", "BOTTLE", "BOW", "BOX", "BRIDGE", "BRUSH", "BUCK", "BUFFALO",
		"BUG", "BUGLE", "BUTTON", "CALF", "CANADA", "CAP", "CAPITAL", "CAR", "CARD", "CARROT",
		"CASINO", "CAST", "CAT", "CELL", "CENTAUR", "CENTER", "CHAIR", "CHANGE", "CHARGE", "CHECK",
		"CHEST", "CHICK", "CHINA", "CHOCOLATE", "CHURCH", "CIRCLE", "CLIFF", "CLOAK", "CLUB", "CODE",
		"COLD", "COMIC", "COMPOUND", "CONCERT", "CONDUCTOR", "CONTRACT", "COOK", "COPPER", "COTTON", "COURT",
		"COVER", "CRANE", "CRASH", "CRICKET", "CROSS", "CROWN", "CYCLE", "CZECH", "DANCE", "DATE",
		"DAY", "DEATH", "DECK", "DEGREE", "DIAMOND", "DICE", "DINOSAUR", "DISEASE", "DOCTOR", "DOG",
		"DRAFT", "DRAGON", "DRESS", "DRILL", "DROP", "DUCK", "DWARF", "EAGLE", "EGYPT", "EMBASSY",
	}

	// Create active words map
	activeWords := make(map[string]bool)
	for _, word := range defaultWords {
		activeWords[word] = true
	}

	return &GameRepository{
		games:       make(map[string]*game.GameState),
		words:       defaultWords,
		activeWords: activeWords,
		mutex:       sync.RWMutex{},
	}
}

// Create stores a new game
func (r *GameRepository) Create(gameState *game.GameState) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.games[gameState.ID]; exists {
		return errors.New("game with this ID already exists")
	}

	r.games[gameState.ID] = gameState
	return nil
}

// FindByID retrieves a game by ID
func (r *GameRepository) FindByID(id string) (*game.GameState, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	gameState, exists := r.games[id]
	if !exists {
		return nil, errors.New("game not found")
	}

	return gameState, nil
}

// FindAll retrieves all games
func (r *GameRepository) FindAll() ([]*game.GameState, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var allGames []*game.GameState
	for _, gameState := range r.games {
		allGames = append(allGames, gameState)
	}

	return allGames, nil
}

// Update modifies an existing game
func (r *GameRepository) Update(gameState *game.GameState) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.games[gameState.ID]; !exists {
		return errors.New("game not found")
	}

	gameState.UpdatedAt = time.Now()
	r.games[gameState.ID] = gameState
	return nil
}

// Delete removes a game
func (r *GameRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.games[id]; !exists {
		return errors.New("game not found")
	}

	delete(r.games, id)
	return nil
}

// GetWords returns all active words
func (r *GameRepository) GetWords() ([]string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Filter only active words
	var activeWords []string
	for word, active := range r.activeWords {
		if active {
			activeWords = append(activeWords, word)
		}
	}

	return activeWords, nil
}

// AddWord adds a new word
func (r *GameRepository) AddWord(word string) error {
	word = strings.TrimSpace(strings.ToUpper(word))
	if word == "" {
		return errors.New("word cannot be empty")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Add to words list if not already there
	found := false
	for _, w := range r.words {
		if w == word {
			found = true
			break
		}
	}

	if !found {
		r.words = append(r.words, word)
	}

	// Mark as active
	r.activeWords[word] = true
	return nil
}

// AddWords adds multiple words
func (r *GameRepository) AddWords(words []string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, word := range words {
		word = strings.TrimSpace(strings.ToUpper(word))
		if word == "" {
			continue
		}

		// Add to words list if not already there
		found := false
		for _, w := range r.words {
			if w == word {
				found = true
				break
			}
		}

		if !found {
			r.words = append(r.words, word)
		}

		// Mark as active
		r.activeWords[word] = true
	}

	return nil
}

// DeleteWord deactivates a word
func (r *GameRepository) DeleteWord(word string) error {
	word = strings.TrimSpace(strings.ToUpper(word))
	if word == "" {
		return errors.New("word cannot be empty")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Mark as inactive
	r.activeWords[word] = false
	return nil
}
