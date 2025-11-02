package productHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
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
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			// This covers an empty ID.
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			// The category ID from the URL was not found.
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			// A general internal server error occurred (DB, corrupt data, etc.).
			logger.ErrorContext(ctx, "Handler: get product by id failed",
				"operation", "get_product_by_id",
				"error_context", "internal server error during product retrieval",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			// A catch-all for any other unhandled errors.
			logger.ErrorContext(ctx, "Handler: get product by id failed",
				"operation", "get_product_by_id",
				"error_context", "unexpected error from service during product retrieval",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.InfoContext(ctx, "Handler: get product by id completed",
		"operation", "get_product_by_id",
		"product_id", productID,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, product, http.StatusOK)
}
