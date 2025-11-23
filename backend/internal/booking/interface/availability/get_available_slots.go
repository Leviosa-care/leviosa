package availabilityHandler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Parse query parameters for filtering
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
	if roomIDStr := r.URL.Query().Get("room_id"); roomIDStr != "" {
		if roomID, err := uuid.Parse(roomIDStr); err == nil {
			filter.RoomID = &roomID
		}
	}
	if partnerIDStr := r.URL.Query().Get("partner_id"); partnerIDStr != "" {
		if partnerID, err := uuid.Parse(partnerIDStr); err == nil {
			filter.UserID = &partnerID
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	// Call service to get available slots
	availabilities, err := h.svc.GetAvailableSlots(ctx, filter)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get available slots")
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
			// CurrentBookings:   availability.CurrentBookings,
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
