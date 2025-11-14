package roomHandler

import (
	"context"
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

func (h *handler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
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

	// Extract room ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		logger.WarnContext(ctx, "Handler: Invalid URL path for update room",
			"path", r.URL.Path,
			"operation", "update_room")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid room ID in path"), http.StatusBadRequest)
		return
	}

	roomIDStr := pathParts[1]
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid room ID format",
			"room_id", roomIDStr,
			"error", err,
			"operation", "update_room")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid room ID format"), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing update room request",
		"room_id", roomID,
		"operation", "update_room",
		"method", r.Method,
		"path", r.URL.Path)

	// Parse request body
	var request domain.UpdateRoomRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"room_id", roomID,
			"operation", "update_room")
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	request.ID = roomID

	// Call service to update room
	room, err := h.svc.UpdateRoom(ctx, &request)
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
		case errors.Is(err, errs.ErrConflict):
			statusCode = http.StatusConflict
			errorContext = "room conflict"
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
			logger.ErrorContext(ctx, "Handler: Update room failed",
				"error", err,
				"room_id", roomID,
				"operation", "update_room",
				"context", errorContext,
				"status_code", statusCode)
		} else {
			logger.WarnContext(ctx, "Handler: Update room failed",
				"error", err,
				"room_id", roomID,
				"operation", "update_room",
				"context", errorContext,
				"status_code", statusCode)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	logger.InfoContext(ctx, "Handler: Room updated successfully",
		"room_id", roomID,
		"room_name", room.Name,
		"operation", "update_room")

	httpx.RespondWithJSON(w, room, http.StatusOK)
}
