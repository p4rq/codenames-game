package api

import (
    "net/http"
    "github.com/gorilla/mux"
    "codenames-game/internal/usecase/game"
)

type GameHandler struct {
    gameService game.Service
}

func NewGameHandler(gs game.Service) *GameHandler {
    return &GameHandler{gameService: gs}
}

func (h *GameHandler) StartGame(w http.ResponseWriter, r *http.Request) {
    // Implementation for starting a game
}

func (h *GameHandler) GetGameState(w http.ResponseWriter, r *http.Request) {
    // Implementation for getting the current game state
}

func (h *GameHandler) RegisterRoutes(r *mux.Router) {
    r.HandleFunc("/api/game/start", h.StartGame).Methods("POST")
    r.HandleFunc("/api/game/state", h.GetGameState).Methods("GET")
}