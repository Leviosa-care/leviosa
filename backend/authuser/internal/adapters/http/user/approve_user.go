package userHandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
)

func (h *handler) ApproveUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing approve user request",
		"operation", "approve_user",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Parse request body
	var request domain.ApproveUserRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "approve_user",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call service to approve user
	err = h.svc.ApproveUser(ctx, &request)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		var statusCode int

		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			logLevel = "warn"
			errorContext = "invalid request validation"
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrDomainNotFound):
			logLevel = "warn"
			errorContext = "user not found or not pending"
			statusCode = http.StatusNotFound
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
		case errors.Is(err, errs.ErrNotEncrypted):
			logLevel = "error"
			errorContext = "encryption failure"
			statusCode = http.StatusInternalServerError
		case errors.Is(err, errs.ErrConflict):
			logLevel = "warn"
			errorContext = "user approval conflict"
			statusCode = http.StatusConflict
		case errors.Is(err, errs.ErrExternalService):
			logLevel = "error"
			errorContext = "external service unavailable"
			statusCode = http.StatusServiceUnavailable
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
			"operation", "approve_user",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", statusCode,
			"error", err,
		}

		switch logLevel {
		case "warn":
			logger.WarnContext(ctx, "Handler: Approve user failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Approve user failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Approve user completed",
		"operation", "approve_user",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"user_id", request.UserID,
		"role", request.Role)

	// Return success response
	httpx.RespondWithJSON(w, map[string]string{"message": "User approved successfully"}, http.StatusOK)
}
