package bookingHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

type claimBookingsRequest struct {
	ClientID string `json:"client_id"`
	Email    string `json:"email"`
}

type claimBookingsResponse struct {
	Claimed int `json:"claimed"`
}

// ClaimBookings handles POST /bookings/claim.
// It links all guest bookings matching the given email to the specified client.
// This endpoint is called internally by authuser after account creation.
func (h *handler) ClaimBookings(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httpx.RespondWithError(w, errors.New("unsupported media type"), http.StatusUnsupportedMediaType)
		return
	}

	ctx := r.Context()

	var payload claimBookingsRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	clientID, err := uuid.Parse(payload.ClientID)
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid client_id format"), http.StatusBadRequest)
		return
	}

	if payload.Email == "" {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("email is required"), http.StatusBadRequest)
		return
	}

	claimed, err := h.svc.ClaimBookings(ctx, clientID, payload.Email)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	httpx.RespondWithJSON(w, claimBookingsResponse{Claimed: claimed}, http.StatusOK)
}
