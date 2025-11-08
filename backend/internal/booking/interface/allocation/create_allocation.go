package allocationHandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) CreateSharedAllocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing create shared allocation request",
		"operation", "create_shared_allocation",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Parse request body
	var request domain.CreateSharedAllocationRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "create_shared_allocation",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call service to create shared allocation
	allocation, err := h.svc.CreateSharedAllocation(ctx, request.RoomID, request.PartnerID)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		var statusCode int

		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			logLevel = "warn"
			errorContext = "invalid request validation"
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrRepositoryNotFound):
			logLevel = "warn"
			errorContext = "room not found"
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrUniqueViolation):
			logLevel = "warn"
			errorContext = "allocation conflict"
			statusCode = http.StatusConflict
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			logLevel = "error"
			errorContext = "database connection failure"
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrResourceExhausted):
			logLevel = "error"
			errorContext = "database resource exhaustion"
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			logLevel = "warn"
			errorContext = "request cancelled"
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, context.DeadlineExceeded):
			logLevel = "warn"
			errorContext = "request timeout"
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			logLevel = "error"
			errorContext = "transaction conflict"
			statusCode = http.StatusServiceUnavailable
		default:
			logLevel = "error"
			errorContext = "unexpected error"
			statusCode = http.StatusInternalServerError
		}

		if logLevel == "error" {
			logger.ErrorContext(ctx, "Handler: Create shared allocation failed",
				"error", err,
				"operation", "create_shared_allocation",
				"context", errorContext,
				"room_id", request.RoomID,
				"partner_id", request.PartnerID,
				"status_code", statusCode)
		} else {
			logger.WarnContext(ctx, "Handler: Create shared allocation failed",
				"error", err,
				"operation", "create_shared_allocation",
				"context", errorContext,
				"room_id", request.RoomID,
				"partner_id", request.PartnerID,
				"status_code", statusCode)
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

	logger.InfoContext(ctx, "Handler: Shared allocation created successfully",
		"allocation_id", allocation.ID,
		"room_id", allocation.RoomID,
		"partner_id", allocation.PartnerID,
		"operation", "create_shared_allocation")

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}

func (h *handler) CreateDedicatedAllocation(w http.ResponseWriter, r *http.Request) {
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
	allocation, err := h.svc.CreateDedicatedAllocation(ctx, request.RoomID, request.PartnerID, request.StartDate, request.EndDate)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		var statusCode int

		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			logLevel = "warn"
			errorContext = "invalid request validation"
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrRepositoryNotFound):
			logLevel = "warn"
			errorContext = "room not found"
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrUniqueViolation):
			logLevel = "warn"
			errorContext = "allocation conflict"
			statusCode = http.StatusConflict
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			logLevel = "error"
			errorContext = "database connection failure"
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrResourceExhausted):
			logLevel = "error"
			errorContext = "database resource exhaustion"
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			logLevel = "warn"
			errorContext = "request cancelled"
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, context.DeadlineExceeded):
			logLevel = "warn"
			errorContext = "request timeout"
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			logLevel = "error"
			errorContext = "transaction conflict"
			statusCode = http.StatusServiceUnavailable
		default:
			logLevel = "error"
			errorContext = "unexpected error"
			statusCode = http.StatusInternalServerError
		}

		if logLevel == "error" {
			logger.ErrorContext(ctx, "Handler: Create dedicated allocation failed",
				"error", err,
				"operation", "create_dedicated_allocation",
				"context", errorContext,
				"room_id", request.RoomID,
				"partner_id", request.PartnerID,
				"status_code", statusCode)
		} else {
			logger.WarnContext(ctx, "Handler: Create dedicated allocation failed",
				"error", err,
				"operation", "create_dedicated_allocation",
				"context", errorContext,
				"room_id", request.RoomID,
				"partner_id", request.PartnerID,
				"status_code", statusCode)
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

	logger.InfoContext(ctx, "Handler: Dedicated allocation created successfully",
		"allocation_id", allocation.ID,
		"room_id", allocation.RoomID,
		"partner_id", allocation.PartnerID,
		"start_date", request.StartDate,
		"end_date", request.EndDate,
		"operation", "create_dedicated_allocation")

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}
