package price

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// UpdatePrice updates an existing price in your database and potentially in Stripe.
func (s *PriceService) UpdatePrice(ctx context.Context, priceID string, input domain.UpdatePriceRequest) (*domain.Price, error) {
	if err := uuid.Validate(priceID); err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("price ID is invalid: %s", err.Error()))
	}

	// 1. Get existing price from DB to get StripePriceID
	existingPrice, err := s.repo.GetPrice(ctx, priceID)
	if err != nil {
		return nil, fmt.Errorf("get existing price: %w", err)
	}

	// 2. Prepare Stripe update request
	// stripeUpdateReq := domain.UpdateStripePriceRequest{
	// 	Active:   input.Active,
	// 	Metadata: input.Metadata,
	// 	Nickname: input.Nickname,
	// }

	// stripeUpdateReq := domain.UpdatePriceInput{
	// 	Active:   input.Active,
	// 	Metadata: input.Metadata,
	// 	Nickname: input.Nickname,
	// }

	stripeUpdateReq := domain.UpdateStripePriceRequest{
		Active:   input.Active,
		Metadata: input.Metadata,
		Nickname: input.Nickname,
	}

	// 3. Update in Stripe (only updatable fields like 'active' can be changed)
	// IMPORTANT: Only call Stripe if active status is actually changing or other updatable fields are present.
	// Stripe.UpdatePrice takes `UpdateStripePriceRequest` which has pointers, so it handles "not provided".
	paymentPrice, err := s.stripe.UpdatePrice(ctx, existingPrice.StripePriceID, stripeUpdateReq)
	_ = paymentPrice
	if err != nil {
		return nil, errs.NewExternalServiceErr(fmt.Errorf("failed to update price in Stripe: %w", err), "Stripe API")
	}

	// 4. Update in local database
	// The `errs.UpdatePriceInput` is designed for local partial update.
	if err := s.repo.UpdatePrice(ctx, priceID, &input); err != nil {
		return nil, fmt.Errorf("update price in database: %w", err)
	}

	// Re-fetch the updated price from DB or construct it from existing + input
	updatedPrice, err := s.repo.GetPrice(ctx, priceID) // Simplest way to get latest state
	if err != nil {
		// This is a critical error after a successful update. Log it.
		return nil, fmt.Errorf("retrieve updated price from DB: %w", err)
	}

	return updatedPrice, nil
}
