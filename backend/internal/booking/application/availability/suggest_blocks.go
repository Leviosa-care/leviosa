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
