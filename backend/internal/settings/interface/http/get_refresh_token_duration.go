package http

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetRefreshTokenDuration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get refresh token duration request",
		"operation", "get_refresh_token_duration",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	response, err := h.svc.GetRefreshTokenDuration(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get refresh token duration")
		return
	}

	logger.InfoContext(ctx, "Handler: Get refresh token duration completed",
		"operation", "get_refresh_token_duration",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"duration_hours", response.Duration)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
