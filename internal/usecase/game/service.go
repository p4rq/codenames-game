package game

import (
    "codenames-game/internal/domain/game"
)

type GameService struct {
    repository game.GameRepository
}

func NewGameService(repo game.GameRepository) *GameService {
    return &GameService{repository: repo}
}

func (s *GameService) StartNewGame() (*game.Game, error) {
    newGame := game.NewGame() // Assuming NewGame initializes a new game entity
    err := s.repository.Save(newGame)
    if err != nil {
        return nil, err
    }
    return newGame, nil
}

func (s *GameService) GetGameState(gameID string) (*game.Game, error) {
    return s.repository.FindByID(gameID)
}