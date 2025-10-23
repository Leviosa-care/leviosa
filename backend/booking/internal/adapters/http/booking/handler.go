package bookingHandler

import (
	"net/http"

	"github.com/Leviosa-care/booking/internal/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	CreateBooking(w http.ResponseWriter, r *http.Request)
	GetBooking(w http.ResponseWriter, r *http.Request)
	GetClientBookings(w http.ResponseWriter, r *http.Request)
	GetPartnerBookings(w http.ResponseWriter, r *http.Request)
	GetUpcomingBookings(w http.ResponseWriter, r *http.Request)
	UpdateBookingNotes(w http.ResponseWriter, r *http.Request)
	CancelBooking(w http.ResponseWriter, r *http.Request)
	CompleteBooking(w http.ResponseWriter, r *http.Request)
	MarkNoShow(w http.ResponseWriter, r *http.Request)
	ProcessPayment(w http.ResponseWriter, r *http.Request)
	RefundBooking(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc    ports.BookingService
	authmw auth.AuthMiddleware
}

func New(svc ports.BookingService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}