package partnerHandler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

// ValidateSpecializationsRequest represents the request body for validating specializations
type ValidateSpecializationsRequest struct {
	SpecializationIDs []string `json:"specialization_ids"`
}

// ValidateSpecializationsResponse represents the response for validating specializations
type ValidateSpecializationsResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

// ValidateProductsRequest represents the request body for validating products
type ValidateProductsRequest struct {
	ProductIDs []string `json:"product_ids"`
}

// ValidateProductsResponse represents the response for validating products
type ValidateProductsResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

// ValidatePartnerSpecializations validates that all specialization IDs exist in the catalog cache
func (h *handler) ValidatePartnerSpecializations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.RespondWithError(w, http.StatusMethodNotAllowed, httpx.ErrMethodNotAllowed)
		return
	}

	var req ValidateSpecializationsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.RespondWithError(w, http.StatusBadRequest, httpx.ErrJSONDecoding)
		return
	}

	// Convert string IDs to UUIDs
	specializationIDs := make([]uuid.UUID, 0, len(req.SpecializationIDs))
	for _, idStr := range req.SpecializationIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			response := ValidateSpecializationsResponse{
				Valid: false,
				Error: "invalid specialization ID format: " + idStr,
			}
			httpx.RespondWithJSON(w, http.StatusBadRequest, response)
			return
		}
		specializationIDs = append(specializationIDs, id)
	}

	// Validate specializations through the service
	if err := h.svc.ValidatePartnerSpecializations(r.Context(), specializationIDs); err != nil {
		response := ValidateSpecializationsResponse{
			Valid: false,
			Error: err.Error(),
		}
		httpx.RespondWithJSON(w, http.StatusBadRequest, response)
		return
	}

	// All validations passed
	response := ValidateSpecializationsResponse{
		Valid: true,
	}
	httpx.RespondWithJSON(w, http.StatusOK, response)
}

// ValidatePartnerProducts validates that all product IDs exist in the catalog cache
func (h *handler) ValidatePartnerProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.RespondWithError(w, http.StatusMethodNotAllowed, httpx.ErrMethodNotAllowed)
		return
	}

	var req ValidateProductsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.RespondWithError(w, http.StatusBadRequest, httpx.ErrJSONDecoding)
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
			httpx.RespondWithJSON(w, http.StatusBadRequest, response)
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
		httpx.RespondWithJSON(w, http.StatusBadRequest, response)
		return
	}

	// All validations passed
	response := ValidateProductsResponse{
		Valid: true,
	}
	httpx.RespondWithJSON(w, http.StatusOK, response)
}