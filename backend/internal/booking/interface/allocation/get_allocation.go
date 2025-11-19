package allocationHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) GetAllocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	_ = logger
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract allocation ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid allocation ID in path"), http.StatusBadRequest)
		return
	}

	allocationID, err := uuid.Parse(pathParts[1])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid allocation ID format"), http.StatusBadRequest)
		return
	}

	// Call service to get allocation
	allocation, err := h.svc.GetAllocation(ctx, allocationID)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
		default:
			statusCode = http.StatusInternalServerError
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Convert to response DTO
	response := domain.RoomAllocationResponse{
		ID:             allocation.ID,
		RoomID:         allocation.RoomID,
		UserID:         allocation.UserID,
		AllocationType: allocation.AllocationType,
		StartDate:      allocation.StartDate,
		EndDate:        allocation.EndDate,
		IsActive:       allocation.IsActive,
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
