package http

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetOTPLength(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get OTP length request",
		"operation", "get_otp_length",
		"method", r.Method,
		"path", r.URL.Path)

	response, err := h.svc.GetOTPLength(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get OTP length")
		return
	}

	logger.InfoContext(ctx, "Handler: Get OTP length completed",
		"operation", "get_otp_length",
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
