package roomHandler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"

	"github.com/google/uuid"
)

func (h *handler) GetRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract room ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		logger.WarnContext(ctx, "Handler: Invalid URL path for get room",
			"path", r.URL.Path,
			"operation", "get_room")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid room ID in path"), http.StatusBadRequest)
		return
	}

	roomIDStr := pathParts[1]
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid room ID format",
			"room_id", roomIDStr,
			"error", err,
			"operation", "get_room")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid room ID format"), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get room request",
		"room_id", roomID,
		"operation", "get_room",
		"method", r.Method,
		"path", r.URL.Path)

	// Call service to get room
	room, err := h.svc.GetRoom(ctx, roomID)
	if err != nil {
		var statusCode int
		var errorContext string

		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			statusCode = http.StatusNotFound
			errorContext = "room not found"
		case errors.Is(err, errs.ErrInvalidValue):
			statusCode = http.StatusBadRequest
			errorContext = "invalid request validation"
		case errors.Is(err, errs.ErrNotEncrypted), errors.Is(err, errs.ErrNotDecrypted):
			statusCode = http.StatusInternalServerError
			errorContext = "encryption error"
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			statusCode = http.StatusInternalServerError
			errorContext = "database operation failed"
		case errors.Is(err, context.Canceled):
			statusCode = http.StatusRequestTimeout
			errorContext = "request cancelled"
		case errors.Is(err, context.DeadlineExceeded):
			statusCode = http.StatusRequestTimeout
			errorContext = "request timeout"
		default:
			statusCode = http.StatusInternalServerError
			errorContext = "unexpected error"
		}

		if statusCode >= 500 {
			logger.ErrorContext(ctx, "Handler: Get room failed",
				"error", err,
				"room_id", roomID,
				"operation", "get_room",
				"context", errorContext,
				"status_code", statusCode)
		} else {
			logger.WarnContext(ctx, "Handler: Get room failed",
				"error", err,
				"room_id", roomID,
				"operation", "get_room",
				"context", errorContext,
				"status_code", statusCode)
		}

		httpx.RespondWithError(w, err, statusCode)
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
		IsActive:    room.IsActive,
	}

	logger.InfoContext(ctx, "Handler: Room retrieved successfully",
		"room_id", roomID,
		"room_name", room.Name,
		"building_id", room.BuildingID,
		"operation", "get_room")

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
