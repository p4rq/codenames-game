package persistence

import (
	"codenames-game/internal/domain/game"
	"errors"
	"sync"
)

// GameRepository implements an in-memory repository for games
type GameRepository struct {
	games map[string]*game.GameState
	mutex sync.RWMutex
}

// NewGameRepository creates a new in-memory game repository
func NewGameRepository() *GameRepository {
	return &GameRepository{
		games: make(map[string]*game.GameState),
	}
}

// Create stores a new game
func (r *GameRepository) Create(g *game.GameState) error {
	if g.ID == "" {
		return errors.New("game ID cannot be empty")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.games[g.ID] = g
	return nil
}

// FindByID retrieves a game by ID
func (r *GameRepository) FindByID(id string) (*game.GameState, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	g, exists := r.games[id]
	if !exists {
		return nil, errors.New("game not found")
	}
	return g, nil
}

// Update updates an existing game
func (r *GameRepository) Update(g *game.GameState) error {
	if g.ID == "" {
		return errors.New("game ID cannot be empty")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.games[g.ID]
	if !exists {
		return errors.New("game not found")
	}

	r.games[g.ID] = g
	return nil
}

// Delete removes a game
func (r *GameRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.games[id]
	if !exists {
		return errors.New("game not found")
	}

	delete(r.games, id)
	return nil
}
