package http

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAccessTokenDuration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get access token duration request",
		"operation", "get_access_token_duration",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	response, err := h.svc.GetAccessTokenDuration(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get access token duration")
		return
	}

	logger.InfoContext(ctx, "Handler: Get access token duration completed",
		"operation", "get_access_token_duration",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"duration_minutes", response.Duration)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
