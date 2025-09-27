package specializationHandler

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
	"github.com/google/uuid"
)

func (h *handler) UpdateSpecialization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract specialization ID from URL path
	specializationIDStr := r.PathValue("id")
	specializationID, err := uuid.Parse(specializationIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid specialization ID format",
			"error", err,
			"operation", "update_specialization",
			"method", r.Method,
			"path", r.URL.Path,
			"specialization_id", specializationIDStr)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid specialization ID format"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing update specialization request",
		"operation", "update_specialization",
		"method", r.Method,
		"path", r.URL.Path,
		"specialization_id", specializationID,
		"user_agent", r.Header.Get("User-Agent"))

	// Parse request body
	var request domain.UpdateSpecializationRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "update_specialization",
			"method", r.Method,
			"path", r.URL.Path,
			"specialization_id", specializationID)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call service to update specialization
	specialization, err := h.svc.UpdateSpecialization(ctx, specializationID, &request)
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
		case errors.Is(err, errs.ErrRepositoryNotFound):
			logLevel = "warn"
			errorContext = "specialization not found"
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrUniqueViolation):
			logLevel = "warn"
			errorContext = "specialization name already exists"
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
			"operation", "update_specialization",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"specialization_id", specializationID,
			"status_code", statusCode,
			"error", err,
		}

		switch logLevel {
		case "warn":
			logger.WarnContext(ctx, "Handler: Update specialization failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Update specialization failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Update specialization completed",
		"operation", "update_specialization",
		"method", r.Method,
		"path", r.URL.Path,
		"specialization_id", specializationID,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, specialization, http.StatusOK)
}

func (h *handler) DeleteSpecialization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract specialization ID from URL path
	specializationIDStr := r.PathValue("id")
	specializationID, err := uuid.Parse(specializationIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid specialization ID format",
			"error", err,
			"operation", "delete_specialization",
			"method", r.Method,
			"path", r.URL.Path,
			"specialization_id", specializationIDStr)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid specialization ID format"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing delete specialization request",
		"operation", "delete_specialization",
		"method", r.Method,
		"path", r.URL.Path,
		"specialization_id", specializationID,
		"user_agent", r.Header.Get("User-Agent"))

	// Call service to delete specialization
	err = h.svc.DeleteSpecialization(ctx, specializationID)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		var statusCode int

		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			logLevel = "warn"
			errorContext = "specialization not found"
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrForeignKeyViolation):
			logLevel = "warn"
			errorContext = "specialization is still in use by partners"
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
			"operation", "delete_specialization",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"specialization_id", specializationID,
			"status_code", statusCode,
			"error", err,
		}

		switch logLevel {
		case "warn":
			logger.WarnContext(ctx, "Handler: Delete specialization failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Delete specialization failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Delete specialization completed",
		"operation", "delete_specialization",
		"method", r.Method,
		"path", r.URL.Path,
		"specialization_id", specializationID,
		"status_code", http.StatusNoContent)

	w.WriteHeader(http.StatusNoContent)
}