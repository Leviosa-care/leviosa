package categoryHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *handler) Healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hit server")
	message := struct {
		Message string `json:"message"`
	}{
		Message: "Hello world",
	}
	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
