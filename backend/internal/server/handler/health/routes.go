package healthHandler

import "net/http"

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /health", h.CheckHealth)
}
