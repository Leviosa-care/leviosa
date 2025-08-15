package settingsHandler

import "net/http"

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /settings/logo", h.AddLogo)
}
