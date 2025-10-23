package allocationHandler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) GetAllocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
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
		PartnerID:      allocation.PartnerID,
		AllocationType: allocation.AllocationType,
		StartDate:      allocation.StartDate,
		EndDate:        allocation.EndDate,
		IsActive:       allocation.IsActive,
		CreatedAt:      allocation.CreatedAt,
		UpdatedAt:      allocation.UpdatedAt,
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

func (h *handler) GetPartnerAllocations(w http.ResponseWriter, r *http.Request) {
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

	partnerID, err := uuid.Parse(pathParts[1])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
		return
	}

	// Parse query parameters
	activeOnly := r.URL.Query().Get("active_only") != "false" // Default to active only

	// Call service to get partner allocations
	allocations, err := h.svc.GetPartnerAllocations(ctx, partnerID, activeOnly)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
		default:
			statusCode = http.StatusInternalServerError
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Convert to response DTOs
	var responses []domain.RoomAllocationResponse
	for _, allocation := range allocations {
		responses = append(responses, domain.RoomAllocationResponse{
			ID:             allocation.ID,
			RoomID:         allocation.RoomID,
			PartnerID:      allocation.PartnerID,
			AllocationType: allocation.AllocationType,
			StartDate:      allocation.StartDate,
			EndDate:        allocation.EndDate,
			IsActive:       allocation.IsActive,
			CreatedAt:      allocation.CreatedAt,
			UpdatedAt:      allocation.UpdatedAt,
		})
	}

	httpx.RespondWithJSON(w, responses, http.StatusOK)
}

func (h *handler) GetRoomAllocations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract room ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid room ID in path"), http.StatusBadRequest)
		return
	}

	roomID, err := uuid.Parse(pathParts[1])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid room ID format"), http.StatusBadRequest)
		return
	}

	// Parse query parameters
	activeOnly := r.URL.Query().Get("active_only") != "false" // Default to active only

	// Call service to get room allocations
	allocations, err := h.svc.GetRoomAllocations(ctx, roomID, activeOnly)
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

	// Convert to response DTOs
	var responses []domain.RoomAllocationResponse
	for _, allocation := range allocations {
		responses = append(responses, domain.RoomAllocationResponse{
			ID:             allocation.ID,
			RoomID:         allocation.RoomID,
			PartnerID:      allocation.PartnerID,
			AllocationType: allocation.AllocationType,
			StartDate:      allocation.StartDate,
			EndDate:        allocation.EndDate,
			IsActive:       allocation.IsActive,
			CreatedAt:      allocation.CreatedAt,
			UpdatedAt:      allocation.UpdatedAt,
		})
	}

	httpx.RespondWithJSON(w, responses, http.StatusOK)
}

func (h *handler) CheckPartnerRoomAccess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract partner ID and room ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 5 {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid URL path"), http.StatusBadRequest)
		return
	}

	partnerID, err := uuid.Parse(pathParts[1])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
		return
	}

	roomID, err := uuid.Parse(pathParts[3])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid room ID format"), http.StatusBadRequest)
		return
	}

	// Parse query parameter for time check (default to now)
	checkTime := time.Now()
	if timeStr := r.URL.Query().Get("at"); timeStr != "" {
		if parsedTime, err := time.Parse(time.RFC3339, timeStr); err == nil {
			checkTime = parsedTime
		}
	}

	// Call service to check access
	hasAccess, err := h.svc.CheckPartnerRoomAccess(ctx, partnerID, roomID, checkTime)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
		default:
			statusCode = http.StatusInternalServerError
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	response := struct {
		HasAccess bool      `json:"has_access"`
		CheckedAt time.Time `json:"checked_at"`
	}{
		HasAccess: hasAccess,
		CheckedAt: checkTime,
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}