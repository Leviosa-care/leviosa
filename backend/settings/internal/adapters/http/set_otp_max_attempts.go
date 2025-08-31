package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/middleware"
)

func (h *handler) SetOTPMaxAttempts(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		middleware.RespondWithError(w, errors.New("unsupported media type: please send 'application/json'"), http.StatusUnsupportedMediaType)
		return
	}

	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		middleware.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	var request domain.SetOTPMaxAttemptsRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.DebugContext(ctx, fmt.Sprintf("Handler: Error decoding JSON body: %v", err))
		middleware.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	response, err := h.svc.SetOTPMaxAttempts(ctx, &request)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			middleware.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.ErrorContext(ctx, fmt.Sprintf("Handler: Internal server error during OTP max attempts update: %v", err))
			middleware.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.ErrorContext(ctx, fmt.Sprintf("Handler: Unhandled error from service during OTP max attempts update: %v", err))
			middleware.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	middleware.RespondWithJSON(w, response, http.StatusOK)
}

