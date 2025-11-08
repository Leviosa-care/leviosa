package roomHandler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
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
		PriceCents:  room.PriceCents,
		Currency:    room.Currency,
		IsActive:    room.IsActive,
		CreatedAt:   room.CreatedAt,
		UpdatedAt:   room.UpdatedAt,
	}

	logger.InfoContext(ctx, "Handler: Room retrieved successfully",
		"room_id", roomID,
		"room_name", room.Name,
		"building_id", room.BuildingID,
		"operation", "get_room")

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

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
	var capacityFilter *int
	if capacityStr := r.URL.Query().Get("min_capacity"); capacityStr != "" {
		if capacity, err := strconv.Atoi(capacityStr); err == nil {
			capacityFilter = &capacity
		}
	}

	// Call service to get rooms by building
	rooms, err := h.svc.GetRoomsByBuilding(ctx, buildingID, activeOnly, capacityFilter)
	if err != nil {
		var statusCode int
		var errorContext string

		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusNotFound
			errorContext = "building not found"
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

	// Convert to response DTOs
	var responses []domain.RoomResponse
	for _, room := range rooms {
		responses = append(responses, domain.RoomResponse{
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
		})
	}

	logger.InfoContext(ctx, "Handler: Rooms retrieved successfully",
		"building_id", buildingID,
		"room_count", len(rooms),
		"active_only", activeOnly,
		"operation", "get_rooms_by_building")

	httpx.RespondWithJSON(w, responses, http.StatusOK)
}

func (h *handler) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get all rooms request",
		"operation", "get_all_rooms",
		"method", r.Method,
		"path", r.URL.Path)

	// Parse query parameters for filtering
	activeOnly := r.URL.Query().Get("active_only") == "true"
	var capacityFilter *int
	if capacityStr := r.URL.Query().Get("min_capacity"); capacityStr != "" {
		if capacity, err := strconv.Atoi(capacityStr); err == nil {
			capacityFilter = &capacity
		}
	}

	// Call service to get all rooms
	rooms, err := h.svc.GetAllRooms(ctx, activeOnly, capacityFilter)
	if err != nil {
		var statusCode int
		var errorContext string

		switch {
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
		default:
			statusCode = http.StatusInternalServerError
			errorContext = "unexpected error"
		}

		logger.ErrorContext(ctx, "Handler: Get all rooms failed",
			"error", err,
			"operation", "get_all_rooms",
			"context", errorContext,
			"status_code", statusCode)

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Convert to response DTOs
	var responses []domain.RoomResponse
	for _, room := range rooms {
		responses = append(responses, domain.RoomResponse{
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
		})
	}

	logger.InfoContext(ctx, "Handler: All rooms retrieved successfully",
		"room_count", len(rooms),
		"active_only", activeOnly,
		"operation", "get_all_rooms")

	httpx.RespondWithJSON(w, responses, http.StatusOK)
}
