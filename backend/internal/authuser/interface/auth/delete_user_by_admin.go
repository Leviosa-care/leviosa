package aggregatorHandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) DeleteUserByAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract logger from context
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract user ID from URL path
	userIDStr := r.PathValue("id")
	if userIDStr == "" {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("user ID is required"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing delete user by admin request",
		"operation", "delete_user_by_admin",
		"user_id", userIDStr,
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid user ID format: %v", err)), http.StatusBadRequest)
		return
	}

	// Call the aggregator service to delete the user
	err = h.svc.DeleteUserByAdmin(ctx, userID)
	if err != nil {
		var statusCode int
		var errorContext string
		var logLevel string

		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			statusCode = http.StatusNotFound
			errorContext = "user not found"
			logLevel = "warn"
		case errors.Is(err, errs.ErrInvalidValue):
			statusCode = http.StatusBadRequest
			errorContext = "validation failed"
			logLevel = "warn"
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
			errorContext = "database connection issue"
			logLevel = "error"
		case errors.Is(err, errs.ErrResourceExhausted):
			statusCode = http.StatusServiceUnavailable
			errorContext = "resource exhausted"
			logLevel = "error"
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			statusCode = http.StatusRequestTimeout
			errorContext = "request cancelled"
			logLevel = "warn"
		case errors.Is(err, context.DeadlineExceeded):
			statusCode = http.StatusRequestTimeout
			errorContext = "request timeout"
			logLevel = "warn"
		case errors.Is(err, errs.ErrTransactionFailure):
			statusCode = http.StatusServiceUnavailable
			errorContext = "transaction failure"
			logLevel = "error"
		default:
			statusCode = http.StatusInternalServerError
			errorContext = "internal server error"
			logLevel = "error"
		}

		logFields := []any{
			"operation", "delete_user_by_admin",
			"user_id", userID.String(),
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", statusCode,
			"error", err,
		}

		switch logLevel {
		case "warn":
			logger.WarnContext(ctx, "Handler: Delete user by admin failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Delete user by admin failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful completion
	logger.InfoContext(ctx, "Handler: Delete user by admin completed",
		"operation", "delete_user_by_admin",
		"user_id", userID.String(),
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Respond with success message
	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
	}{
		Message: "User deleted successfully",
	}, http.StatusOK)
}
