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

// GetPartnerMetrics handles GET /partners/metrics/{partner_id}?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func (h *Handler) GetPartnerMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Getting partner metrics",
		"operation", "get_partner_metrics")

	// Parse partner_id from URL path
	partnerIDStr := r.PathValue("partner_id")
	if partnerIDStr == "" {
		logger.WarnContext(ctx, "Handler: Missing partner_id parameter",
			"operation", "get_partner_metrics")
		httpx.RespondWithError(w, errors.New("partner_id is required"), http.StatusBadRequest)
		return
	}

	partnerID, err := uuid.Parse(partnerIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid partner_id format",
			"partner_id", partnerIDStr,
			"error", err.Error(),
			"operation", "get_partner_metrics")
		httpx.RespondWithError(w, errors.New("invalid partner_id format"), http.StatusBadRequest)
		return
	}

	// Parse date range from query parameters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		logger.WarnContext(ctx, "Handler: Missing date parameters",
			"partner_id", partnerID,
			"operation", "get_partner_metrics")
		httpx.RespondWithError(w, errors.New("start_date and end_date query parameters are required (format: YYYY-MM-DD)"), http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid start_date format",
			"partner_id", partnerID,
			"start_date", startDateStr,
			"error", err.Error(),
			"operation", "get_partner_metrics")
		httpx.RespondWithError(w, errors.New("invalid start_date format, use YYYY-MM-DD"), http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid end_date format",
			"partner_id", partnerID,
			"end_date", endDateStr,
			"error", err.Error(),
			"operation", "get_partner_metrics")
		httpx.RespondWithError(w, errors.New("invalid end_date format, use YYYY-MM-DD"), http.StatusBadRequest)
		return
	}

	// Call service
	response, err := h.svc.GetPartnerUtilization(ctx, partnerID, startDate, endDate)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get partner utilization metrics")
		return
	}

	logger.InfoContext(ctx, "Handler: Successfully retrieved partner metrics",
		"partner_id", partnerID,
		"start_date", startDate.Format("2006-01-02"),
		"end_date", endDate.Format("2006-01-02"),
		"rooms_count", len(response.RoomMetrics),
		"days_analyzed", response.Summary.DaysAnalyzed,
		"operation", "get_partner_metrics")

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.ErrorContext(ctx, "Handler: Failed to encode response",
			"error", err.Error(),
			"operation", "get_partner_metrics")
	}
}
