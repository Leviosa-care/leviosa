package productHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) ModifyProduct(w http.ResponseWriter, r *http.Request) {
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

	logger.InfoContext(ctx, "Handler: Processing modify product",
		"operation", "modify_product",
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

	var req domain.UpdateProductRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevent clients from sending unexpected fields
	if err := decoder.Decode(&req); err != nil {
		logger.ErrorContext(ctx, "Handler: modify product failed",
			"operation", "modify_product",
			"error_context", "invalid JSON request body",
			"status_code", http.StatusBadRequest,
			"error", err)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	if req.Name == nil && req.Description == nil && req.Status == nil && req.Metadata == nil {
		httpx.RespondWithError(w, errors.New("no updatable fields provided in request body"), http.StatusBadRequest)
		return
	}

	if err := h.productService.UpdateProduct(ctx, productID, &req); err != nil {
		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrUnexpectedError):
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
		case errors.Is(err, errs.ErrDomainNotUpdated):
			// TODO: change that status, this is bad
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrExternalService):
			httpx.RespondWithError(w, errors.New("failed to delete product images due to external storage issue"), http.StatusServiceUnavailable)
		case errors.Is(err, errs.ErrQueryFailed):
			logger.ErrorContext(ctx, "Handler: modify product failed",
				"operation", "modify_product",
				"error_context", "internal server error during product update",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.ErrorContext(ctx, "Handler: modify product failed",
				"operation", "modify_product",
				"error_context", "internal server error",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("internal server occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.InfoContext(ctx, "Handler: modify product completed",
		"operation", "modify_product",
		"product_id", productID,
		"status_code", http.StatusOK)

	w.WriteHeader(http.StatusNoContent)
}
