package api

import (
	"encoding/json"
	"net/http"

	"codenames-game/internal/domain/game"
	gameservice "codenames-game/internal/usecase/game"

	"github.com/gorilla/mux"
)

type GameHandler struct {
	gameService gameservice.Service
}

func NewGameHandler(gs gameservice.Service) *GameHandler {
	return &GameHandler{gameService: gs}
}

func (h *GameHandler) StartGame(w http.ResponseWriter, r *http.Request) {
	var req game.CreateGameRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gameState, err := h.gameService.CreateGame(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(gameState)
}

func (h *GameHandler) GetGameState(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("id")
	if gameID == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}

	gameState, err := h.gameService.GetGame(gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameState)
}

func (h *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	var req game.JoinGameRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gameState, err := h.gameService.JoinGame(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameState)
}

func (h *GameHandler) RevealCard(w http.ResponseWriter, r *http.Request) {
	var req game.RevealCardRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gameState, err := h.gameService.RevealCard(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameState)
}

func (h *GameHandler) SetSpymaster(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("game_id")
	playerID := r.URL.Query().Get("player_id")

	if gameID == "" || playerID == "" {
		http.Error(w, "Game ID and Player ID are required", http.StatusBadRequest)
		return
	}

	gameState, err := h.gameService.SetSpymaster(gameID, playerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameState)
}

func (h *GameHandler) EndTurn(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("game_id")
	playerID := r.URL.Query().Get("player_id")

	if gameID == "" || playerID == "" {
		http.Error(w, "Game ID and Player ID are required", http.StatusBadRequest)
		return
	}

	gameState, err := h.gameService.EndTurn(gameID, playerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameState)
}

func (h *GameHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/game/start", h.StartGame).Methods("POST")
	r.HandleFunc("/api/game/state", h.GetGameState).Methods("GET")
	r.HandleFunc("/api/game/join", h.JoinGame).Methods("POST")
	r.HandleFunc("/api/game/reveal", h.RevealCard).Methods("POST")
	r.HandleFunc("/api/game/set-spymaster", h.SetSpymaster).Methods("POST")
	r.HandleFunc("/api/game/end-turn", h.EndTurn).Methods("POST")
}
