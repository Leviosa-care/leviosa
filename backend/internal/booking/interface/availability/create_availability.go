package availabilityHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	availabilitySvc "github.com/Leviosa-care/leviosa/backend/internal/booking/application/availability"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) CreateAvailability(w http.ResponseWriter, r *http.Request) {
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

	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if !ok {
		logger.WarnContext(ctx, "Handler: Missing session info",
			"error", err,
			"operation", "complete_partner",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("session information from context required"), http.StatusUnauthorized)
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

	request.UserID = sessionInfo.UserID

	// Call service to create availability
	availability, err := h.svc.CreateAvailability(ctx, &request)
	if err != nil {
		// Check if error is InvalidDurationError and return structured response
		var invalidDurationErr *availabilitySvc.InvalidDurationError
		if errors.As(err, &invalidDurationErr) {
			logger.WarnContext(ctx, "Handler: Invalid availability duration",
				"requested_duration", invalidDurationErr.RequestedDuration,
				"valid_blocks_count", len(invalidDurationErr.ValidBlocks),
				"operation", "create_availability")
			httpx.RespondWithJSON(w, invalidDurationErr.ToJSON(), http.StatusBadRequest)
			return
		}

		httpx.RespondWithServiceError(w, logger, ctx, err, "create availability")
		return
	}

	// Set additional fields if provided
	if request.ServiceType != "" || request.PriceCents != nil || request.Notes != "" {
		availability.SetServiceDetails(request.ServiceType, request.PriceCents, request.Notes)

		// Update availability with service details
		updateRequest := &domain.UpdateAvailabilityRequest{
			ID:          availability.ID,
			ServiceType: &request.ServiceType,
			PriceCents:  request.PriceCents,
			Notes:       &request.Notes,
		}
		availability, err = h.svc.UpdateAvailability(ctx, updateRequest)
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
		ID:          availability.ID,
		UserID:      availability.UserID,
		RoomID:      availability.RoomID,
		StartTime:   availability.StartTime,
		EndTime:     availability.EndTime,
		MaxCapacity: availability.MaxCapacity,
		// CurrentBookings: availability.CurrentBookings,
		Status:      availability.Status,
		ServiceType: availability.ServiceType,
		PriceCents:  availability.PriceCents,
		Notes:       availability.Notes,
		// RecurrencePattern: availability.RecurrencePattern,
		// ParentID:  availability.ParentID,
		CreatedAt: availability.CreatedAt,
		UpdatedAt: availability.UpdatedAt,
	}

	logger.InfoContext(ctx, "Handler: Availability created successfully",
		"availability_id", availability.ID,
		"partner_id", availability.UserID,
		"room_id", availability.RoomID,
		"operation", "create_availability")

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}
