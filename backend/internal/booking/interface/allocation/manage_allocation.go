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

	// Extract allocation ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid allocation ID in path"), http.StatusBadRequest)
		return
	}

	allocationID, err := uuid.Parse(pathParts[1])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid allocation ID format"), http.StatusBadRequest)
		return
	}

	// Parse request body
	var request domain.UpdateDedicatedAllocationRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call service to update dedicated period
	allocation, err := h.svc.UpdateDedicatedPeriod(ctx, allocationID, request.StartDate, request.EndDate)
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
		CreatedAt:      allocation.CreatedAt,
		UpdatedAt:      allocation.UpdatedAt,
	}

	logger.InfoContext(ctx, "Handler: Dedicated period updated successfully",
		"allocation_id", allocationID,
		"operation", "update_dedicated_period")

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

func (h *handler) DeactivateAllocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract allocation ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid allocation ID in path"), http.StatusBadRequest)
		return
	}

	allocationID, err := uuid.Parse(pathParts[1])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid allocation ID format"), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing deactivate allocation request",
		"allocation_id", allocationID,
		"operation", "deactivate_allocation")

	// Call service to deactivate allocation
	err = h.svc.DeactivateAllocation(ctx, allocationID)
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

		if statusCode >= 500 {
			logger.ErrorContext(ctx, "Handler: Deactivate allocation failed",
				"error", err,
				"allocation_id", allocationID,
				"operation", "deactivate_allocation")
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Get updated allocation to return
	allocation, err := h.svc.GetAllocation(ctx, allocationID)
	if err != nil {
		logger.ErrorContext(ctx, "Handler: Failed to retrieve deactivated allocation",
			"error", err,
			"allocation_id", allocationID,
			"operation", "deactivate_allocation")
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
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
		CreatedAt:      allocation.CreatedAt,
		UpdatedAt:      allocation.UpdatedAt,
	}

	logger.InfoContext(ctx, "Handler: Allocation deactivated successfully",
		"allocation_id", allocationID,
		"operation", "deactivate_allocation")

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
