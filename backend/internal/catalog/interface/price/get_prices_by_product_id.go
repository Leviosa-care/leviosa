package priceHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

// ListPricesForProduct handles GET /admin/products/{id}/prices
func (h *handler) GetPricesByProductID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get prices by product id",
		"operation", "get_prices_by_product_id",
		"method", r.Method,
		"path", r.URL.Path)

	productID := strings.Split(r.URL.Path, "/")[3] // Extract internal product ID
	if productID == "" {
		httpx.RespondWithError(w, errors.New("product ID is missing from URL"), http.StatusBadRequest)
		return
	}

	prices, err := h.svc.GetPricesByProductID(ctx, productID) // Service might take options or always active
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrConflict):
			httpx.RespondWithError(w, err, http.StatusConflict)
		case errors.Is(err, errs.ErrDomainNotFound):
			// If the price isn't found
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.ErrorContext(ctx, "Handler: get prices by product id failed",
				"operation", "get_prices_by_product_id",
				"error_context", "internal server error listing prices",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.ErrorContext(ctx, "Handler: get prices by product id failed",
				"operation", "get_prices_by_product_id",
				"error_context", "unhandled error listing prices",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	count := len(prices)
	logger.InfoContext(ctx, "Handler: get prices by product id completed",
		"operation", "get_prices_by_product_id",
		"product_id", productID,
		"count", count,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, prices, http.StatusOK)
}
