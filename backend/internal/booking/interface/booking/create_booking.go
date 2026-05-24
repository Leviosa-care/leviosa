package bookingHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/Leviosa-care/leviosa/backend/internal/common/validation"
)

func (h *handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httpx.RespondWithError(w, errors.New("unsupported media type: please send 'application/json'"), http.StatusUnsupportedMediaType)
		return
	}

	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Parse request body
	var request domain.CreateBookingRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Validate guest contact format when provided
	if request.GuestEmail != "" {
		if err := validation.ValidateEmail(request.GuestEmail); err != nil {
			httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("guest_email: %v", err)), http.StatusBadRequest)
			return
		}
	}
	if request.GuestPhone != "" {
		if err := validation.ValidatePhone(request.GuestPhone); err != nil {
			httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("guest_phone: %v", err)), http.StatusBadRequest)
			return
		}
	}

	// Call service to create booking
	booking, err := h.svc.CreateBooking(ctx, request.AvailabilityID, request.ClientID, request.ProductID, request.SlotStartTime, request.ClientNotes, request.GuestFirstName, request.GuestLastName, request.GuestEmail, request.GuestPhone)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrUniqueViolation):
			statusCode = http.StatusConflict
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

	logger.InfoContext(ctx, "Handler: Booking created successfully",
		"booking_id", booking.ID,
		"availability_id", booking.AvailabilityID,
		"is_guest", booking.IsGuestBooking(),
		"operation", "create_booking")

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}