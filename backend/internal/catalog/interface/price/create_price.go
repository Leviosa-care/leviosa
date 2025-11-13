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

// CreatePrice handles POST /admin/products/{id}/prices
func (h *handler) CreatePrice(w http.ResponseWriter, r *http.Request) {
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

	logger.InfoContext(ctx, "Handler: Processing create price",
		"operation", "create_price",
		"method", r.Method,
		"path", r.URL.Path)

	productID := strings.Split(r.URL.Path, "/")[3] // Extract internal product ID
	if productID == "" {
		httpx.RespondWithError(w, errors.New("product ID is missing from URL"), http.StatusBadRequest)
		return
	}

	var request domain.CreatePriceRequest // Use your domain input struct
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.ErrorContext(ctx, "Handler: create price failed",
			"operation", "create_price",
			"error_context", "invalid JSON request body",
			"status_code", http.StatusBadRequest,
			"error", err)
		httpx.RespondWithError(w, errors.New("invalid request payload"), http.StatusBadRequest)
		return
	}

	newPrice, err := h.svc.CreatePrice(ctx, productID, &request)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound) // Product not found
		case errors.Is(err, errs.ErrConflict): // If Stripe price ID unique constraint check fails
			httpx.RespondWithError(w, err, http.StatusConflict)
		case errors.Is(err, errs.ErrExternalService):
			logger.ErrorContext(ctx, "Handler: create price failed",
				"operation", "create_price",
				"error_context", "external service error creating price",
				"status_code", http.StatusServiceUnavailable,
				"error", err)
			httpx.RespondWithError(w, errors.New("failed to create price due to external service issue"), http.StatusServiceUnavailable)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.ErrorContext(ctx, "Handler: create price failed",
				"operation", "create_price",
				"error_context", "internal server error creating price",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.ErrorContext(ctx, "Handler: create price failed",
				"operation", "create_price",
				"error_context", "unhandled error creating price",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.InfoContext(ctx, "Handler: create price completed",
		"operation", "create_price",
		"price_id", newPrice,
		"status_code", http.StatusCreated)

	httpx.RespondWithJSON(w, newPrice, http.StatusCreated)
}
