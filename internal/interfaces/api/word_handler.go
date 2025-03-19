package api

import (
	gameservice "codenames-game/internal/usecase/game"
	"encoding/json"
	"net/http"
)

// WordHandler handles HTTP requests for word operations
type WordHandler struct {
	gameService gameservice.Service
}

// NewWordHandler creates a new word handler
func NewWordHandler(gs gameservice.Service) *WordHandler {
	return &WordHandler{
		gameService: gs,
	}
}

// GetWords returns all words
func (h *WordHandler) GetWords(w http.ResponseWriter, r *http.Request) {
	words, err := h.gameService.GetAllWords()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Words []string `json:"words"`
		Count int      `json:"count"`
	}{
		Words: words,
		Count: len(words),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AddWord adds a new word
func (h *WordHandler) AddWord(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Word string `json:"word"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Word == "" {
		http.Error(w, "Word cannot be empty", http.StatusBadRequest)
		return
	}

	if err := h.gameService.AddNewWord(req.Word); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"word":   req.Word,
	})
}

// DeleteWord deactivates a word
func (h *WordHandler) DeleteWord(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Word string `json:"word"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Word == "" {
		http.Error(w, "Word cannot be empty", http.StatusBadRequest)
		return
	}

	if err := h.gameService.DeleteExistingWord(req.Word); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"word":   req.Word,
	})
}
