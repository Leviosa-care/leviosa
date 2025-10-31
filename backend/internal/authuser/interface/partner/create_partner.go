package partnerHandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) CreatePartner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing create partner request",
		"operation", "create_partner",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Parse request body
	var request domain.CreatePartnerRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "create_partner",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body, user ID wrong formatted",
			"error", err,
			"operation", "create_partner",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Call service to create partner
	partner, err := h.svc.CreatePartner(ctx,
		userID,
		request.Bio,
		request.Experience,
		// request.Certifications,
		request.CategoryIDs,
		request.ProductIDs,
	)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		var statusCode int

		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			logLevel = "warn"
			errorContext = "invalid request validation"
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrUniqueViolation):
			logLevel = "warn"
			errorContext = "user email already exists"
			statusCode = http.StatusConflict
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
			errorContext = "permission denied"
			statusCode = http.StatusForbidden
		case errors.Is(err, errs.ErrInvalidValue):
			logLevel = "error"
			errorContext = "encryption failure"
			statusCode = http.StatusInternalServerError
		case errors.Is(err, errs.ErrDatabase):
			logLevel = "error"
			errorContext = "general database error"
			statusCode = http.StatusInternalServerError
		default:
			logLevel = "error"
			errorContext = "unexpected error"
			statusCode = http.StatusInternalServerError
		}

		logFields := []any{
			"operation", "create_partner",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", statusCode,
			"user_id", request.UserID,
			"error", err,
		}

		switch logLevel {
		case "warn":
			logger.WarnContext(ctx, "Handler: Create partner failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Create partner failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Create partner completed",
		"operation", "create_partner",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusCreated,
		"user_id", partner.UserID)

	httpx.RespondWithJSON(w, partner, http.StatusCreated)
}
