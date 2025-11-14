package buildingHandler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) GetBuildingByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get building by ID request",
		"operation", "get_building_by_id",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Extract building ID from URL path parameter
	idStr := r.PathValue("id")
	if idStr == "" {
		logger.WarnContext(ctx, "Handler: Missing building ID in path",
			"operation", "get_building_by_id",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("building ID is required"), http.StatusBadRequest)
		return
	}

	buildingID, err := uuid.Parse(idStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid building ID format",
			"error", err,
			"building_id", idStr,
			"operation", "get_building_by_id",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid building ID format"), http.StatusBadRequest)
		return
	}

	// Call service to get building
	building, err := h.svc.GetBuildingByID(ctx, buildingID)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		var statusCode int

		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			logLevel = "warn"
			errorContext = "building not found"
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrInvalidValue):
			logLevel = "warn"
			errorContext = "invalid building ID"
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrNotEncrypted), errors.Is(err, errs.ErrNotDecrypted):
			logLevel = "error"
			errorContext = "encryption error"
			statusCode = http.StatusInternalServerError
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logLevel = "error"
			errorContext = "database operation failed"
			statusCode = http.StatusInternalServerError
		case errors.Is(err, context.Canceled):
			logLevel = "warn"
			errorContext = "request cancelled"
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, context.DeadlineExceeded):
			logLevel = "warn"
			errorContext = "request timeout"
			statusCode = http.StatusRequestTimeout
		default:
			logLevel = "error"
			errorContext = "unexpected error"
			statusCode = http.StatusInternalServerError
		}

		if logLevel == "error" {
			logger.ErrorContext(ctx, "Handler: Get building by ID failed",
				"error", err,
				"building_id", buildingID,
				"operation", "get_building_by_id",
				"context", errorContext,
				"status_code", statusCode)
		} else {
			logger.WarnContext(ctx, "Handler: Get building by ID failed",
				"error", err,
				"building_id", buildingID,
				"operation", "get_building_by_id",
				"context", errorContext,
				"status_code", statusCode)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	logger.InfoContext(ctx, "Handler: Building retrieved successfully",
		"building_id", building.ID,
		"building_name", building.Name,
		"operation", "get_building_by_id")

	httpx.RespondWithJSON(w, building, http.StatusOK)
}
