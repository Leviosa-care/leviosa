package availabilityHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
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

	// Extract availability ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid availability ID in path"), http.StatusBadRequest)
		return
	}

	availabilityID, err := uuid.Parse(pathParts[1])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid availability ID format"), http.StatusBadRequest)
		return
	}

	// Parse request body
	var request domain.UpdateAvailabilityRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	request.ID = availabilityID

	// Call service to update availability
	availability, err := h.svc.UpdateAvailability(ctx, &request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "update availability")
		return
	}

	// Convert to response DTO
	response := domain.AvailabilityResponse{
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
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
