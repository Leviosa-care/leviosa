package bookingHandler

import (
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) GetPartnerBookings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract partner ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID in path"), http.StatusBadRequest)
		return
	}

	partnerID, err := uuid.Parse(pathParts[1])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
		return
	}

	// Parse query parameters
	filter := parseBookingFilter(r)

	// Call service to get partner bookings
	bookings, err := h.svc.GetPartnerBookings(ctx, partnerID, filter)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	responses := convertBookingsToResponses(bookings)
	httpx.RespondWithJSON(w, responses, http.StatusOK)
}