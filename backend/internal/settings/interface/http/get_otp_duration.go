package http

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetOTPDuration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get OTP duration request",
		"operation", "get_otp_duration",
		"method", r.Method,
		"path", r.URL.Path)

	response, err := h.svc.GetOTPDuration(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get OTP duration")
		return
	}

	logger.InfoContext(ctx, "Handler: Get OTP duration completed",
		"operation", "get_otp_duration",
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
