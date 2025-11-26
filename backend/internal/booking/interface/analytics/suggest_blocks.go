package analytics

// TEMPORARILY DISABLED - ARCHITECTURAL MISMATCH
//
// HTTP handler for SuggestBlocks endpoint. See application/analytics/suggest_blocks.go
// for explanation of why this feature is disabled.

/*
package availabilityHandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

// SuggestBlocks handles GET /partners/{partner_id}/rooms/{room_id}/suggest-blocks
func (h *handler) SuggestBlocks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Getting availability block suggestions",
		"operation", "suggest_blocks")

	// Parse partner_id from URL path
	partnerIDStr := r.PathValue("partner_id")
	if partnerIDStr == "" {
		logger.WarnContext(ctx, "Handler: Missing partner_id parameter",
			"operation", "suggest_blocks")
		httpx.RespondWithError(w, errors.New("partner_id is required"), http.StatusBadRequest)
		return
	}

	partnerID, err := uuid.Parse(partnerIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid partner_id format",
			"partner_id", partnerIDStr,
			"error", err.Error(),
			"operation", "suggest_blocks")
		httpx.RespondWithError(w, errors.New("invalid partner_id format"), http.StatusBadRequest)
		return
	}

	// Parse room_id from URL path
	roomIDStr := r.PathValue("room_id")
	if roomIDStr == "" {
		logger.WarnContext(ctx, "Handler: Missing room_id parameter",
			"partner_id", partnerID,
			"operation", "suggest_blocks")
		httpx.RespondWithError(w, errors.New("room_id is required"), http.StatusBadRequest)
		return
	}

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid room_id format",
			"partner_id", partnerID,
			"room_id", roomIDStr,
			"error", err.Error(),
			"operation", "suggest_blocks")
		httpx.RespondWithError(w, errors.New("invalid room_id format"), http.StatusBadRequest)
		return
	}

	// Call service
	suggestions, err := h.svc.SuggestAvailabilityBlocks(ctx, partnerID, roomID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get availability block suggestions")
		return
	}

	// Convert to response DTO
	response := convertToSuggestionsResponse(suggestions)

	logger.InfoContext(ctx, "Handler: Successfully retrieved block suggestions",
		"partner_id", partnerID,
		"room_id", roomID,
		"suggestions_count", len(response.RecommendedBlocks),
		"allocation_type", response.AllocationType,
		"operation", "suggest_blocks")

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.ErrorContext(ctx, "Handler: Failed to encode response",
			"error", err.Error(),
			"operation", "suggest_blocks")
	}
}

// convertToSuggestionsResponse converts domain model to API response DTO
func convertToSuggestionsResponse(suggestions *domain.AvailabilitySuggestions) domain.GetAvailabilitySuggestionsResponse {
	blocks := make([]domain.BlockSuggestionResponse, 0, len(suggestions.RecommendedBlocks))

	for _, block := range suggestions.RecommendedBlocks {
		combos := make([]domain.ProductComboResponse, 0, len(block.ProductCombinations))

		for _, combo := range block.ProductCombinations {
			products := make([]domain.ProductInfoResponse, 0, len(combo.Products))

			for _, product := range combo.Products {
				products = append(products, domain.ProductInfoResponse{
					ID:         product.ID,
					Name:       product.Name,
					Duration:   product.Duration,
					BufferTime: product.BufferTime,
				})
			}

			combos = append(combos, domain.ProductComboResponse{
				Products:      products,
				TotalDuration: combo.TotalDuration,
				SessionCount:  combo.SessionCount,
			})
		}

		blocks = append(blocks, domain.BlockSuggestionResponse{
			DurationMinutes:     block.DurationMinutes,
			Rationale:           block.Rationale,
			ProductCombinations: combos,
			Priority:            block.Priority,
		})
	}

	return domain.GetAvailabilitySuggestionsResponse{
		PartnerID:         suggestions.PartnerID,
		RoomID:            suggestions.RoomID,
		AllocationType:    string(suggestions.AllocationType),
		RecommendedBlocks: blocks,
	}
}
*/
