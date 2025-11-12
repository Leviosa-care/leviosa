package roomHandler

// import (
// 	"context"
// 	"errors"
// 	"net/http"
// 	"strconv"
//
// 	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
// 	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
// 	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
// 	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
// )
//
// func (h *handler) GetAllRooms(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
//
// 	logger, err := ctxutil.GetLoggerFromContext(ctx)
// 	if err != nil {
// 		httpx.RespondWithError(w, err, http.StatusInternalServerError)
// 		return
// 	}
//
// 	logger.InfoContext(ctx, "Handler: Processing get all rooms request",
// 		"operation", "get_all_rooms",
// 		"method", r.Method,
// 		"path", r.URL.Path)
//
// 	// Parse query parameters for filtering
// 	activeOnly := r.URL.Query().Get("active_only") == "true"
// 	var capacityFilter *int
// 	if capacityStr := r.URL.Query().Get("min_capacity"); capacityStr != "" {
// 		if capacity, err := strconv.Atoi(capacityStr); err == nil {
// 			capacityFilter = &capacity
// 		}
// 	}
//
// 	// Call service to get all rooms
// 	rooms, err := h.svc.GetAllRooms(ctx, activeOnly, capacityFilter)
// 	if err != nil {
// 		var statusCode int
// 		var errorContext string
//
// 		switch {
// 		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
// 			statusCode = http.StatusServiceUnavailable
// 			errorContext = "database connection failure"
// 		case errors.Is(err, errs.ErrResourceExhausted):
// 			statusCode = http.StatusServiceUnavailable
// 			errorContext = "database resource exhaustion"
// 		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
// 			statusCode = http.StatusRequestTimeout
// 			errorContext = "request cancelled"
// 		case errors.Is(err, context.DeadlineExceeded):
// 			statusCode = http.StatusRequestTimeout
// 			errorContext = "request timeout"
// 		default:
// 			statusCode = http.StatusInternalServerError
// 			errorContext = "unexpected error"
// 		}
//
// 		logger.ErrorContext(ctx, "Handler: Get all rooms failed",
// 			"error", err,
// 			"operation", "get_all_rooms",
// 			"context", errorContext,
// 			"status_code", statusCode)
//
// 		httpx.RespondWithError(w, err, statusCode)
// 		return
// 	}
//
// 	// Convert to response DTOs
// 	var responses []domain.RoomResponse
// 	for _, room := range rooms {
// 		responses = append(responses, domain.RoomResponse{
// 			ID:          room.ID,
// 			BuildingID:  room.BuildingID,
// 			Name:        room.Name,
// 			Description: room.Description,
// 			Capacity:    room.Capacity,
// 			Equipment:   room.Equipment,
// 			PriceCents:  room.PriceCents,
// 			Currency:    room.Currency,
// 			IsActive:    room.IsActive,
// 			CreatedAt:   room.CreatedAt,
// 			UpdatedAt:   room.UpdatedAt,
// 		})
// 	}
//
// 	logger.InfoContext(ctx, "Handler: All rooms retrieved successfully",
// 		"room_count", len(rooms),
// 		"active_only", activeOnly,
// 		"operation", "get_all_rooms")
//
// 	httpx.RespondWithJSON(w, responses, http.StatusOK)
// }
