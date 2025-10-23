package promotionCodeHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// Public routes for promotion code validation and lookup
	router.HandleFunc("POST /promotion-codes/validate", middleware.EnableCORS(h.ValidatePromotionCode))
	router.HandleFunc("GET /promotion-codes/code/{code}", middleware.EnableCORS(h.GetPromotionCodeWithCoupon))

	// Admin routes for promotion code management
	router.HandleFunc("GET /admin/promotion-codes", middleware.EnableCORS(h.GetAllPromotionCodes))
	router.HandleFunc("GET /admin/promotion-codes/active", middleware.EnableCORS(h.GetActivePromotionCodes))
	router.HandleFunc("GET /admin/promotion-codes/{id}", middleware.EnableCORS(h.GetPromotionCodeByID))
	router.HandleFunc("GET /admin/promotion-codes/code/{code}", middleware.EnableCORS(h.GetPromotionCodeByCode))
	router.HandleFunc("POST /admin/promotion-codes", middleware.EnableCORS(h.CreatePromotionCode))
	router.HandleFunc("PATCH /admin/promotion-codes/{id}", middleware.EnableCORS(h.UpdatePromotionCode))
	router.HandleFunc("POST /admin/promotion-codes/{id}/deactivate", middleware.EnableCORS(h.DeactivatePromotionCode))
	router.HandleFunc("DELETE /admin/promotion-codes/{id}", middleware.EnableCORS(h.DeletePromotionCode))
}
