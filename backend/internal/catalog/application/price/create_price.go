package price

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// CreatePrice creates a new Stripe Price and stores its details locally.
func (s *PriceService) CreatePrice(ctx context.Context, productIDStr string, request *domain.CreatePriceRequest) (string, error) {
	// 1. Validate internal product ID
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return "", errs.NewInvalidValueErr(fmt.Sprintf("internal product ID is invalid: %s", err.Error()))
	}

	// 2. Validate incoming price data (errs.CreatePriceInput is analogous to errs.Price here)
	// You might define a `CreatePriceInput` struct in `errs` if different from `Price` directly.
	// For simplicity, let's assume `input` has relevant fields for `Price.Valid` check.
	// A new `errs.Price` will be constructed for `Valid` call.
	tempPrice := domain.Price{
		ProductID: productID,
		Amount:    request.Amount,
		Currency:  request.Currency,
		Interval:  domain.Interval(request.Interval),
	}

	if err := tempPrice.Valid(ctx); err != nil {
		return "", errs.NewInvalidValueErr(fmt.Errorf("price validation failed: %w", err).Error())
	}

	// 3. Get Stripe Product ID for the given internal product ID.
	stripeProductID, _, err := s.sharedRepo.GetStripeProductAndPriceIDs(ctx, productID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return "", errs.NewNotFoundErr(err, "product for price creation")
		case errors.Is(err, errs.ErrDBQuery):
			return "", errs.NewQueryFailedErr(fmt.Errorf("failed to get Stripe Product ID: %w", err))
		default:
			return "", errs.NewUnexpectedError(fmt.Errorf("unhandled error getting Stripe Product ID: %w", err))
		}
	}

	// 4. Create Price in Stripe.
	stripeReq := domain.CreateStripePriceRequest{
		ProductID: stripeProductID, // Use Stripe's product ID
		Amount:    int64(request.Amount),
		Currency:  request.Currency,
		Interval:  request.Interval,
		Active:    true, // Default to active when creating prices
		Nickname:  "",   // Set from request.Nickname if needed
		Metadata:  request.Metadata,
	}

	// Set optional fields if provided
	if request.Nickname != nil {
		stripeReq.Nickname = *request.Nickname
	}
	// paymentPrice, err := s.stripe.CreatePrice(ctx, stripeReq)
	paymentPrice, err := s.stripe.CreatePrice(ctx, stripeReq)
	if err != nil {
		return "", errs.NewExternalServiceErr(fmt.Errorf("failed to create price in Stripe: %w", err), "Stripe API")
	}

	// 5. Store new Price details in your database.
	now := time.Now()
	newPrice := &domain.Price{
		ID:            uuid.New(),      // Generate a new internal ID
		ProductID:     productID,       // Your internal product ID
		StripePriceID: paymentPrice.ID, // Stripe's price ID
		Amount:        int(paymentPrice.Amount),
		Currency:      paymentPrice.Currency,
		Interval:      domain.Interval(paymentPrice.Interval),
		IsActive:      paymentPrice.Active,
		CreatedAt:     now, // DB will set this, but for returning, use now or Stripe's created time
		UpdatedAt:     now,
	}

	if err := s.repo.CreatePrice(ctx, newPrice); err != nil {
		// IMPORTANT: If DB creation fails, consider rolling back Stripe price creation
		// or flagging it for manual review/cleanup. This adds complexity.
		// For now, we'll just return the error.
		return "", errs.NewQueryFailedErr(fmt.Errorf("failed to save price to database: %w", err))
	}

	return newPrice.ID.String(), nil
}
