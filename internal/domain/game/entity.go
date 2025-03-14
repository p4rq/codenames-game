package game

import "time"

type Game struct {
	ID        string
	Name      string
	Players   []string
	State     string
	CreatedAt string
	UpdatedAt string
}

func NewGame(id, name string, players []string) *Game {
	return &Game{
		ID:        id,
		Name:      name,
		Players:   players,
		State:     "waiting",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
}

func (g *Game) UpdateState(state string) {
	g.State = state
	g.UpdatedAt = time.Now().Format(time.RFC3339)
}
