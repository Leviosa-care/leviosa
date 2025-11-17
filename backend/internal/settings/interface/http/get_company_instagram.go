package http

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetCompanyInstagram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get company instagram request",
		"operation", "get_company_instagram",
		"method", r.Method,
		"path", r.URL.Path)

	response, err := h.svc.GetCompanyInstagram(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get company instagram")
		return
	}

	logger.InfoContext(ctx, "Handler: Get company instagram completed",
		"operation", "get_company_instagram",
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
