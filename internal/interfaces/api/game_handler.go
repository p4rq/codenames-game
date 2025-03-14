package api

import (
	"encoding/json"
	"log"
	"net/http"

	"codenames-game/internal/domain/game"
	gameservice "codenames-game/internal/usecase/game"

	"github.com/gorilla/mux"
)

// GameHandler handles HTTP requests related to game operations
type GameHandler struct {
	gameService gameservice.Service
}

// NewGameHandler creates a new game handler
func NewGameHandler(gs gameservice.Service) *GameHandler {
	return &GameHandler{
		gameService: gs,
	}
}

// StartGame handles the request to create a new game
func (h *GameHandler) StartGame(w http.ResponseWriter, r *http.Request) {
	log.Println("StartGame handler called")

	var req struct {
		CreatorID string `json:"creator_id"`
		Username  string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Request received: %+v", req)

	createReq := game.CreateGameRequest{
		CreatorID: req.CreatorID,
		Username:  req.Username,
	}

	gameState, err := h.gameService.CreateGame(createReq)
	if err != nil {
		log.Printf("Error creating game: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Game created with ID: %s", gameState.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameState)
}

// JoinGame handles the request to join an existing game
func (h *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GameID   string `json:"game_id"`
		PlayerID string `json:"player_id"`
		Username string `json:"username"`
		Team     string `json:"team"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	joinReq := game.JoinGameRequest{
		GameID:   req.GameID,
		PlayerID: req.PlayerID,
		Username: req.Username,
		Team:     game.Team(req.Team),
	}

	gameState, err := h.gameService.JoinGame(joinReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameState)
}

// GetGameState handles the request to get the current state of a game
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

// RevealCard handles the request to reveal a card
func (h *GameHandler) RevealCard(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GameID   string `json:"game_id"`
		CardID   string `json:"card_id"`
		PlayerID string `json:"player_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	revealReq := game.RevealCardRequest{
		GameID:   req.GameID,
		CardID:   req.CardID,
		PlayerID: req.PlayerID,
	}

	gameState, err := h.gameService.RevealCard(revealReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameState)
}

// SetSpymaster handles the request to set a player as a spymaster
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

// EndTurn handles the request to end the current team's turn
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

// RegisterRoutes registers all game routes
func (h *GameHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/game/start", h.StartGame).Methods("POST")
	r.HandleFunc("/api/game/state", h.GetGameState).Methods("GET")
	r.HandleFunc("/api/game/join", h.JoinGame).Methods("POST")
	r.HandleFunc("/api/game/reveal", h.RevealCard).Methods("POST")
	r.HandleFunc("/api/game/set-spymaster", h.SetSpymaster).Methods("POST")
	r.HandleFunc("/api/game/end-turn", h.EndTurn).Methods("POST")
}
