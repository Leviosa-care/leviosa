package availabilityHandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

// GetRoomGaps handles GET /availabilities/rooms/:room_id/gaps?date=YYYY-MM-DD
func (h *handler) GetRoomGaps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Getting room gaps",
		"operation", "get_room_gaps")

	// Parse room_id from URL path parameter
	roomIDStr := r.PathValue("room_id")
	if roomIDStr == "" {
		logger.WarnContext(ctx, "Handler: Missing room_id parameter",
			"operation", "get_room_gaps")
		httpx.RespondWithError(w, errors.New("room_id is required"), http.StatusBadRequest)
		return
	}

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid room_id format",
			"room_id", roomIDStr,
			"error", err.Error(),
			"operation", "get_room_gaps")
		httpx.RespondWithError(w, errors.New("invalid room_id format"), http.StatusBadRequest)
		return
	}

	// Parse date query parameter
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		logger.WarnContext(ctx, "Handler: Missing date parameter",
			"room_id", roomID,
			"operation", "get_room_gaps")
		httpx.RespondWithError(w, errors.New("date query parameter is required (format: YYYY-MM-DD)"), http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid date format",
			"room_id", roomID,
			"date_string", dateStr,
			"error", err.Error(),
			"operation", "get_room_gaps")
		httpx.RespondWithError(w, errors.New("invalid date format, use YYYY-MM-DD"), http.StatusBadRequest)
		return
	}

	// Build request
	request := domain.GetRoomGapsRequest{
		RoomID: roomID,
		Date:   date,
	}

	// Call service
	response, err := h.svc.GetRoomGaps(ctx, request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get room gaps")
		return
	}

	logger.InfoContext(ctx, "Handler: Successfully retrieved room gaps",
		"room_id", roomID,
		"date", date.Format("2006-01-02"),
		"gaps_count", len(response.Gaps),
		"total_gap_minutes", response.TotalGapMinutes,
		"operation", "get_room_gaps")

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.ErrorContext(ctx, "Handler: Failed to encode response",
			"error", err.Error(),
			"operation", "get_room_gaps")
	}
}
