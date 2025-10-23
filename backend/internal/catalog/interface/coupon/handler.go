package couponHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	CreateCoupon(w http.ResponseWriter, r *http.Request)
	GetCouponByID(w http.ResponseWriter, r *http.Request)
	GetCouponByStripeID(w http.ResponseWriter, r *http.Request)
	GetAllCoupons(w http.ResponseWriter, r *http.Request)
	GetValidCoupons(w http.ResponseWriter, r *http.Request)
	UpdateCoupon(w http.ResponseWriter, r *http.Request)
	DeactivateCoupon(w http.ResponseWriter, r *http.Request)
	DeleteCoupon(w http.ResponseWriter, r *http.Request)
	ValidateCoupon(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc ports.CouponService
}

func New(couponService ports.CouponService) Handler {
	return &handler{
		svc: couponService,
	}
}