package roomHandler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
	"github.com/google/uuid"
)

func (h *handler) ActivateRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract room ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		logger.WarnContext(ctx, "Handler: Invalid URL path for activate room",
			"path", r.URL.Path,
			"operation", "activate_room")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid room ID in path"), http.StatusBadRequest)
		return
	}

	roomIDStr := pathParts[1]
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid room ID format",
			"room_id", roomIDStr,
			"error", err,
			"operation", "activate_room")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid room ID format"), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing activate room request",
		"room_id", roomID,
		"operation", "activate_room",
		"method", r.Method,
		"path", r.URL.Path)

	// Call service to activate room
	err = h.svc.ActivateRoom(ctx, roomID)
	if err != nil {
		var statusCode int
		var errorContext string

		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusNotFound
			errorContext = "room not found"
		case errors.Is(err, errs.ErrInvalidInput):
			statusCode = http.StatusBadRequest
			errorContext = "cannot activate room (building may be inactive)"
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
			errorContext = "database connection failure"
		case errors.Is(err, errs.ErrResourceExhausted):
			statusCode = http.StatusServiceUnavailable
			errorContext = "database resource exhaustion"
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			statusCode = http.StatusRequestTimeout
			errorContext = "request cancelled"
		case errors.Is(err, context.DeadlineExceeded):
			statusCode = http.StatusRequestTimeout
			errorContext = "request timeout"
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			statusCode = http.StatusServiceUnavailable
			errorContext = "transaction conflict"
		default:
			statusCode = http.StatusInternalServerError
			errorContext = "unexpected error"
		}

		if statusCode >= 500 {
			logger.ErrorContext(ctx, "Handler: Activate room failed",
				"error", err,
				"room_id", roomID,
				"operation", "activate_room",
				"context", errorContext,
				"status_code", statusCode)
		} else {
			logger.WarnContext(ctx, "Handler: Activate room failed",
				"error", err,
				"room_id", roomID,
				"operation", "activate_room",
				"context", errorContext,
				"status_code", statusCode)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Get updated room to return
	room, err := h.svc.GetRoom(ctx, roomID)
	if err != nil {
		logger.ErrorContext(ctx, "Handler: Failed to retrieve activated room",
			"error", err,
			"room_id", roomID,
			"operation", "activate_room")
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Convert to response DTO
	response := domain.RoomResponse{
		ID:          room.ID,
		BuildingID:  room.BuildingID,
		Name:        room.Name,
		Description: room.Description,
		Capacity:    room.Capacity,
		Equipment:   room.Equipment,
		PriceCents:  room.PriceCents,
		Currency:    room.Currency,
		IsActive:    room.IsActive,
		CreatedAt:   room.CreatedAt,
		UpdatedAt:   room.UpdatedAt,
	}

	logger.InfoContext(ctx, "Handler: Room activated successfully",
		"room_id", roomID,
		"operation", "activate_room")

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

func (h *handler) DeactivateRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract room ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		logger.WarnContext(ctx, "Handler: Invalid URL path for deactivate room",
			"path", r.URL.Path,
			"operation", "deactivate_room")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid room ID in path"), http.StatusBadRequest)
		return
	}

	roomIDStr := pathParts[1]
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid room ID format",
			"room_id", roomIDStr,
			"error", err,
			"operation", "deactivate_room")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid room ID format"), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing deactivate room request",
		"room_id", roomID,
		"operation", "deactivate_room",
		"method", r.Method,
		"path", r.URL.Path)

	// Call service to deactivate room
	err = h.svc.DeactivateRoom(ctx, roomID)
	if err != nil {
		var statusCode int
		var errorContext string

		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusNotFound
			errorContext = "room not found"
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
			errorContext = "database connection failure"
		case errors.Is(err, errs.ErrResourceExhausted):
			statusCode = http.StatusServiceUnavailable
			errorContext = "database resource exhaustion"
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			statusCode = http.StatusRequestTimeout
			errorContext = "request cancelled"
		case errors.Is(err, context.DeadlineExceeded):
			statusCode = http.StatusRequestTimeout
			errorContext = "request timeout"
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			statusCode = http.StatusServiceUnavailable
			errorContext = "transaction conflict"
		default:
			statusCode = http.StatusInternalServerError
			errorContext = "unexpected error"
		}

		if statusCode >= 500 {
			logger.ErrorContext(ctx, "Handler: Deactivate room failed",
				"error", err,
				"room_id", roomID,
				"operation", "deactivate_room",
				"context", errorContext,
				"status_code", statusCode)
		} else {
			logger.WarnContext(ctx, "Handler: Deactivate room failed",
				"error", err,
				"room_id", roomID,
				"operation", "deactivate_room",
				"context", errorContext,
				"status_code", statusCode)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Get updated room to return
	room, err := h.svc.GetRoom(ctx, roomID)
	if err != nil {
		logger.ErrorContext(ctx, "Handler: Failed to retrieve deactivated room",
			"error", err,
			"room_id", roomID,
			"operation", "deactivate_room")
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Convert to response DTO
	response := domain.RoomResponse{
		ID:          room.ID,
		BuildingID:  room.BuildingID,
		Name:        room.Name,
		Description: room.Description,
		Capacity:    room.Capacity,
		Equipment:   room.Equipment,
		PriceCents:  room.PriceCents,
		Currency:    room.Currency,
		IsActive:    room.IsActive,
		CreatedAt:   room.CreatedAt,
		UpdatedAt:   room.UpdatedAt,
	}

	logger.InfoContext(ctx, "Handler: Room deactivated successfully",
		"room_id", roomID,
		"operation", "deactivate_room")

	httpx.RespondWithJSON(w, response, http.StatusOK)
}