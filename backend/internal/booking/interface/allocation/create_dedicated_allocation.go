package allocationHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) CreateDedicatedAllocation(w http.ResponseWriter, r *http.Request) {
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
	logger.InfoContext(ctx, "Handler: Processing create dedicated allocation request",
		"operation", "create_dedicated_allocation",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Parse request body
	var request domain.CreateDedicatedAllocationRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "create_dedicated_allocation",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call service to create dedicated allocation
	allocation, err := h.svc.CreateDedicatedAllocation(ctx, &request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "create dedicated allocation for partner")
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

	logger.InfoContext(ctx, "Handler: Dedicated allocation created successfully",
		"allocation_id", allocation.ID,
		"room_id", allocation.RoomID,
		"partner_id", allocation.UserID,
		"start_date", request.StartDate,
		"end_date", request.EndDate,
		"operation", "create_dedicated_allocation")

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}
