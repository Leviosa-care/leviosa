package bookingHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) UpdateBookingNotes(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httpx.RespondWithError(w, errors.New("unsupported media type: please send 'application/json'"), http.StatusUnsupportedMediaType)
		return
	}

	ctx := r.Context()

	// Extract booking ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid booking ID in path"), http.StatusBadRequest)
		return
	}

	bookingID, err := uuid.Parse(pathParts[1])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid booking ID format"), http.StatusBadRequest)
		return
	}

	// Parse request body
	var request domain.UpdateBookingNotesRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call service to update booking notes
	booking, err := h.svc.UpdateBookingNotes(ctx, bookingID, request.ClientNotes, request.PartnerNotes)
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

	// Convert to response DTO
	var cancellationReason *string
	if booking.CancellationReason != "" {
		cancellationReason = &booking.CancellationReason
	}

	response := domain.BookingResponse{
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
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}