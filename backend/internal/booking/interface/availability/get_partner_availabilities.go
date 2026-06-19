package availabilityHandler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"

	"github.com/google/uuid"
)

func (h *handler) GetPartnerAvailabilities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract partner ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID in path"), http.StatusBadRequest)
		return
	}

	partnerID, err := uuid.Parse(pathParts[2])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
		return
	}

	// Parse query parameters
	filter := ports.AvailabilityFilter{}
	if startStr := r.URL.Query().Get("start_time"); startStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startStr); err == nil {
			filter.StartTime = &startTime
		}
	}
	if endStr := r.URL.Query().Get("end_time"); endStr != "" {
		if endTime, err := time.Parse(time.RFC3339, endStr); err == nil {
			filter.EndTime = &endTime
		}
	}
	if statusParams := r.URL.Query()["status"]; len(statusParams) > 0 {
		statuses := make([]domain.AvailabilityStatus, 0, len(statusParams))
		for _, statusStr := range statusParams {
			if statusStr != "" {
				statuses = append(statuses, domain.AvailabilityStatus(statusStr))
			}
		}
		if len(statuses) > 0 {
			filter.Status = statuses
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	// Call service to get partner availabilities
	availabilities, err := h.svc.GetPartnerAvailabilities(ctx, partnerID, filter)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get partner availabilities")
		return
	}

	// Convert to response DTOs
	var responses []domain.AvailabilityResponse
	for _, availability := range availabilities {
		responses = append(responses, domain.AvailabilityResponse{
			ID:          availability.ID,
			UserID:      availability.UserID,
			RoomID:      availability.RoomID,
			StartTime:   availability.StartTime,
			EndTime:     availability.EndTime,
			MaxCapacity: availability.MaxCapacity,
			// CurrentBookings not available on domain.Availability — requires join query
			Status:            availability.Status,
			ServiceType:       availability.ServiceType,
			PriceCents:        availability.PriceCents,
			Notes:             availability.Notes,
			RecurrencePattern: availability.RecurrencePattern,
			// ParentID:          availability.ParentID,
			CreatedAt: availability.CreatedAt,
			UpdatedAt: availability.UpdatedAt,
		})
	}

	httpx.RespondWithJSON(w, responses, http.StatusOK)
}
