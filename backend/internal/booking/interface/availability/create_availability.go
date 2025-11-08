package availabilityHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) CreateAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Parse request body
	var request domain.CreateAvailabilityRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "create_availability")
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call service to create availability
	availability, err := h.svc.CreateAvailability(ctx, request.PartnerID, request.RoomID, request.StartTime, request.EndTime, request.MaxCapacity)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrUniqueViolation):
			statusCode = http.StatusConflict
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
		default:
			statusCode = http.StatusInternalServerError
		}

		logger.ErrorContext(ctx, "Handler: Create availability failed",
			"error", err,
			"partner_id", request.PartnerID,
			"room_id", request.RoomID,
			"operation", "create_availability")

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Set additional fields if provided
	if request.ServiceType != "" || request.PriceCents != nil || request.Notes != "" {
		availability.SetServiceDetails(request.ServiceType, request.PriceCents, request.Notes)

		// Update availability with service details
		availability, err = h.svc.UpdateAvailability(ctx, availability.ID, availability.StartTime, availability.EndTime, availability.ServiceType, availability.PriceCents, availability.Notes)
		if err != nil {
			logger.ErrorContext(ctx, "Handler: Failed to update availability with service details",
				"error", err,
				"availability_id", availability.ID,
				"operation", "create_availability")
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}
	}

	// Convert to response DTO
	response := domain.AvailabilityResponse{
		ID:              availability.ID,
		PartnerID:       availability.PartnerID,
		RoomID:          availability.RoomID,
		StartTime:       availability.StartTime,
		EndTime:         availability.EndTime,
		MaxCapacity:     availability.MaxCapacity,
		CurrentBookings: availability.CurrentBookings,
		Status:          availability.Status,
		ServiceType:     availability.ServiceType,
		PriceCents:      availability.PriceCents,
		Notes:           availability.Notes,
		RecurrenceRule:  availability.RecurrenceRule,
		ParentID:        availability.ParentID,
		CreatedAt:       availability.CreatedAt,
		UpdatedAt:       availability.UpdatedAt,
	}

	logger.InfoContext(ctx, "Handler: Availability created successfully",
		"availability_id", availability.ID,
		"partner_id", availability.PartnerID,
		"room_id", availability.RoomID,
		"operation", "create_availability")

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}

func (h *handler) CreateRecurringAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Parse request body
	var request domain.CreateRecurringAvailabilityRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "create_recurring_availability")
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call service to create recurring availability
	availability, err := h.svc.CreateRecurringAvailability(ctx, request.PartnerID, request.RoomID, request.StartTime, request.EndTime, request.MaxCapacity, request.Pattern)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrUniqueViolation):
			statusCode = http.StatusConflict
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
		default:
			statusCode = http.StatusInternalServerError
		}

		logger.ErrorContext(ctx, "Handler: Create recurring availability failed",
			"error", err,
			"partner_id", request.PartnerID,
			"room_id", request.RoomID,
			"operation", "create_recurring_availability")

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Set additional fields if provided
	if request.ServiceType != "" || request.PriceCents != nil || request.Notes != "" {
		availability.SetServiceDetails(request.ServiceType, request.PriceCents, request.Notes)

		// Update availability with service details
		availability, err = h.svc.UpdateAvailability(ctx, availability.ID, availability.StartTime, availability.EndTime, availability.ServiceType, availability.PriceCents, availability.Notes)
		if err != nil {
			logger.ErrorContext(ctx, "Handler: Failed to update recurring availability with service details",
				"error", err,
				"availability_id", availability.ID,
				"operation", "create_recurring_availability")
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}
	}

	// Convert to response DTO
	response := domain.AvailabilityResponse{
		ID:              availability.ID,
		PartnerID:       availability.PartnerID,
		RoomID:          availability.RoomID,
		StartTime:       availability.StartTime,
		EndTime:         availability.EndTime,
		MaxCapacity:     availability.MaxCapacity,
		CurrentBookings: availability.CurrentBookings,
		Status:          availability.Status,
		ServiceType:     availability.ServiceType,
		PriceCents:      availability.PriceCents,
		Notes:           availability.Notes,
		RecurrenceRule:  availability.RecurrenceRule,
		ParentID:        availability.ParentID,
		CreatedAt:       availability.CreatedAt,
		UpdatedAt:       availability.UpdatedAt,
	}

	logger.InfoContext(ctx, "Handler: Recurring availability created successfully",
		"availability_id", availability.ID,
		"partner_id", availability.PartnerID,
		"room_id", availability.RoomID,
		"pattern", request.Pattern,
		"operation", "create_recurring_availability")

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}
