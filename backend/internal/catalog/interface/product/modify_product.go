package productHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) ModifyProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
		log.Printf("Handler: Error decoding JSON body: %v", err)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	if req.Name == nil && req.Description == nil && req.Status == nil && req.Metadata == nil {
		httpx.RespondWithError(w, errors.New("no updatable fields provided in request body"), http.StatusBadRequest)
		return
	}

	if err := h.productService.UpdateProduct(ctx, productID, &req); err != nil {
		switch {
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
			log.Printf("Handler: Internal server error during product deletion: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			httpx.RespondWithError(w, errors.New("internal server occurred"), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
	log.Println("Handler: Product updated successfully.")
}
