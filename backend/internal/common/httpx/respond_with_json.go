package httpx

import (
	"encoding/json"
	"log"
	"net/http"
)

// Helper to respond with JSON data (e.g., for validation errors from errsx.Map)
func RespondWithJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
