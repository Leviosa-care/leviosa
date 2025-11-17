package priceHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

// GetPrice handles GET /admin/prices/{id}
func (h *handler) GetPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get price",
		"operation", "get_price",
		"method", r.Method,
		"path", r.URL.Path)

	priceID := strings.Split(r.URL.Path, "/")[3] // Extract internal price ID
	if priceID == "" {
		httpx.RespondWithError(w, errors.New("price ID is missing from URL"), http.StatusBadRequest)
		return
	}

	price, err := h.svc.GetPrice(ctx, priceID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get price")
		return
	}

	logger.InfoContext(ctx, "Handler: get price completed",
		"operation", "get_price",
		"price_id", priceID,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, price, http.StatusOK)
}
