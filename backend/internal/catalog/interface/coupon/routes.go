package couponHandler

import (
	"net/http"

	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin

	// === Public Endpoints (no authentication required) ===

	// Validate coupon code during checkout (public access)
	router.HandleFunc("POST "+ValidateCouponEndpoint, mw.EnableCORS(h.ValidateCoupon))

	// Get all valid/active coupons (public access)
	router.HandleFunc("GET "+GetValidCouponsEndpoint, mw.EnableCORS(h.GetValidCoupons))

	// === Admin-Only Endpoints ===

	// Get all coupons (admin only)
	router.HandleFunc("GET "+GetAllCouponsEndpoint, RequireAdmin(mw.EnableCORS(h.GetAllCoupons)))

	// Get coupon by ID (admin only)
	router.HandleFunc("GET "+GetCouponByIDEndpoint, RequireAdmin(mw.EnableCORS(h.GetCouponByID)))

	// Get coupon by Stripe ID (admin only)
	router.HandleFunc("GET "+GetCouponByStripeIDEndpoint, RequireAdmin(mw.EnableCORS(h.GetCouponByStripeID)))

	// Create coupon (admin only)
	router.HandleFunc("POST "+CreateCouponEndpoint, RequireAdmin(mw.EnableCORS(h.CreateCoupon)))

	// Update coupon (admin only)
	router.HandleFunc("PATCH "+UpdateCouponEndpoint, RequireAdmin(mw.EnableCORS(h.UpdateCoupon)))

	// Deactivate coupon (admin only)
	router.HandleFunc("POST "+DeactivateCouponEndpoint, RequireAdmin(mw.EnableCORS(h.DeactivateCoupon)))

	// Delete coupon (admin only)
	router.HandleFunc("DELETE "+DeleteCouponEndpoint, RequireAdmin(mw.EnableCORS(h.DeleteCoupon)))
}

