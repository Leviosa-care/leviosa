package productHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) RemoveProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing remove product",
		"operation", "remove_product",
		"method", r.Method,
		"path", r.URL.Path)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 || parts[0] != "" || parts[1] != "admin" || parts[2] != "products" {
		httpx.RespondWithError(w, errors.New("invalid URL path format. Expected /admin/products/{id}"), http.StatusBadRequest)
		return
	}
	productID := parts[3] // The ID should be the last part
	if productID == "" {
		httpx.RespondWithError(w, errors.New("product ID is missing from the URL"), http.StatusBadRequest)
		return
	}

	if err := h.productService.RemoveProduct(ctx, productID); err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "remove product")
		return
	}

	logger.InfoContext(ctx, "Handler: remove product completed",
		"operation", "remove_product",
		"product_id", productID,
		"status_code", http.StatusOK)

	w.WriteHeader(http.StatusNoContent)
}
