package bookingHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) MarkNoShow(w http.ResponseWriter, r *http.Request) {
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

	// Call service to mark as no-show
	booking, err := h.svc.MarkNoShow(ctx, bookingID)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrInvalidInput),
			errors.Is(err, domain.ErrCannotMarkCancelledAsNoShow),
			errors.Is(err, domain.ErrCannotMarkCompletedAsNoShow):
			statusCode = http.StatusBadRequest
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
		Token:              booking.Token,
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}