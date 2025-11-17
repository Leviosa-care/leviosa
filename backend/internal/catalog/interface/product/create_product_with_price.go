package productHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) CreateProductWithPrice(w http.ResponseWriter, r *http.Request) {
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

	logger.InfoContext(ctx, "Handler: Processing create product with price",
		"operation", "create_product_with_price",
		"method", r.Method,
		"path", r.URL.Path)

	var request domain.CreateProductWithPriceRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevent clients from sending unexpected fields
	if err := decoder.Decode(&request); err != nil {
		logger.ErrorContext(ctx, "Handler: create product with price failed",
			"operation", "create_product_with_price",
			"error_context", "invalid JSON request body",
			"status_code", http.StatusBadRequest,
			"error", err)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// productID, priceID, err := h.productPriceService.CreateProductWithPrice(ctx, &request)
	productID, priceID, err := h.aggr.CreateProductWithPrice(ctx, &request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "create product with price")
		return
	}

	logger.InfoContext(ctx, "Handler: create product with price completed",
		"operation", "create_product_with_price",
		"product_id", productID,
		"status_code", http.StatusCreated)

	httpx.RespondWithJSON(
		w,
		struct {
			ProductID string `json:"product_id"`
			PriceID   string `json:"price_id"`
			Message   string `json:"message"`
		}{
			ProductID: productID,
			PriceID:   priceID,
			Message:   "Product created successfully!",
		},
		http.StatusCreated,
	)
}
