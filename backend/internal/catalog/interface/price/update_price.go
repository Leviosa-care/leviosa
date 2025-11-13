package priceHandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

// UpdatePrice handles PATCH /admin/prices/{id}
func (h *handler) UpdatePrice(w http.ResponseWriter, r *http.Request) {
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

	logger.InfoContext(ctx, "Handler: Processing update price",
		"operation", "update_price",
		"method", r.Method,
		"path", r.URL.Path)

	priceID := strings.Split(r.URL.Path, "/")[3] // Extract internal price ID
	if priceID == "" {
		httpx.RespondWithError(w, errors.New("price ID is missing from URL"), http.StatusBadRequest)
		return
	}

	var req domain.UpdatePriceRequest // Use your errs input struct with pointers
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorContext(ctx, "Handler: update price failed",
			"operation", "update_price",
			"error_context", "invalid JSON request body",
			"status_code", http.StatusBadRequest,
			"error", err)
		httpx.RespondWithError(w, errors.New("invalid request payload"), http.StatusBadRequest)
		return
	}

	// You might add an early check if req is empty (no updatable fields provided)
	if req.Active == nil && req.Metadata == nil && req.Nickname == nil {
		httpx.RespondWithError(w, errors.New("no updatable fields provided in request body"), http.StatusBadRequest)
		return
	}

	updatedPrice, err := h.svc.UpdatePrice(ctx, priceID, req)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrExternalService):
			logger.ErrorContext(ctx, "Handler: update price failed",
				"operation", "update_price",
				"error_context", "external service error updating price",
				"status_code", http.StatusServiceUnavailable,
				"error", err)
			httpx.RespondWithError(w, errors.New("failed to update price due to external service issue"), http.StatusServiceUnavailable)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.ErrorContext(ctx, "Handler: update price failed",
				"operation", "update_price",
				"error_context", "internal server error updating price",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.ErrorContext(ctx, "Handler: update price failed",
				"operation", "update_price",
				"error_context", "unhandled error updating price",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.InfoContext(ctx, "Handler: update price completed",
		"operation", "update_price",
		"price_id", priceID,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, updatedPrice, http.StatusOK)
}
