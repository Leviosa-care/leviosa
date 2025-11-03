package imageHandler

import (
	"net/http"

	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin

	// === Admin-Only Endpoints ===
	// (All image management endpoints are admin-only)

	// Upload a new image (admin only)
	router.HandleFunc("POST "+UploadImageEndpoint, RequireAdmin(mw.EnableCORS(h.UploadImage)))

	// Remove an image (admin only)
	router.HandleFunc("DELETE "+RemoveImageEndpoint, RequireAdmin(mw.EnableCORS(h.RemoveImage)))

	// Set active image for a resource (admin only)
	router.HandleFunc("POST "+SetActiveImageEndpoint, RequireAdmin(mw.EnableCORS(h.SetActiveImage)))
}
