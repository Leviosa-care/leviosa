package categoryHandler

import (
	"net/http"

	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin

	// === Public Endpoints (no authentication required) ===

	// Retrieves all published categories visible to the public
	router.HandleFunc("GET "+GetAllPublishedCategoriesEndpoint, mw.EnableCORS(h.GetAllPublishedCategories))

	// Retrieves a specific category by ID (public access)
	router.HandleFunc("GET "+GetCategoryByIDEndpoint, mw.EnableCORS(h.GetCategoryByID))

	// Workaround for strict path matching
	router.HandleFunc("GET "+CategoriesBasePath+"/", mw.EnableCORS(h.GetCategoryByID))

	// === Admin-Only Endpoints ===

	// Retrieves all categories including unpublished/draft categories (admin only)
	router.HandleFunc("GET "+GetAdminAllCategoriesEndpoint, RequireAdmin(mw.EnableCORS(h.GetAdminAllCategories)))

	// Creates a new category (admin only)
	router.HandleFunc("POST "+CreateCategoryEndpoint, RequireAdmin(mw.EnableCORS(h.CreateCategory)))

	// Modifies an existing category by ID (admin only)
	router.HandleFunc("PATCH "+ModifyCategoryEndpoint, RequireAdmin(mw.EnableCORS(h.ModifyCategory)))

	// Removes a category by ID (admin only)
	router.HandleFunc("DELETE "+RemoveCategoryEndpoint, RequireAdmin(mw.EnableCORS(h.RemoveCategory)))
}
