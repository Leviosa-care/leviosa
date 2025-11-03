package promotionCodeHandler

import (
	"net/http"

	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin

	// === Public Endpoints (no authentication required) ===

	// Validate promotion code during checkout (public access)
	router.HandleFunc("POST "+ValidatePromotionCodeEndpoint, mw.EnableCORS(h.ValidatePromotionCode))

	// Get promotion code with associated coupon by code (public access)
	router.HandleFunc("GET "+GetPromotionCodeWithCouponEndpoint, mw.EnableCORS(h.GetPromotionCodeWithCoupon))

	// === Admin-Only Endpoints ===

	// Get all promotion codes (admin only)
	router.HandleFunc("GET "+GetAllPromotionCodesEndpoint, RequireAdmin(mw.EnableCORS(h.GetAllPromotionCodes)))

	// Get active promotion codes (admin only)
	router.HandleFunc("GET "+GetActivePromotionCodesEndpoint, RequireAdmin(mw.EnableCORS(h.GetActivePromotionCodes)))

	// Get promotion code by ID (admin only)
	router.HandleFunc("GET "+GetPromotionCodeByIDEndpoint, RequireAdmin(mw.EnableCORS(h.GetPromotionCodeByID)))

	// Get promotion code by code string (admin only)
	router.HandleFunc("GET "+GetPromotionCodeByCodeEndpoint, RequireAdmin(mw.EnableCORS(h.GetPromotionCodeByCode)))

	// Create promotion code (admin only)
	router.HandleFunc("POST "+CreatePromotionCodeEndpoint, RequireAdmin(mw.EnableCORS(h.CreatePromotionCode)))

	// Update promotion code (admin only)
	router.HandleFunc("PATCH "+UpdatePromotionCodeEndpoint, RequireAdmin(mw.EnableCORS(h.UpdatePromotionCode)))

	// Deactivate promotion code (admin only)
	router.HandleFunc("POST "+DeactivatePromotionCodeEndpoint, RequireAdmin(mw.EnableCORS(h.DeactivatePromotionCode)))

	// Delete promotion code (admin only)
	router.HandleFunc("DELETE "+DeletePromotionCodeEndpoint, RequireAdmin(mw.EnableCORS(h.DeletePromotionCode)))
}
