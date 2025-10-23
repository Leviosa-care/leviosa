package priceHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// done
	// todo
	// Create Price for a specific product
	router.HandleFunc("POST /admin/products/{id}/prices", middleware.EnableCORS(h.CreatePrice))
	// List Prices for a specific product
	// GET /admin/products/{productId}/prices
	router.HandleFunc("GET /admin/products/{id}/prices", middleware.EnableCORS(h.GetPricesByProductID))
	// Get (Retrieve) a specific Price by its ID (could be top-level or nested)
	// GET /admin/prices/{priceId}
	router.HandleFunc("GET /admin/prices/{id}", middleware.EnableCORS(h.GetPrice))
	// Update a specific Price (e.g., set inactive, update metadata)
	// PATCH /admin/prices/{priceId}
	router.HandleFunc("PATCH /admin/prices/{id}", middleware.EnableCORS(h.UpdatePrice))
	// to remove
}

// NOTE:
// A Product can have multiple Prices (e.g. one-time purchase, monthly subscription, discounted price,
// region-specific price, etc.).
