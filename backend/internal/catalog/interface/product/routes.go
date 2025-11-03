package productHandler

import (
	"net/http"

	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin

	// === Public Endpoints (no authentication required) ===

	// Retrieves all published products visible to the public
	router.HandleFunc("GET "+GetAllPublishedProductsEndpoint, mw.EnableCORS(h.GetAllPublishedProducts))

	// Retrieves a specific product by ID (public access)
	router.HandleFunc("GET "+GetProductByIDEndpoint, mw.EnableCORS(h.GetProductByID))

	// === Admin-Only Endpoints ===

	// Retrieves all products including unpublished/draft products (admin only)
	router.HandleFunc("GET "+GetAdminAllProductsEndpoint, RequireAdmin(mw.EnableCORS(h.GetAdminAllProducts)))

	// Creates a new product with price (admin only)
	router.HandleFunc("POST "+CreateProductWithPriceEndpoint, RequireAdmin(mw.EnableCORS(h.CreateProductWithPrice)))

	// Modifies an existing product by ID (admin only)
	router.HandleFunc("PATCH "+ModifyProductEndpoint, RequireAdmin(mw.EnableCORS(h.ModifyProduct)))

	// Removes a product by ID (admin only)
	router.HandleFunc("DELETE "+RemoveProductEndpoint, RequireAdmin(mw.EnableCORS(h.RemoveProduct)))
}
