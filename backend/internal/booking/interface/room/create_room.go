package roomHandler

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

func (h *handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing create room request",
		"operation", "create_room",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Parse request body
	var request domain.CreateRoomRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "create_room",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call service to create room
	room, err := h.svc.CreateRoom(ctx, request.BuildingID, request.Name, request.Description, request.Capacity)
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
			errorContext = "building not found"
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrUniqueViolation):
			logLevel = "warn"
			errorContext = "room name already exists in building"
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
			logger.ErrorContext(ctx, "Handler: Create room failed",
				"error", err,
				"operation", "create_room",
				"context", errorContext,
				"building_id", request.BuildingID,
				"room_name", request.Name,
				"status_code", statusCode)
		} else {
			logger.WarnContext(ctx, "Handler: Create room failed",
				"error", err,
				"operation", "create_room",
				"context", errorContext,
				"building_id", request.BuildingID,
				"room_name", request.Name,
				"status_code", statusCode)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Set additional fields if provided
	if len(request.Equipment) > 0 {
		room.SetEquipment(request.Equipment)
	}
	if request.PriceCents != nil && request.Currency != "" {
		room.SetPricing(*request.PriceCents, request.Currency)
	}

	// Update room with additional fields if any were set
	if len(request.Equipment) > 0 || request.PriceCents != nil {
		room, err = h.svc.UpdateRoom(ctx, room.ID, room.Name, room.Description, room.Capacity, room.Equipment)
		if err != nil {
			logger.ErrorContext(ctx, "Handler: Failed to update room with additional fields",
				"error", err,
				"room_id", room.ID,
				"operation", "create_room")
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		if request.PriceCents != nil && request.Currency != "" {
			room, err = h.svc.UpdateRoomPricing(ctx, room.ID, request.PriceCents, request.Currency)
			if err != nil {
				logger.ErrorContext(ctx, "Handler: Failed to update room pricing",
					"error", err,
					"room_id", room.ID,
					"operation", "create_room")
				httpx.RespondWithError(w, err, http.StatusInternalServerError)
				return
			}
		}
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

	logger.InfoContext(ctx, "Handler: Room created successfully",
		"room_id", room.ID,
		"room_name", room.Name,
		"building_id", room.BuildingID,
		"operation", "create_room")

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}
