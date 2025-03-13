package game

type GameRepository interface {
    Save(game *Game) error
    FindByID(id string) (*Game, error)
    FindAll() ([]*Game, error)
    Update(game *Game) error
    Delete(id string) error
}