package bookingHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetUpcomingBookings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	filter := parseBookingFilter(r)

	// Call service to get upcoming bookings
	bookings, err := h.svc.GetUpcomingBookings(ctx, filter)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	responses := convertBookingsToResponses(bookings)
	httpx.RespondWithJSON(w, responses, http.StatusOK)
}