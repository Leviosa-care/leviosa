package aggregatorHandler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Leviosa-care/core/auth/cookies"
	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
)

func (h *handler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Get provider from URL path
	provider := r.PathValue("provider")
	if provider == "" {
		logger.WarnContext(ctx, "Handler: Missing provider in OAuth callback request",
			"operation", "oauth_callback",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("provider parameter is required"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing OAuth callback request",
		"provider", provider,
		"operation", "oauth_callback",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Call service
	response, err := h.svc.OAuthCallback(ctx, w, r, provider)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			logLevel = "warn"
			errorContext = "invalid OAuth callback data or validation error"
		case errors.Is(err, errs.ErrUnauthorized):
			logLevel = "warn"
			errorContext = "OAuth authentication failure"
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			logLevel = "error"
			errorContext = "infrastructure connection failure"
		case errors.Is(err, errs.ErrResourceExhausted):
			logLevel = "error"
			errorContext = "infrastructure resource exhaustion"
		case errors.Is(err, errs.ErrExternalService):
			logLevel = "error"
			errorContext = "external OAuth service failure"
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			logLevel = "warn"
			errorContext = "request cancelled"
		case errors.Is(err, context.DeadlineExceeded):
			logLevel = "warn"
			errorContext = "request timeout"
		default:
			logLevel = "error"
			errorContext = "unexpected error"
		}

		logFields := []any{
			"provider", provider,
			"operation", "oauth_callback",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", getStatusCodeForError(err),
			"error", err,
		}

		switch logLevel {
		case "warn":
			logger.WarnContext(ctx, "Handler: OAuth callback request failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: OAuth callback request failed", logFields...)
		}

		var statusCode int
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrUnauthorized):
			statusCode = http.StatusUnauthorized
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrResourceExhausted):
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrExternalService):
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, context.DeadlineExceeded):
			statusCode = http.StatusRequestTimeout
		default:
			statusCode = http.StatusInternalServerError
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: OAuth callback request completed successfully",
		"provider", provider,
		"is_new_user", response.IsNewUser,
		"operation", "oauth_callback",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusCreated)

	// Set dual token cookies
	cookies.SetTokenCookies(w, response.AccessToken, response.RefreshToken,
		time.Unix(response.AccessTokenExpiry, 0), time.Unix(response.RefreshTokenExpiry, 0))

	httpx.RespondWithJSON(w, struct {
		Message   string `json:"message"`
		Status    string `json:"status"`
		IsNewUser bool   `json:"is_new_user"`
	}{
		Message:   "OAuth login successful",
		Status:    "created",
		IsNewUser: response.IsNewUser,
	}, http.StatusCreated)
}

