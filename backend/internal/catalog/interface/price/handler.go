package priceHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	CreatePrice(w http.ResponseWriter, r *http.Request)
	GetPrice(w http.ResponseWriter, r *http.Request)
	GetPricesByProductID(w http.ResponseWriter, r *http.Request)
	GetPublicProductPrices(w http.ResponseWriter, r *http.Request)
	UpdatePrice(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc    ports.PriceService
	authmw auth.AuthMiddleware
}

func New(services ports.PriceService, authmw auth.AuthMiddleware) Handler {
	return &handler{
		svc:    services,
		authmw: authmw,
	}
}
