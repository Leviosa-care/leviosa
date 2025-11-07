package partnerHandler

import (
	"context"
	"errors"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAllPartnersByCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract category IDs from query parameters
	categoryIDs := r.URL.Query()["category_id"]
	if len(categoryIDs) == 0 {
		logger.WarnContext(ctx, "Handler: Missing category IDs in request",
			"operation", "get_partners_by_categories",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("at least one category_id query parameter is required"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get partners by categories request",
		"operation", "get_partners_by_categories",
		"method", r.Method,
		"path", r.URL.Path,
		"category_ids", categoryIDs,
		"category_count", len(categoryIDs),
		"user_agent", r.Header.Get("User-Agent"))

	partners, err := h.svc.GetAllPartnersByCategories(ctx, categoryIDs)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		var statusCode int

		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			logLevel = "warn"
			errorContext = "invalid category ID format"
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
			"operation", "get_partners_by_categories",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"category_ids", categoryIDs,
			"category_count", len(categoryIDs),
			"status_code", statusCode,
			"error", err,
		}

		switch logLevel {
		case "warn":
			logger.WarnContext(ctx, "Handler: Get partners by categories failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Get partners by categories failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get partners by categories completed",
		"operation", "get_partners_by_categories",
		"method", r.Method,
		"path", r.URL.Path,
		"category_ids", categoryIDs,
		"category_count", len(categoryIDs),
		"status_code", http.StatusOK,
		"partner_count", len(partners))

	httpx.RespondWithJSON(w, partners, http.StatusOK)
}
