package smsHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /categories", middleware.EnableCORS(h.WelcomeUser))
}
