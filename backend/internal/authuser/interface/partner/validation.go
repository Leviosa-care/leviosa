package partnerHandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

// ValidateProductsRequest represents the request body for validating products
type ValidateProductsRequest struct {
	ProductIDs []string `json:"product_ids"`
}

// ValidateProductsResponse represents the response for validating products
type ValidateProductsResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

// ValidatePartnerProducts validates that all product IDs exist in the catalog cache
func (h *handler) ValidatePartnerProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.RespondWithError(w, errors.New(""), http.StatusMethodNotAllowed)
		return
	}

	var req ValidateProductsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.RespondWithError(w, errors.New(""), http.StatusBadRequest)
		return
	}

	// Convert string IDs to UUIDs
	productIDs := make([]uuid.UUID, 0, len(req.ProductIDs))
	for _, idStr := range req.ProductIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			response := ValidateProductsResponse{
				Valid: false,
				Error: "invalid product ID format: " + idStr,
			}
			httpx.RespondWithJSON(w, errors.New(response.Error), http.StatusBadRequest)
			return
		}
		productIDs = append(productIDs, id)
	}

	// Validate products through the service
	if err := h.svc.ValidatePartnerProducts(r.Context(), productIDs); err != nil {
		response := ValidateProductsResponse{
			Valid: false,
			Error: err.Error(),
		}
		httpx.RespondWithJSON(w, errors.New(response.Error), http.StatusBadRequest)
		return
	}

	// All validations passed
	response := ValidateProductsResponse{
		Valid: true,
	}
	httpx.RespondWithJSON(w, errors.New(response.Error), http.StatusOK)
}

