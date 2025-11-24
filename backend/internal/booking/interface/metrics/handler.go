package metricsHandler

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
)

// Handler handles HTTP requests for metrics
type Handler struct {
	svc    ports.MetricsService
	authmw auth.AuthMiddleware
}

// New creates a new metrics handler
func New(svc ports.MetricsService, authmw auth.AuthMiddleware) *Handler {
	return &Handler{
		svc:    svc,
		authmw: authmw,
	}
}
