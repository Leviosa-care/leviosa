package userHandler

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

func (h *handler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httpx.RespondWithError(w, errors.New("unsupported media type: please send 'application/json'"), http.StatusUnsupportedMediaType)
		return
	}

	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract user ID from path parameter
	userIDStr := r.PathValue("id")
	if userIDStr == "" {
		logger.WarnContext(ctx, "Handler: Missing user ID in path",
			"operation", "update_user_role",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Parse user ID as UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid user ID format",
			"operation", "update_user_role",
			"method", r.Method,
			"path", r.URL.Path,
			"user_id", userIDStr,
			"error", err)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing update user role request",
		"operation", "update_user_role",
		"method", r.Method,
		"path", r.URL.Path,
		"user_id", userID,
		"user_agent", r.Header.Get("User-Agent"))

	// Parse request body
	var requestBody struct {
		Role string `json:"role"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&requestBody); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "update_user_role",
			"method", r.Method,
			"path", r.URL.Path,
			"user_id", userID)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Create service request
	request := &domain.UpdateUserRoleRequest{
		UserID: userID,
		Role:   requestBody.Role,
	}

	// Call service to update user role
	err = h.svc.UpdateUserRole(ctx, request)
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
			errorContext = "user not found"
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
		case errors.Is(err, errs.ErrNotEncrypted), errors.Is(err, errs.ErrNotDecrypted):
			logLevel = "error"
			errorContext = "encryption/decryption failure"
			statusCode = http.StatusInternalServerError
		case errors.Is(err, errs.ErrConflict):
			logLevel = "warn"
			errorContext = "user role update conflict"
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
			"operation", "update_user_role",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"user_id", userID,
			"role", requestBody.Role,
			"status_code", statusCode,
			"error", err,
		}

		switch logLevel {
		case "warn":
			logger.WarnContext(ctx, "Handler: Update user role failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Update user role failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Update user role completed",
		"operation", "update_user_role",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"user_id", userID,
		"role", requestBody.Role)

	// Return success response
	httpx.RespondWithJSON(w, map[string]string{"message": "User role updated successfully"}, http.StatusOK)
}
