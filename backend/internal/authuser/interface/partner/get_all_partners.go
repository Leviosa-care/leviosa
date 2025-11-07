package partnerHandler

import (
	"context"
	"errors"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAllPartners(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get all partners request",
		"operation", "get_all_partners",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	partners, err := h.svc.GetAllPartners(ctx)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		var statusCode int

		switch {
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			logLevel = "error"
			errorContext = "database connection failure"
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrResourceExhausted):
			logLevel = "error"
			errorContext = "database resource exhaustion"
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			logLevel = "warn"
			errorContext = "request cancelled"
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, context.DeadlineExceeded):
			logLevel = "warn"
			errorContext = "request timeout"
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			logLevel = "error"
			errorContext = "database transaction failure"
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrPermissionDenied):
			logLevel = "error"
			errorContext = "database permission denied"
			statusCode = http.StatusInternalServerError
		case errors.Is(err, errs.ErrDatabase):
			logLevel = "error"
			errorContext = "general database error"
			statusCode = http.StatusInternalServerError
		case errors.Is(err, errs.ErrNotDecrypted):
			logLevel = "error"
			errorContext = "data decryption failure"
			statusCode = http.StatusInternalServerError
		default:
			logLevel = "error"
			errorContext = "unexpected error"
			statusCode = http.StatusInternalServerError
		}

		logFields := []any{
			"operation", "get_all_partners",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", statusCode,
			"error", err,
		}

		switch logLevel {
		case "warn":
			logger.WarnContext(ctx, "Handler: Get all partners failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Get all partners failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get all partners completed",
		"operation", "get_all_partners",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"partner_count", len(partners))

	httpx.RespondWithJSON(w, partners, http.StatusOK)
}
