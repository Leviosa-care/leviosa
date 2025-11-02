package productHandler

import (
	"errors"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
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
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			logger.ErrorContext(ctx, "Handler: get all published products failed",
				"operation", "get_all_published_products",
				"error_context", "invalid value error",
				"status_code", http.StatusBadRequest,
				"error", err)
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrQueryFailed):
			logger.ErrorContext(ctx, "Handler: get all published products failed",
				"operation", "get_all_published_products",
				"error_context", "query failed",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
		default:
			logger.ErrorContext(ctx, "Handler: get all published products failed",
				"operation", "get_all_published_products",
				"error_context", "internal server error",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("internal server occurred"), http.StatusInternalServerError)
		}
		return
	}

	count := len(products)
	logger.InfoContext(ctx, "Handler: get all published products completed",
		"operation", "get_all_published_products",
		"count", count,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, products, http.StatusOK)
}
