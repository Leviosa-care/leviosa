package metricsHandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

// GetRoomMetrics handles GET /rooms/{room_id}/metrics?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func (h *Handler) GetRoomMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Getting room metrics",
		"operation", "get_room_metrics")

	// Parse room_id from URL path
	roomIDStr := r.PathValue("room_id")
	if roomIDStr == "" {
		logger.WarnContext(ctx, "Handler: Missing room_id parameter",
			"operation", "get_room_metrics")
		httpx.RespondWithError(w, errors.New("room_id is required"), http.StatusBadRequest)
		return
	}

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid room_id format",
			"room_id", roomIDStr,
			"error", err.Error(),
			"operation", "get_room_metrics")
		httpx.RespondWithError(w, errors.New("invalid room_id format"), http.StatusBadRequest)
		return
	}

	// Parse date range from query parameters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		logger.WarnContext(ctx, "Handler: Missing date parameters",
			"room_id", roomID,
			"operation", "get_room_metrics")
		httpx.RespondWithError(w, errors.New("start_date and end_date query parameters are required (format: YYYY-MM-DD)"), http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid start_date format",
			"room_id", roomID,
			"start_date", startDateStr,
			"error", err.Error(),
			"operation", "get_room_metrics")
		httpx.RespondWithError(w, errors.New("invalid start_date format, use YYYY-MM-DD"), http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid end_date format",
			"room_id", roomID,
			"end_date", endDateStr,
			"error", err.Error(),
			"operation", "get_room_metrics")
		httpx.RespondWithError(w, errors.New("invalid end_date format, use YYYY-MM-DD"), http.StatusBadRequest)
		return
	}

	// Call service
	response, err := h.svc.GetRoomUtilization(ctx, roomID, startDate, endDate)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get room utilization metrics")
		return
	}

	logger.InfoContext(ctx, "Handler: Successfully retrieved room metrics",
		"room_id", roomID,
		"start_date", startDate.Format("2006-01-02"),
		"end_date", endDate.Format("2006-01-02"),
		"days_analyzed", response.Summary.DaysAnalyzed,
		"operation", "get_room_metrics")

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.ErrorContext(ctx, "Handler: Failed to encode response",
			"error", err.Error(),
			"operation", "get_room_metrics")
	}
}
