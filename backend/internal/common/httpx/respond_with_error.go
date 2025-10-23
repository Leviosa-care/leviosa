package httpx

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// Helper to respond with structured error (assuming errsx.Map can be marshaled)
func RespondWithError(w http.ResponseWriter, err error, status int) {
	errorResponse := map[string]string{"error": err.Error()}

	// Use a buffer to capture JSON output before writing to the response writer.
	// This allows handling encoding errors before headers are sent.
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if encodeErr := encoder.Encode(errorResponse); encodeErr != nil { // Corrected: encodeErr != nil
		log.Printf("Error encoding error response: %v", encodeErr)
		// If encoding fails, fallback to a plain text error without setting headers again.
		// http.Error will handle setting the status and writing the body.
		http.Error(w, "Internal server error: Failed to encode error response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)               // Set the intended status code
	_, writeErr := w.Write(buf.Bytes()) // Write the buffered JSON
	if writeErr != nil {
		log.Printf("Error writing JSON response: %v", writeErr)
		// At this point, headers are already sent, so we can only log.
	}
}
