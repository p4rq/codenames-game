package persistence

import (
    "errors"
    "codenames-game/internal/domain/game"
)

type InMemoryGameRepository struct {
    games map[string]*game.Game
}

func NewInMemoryGameRepository() *InMemoryGameRepository {
    return &InMemoryGameRepository{
        games: make(map[string]*game.Game),
    }
}

func (r *InMemoryGameRepository) Save(game *game.Game) error {
    if game == nil {
        return errors.New("game cannot be nil")
    }
    r.games[game.ID] = game
    return nil
}

func (r *InMemoryGameRepository) FindByID(id string) (*game.Game, error) {
    game, exists := r.games[id]
    if !exists {
        return nil, errors.New("game not found")
    }
    return game, nil
}

func (r *InMemoryGameRepository) FindAll() ([]*game.Game, error) {
    var allGames []*game.Game
    for _, game := range r.games {
        allGames = append(allGames, game)
    }
    return allGames, nil
}