package aggregatorHandler

import (
	"context"
	"errors"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) OAuthStart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Get provider from URL path
	provider := r.PathValue("provider")
	if provider == "" {
		logger.WarnContext(ctx, "Handler: Missing provider in OAuth start request",
			"operation", "oauth_start",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("provider parameter is required"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing OAuth start request",
		"provider", provider,
		"operation", "oauth_start",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Create request
	request := &domain.OAuthStartRequest{
		Provider: provider,
	}

	// Call service
	response, err := h.svc.OAuthStart(ctx, request)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			logLevel = "warn"
			errorContext = "invalid provider or validation error"
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
			"operation", "oauth_start",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", getStatusCodeForError(err),
			"error", err,
		}

		switch logLevel {
		case "warn":
			logger.WarnContext(ctx, "Handler: OAuth start request failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: OAuth start request failed", logFields...)
		}

		var statusCode int
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			statusCode = http.StatusBadRequest
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
	logger.InfoContext(ctx, "Handler: OAuth start request completed successfully",
		"provider", provider,
		"operation", "oauth_start",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusFound,
		"redirect_url", response.AuthorizationURL)

	// Redirect to OAuth provider's authorization URL
	http.Redirect(w, r, response.AuthorizationURL, http.StatusFound)
}