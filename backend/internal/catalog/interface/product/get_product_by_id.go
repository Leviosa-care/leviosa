package productHandler

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
			log.Printf("Handler: Internal server error during category retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			// A catch-all for any other unhandled errors.
			log.Printf("Handler: Unhandled error from service during category retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	httpx.RespondWithJSON(w, product, http.StatusOK)
}
