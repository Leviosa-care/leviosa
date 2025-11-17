package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"
)

func (h *handler) SetOTPMaxAttempts(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httpx.RespondWithError(w, errors.New("unsupported media type: please send 'application/json'"), http.StatusUnsupportedMediaType)
		return
	}

	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing set OTP max attempts request",
		"operation", "set_otp_max_attempts",
		"method", r.Method,
		"path", r.URL.Path)

	var request domain.SetOTPMaxAttemptsRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.DebugContext(ctx, fmt.Sprintf("Handler: Error decoding JSON body: %v", err))
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	response, err := h.svc.SetOTPMaxAttempts(ctx, &request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "set OTP max attempts")
		return
	}

	logger.InfoContext(ctx, "Handler: Set OTP max attempts completed",
		"operation", "set_otp_max_attempts",
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
