package api

import (
	"net/http"
)

// RegisterHandlers registers all API handlers to the given router
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/api/health", healthCheckHandler)
	// More handlers will be added later
}

// healthCheckHandler provides a simple health check endpoint
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
