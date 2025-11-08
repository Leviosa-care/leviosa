package bookingHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
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

	// Call service to create booking
	booking, err := h.svc.CreateBooking(ctx, request.AvailabilityID, request.ClientID, request.ClientNotes)
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
	response := domain.BookingResponse{
		ID:                 booking.ID,
		AvailabilityID:     booking.AvailabilityID,
		ClientID:           booking.ClientID,
		PartnerID:          booking.PartnerID,
		RoomID:             booking.RoomID,
		Status:             booking.Status,
		TotalPriceCents:    booking.TotalPriceCents,
		Currency:           booking.Currency,
		PaymentStatus:      booking.PaymentStatus,
		PaymentIntentID:    booking.PaymentIntentID,
		ClientNotes:        booking.ClientNotes,
		PartnerNotes:       booking.PartnerNotes,
		CancellationReason: booking.CancellationReason,
		CancelledAt:        booking.CancelledAt,
		CompletedAt:        booking.CompletedAt,
		CreatedAt:          booking.CreatedAt,
		UpdatedAt:          booking.UpdatedAt,
	}

	logger.InfoContext(ctx, "Handler: Booking created successfully",
		"booking_id", booking.ID,
		"availability_id", booking.AvailabilityID,
		"client_id", booking.ClientID,
		"operation", "create_booking")

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}

func (h *handler) GetBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract booking ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid booking ID in path"), http.StatusBadRequest)
		return
	}

	bookingID, err := uuid.Parse(pathParts[1])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid booking ID format"), http.StatusBadRequest)
		return
	}

	// Call service to get booking
	booking, err := h.svc.GetBooking(ctx, bookingID)
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
	response := domain.BookingResponse{
		ID:                 booking.ID,
		AvailabilityID:     booking.AvailabilityID,
		ClientID:           booking.ClientID,
		PartnerID:          booking.PartnerID,
		RoomID:             booking.RoomID,
		Status:             booking.Status,
		TotalPriceCents:    booking.TotalPriceCents,
		Currency:           booking.Currency,
		PaymentStatus:      booking.PaymentStatus,
		PaymentIntentID:    booking.PaymentIntentID,
		ClientNotes:        booking.ClientNotes,
		PartnerNotes:       booking.PartnerNotes,
		CancellationReason: booking.CancellationReason,
		CancelledAt:        booking.CancelledAt,
		CompletedAt:        booking.CompletedAt,
		CreatedAt:          booking.CreatedAt,
		UpdatedAt:          booking.UpdatedAt,
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

func (h *handler) GetClientBookings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract client ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid client ID in path"), http.StatusBadRequest)
		return
	}

	clientID, err := uuid.Parse(pathParts[1])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid client ID format"), http.StatusBadRequest)
		return
	}

	// Parse query parameters
	filter := parseBookingFilter(r)

	// Call service to get client bookings
	bookings, err := h.svc.GetClientBookings(ctx, clientID, filter)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	responses := convertBookingsToResponses(bookings)
	httpx.RespondWithJSON(w, responses, http.StatusOK)
}

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

func (h *handler) UpdateBookingNotes(w http.ResponseWriter, r *http.Request) {
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
	response := domain.BookingResponse{
		ID:                 booking.ID,
		AvailabilityID:     booking.AvailabilityID,
		ClientID:           booking.ClientID,
		PartnerID:          booking.PartnerID,
		RoomID:             booking.RoomID,
		Status:             booking.Status,
		TotalPriceCents:    booking.TotalPriceCents,
		Currency:           booking.Currency,
		PaymentStatus:      booking.PaymentStatus,
		PaymentIntentID:    booking.PaymentIntentID,
		ClientNotes:        booking.ClientNotes,
		PartnerNotes:       booking.PartnerNotes,
		CancellationReason: booking.CancellationReason,
		CancelledAt:        booking.CancelledAt,
		CompletedAt:        booking.CompletedAt,
		CreatedAt:          booking.CreatedAt,
		UpdatedAt:          booking.UpdatedAt,
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

func (h *handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
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
	var request domain.CancelBookingRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call service to cancel booking
	booking, err := h.svc.CancelBooking(ctx, bookingID, request.Reason)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrInvalidInput):
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Convert to response DTO
	response := domain.BookingResponse{
		ID:                 booking.ID,
		AvailabilityID:     booking.AvailabilityID,
		ClientID:           booking.ClientID,
		PartnerID:          booking.PartnerID,
		RoomID:             booking.RoomID,
		Status:             booking.Status,
		TotalPriceCents:    booking.TotalPriceCents,
		Currency:           booking.Currency,
		PaymentStatus:      booking.PaymentStatus,
		PaymentIntentID:    booking.PaymentIntentID,
		ClientNotes:        booking.ClientNotes,
		PartnerNotes:       booking.PartnerNotes,
		CancellationReason: booking.CancellationReason,
		CancelledAt:        booking.CancelledAt,
		CompletedAt:        booking.CompletedAt,
		CreatedAt:          booking.CreatedAt,
		UpdatedAt:          booking.UpdatedAt,
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

func (h *handler) CompleteBooking(w http.ResponseWriter, r *http.Request) {
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

	// Call service to complete booking
	booking, err := h.svc.CompleteBooking(ctx, bookingID)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrInvalidInput):
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Convert to response DTO
	response := domain.BookingResponse{
		ID:                 booking.ID,
		AvailabilityID:     booking.AvailabilityID,
		ClientID:           booking.ClientID,
		PartnerID:          booking.PartnerID,
		RoomID:             booking.RoomID,
		Status:             booking.Status,
		TotalPriceCents:    booking.TotalPriceCents,
		Currency:           booking.Currency,
		PaymentStatus:      booking.PaymentStatus,
		PaymentIntentID:    booking.PaymentIntentID,
		ClientNotes:        booking.ClientNotes,
		PartnerNotes:       booking.PartnerNotes,
		CancellationReason: booking.CancellationReason,
		CancelledAt:        booking.CancelledAt,
		CompletedAt:        booking.CompletedAt,
		CreatedAt:          booking.CreatedAt,
		UpdatedAt:          booking.UpdatedAt,
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

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
		case errors.Is(err, errs.ErrInvalidInput):
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Convert to response DTO
	response := domain.BookingResponse{
		ID:                 booking.ID,
		AvailabilityID:     booking.AvailabilityID,
		ClientID:           booking.ClientID,
		PartnerID:          booking.PartnerID,
		RoomID:             booking.RoomID,
		Status:             booking.Status,
		TotalPriceCents:    booking.TotalPriceCents,
		Currency:           booking.Currency,
		PaymentStatus:      booking.PaymentStatus,
		PaymentIntentID:    booking.PaymentIntentID,
		ClientNotes:        booking.ClientNotes,
		PartnerNotes:       booking.PartnerNotes,
		CancellationReason: booking.CancellationReason,
		CancelledAt:        booking.CancelledAt,
		CompletedAt:        booking.CompletedAt,
		CreatedAt:          booking.CreatedAt,
		UpdatedAt:          booking.UpdatedAt,
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

func (h *handler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
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
	var request domain.ProcessPaymentRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call service to process payment
	booking, err := h.svc.ProcessPayment(ctx, bookingID, request.PaymentIntentID)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrInvalidInput):
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Convert to response DTO
	response := domain.BookingResponse{
		ID:                 booking.ID,
		AvailabilityID:     booking.AvailabilityID,
		ClientID:           booking.ClientID,
		PartnerID:          booking.PartnerID,
		RoomID:             booking.RoomID,
		Status:             booking.Status,
		TotalPriceCents:    booking.TotalPriceCents,
		Currency:           booking.Currency,
		PaymentStatus:      booking.PaymentStatus,
		PaymentIntentID:    booking.PaymentIntentID,
		ClientNotes:        booking.ClientNotes,
		PartnerNotes:       booking.PartnerNotes,
		CancellationReason: booking.CancellationReason,
		CancelledAt:        booking.CancelledAt,
		CompletedAt:        booking.CompletedAt,
		CreatedAt:          booking.CreatedAt,
		UpdatedAt:          booking.UpdatedAt,
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

func (h *handler) RefundBooking(w http.ResponseWriter, r *http.Request) {
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

	// Call service to refund booking
	booking, err := h.svc.RefundBooking(ctx, bookingID)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrInvalidInput):
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Convert to response DTO
	response := domain.BookingResponse{
		ID:                 booking.ID,
		AvailabilityID:     booking.AvailabilityID,
		ClientID:           booking.ClientID,
		PartnerID:          booking.PartnerID,
		RoomID:             booking.RoomID,
		Status:             booking.Status,
		TotalPriceCents:    booking.TotalPriceCents,
		Currency:           booking.Currency,
		PaymentStatus:      booking.PaymentStatus,
		PaymentIntentID:    booking.PaymentIntentID,
		ClientNotes:        booking.ClientNotes,
		PartnerNotes:       booking.PartnerNotes,
		CancellationReason: booking.CancellationReason,
		CancelledAt:        booking.CancelledAt,
		CompletedAt:        booking.CompletedAt,
		CreatedAt:          booking.CreatedAt,
		UpdatedAt:          booking.UpdatedAt,
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

// Helper functions
func parseBookingFilter(r *http.Request) ports.BookingFilter {
	filter := ports.BookingFilter{}

	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		status := domain.BookingStatus(statusStr)
		filter.Status = &status
	}

	if paymentStatusStr := r.URL.Query().Get("payment_status"); paymentStatusStr != "" {
		paymentStatus := domain.PaymentStatus(paymentStatusStr)
		filter.PaymentStatus = &paymentStatus
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = &limit
		}
	}

	return filter
}

func convertBookingsToResponses(bookings []*domain.Booking) []domain.BookingResponse {
	var responses []domain.BookingResponse
	for _, booking := range bookings {
		responses = append(responses, domain.BookingResponse{
			ID:                 booking.ID,
			AvailabilityID:     booking.AvailabilityID,
			ClientID:           booking.ClientID,
			PartnerID:          booking.PartnerID,
			RoomID:             booking.RoomID,
			Status:             booking.Status,
			TotalPriceCents:    booking.TotalPriceCents,
			Currency:           booking.Currency,
			PaymentStatus:      booking.PaymentStatus,
			PaymentIntentID:    booking.PaymentIntentID,
			ClientNotes:        booking.ClientNotes,
			PartnerNotes:       booking.PartnerNotes,
			CancellationReason: booking.CancellationReason,
			CancelledAt:        booking.CancelledAt,
			CompletedAt:        booking.CompletedAt,
			CreatedAt:          booking.CreatedAt,
			UpdatedAt:          booking.UpdatedAt,
		})
	}
	return responses
}
