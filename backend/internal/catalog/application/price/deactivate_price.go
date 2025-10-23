package price

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// DeactivatePrice deactivates a price in Stripe and marks it inactive in your database.
// This can be simplified by calling UpdatePrice with `Active: stripe.Bool(false)`.
func (s *PriceService) DeactivatePrice(ctx context.Context, priceID string) error {
	if err := uuid.Validate(priceID); err != nil {
		return errs.NewInvalidValueErr(fmt.Sprintf("price ID is invalid: %s", err.Error()))
	}

	// Use UpdatePrice to set active to false
	input := domain.UpdatePriceRequest{
		Active: Bool(false), // Assuming errs.Bool helper if not using aws.Bool
	}
	_, err := s.UpdatePrice(ctx, priceID, input)
	return err // Propagate errors from UpdatePrice
}

// Helper to convert bool to *bool
func Bool(b bool) *bool {
	return &b
}
