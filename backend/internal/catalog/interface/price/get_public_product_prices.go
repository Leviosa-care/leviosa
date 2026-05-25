package priceHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

// GetPublicProductPrices handles GET /products/{id}/prices
// Returns only active prices for a product (public access, no authentication required).
func (h *handler) GetPublicProductPrices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get public product prices",
		"operation", "get_public_product_prices",
		"method", r.Method,
		"path", r.URL.Path)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		httpx.RespondWithError(w, errors.New("product ID is missing from URL"), http.StatusBadRequest)
		return
	}
	productID := parts[2] // /products/{id}/prices
	if productID == "" {
		httpx.RespondWithError(w, errors.New("product ID is missing from URL"), http.StatusBadRequest)
		return
	}

	prices, err := h.svc.GetPricesByProductID(ctx, productID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get public product prices")
		return
	}

	// Filter to only return active prices
	activePrices := make([]PublicPriceResponse, 0)
	for _, p := range prices {
		if p.IsActive {
			activePrices = append(activePrices, PublicPriceResponse{
				ID:       p.ID.String(),
				Amount:   p.Amount,
				Currency: p.Currency,
				Interval: string(p.Interval),
			})
		}
	}

	logger.InfoContext(ctx, "Handler: get public product prices completed",
		"operation", "get_public_product_prices",
		"product_id", productID,
		"count", len(activePrices),
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, activePrices, http.StatusOK)
}

// PublicPriceResponse is the DTO returned by the public price endpoint.
// Only exposes fields safe for unauthenticated access (no Stripe IDs).
type PublicPriceResponse struct {
	ID       string `json:"id"`
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
	Interval string `json:"interval"`
}
