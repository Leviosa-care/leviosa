package roomHandler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"

	"github.com/google/uuid"
)

func (h *handler) GetRoomsByBuilding(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract building ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		logger.WarnContext(ctx, "Handler: Invalid URL path for get rooms by building",
			"path", r.URL.Path,
			"operation", "get_rooms_by_building")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid building ID in path"), http.StatusBadRequest)
		return
	}

	buildingIDStr := pathParts[1]
	buildingID, err := uuid.Parse(buildingIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid building ID format",
			"building_id", buildingIDStr,
			"error", err,
			"operation", "get_rooms_by_building")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid building ID format"), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get rooms by building request",
		"building_id", buildingID,
		"operation", "get_rooms_by_building",
		"method", r.Method,
		"path", r.URL.Path)

	// Parse query parameters for filtering
	activeOnly := r.URL.Query().Get("active_only") == "true"

	// Call service to get rooms by building
	rooms, err := h.svc.GetRoomsByBuilding(ctx, buildingID, activeOnly)
	if err != nil {
		var statusCode int
		var errorContext string

		switch {
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

		logger.ErrorContext(ctx, "Handler: Get rooms by building failed",
			"error", err,
			"building_id", buildingID,
			"operation", "get_rooms_by_building",
			"context", errorContext,
			"status_code", statusCode)

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	logger.InfoContext(ctx, "Handler: Rooms retrieved successfully",
		"building_id", buildingID,
		"room_count", len(rooms),
		"active_only", activeOnly,
		"operation", "get_rooms_by_building")

	httpx.RespondWithJSON(w, rooms, http.StatusOK)
}
