package promotionCodeHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
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
	svc ports.PromotionCodeService
}

func New(promotionCodeService ports.PromotionCodeService) Handler {
	return &handler{
		svc: promotionCodeService,
	}
}

