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

func (h *handler) AddPartnerSpecialization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract partner ID and specialization ID from URL path
	partnerIDStr := r.PathValue("id")
	partnerID, err := uuid.Parse(partnerIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid partner ID format",
			"error", err,
			"operation", "add_partner_specialization",
			"method", r.Method,
			"path", r.URL.Path,
			"partner_id", partnerIDStr)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
		return
	}

	specializationIDStr := r.PathValue("specializationID")
	specializationID, err := uuid.Parse(specializationIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid specialization ID format",
			"error", err,
			"operation", "add_partner_specialization",
			"method", r.Method,
			"path", r.URL.Path,
			"partner_id", partnerID,
			"specialization_id", specializationIDStr)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid specialization ID format"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing add partner specialization request",
		"operation", "add_partner_specialization",
		"method", r.Method,
		"path", r.URL.Path,
		"partner_id", partnerID,
		"specialization_id", specializationID,
		"user_agent", r.Header.Get("User-Agent"))

	// Call service to add specialization
	err = h.svc.AddPartnerSpecialization(ctx, partnerID, specializationID)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		var statusCode int

		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			logLevel = "warn"
			errorContext = "partner or specialization not found"
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrInvalidInput):
			logLevel = "warn"
			errorContext = "invalid request validation"
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrUniqueViolation):
			logLevel = "warn"
			errorContext = "specialization already assigned to partner"
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
			"operation", "add_partner_specialization",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"partner_id", partnerID,
			"specialization_id", specializationID,
			"status_code", statusCode,
			"error", err,
		}

		switch logLevel {
		case "warn":
			logger.WarnContext(ctx, "Handler: Add partner specialization failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Add partner specialization failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Add partner specialization completed",
		"operation", "add_partner_specialization",
		"method", r.Method,
		"path", r.URL.Path,
		"partner_id", partnerID,
		"specialization_id", specializationID,
		"status_code", http.StatusCreated)

	httpx.RespondWithJSON(w, map[string]string{"message": "Specialization added to partner successfully"}, http.StatusCreated)
}

func (h *handler) RemovePartnerSpecialization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract partner ID and specialization ID from URL path
	partnerIDStr := r.PathValue("id")
	partnerID, err := uuid.Parse(partnerIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid partner ID format",
			"error", err,
			"operation", "remove_partner_specialization",
			"method", r.Method,
			"path", r.URL.Path,
			"partner_id", partnerIDStr)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
		return
	}

	specializationIDStr := r.PathValue("specializationID")
	specializationID, err := uuid.Parse(specializationIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid specialization ID format",
			"error", err,
			"operation", "remove_partner_specialization",
			"method", r.Method,
			"path", r.URL.Path,
			"partner_id", partnerID,
			"specialization_id", specializationIDStr)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid specialization ID format"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing remove partner specialization request",
		"operation", "remove_partner_specialization",
		"method", r.Method,
		"path", r.URL.Path,
		"partner_id", partnerID,
		"specialization_id", specializationID,
		"user_agent", r.Header.Get("User-Agent"))

	// Call service to remove specialization
	err = h.svc.RemovePartnerSpecialization(ctx, partnerID, specializationID)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		var statusCode int

		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			logLevel = "warn"
			errorContext = "partner or specialization association not found"
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
			"operation", "remove_partner_specialization",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"partner_id", partnerID,
			"specialization_id", specializationID,
			"status_code", statusCode,
			"error", err,
		}

		switch logLevel {
		case "warn":
			logger.WarnContext(ctx, "Handler: Remove partner specialization failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Remove partner specialization failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Remove partner specialization completed",
		"operation", "remove_partner_specialization",
		"method", r.Method,
		"path", r.URL.Path,
		"partner_id", partnerID,
		"specialization_id", specializationID,
		"status_code", http.StatusNoContent)

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) GetPartnerSpecializations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract partner ID from URL path
	partnerIDStr := r.PathValue("id")
	partnerID, err := uuid.Parse(partnerIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid partner ID format",
			"error", err,
			"operation", "get_partner_specializations",
			"method", r.Method,
			"path", r.URL.Path,
			"partner_id", partnerIDStr)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get partner specializations request",
		"operation", "get_partner_specializations",
		"method", r.Method,
		"path", r.URL.Path,
		"partner_id", partnerID,
		"user_agent", r.Header.Get("User-Agent"))

	// Call service to get specializations
	specializations, err := h.svc.GetPartnerSpecializations(ctx, partnerID)
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
			"operation", "get_partner_specializations",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"partner_id", partnerID,
			"status_code", statusCode,
			"error", err,
		}

		switch logLevel {
		case "warn":
			logger.WarnContext(ctx, "Handler: Get partner specializations failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Get partner specializations failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get partner specializations completed",
		"operation", "get_partner_specializations",
		"method", r.Method,
		"path", r.URL.Path,
		"partner_id", partnerID,
		"status_code", http.StatusOK,
		"specialization_count", len(specializations.Specializations))

	httpx.RespondWithJSON(w, specializations, http.StatusOK)
}