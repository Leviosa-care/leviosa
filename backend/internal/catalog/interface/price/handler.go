package priceHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	CreatePrice(w http.ResponseWriter, r *http.Request)
	GetPrice(w http.ResponseWriter, r *http.Request)
	GetPricesByProductID(w http.ResponseWriter, r *http.Request)
	UpdatePrice(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc ports.PriceService
}

func New(services ports.PriceService) *handler {
	return &handler{
		svc: services,
	}
}
