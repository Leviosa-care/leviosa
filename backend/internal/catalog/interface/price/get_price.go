package priceHandler

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

// GetPrice handles GET /admin/prices/{id}
func (h *handler) GetPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	priceID := strings.Split(r.URL.Path, "/")[3] // Extract internal price ID
	if priceID == "" {
		httpx.RespondWithError(w, errors.New("price ID is missing from URL"), http.StatusBadRequest)
		return
	}

	price, err := h.svc.GetPrice(ctx, priceID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			log.Printf("Handler: Internal server error getting price: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			log.Printf("Handler: Unhandled error getting price: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	httpx.RespondWithJSON(w, price, http.StatusOK)
}
