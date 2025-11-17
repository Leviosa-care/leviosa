package productHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAdminAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get admin all products",
		"operation", "get_admin_all_products",
		"method", r.Method,
		"path", r.URL.Path)

	products, err := h.aggr.GetAdminAllProducts(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get admin all products")
		return
	}

	count := len(products)
	logger.InfoContext(ctx, "Handler: get admin all products completed",
		"operation", "get_admin_all_products",
		"count", count,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, products, http.StatusOK)
}
