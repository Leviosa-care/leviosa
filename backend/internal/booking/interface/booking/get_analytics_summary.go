package bookingHandler

import (
	"net/http"
	"strconv"

	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAnalyticsSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse optional ?months query parameter (default 6, max 12)
	months := 6
	if monthsStr := r.URL.Query().Get("months"); monthsStr != "" {
		if m, err := strconv.Atoi(monthsStr); err == nil && m >= 1 && m <= 12 {
			months = m
		}
	}

	result, err := h.svc.GetAnalyticsSummary(ctx, months)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	httpx.RespondWithJSON(w, result, http.StatusOK)
}
