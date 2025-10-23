package priceHandler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

// CreatePrice handles POST /admin/products/{id}/prices
func (h *handler) CreatePrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	productID := strings.Split(r.URL.Path, "/")[3] // Extract internal product ID
	if productID == "" {
		httpx.RespondWithError(w, errors.New("product ID is missing from URL"), http.StatusBadRequest)
		return
	}

	var request domain.CreatePriceRequest // Use your domain input struct
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		httpx.RespondWithError(w, errors.New("invalid request payload"), http.StatusBadRequest)
		return
	}

	newPrice, err := h.svc.CreatePrice(ctx, productID, &request)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound) // Product not found
		case errors.Is(err, errs.ErrConflict): // If Stripe price ID unique constraint check fails
			httpx.RespondWithError(w, err, http.StatusConflict)
		case errors.Is(err, errs.ErrExternalService):
			log.Printf("Handler: External service error creating price: %v", err)
			httpx.RespondWithError(w, errors.New("failed to create price due to external service issue"), http.StatusServiceUnavailable)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			log.Printf("Handler: Internal server error creating price: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			log.Printf("Handler: Unhandled error creating price: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	httpx.RespondWithJSON(w, newPrice, http.StatusCreated)
}
