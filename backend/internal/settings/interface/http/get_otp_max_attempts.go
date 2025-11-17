package http

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetOTPMaxAttempts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get OTP max attempts request",
		"operation", "get_otp_max_attempts",
		"method", r.Method,
		"path", r.URL.Path)

	response, err := h.svc.GetOTPMaxAttempts(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get OTP max attempts")
		return
	}

	logger.InfoContext(ctx, "Handler: Get OTP max attempts completed",
		"operation", "get_otp_max_attempts",
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
