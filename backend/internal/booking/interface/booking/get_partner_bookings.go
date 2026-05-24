package bookingHandler

import (
	"errors"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) GetPartnerBookings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	partnerID, err := uuid.Parse(r.PathValue("partnerId"))
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
		return
	}

	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if !ok {
		httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	if sessionInfo.UserID != partnerID && !sessionInfo.Role.IsAtLeast(identity.Administrator) {
		httpx.RespondWithError(w, errs.ErrForbidden, http.StatusForbidden)
		return
	}

	filter := parseBookingFilter(r)

	bookings, err := h.svc.GetPartnerBookingsEnriched(ctx, partnerID, filter)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	httpx.RespondWithJSON(w, bookings, http.StatusOK)
}