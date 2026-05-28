package bookingHandler

import (
	"errors"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

// LookupBooking handles unauthenticated public booking lookups via two paths:
//   - ?token=xxx  — verifies the signed booking token
//   - ?ref=<id>&email=<email> or ?ref=<id>&phone=<phone> — manual fallback
//
// The response never exposes decrypted guest contact fields.
func (h *handler) LookupBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := httpx.FormString(r, "token")
	ref := httpx.FormString(r, "ref")
	email := httpx.FormString(r, "email")
	phone := httpx.FormString(r, "phone")

	var response any
	var err error

	switch {
	case token != "":
		response, err = h.svc.LookupBookingByToken(ctx, token)
	case ref != "" && (email != "" || phone != ""):
		response, err = h.svc.LookupBookingByRefAndContact(ctx, ref, email, phone)
	default:
		httpx.RespondWithError(w, errs.NewInvalidValueErr("provide either a token, or a booking reference with email or phone"), http.StatusBadRequest)
		return
	}

	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrUnauthorized):
			statusCode = http.StatusUnauthorized
		case errors.Is(err, errs.ErrInvalidValue):
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
