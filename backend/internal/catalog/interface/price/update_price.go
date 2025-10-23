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

// UpdatePrice handles PATCH /admin/prices/{id}
func (h *handler) UpdatePrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	priceID := strings.Split(r.URL.Path, "/")[3] // Extract internal price ID
	if priceID == "" {
		httpx.RespondWithError(w, errors.New("price ID is missing from URL"), http.StatusBadRequest)
		return
	}

	var req domain.UpdatePriceRequest // Use your errs input struct with pointers
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.RespondWithError(w, errors.New("invalid request payload"), http.StatusBadRequest)
		return
	}

	// You might add an early check if req is empty (no updatable fields provided)
	if req.Active == nil && req.Metadata == nil && req.Nickname == nil {
		httpx.RespondWithError(w, errors.New("no updatable fields provided in request body"), http.StatusBadRequest)
		return
	}

	updatedPrice, err := h.svc.UpdatePrice(ctx, priceID, req)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrExternalService):
			log.Printf("Handler: External service error updating price: %v", err)
			httpx.RespondWithError(w, errors.New("failed to update price due to external service issue"), http.StatusServiceUnavailable)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			log.Printf("Handler: Internal server error updating price: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			log.Printf("Handler: Unhandled error updating price: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	httpx.RespondWithJSON(w, updatedPrice, http.StatusOK)
}
