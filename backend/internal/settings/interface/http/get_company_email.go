package http

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetCompanyEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get company email",
		"operation", "get_company_email",
		"method", r.Method,
		"path", r.URL.Path)

	response, err := h.svc.GetCompanyEmail(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get company email")
		return
	}

	logger.InfoContext(ctx, "Handler: Get company email completed",
		"operation", "get_company_email",
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
