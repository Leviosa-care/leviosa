package priceHandler

import (
	"net/http"

	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin

	// === Public Endpoints ===

	// Get all active prices for a product (public access)
	router.HandleFunc("GET "+GetPublicProductPricesEndpoint, mw.EnableCORS(h.GetPublicProductPrices))

	// === Admin-Only Endpoints ===
	// (All price management endpoints are admin-only)

	// Get a specific price by its ID (admin only)
	router.HandleFunc("GET "+GetPriceEndpoint, RequireAdmin(mw.EnableCORS(h.GetPrice)))

	// List all prices for a specific product (admin only)
	router.HandleFunc("GET "+GetPricesByProductIDEndpoint, RequireAdmin(mw.EnableCORS(h.GetPricesByProductID)))

	// Create a new price for a specific product (admin only)
	router.HandleFunc("POST "+CreatePriceEndpoint, RequireAdmin(mw.EnableCORS(h.CreatePrice)))

	// Update a specific price (e.g., set inactive, update metadata) (admin only)
	router.HandleFunc("PATCH "+UpdatePriceEndpoint, RequireAdmin(mw.EnableCORS(h.UpdatePrice)))
}
