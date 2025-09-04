package aggregatorHandler

import (
	"context"
	"errors"
	"net/http"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
	"github.com/Leviosa-care/core/middleware/auth"
)

func (h *handler) RefreshSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract refresh token from cookies
	refreshToken, err := auth.GetRefreshTokenFromCookies(r)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Missing refresh token cookie",
			"error", err,
			"operation", "refresh_session",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewUnauthorizedErr("missing refresh token"), http.StatusUnauthorized)
		return
	}

	request := &domain.RefreshSessionRequest{
		RefreshToken: refreshToken,
	}

	logger.InfoContext(ctx, "Handler: Processing session refresh request",
		"operation", "refresh_session",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	response, err := h.svc.RefreshSession(ctx, request)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			logLevel = "warn"
			errorContext = "validation failed"
		case errors.Is(err, errs.ErrDomainNotFound):
			logLevel = "warn"
			errorContext = "session not found"
		case errors.Is(err, errs.ErrUnauthorized):
			logLevel = "warn"
			errorContext = "invalid session state"
		case errors.Is(err, errs.ErrExpiredToken):
			logLevel = "info"
			errorContext = "refresh token expired"
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			logLevel = "error"
			errorContext = "infrastructure connection failure"
		case errors.Is(err, errs.ErrResourceExhausted):
			logLevel = "error"
			errorContext = "infrastructure resource exhaustion"
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			logLevel = "warn"
			errorContext = "request cancelled"
		case errors.Is(err, context.DeadlineExceeded):
			logLevel = "warn"
			errorContext = "request timeout"
		case errors.Is(err, errs.ErrTransactionFailure):
			logLevel = "error"
			errorContext = "infrastructure transaction failure"
		default:
			logLevel = "error"
			errorContext = "unexpected error"
		}

		logFields := []any{
			"operation", "refresh_session",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"error", err,
		}

		var statusCode int
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrDomainNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrUnauthorized):
			statusCode = http.StatusUnauthorized
		case errors.Is(err, errs.ErrExpiredToken):
			statusCode = http.StatusUnauthorized
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrResourceExhausted):
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, context.DeadlineExceeded):
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, errs.ErrTransactionFailure):
			statusCode = http.StatusServiceUnavailable
		default:
			statusCode = http.StatusInternalServerError
		}

		logFields = append(logFields, "status_code", statusCode)

		switch logLevel {
		case "info":
			logger.InfoContext(ctx, "Handler: Session refresh request result", logFields...)
		case "warn":
			logger.WarnContext(ctx, "Handler: Session refresh request failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Session refresh request failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Session refresh request completed successfully",
		"operation", "refresh_session",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Set new token cookies
	auth.SetTokenCookies(w, response.AccessToken, response.RefreshToken,
		response.AccessTokenExpiry, response.RefreshTokenExpiry)

	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}{
		Message: "Session refreshed successfully",
		Status:  "success",
	}, http.StatusOK)
}

