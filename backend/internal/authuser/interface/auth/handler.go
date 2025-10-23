package aggregatorHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
}

type handler struct {
	svc    ports.AuthAggregatorService
	authmw auth.AuthMiddleware
}

func New(svc ports.AuthAggregatorService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}
