package availabilityHandler

import (
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) GetAvailability(w http.ResponseWriter, r *http.Request) {
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

	// Call service to get availability
	availability, err := h.svc.GetAvailability(ctx, availabilityID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get availability")
		return
	}

	// Convert to response DTO
	response := domain.AvailabilityResponse{
		ID:                availability.ID,
		UserID:            availability.UserID,
		RoomID:            availability.RoomID,
		StartTime:         availability.StartTime,
		EndTime:           availability.EndTime,
		MaxCapacity:       availability.MaxCapacity,
		Status:            availability.Status,
		ServiceType:       availability.ServiceType,
		PriceCents:        availability.PriceCents,
		Notes:             availability.Notes,
		RecurrencePattern: availability.RecurrencePattern,
		CreatedAt:         availability.CreatedAt,
		UpdatedAt:         availability.UpdatedAt,
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
