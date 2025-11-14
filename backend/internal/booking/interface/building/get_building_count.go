package buildingHandler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetBuildingCount(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get building count request",
		"operation", "get_building_count",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Parse query parameters - same filters as GetAllBuildings
	filter := ports.BuildingFilter{}

	// Parse is_active filter
	if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err != nil {
			logger.WarnContext(ctx, "Handler: Invalid is_active parameter",
				"error", err,
				"is_active", isActiveStr,
				"operation", "get_building_count")
			httpx.RespondWithError(w, errs.NewInvalidValueErr("is_active must be a boolean (true/false)"), http.StatusBadRequest)
			return
		}
		filter.IsActive = &isActive
	}

	// Parse city filter
	if city := r.URL.Query().Get("city"); city != "" {
		filter.City = &city
	}

	// Parse country filter
	if country := r.URL.Query().Get("country"); country != "" {
		filter.Country = &country
	}

	// Call service to get building count
	count, err := h.svc.GetBuildingCount(ctx, filter)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		var statusCode int

		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			logLevel = "warn"
			errorContext = "invalid filter parameters"
			statusCode = http.StatusBadRequest
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
			logger.ErrorContext(ctx, "Handler: Get building count failed",
				"error", err,
				"filter", filter,
				"operation", "get_building_count",
				"context", errorContext,
				"status_code", statusCode)
		} else {
			logger.WarnContext(ctx, "Handler: Get building count failed",
				"error", err,
				"filter", filter,
				"operation", "get_building_count",
				"context", errorContext,
				"status_code", statusCode)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	logger.InfoContext(ctx, "Handler: Building count retrieved successfully",
		"count", count,
		"filter", filter,
		"operation", "get_building_count")

	// Return count as JSON response
	response := map[string]int{"count": count}
	httpx.RespondWithJSON(w, response, http.StatusOK)
}
