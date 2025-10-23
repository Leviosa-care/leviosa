package couponHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// Public routes for coupon validation
	router.HandleFunc("POST /coupons/validate", middleware.EnableCORS(h.ValidateCoupon))
	router.HandleFunc("GET /coupons/valid", middleware.EnableCORS(h.GetValidCoupons))
	
	// Admin routes for coupon management
	router.HandleFunc("GET /admin/coupons", middleware.EnableCORS(h.GetAllCoupons))
	router.HandleFunc("GET /admin/coupons/valid", middleware.EnableCORS(h.GetValidCoupons))
	router.HandleFunc("GET /admin/coupons/{id}", middleware.EnableCORS(h.GetCouponByID))
	router.HandleFunc("GET /admin/coupons/stripe/{stripeId}", middleware.EnableCORS(h.GetCouponByStripeID))
	router.HandleFunc("POST /admin/coupons", middleware.EnableCORS(h.CreateCoupon))
	router.HandleFunc("PATCH /admin/coupons/{id}", middleware.EnableCORS(h.UpdateCoupon))
	router.HandleFunc("POST /admin/coupons/{id}/deactivate", middleware.EnableCORS(h.DeactivateCoupon))
	router.HandleFunc("DELETE /admin/coupons/{id}", middleware.EnableCORS(h.DeleteCoupon))
}