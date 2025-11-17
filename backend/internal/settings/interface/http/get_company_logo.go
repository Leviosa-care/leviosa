package http

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetCompanyLogo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get company logo request",
		"operation", "get_company_logo",
		"method", r.Method,
		"path", r.URL.Path)

	response, err := h.svc.GetCompanyLogo(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get company logo")
		return
	}

	logger.InfoContext(ctx, "Handler: Get company logo completed",
		"operation", "get_company_logo",
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
