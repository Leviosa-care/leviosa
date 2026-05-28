package bookingHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

// cancelBookingResponse builds a BookingResponse DTO from a domain Booking.
func cancelBookingResponse(booking *domain.Booking) domain.BookingResponse {
	var cancellationReason *string
	if booking.CancellationReason != "" {
		cancellationReason = &booking.CancellationReason
	}

	return domain.BookingResponse{
		ID:                 booking.ID,
		AvailabilityID:     booking.AvailabilityID,
		ClientID:           booking.ClientID,
		PartnerID:          booking.PartnerID,
		RoomID:             booking.RoomID,
		ProductID:          booking.ProductID,
		SlotStartTime:      booking.SlotStartTime,
		SlotEndTime:        booking.SlotEndTime,
		Status:             booking.Status,
		TotalPriceCents:    booking.TotalPriceCents,
		Currency:           booking.Currency,
		PaymentStatus:      booking.PaymentStatus,
		PaymentIntentID:    booking.PaymentIntentID,
		ClientNotes:        booking.ClientNotes,
		PartnerNotes:       booking.PartnerNotes,
		CancellationReason: cancellationReason,
		CancelledAt:        booking.CancelledAt,
		CompletedAt:        booking.CompletedAt,
		GuestFirstName:     booking.GuestFirstName,
		GuestLastName:      booking.GuestLastName,
		GuestEmail:         booking.GuestEmail,
		GuestPhone:         booking.GuestPhone,
		Token:              booking.Token,
	}
}

// decodeCancelRequest reads and validates the cancel request body.
func decodeCancelRequest(r *http.Request) (domain.CancelBookingRequest, error) {
	var req domain.CancelBookingRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		return req, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err))
	}
	return req, nil
}

// mapCancelError maps cancellation errors to appropriate HTTP status codes.
func mapCancelError(err error) int {
	switch {
	case errors.Is(err, errs.ErrRepositoryNotFound):
		return http.StatusNotFound
	case errors.Is(err, errs.ErrInvalidInput), errors.Is(err, errs.ErrInvalidValue):
		return http.StatusBadRequest
	case errors.Is(err, errs.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, errs.ErrCancellationWindowClosed):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

func (h *handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httpx.RespondWithError(w, errors.New("unsupported media type: please send 'application/json'"), http.StatusUnsupportedMediaType)
		return
	}

	ctx := r.Context()

	bookingID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid booking ID format"), http.StatusBadRequest)
		return
	}

	request, err := decodeCancelRequest(r)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	booking, err := h.svc.CancelBooking(ctx, bookingID, request.Reason)
	if err != nil {
		httpx.RespondWithError(w, err, mapCancelError(err))
		return
	}

	httpx.RespondWithJSON(w, cancelBookingResponse(booking), http.StatusOK)
}

// CancelBookingPublic handles token-based booking cancellation for unauthenticated
// guests. The booking token is passed as a query parameter (?token=xxx).
// The response uses the public DTO (no internal IDs or guest contact fields).
func (h *handler) CancelBookingPublic(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httpx.RespondWithError(w, errors.New("unsupported media type: please send 'application/json'"), http.StatusUnsupportedMediaType)
		return
	}

	token := httpx.FormString(r, "token")
	if token == "" {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("missing booking token"), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	bookingID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid booking ID format"), http.StatusBadRequest)
		return
	}

	request, err := decodeCancelRequest(r)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	booking, err := h.svc.CancelBookingByToken(ctx, bookingID, token, request.Reason)
	if err != nil {
		httpx.RespondWithError(w, err, mapCancelError(err))
		return
	}

	httpx.RespondWithJSON(w, publicCancelResponse(booking), http.StatusOK)
}

// publicCancelResponse converts a cancelled booking to a public response DTO.
func publicCancelResponse(booking *domain.Booking) domain.PublicBookingLookupResponse {
	return domain.PublicBookingLookupResponse{
		ID:              booking.ID,
		SlotStartTime:   booking.SlotStartTime,
		SlotEndTime:     booking.SlotEndTime,
		Status:          booking.Status,
		TotalPriceCents: booking.TotalPriceCents,
		Currency:        booking.Currency,
		PaymentStatus:   booking.PaymentStatus,
	}
}