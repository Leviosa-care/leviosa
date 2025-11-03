package promotionCodeHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	CreatePromotionCode(w http.ResponseWriter, r *http.Request)
	GetPromotionCodeByID(w http.ResponseWriter, r *http.Request)
	GetPromotionCodeByCode(w http.ResponseWriter, r *http.Request)
	GetAllPromotionCodes(w http.ResponseWriter, r *http.Request)
	GetActivePromotionCodes(w http.ResponseWriter, r *http.Request)
	UpdatePromotionCode(w http.ResponseWriter, r *http.Request)
	DeactivatePromotionCode(w http.ResponseWriter, r *http.Request)
	DeletePromotionCode(w http.ResponseWriter, r *http.Request)
	ValidatePromotionCode(w http.ResponseWriter, r *http.Request)
	GetPromotionCodeWithCoupon(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc    ports.PromotionCodeService
	authmw auth.AuthMiddleware
}

func New(promotionCodeService ports.PromotionCodeService, authmw auth.AuthMiddleware) Handler {
	return &handler{
		svc:    promotionCodeService,
		authmw: authmw,
	}
}

