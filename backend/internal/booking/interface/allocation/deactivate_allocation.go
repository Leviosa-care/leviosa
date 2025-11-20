package allocationHandler

import (
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) DeactivateAllocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing deactivate allocation request",
		"operation", "deactivate_allocation",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Extract allocation ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		logger.WarnContext(ctx, "Handler: Invalid allocation ID in path",
			"error", "invalid allocation ID in path",
			"path", r.URL.Path,
			"operation", "deactivate_allocation",
			"method", r.Method)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid allocation ID in path"), http.StatusBadRequest)
		return
	}

	allocationID, err := uuid.Parse(pathParts[1])
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid allocation ID format",
			"error", err,
			"raw_allocation_id", pathParts[1],
			"operation", "deactivate_allocation",
			"method", r.Method)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid allocation ID format"), http.StatusBadRequest)
		return
	}

	// Call service to deactivate allocation
	err = h.svc.DeactivateAllocation(ctx, allocationID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "deactivate allocation for partner")
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
	}

	logger.InfoContext(ctx, "Handler: Allocation deactivated successfully",
		"allocation_id", allocationID,
		"room_id", allocation.RoomID,
		"user_id", allocation.UserID,
		"is_active", allocation.IsActive,
		"operation", "deactivate_allocation")

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
