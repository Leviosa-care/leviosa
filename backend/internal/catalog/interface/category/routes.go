package categoryHandler

import (
	"net/http"

	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /categories", mw.EnableCORS(h.GetAllPublishedCategories))
	router.HandleFunc("GET /categories/{id}", mw.EnableCORS(h.GetCategoryByID))
	router.HandleFunc("GET /categories/", mw.EnableCORS(h.GetCategoryByID)) // just to fit router strict path matching
	router.HandleFunc("GET /admin/categories", mw.EnableCORS(h.GetAdminAllCategories))
	router.HandleFunc("POST /admin/categories", mw.EnableCORS(h.CreateCategory))
	router.HandleFunc("PATCH /admin/categories/{id}", mw.EnableCORS(h.ModifyCategory))
	router.HandleFunc("DELETE /admin/categories/{id}", mw.EnableCORS(h.RemoveCategory))
}
