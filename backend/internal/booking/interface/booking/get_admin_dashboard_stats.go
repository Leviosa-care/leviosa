package bookingHandler

import (
	"errors"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stats, err := h.svc.GetDashboardStats(ctx)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrUnauthorized):
			statusCode = http.StatusUnauthorized
		case errors.Is(err, errs.ErrForbidden):
			statusCode = http.StatusForbidden
		default:
			statusCode = http.StatusInternalServerError
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	httpx.RespondWithJSON(w, stats, http.StatusOK)
}
