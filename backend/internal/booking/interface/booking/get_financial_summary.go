package bookingHandler

import (
	"net/http"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetFinancialSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	from, to, err := parseFinancialDateRange(r)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	result, err := h.svc.GetFinancialSummary(ctx, from, to)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	httpx.RespondWithJSON(w, result, http.StatusOK)
}

// parseFinancialDateRange extracts the from/to query parameters with sensible defaults.
// Defaults: from = start of current month, to = end of today.
// All times are UTC to stay consistent with time.Parse behaviour.
func parseFinancialDateRange(r *http.Request) (from, to time.Time, err error) {
	now := time.Now().UTC()

	// Default: start of current month (UTC)
	from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	// Default: start of tomorrow (exclusive upper bound for today)
	to = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)

	q := r.URL.Query()

	if fromStr := q.Get("from"); fromStr != "" {
		parsed, e := time.ParseInLocation(time.DateOnly, fromStr, time.UTC)
		if e != nil {
			return from, to, errs.NewInvalidValueErr("invalid from date, expected YYYY-MM-DD")
		}
		from = parsed
	}

	if toStr := q.Get("to"); toStr != "" {
		parsed, e := time.ParseInLocation(time.DateOnly, toStr, time.UTC)
		if e != nil {
			return from, to, errs.NewInvalidValueErr("invalid to date, expected YYYY-MM-DD")
		}
		// to is exclusive upper bound → next day at midnight
		to = parsed.AddDate(0, 0, 1)
	}

	if !from.Before(to) {
		return from, to, errs.NewInvalidValueErr("from must be before to")
	}

	return from, to, nil
}
