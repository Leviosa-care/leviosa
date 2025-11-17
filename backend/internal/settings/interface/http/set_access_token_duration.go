package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"
)

func (h *handler) SetAccessTokenDuration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	var payload domain.SetAccessTokenDurationRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "set_access_token_duration",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing set access token duration request",
		"operation", "set_access_token_duration",
		"method", r.Method,
		"path", r.URL.Path,
		"duration_minutes", payload.Duration,
		"user_agent", r.Header.Get("User-Agent"))

	response, err := h.svc.SetAccessTokenDuration(ctx, &payload)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "set access token duration")
		return
	}

	logger.InfoContext(ctx, "Handler: Set access token duration completed",
		"operation", "set_access_token_duration",
		"method", r.Method,
		"path", r.URL.Path,
		"duration_minutes", payload.Duration,
		"status_code", http.StatusOK,
		"success", response.Success)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
