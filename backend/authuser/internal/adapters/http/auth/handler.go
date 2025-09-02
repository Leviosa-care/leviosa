package aggregatorHandler

import (
	"net/http"

	"github.com/Leviosa-care/authuser/internal/ports"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
}

type handler struct {
	svc ports.AuthAggregatorService
}

func New(svc ports.AuthAggregatorService) Handler {
	return &handler{svc: svc}
}
