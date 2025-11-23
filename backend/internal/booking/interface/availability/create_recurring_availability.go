package availabilityHandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) CreateRecurringAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if !ok {
		logger.WarnContext(ctx, "Handler: Missing session info",
			"operation", "create_recurring_availability",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("session information from context required"), http.StatusUnauthorized)
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

	request.UserID = sessionInfo.UserID

	// Call service to create recurring availability
	availability, err := h.svc.CreateRecurringAvailability(ctx, &request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "create recurring availability")
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
		// ParentID:          availability.ParentID,
		CreatedAt: availability.CreatedAt,
		UpdatedAt: availability.UpdatedAt,
	}

	logger.InfoContext(ctx, "Handler: Recurring availability created successfully",
		"availability_id", availability.ID,
		"partner_id", availability.UserID,
		"room_id", availability.RoomID,
		"pattern", request.Pattern,
		"operation", "create_recurring_availability")

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}
