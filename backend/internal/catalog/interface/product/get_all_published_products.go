package productHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAllPublishedProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get all published products",
		"operation", "get_all_published_products",
		"method", r.Method,
		"path", r.URL.Path)

	products, err := h.aggr.GetAllPublishedProducts(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get all published products")
		return
	}

	count := len(products)
	logger.InfoContext(ctx, "Handler: get all published products completed",
		"operation", "get_all_published_products",
		"count", count,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, products, http.StatusOK)
}
