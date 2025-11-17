package http

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetCompanyPhone(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get company phone",
		"operation", "get_company_phone",
		"method", r.Method,
		"path", r.URL.Path)

	response, err := h.svc.GetCompanyTelephone(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get company phone")
		return
	}

	logger.InfoContext(ctx, "Handler: Get company phone completed",
		"operation", "get_company_phone",
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
