package roomHandler

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

	"github.com/google/uuid"
)

func (h *handler) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	print("IN THE HANDLER")
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get all rooms request",
		"operation", "get_all_rooms",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Initialize filter with defaults
	filter := ports.RoomFilter{
		Limit:          20, // Default limit
		Offset:         0,  // Default offset
		OrderBy:        "created_at",
		OrderDirection: "desc",
	}

	// Parse is_active filter
	if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err != nil {
			logger.WarnContext(ctx, "Handler: Invalid is_active parameter",
				"error", err,
				"is_active", isActiveStr,
				"operation", "get_all_rooms")
			httpx.RespondWithError(w, errs.NewInvalidValueErr("is_active must be a boolean (true/false)"), http.StatusBadRequest)
			return
		}
		filter.IsActive = &isActive
	}

	// Parse building_id filter
	if buildingIDStr := r.URL.Query().Get("building_id"); buildingIDStr != "" {
		buildingID, err := uuid.Parse(buildingIDStr)
		if err != nil {
			logger.WarnContext(ctx, "Handler: Invalid building_id parameter",
				"error", err,
				"building_id", buildingIDStr,
				"operation", "get_all_rooms")
			httpx.RespondWithError(w, errs.NewInvalidValueErr("building_id must be a valid UUID"), http.StatusBadRequest)
			return
		}
		filter.BuildingID = &buildingID
	}

	// Parse min_capacity filter
	if minCapacityStr := r.URL.Query().Get("min_capacity"); minCapacityStr != "" {
		minCapacity, err := strconv.Atoi(minCapacityStr)
		if err != nil || minCapacity < 0 {
			logger.WarnContext(ctx, "Handler: Invalid min_capacity parameter",
				"error", err,
				"min_capacity", minCapacityStr,
				"operation", "get_all_rooms")
			httpx.RespondWithError(w, errs.NewInvalidValueErr("min_capacity must be a non-negative integer"), http.StatusBadRequest)
			return
		}
		filter.MinCapacity = &minCapacity
	}

	// Parse max_capacity filter
	if maxCapacityStr := r.URL.Query().Get("max_capacity"); maxCapacityStr != "" {
		maxCapacity, err := strconv.Atoi(maxCapacityStr)
		if err != nil || maxCapacity < 0 {
			logger.WarnContext(ctx, "Handler: Invalid max_capacity parameter",
				"error", err,
				"max_capacity", maxCapacityStr,
				"operation", "get_all_rooms")
			httpx.RespondWithError(w, errs.NewInvalidValueErr("max_capacity must be a non-negative integer"), http.StatusBadRequest)
			return
		}
		filter.MaxCapacity = &maxCapacity
	}

	// Validate capacity range
	if filter.MinCapacity != nil && filter.MaxCapacity != nil && *filter.MinCapacity > *filter.MaxCapacity {
		logger.WarnContext(ctx, "Handler: Invalid capacity range",
			"min_capacity", *filter.MinCapacity,
			"max_capacity", *filter.MaxCapacity,
			"operation", "get_all_rooms")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("min_capacity cannot be greater than max_capacity"), http.StatusBadRequest)
		return
	}

	// Parse name filter (searchable field)
	if name := r.URL.Query().Get("name"); name != "" {
		filter.Name = &name
	}

	// Parse room_number filter (searchable field)
	if roomNumber := r.URL.Query().Get("room_number"); roomNumber != "" {
		filter.RoomNumber = &roomNumber
	}

	// Parse pagination parameters
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 100 {
			logger.WarnContext(ctx, "Handler: Invalid limit parameter",
				"error", err,
				"limit", limitStr,
				"operation", "get_all_rooms")
			httpx.RespondWithError(w, errs.NewInvalidValueErr("limit must be a positive integer between 1 and 100"), http.StatusBadRequest)
			return
		}
		filter.Limit = limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			logger.WarnContext(ctx, "Handler: Invalid offset parameter",
				"error", err,
				"offset", offsetStr,
				"operation", "get_all_rooms")
			httpx.RespondWithError(w, errs.NewInvalidValueErr("offset must be a non-negative integer"), http.StatusBadRequest)
			return
		}
		filter.Offset = offset
	}

	// Parse sorting parameters
	if orderBy := r.URL.Query().Get("order_by"); orderBy != "" {
		// Validate order_by values
		validOrderBy := map[string]bool{
			"name":       true,
			"created_at": true,
			"capacity":   true,
		}
		if !validOrderBy[orderBy] {
			logger.WarnContext(ctx, "Handler: Invalid order_by parameter",
				"order_by", orderBy,
				"operation", "get_all_rooms")
			httpx.RespondWithError(w, errs.NewInvalidValueErr("order_by must be one of: name, created_at, capacity"), http.StatusBadRequest)
			return
		}
		filter.OrderBy = orderBy
	}

	if orderDirection := r.URL.Query().Get("order_direction"); orderDirection != "" {
		// Validate order_direction values
		if orderDirection != "asc" && orderDirection != "desc" {
			logger.WarnContext(ctx, "Handler: Invalid order_direction parameter",
				"order_direction", orderDirection,
				"operation", "get_all_rooms")
			httpx.RespondWithError(w, errs.NewInvalidValueErr("order_direction must be either 'asc' or 'desc'"), http.StatusBadRequest)
			return
		}
		filter.OrderDirection = orderDirection
	}
	println("GETTING TO THE APPLICATION PART")

	// Call service to get rooms
	rooms, err := h.svc.ListRooms(ctx, filter)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		var statusCode int

		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			logLevel = "warn"
			errorContext = "invalid filter parameters"
			statusCode = http.StatusBadRequest
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
			logger.ErrorContext(ctx, "Handler: Get all rooms failed",
				"error", err,
				"filter", filter,
				"operation", "get_all_rooms",
				"context", errorContext,
				"status_code", statusCode)
		} else {
			logger.WarnContext(ctx, "Handler: Get all rooms failed",
				"error", err,
				"filter", filter,
				"operation", "get_all_rooms",
				"context", errorContext,
				"status_code", statusCode)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	logger.InfoContext(ctx, "Handler: Rooms retrieved successfully",
		"count", len(rooms),
		"filter", filter,
		"operation", "get_all_rooms")

	httpx.RespondWithJSON(w, rooms, http.StatusOK)
}
