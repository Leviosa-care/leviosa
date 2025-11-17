package priceHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
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
		httpx.RespondWithServiceError(w, logger, ctx, err, "get prices by product ID")
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
