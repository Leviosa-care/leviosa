package partnerHandler

import (
	"context"
	"errors"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) GetPartnerByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract user ID from URL path
	partnerIDStr := r.PathValue("id")
	partnerID, err := uuid.Parse(partnerIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid partner ID format",
			"error", err,
			"operation", "get_partner_by_id",
			"method", r.Method,
			"path", r.URL.Path,
			"user_id", partnerIDStr)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get partner by ID request",
		"operation", "get_partner_by_id",
		"method", r.Method,
		"path", r.URL.Path,
		"user_id", partnerID,
		"user_agent", r.Header.Get("User-Agent"))

	// partner, err := h.svc.GetPartnerByUserID(ctx, userID)
	partner, err := h.svc.GetPartnerByID(ctx, partnerID)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		var statusCode int

		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			logLevel = "warn"
			errorContext = "partner not found"
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrInvalidValue):
			logLevel = "warn"
			errorContext = "invalid request validation"
			statusCode = http.StatusBadRequest
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
		case errors.Is(err, errs.ErrInvalidInput):
			logLevel = "warn"
			errorContext = "invalid input data"
			statusCode = http.StatusBadRequest
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
			"operation", "get_partner_by_id",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"user_id", partnerID,
			"status_code", statusCode,
			"error", err,
		}

		switch logLevel {
		case "warn":
			logger.WarnContext(ctx, "Handler: Get partner by ID failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Get partner by ID failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get partner by ID completed",
		"operation", "get_partner_by_id",
		"method", r.Method,
		"path", r.URL.Path,
		"user_id", partnerID,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, partner, http.StatusOK)
}
