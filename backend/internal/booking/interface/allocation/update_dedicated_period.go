package allocationHandler

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

func (h *handler) UpdateDedicatedPeriod(w http.ResponseWriter, r *http.Request) {
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

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing update dedicated period request",
		"operation", "update_dedicated_period",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Extract allocation ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		logger.WarnContext(ctx, "Handler: Invalid allocation ID in path",
			"error", "invalid allocation ID in path",
			"path", r.URL.Path,
			"operation", "update_dedicated_period",
			"method", r.Method)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid allocation ID in path"), http.StatusBadRequest)
		return
	}

	allocationID, err := uuid.Parse(pathParts[1])
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid allocation ID format",
			"error", err,
			"raw_allocation_id", pathParts[1],
			"operation", "update_dedicated_period",
			"method", r.Method)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid allocation ID format"), http.StatusBadRequest)
		return
	}

	// Parse request body
	var request domain.UpdateDedicatedAllocationRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "update_dedicated_period",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	request.ID = allocationID

	// Call service to update dedicated period
	allocation, err := h.svc.UpdateDedicatedPeriod(ctx, &request)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrInvalidInput):
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrUniqueViolation):
			statusCode = http.StatusConflict
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
		default:
			statusCode = http.StatusInternalServerError
		}

		if statusCode >= 500 {
			logger.ErrorContext(ctx, "Handler: Update dedicated period failed",
				"error", err,
				"allocation_id", allocationID,
				"operation", "update_dedicated_period")
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

	logger.InfoContext(ctx, "Handler: Dedicated period updated successfully",
		"allocation_id", allocationID,
		"room_id", allocation.RoomID,
		"user_id", allocation.UserID,
		"start_date", request.StartDate,
		"end_date", request.EndDate,
		"operation", "update_dedicated_period")

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
