package imageHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// done
	router.HandleFunc("POST /admin/images", middleware.EnableCORS(h.UploadImage))
	router.HandleFunc("DELETE /admin/images", middleware.EnableCORS(h.RemoveImage))
	router.HandleFunc("POST /admin/images/set-active", middleware.EnableCORS(h.SetActiveImage))
	// todo
	// router.HandleFunc("GET /admin/images", middleware.EnableCORS(h.UploadImage))
	// router.HandleFunc("PATCH /admin/images", middleware.EnableCORS(h.UploadImage))
}
