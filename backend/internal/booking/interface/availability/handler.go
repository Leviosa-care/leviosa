package availabilityHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	CreateAvailability(w http.ResponseWriter, r *http.Request)
	CreateRecurringAvailability(w http.ResponseWriter, r *http.Request)
	GetAvailability(w http.ResponseWriter, r *http.Request)
	GetPartnerAvailabilities(w http.ResponseWriter, r *http.Request)
	GetAvailableSlots(w http.ResponseWriter, r *http.Request)
	UpdateAvailability(w http.ResponseWriter, r *http.Request)
	CancelAvailability(w http.ResponseWriter, r *http.Request)
	BlockAvailability(w http.ResponseWriter, r *http.Request)
	CheckAvailabilityConflict(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc    ports.AvailabilityService
	authmw auth.AuthMiddleware
}

func New(svc ports.AvailabilityService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}
