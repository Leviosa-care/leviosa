package productHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get product by id",
		"operation", "get_product_by_id",
		"method", r.Method,
		"path", r.URL.Path)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[0] != "" || parts[1] != "products" {
		httpx.RespondWithError(w, errors.New("invalid URL path format. Expected /products/{id}"), http.StatusBadRequest)
		return
	}
	productID := parts[2]
	if productID == "" {
		httpx.RespondWithError(w, errors.New("product ID is missing from the URL"), http.StatusBadRequest)
		return
	}

	product, err := h.aggr.GetProductByID(ctx, productID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get product by ID")
		return
	}

	logger.InfoContext(ctx, "Handler: get product by id completed",
		"operation", "get_product_by_id",
		"product_id", productID,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, product, http.StatusOK)
}
