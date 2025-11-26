package analytics

// TEMPORARILY DISABLED - ARCHITECTURAL MISMATCH
//
// This feature (along with suggestion_algorithm.go) has been moved here from the
// availability service because it violates the separation of concerns.
//
// ISSUE:
// SuggestAvailabilityBlocks queries the product catalog and generates availability
// block recommendations based on product durations. This couples the availability
// layer to products, which should be separate concerns.
//
// The availability layer should be product-agnostic - partners declare time blocks,
// clients book products within those blocks. Product-based scheduling optimization
// belongs elsewhere (booking analytics, partner dashboard, etc.).
//
// TODO: Decide whether to:
// - Reimplement as a UI/dashboard feature (not a backend service)
// - Move to a separate scheduling optimization service
// - Delete if not needed
//
// For now, commented out to prevent compilation errors.

/*
package availability

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

// SuggestAvailabilityBlocks recommends optimal availability block durations
// based on the partner's products and room allocation type
func (s *AvailabilityService) SuggestAvailabilityBlocks(
	ctx context.Context,
	partnerID uuid.UUID,
	roomID uuid.UUID,
) (*domain.AvailabilitySuggestions, error) {
	// 1. Hash partner ID for querying
	userIDBytes, err := encx.SerializeValue(partnerID)
	if err != nil {
		return nil, fmt.Errorf("serialize partner ID for hashing: %w", err)
	}
	userIDHash := s.crypto.HashBasic(ctx, userIDBytes)

	// 2. Get room allocation to determine allocation type
	// Use current time to check if allocation is active
	roomAllocationEncx, err := s.allocationRepo.GetActiveAllocationForPartnerAndRoom(ctx, userIDHash, roomID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("get room allocation: %w", err)
	}

	// Decrypt allocation to access allocation type
	allocation, err := domain.DecryptRoomAllocationEncx(ctx, s.crypto, roomAllocationEncx)
	if err != nil {
		return nil, fmt.Errorf("decrypt room allocation: %w", err)
	}

	// 3. Get partner's products
	products, err := s.productService.GetAllPublishedProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("get partner products: %w", err)
	}

	// 4. Generate block suggestions based on products and allocation type
	suggestions := generateBlockSuggestions(products, allocation.AllocationType)

	// 5. Rank suggestions by priority
	rankSuggestions(suggestions)

	return &domain.AvailabilitySuggestions{
		PartnerID:         partnerID,
		RoomID:            roomID,
		AllocationType:    allocation.AllocationType,
		RecommendedBlocks: suggestions,
	}, nil
}
*/
