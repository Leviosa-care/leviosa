package product

import (
	"context"
	"fmt"
	"log"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *ProductService) RemoveProduct(ctx context.Context, productIDStr string) error {
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return errs.NewInvalidValueErr("product ID is required")
	}

	stripeProductID, stripePriceIDs, err := s.sharedRepo.GetStripeProductAndPriceIDs(ctx, productID)
	if err != nil {
		return fmt.Errorf("get stripe product and prices IDs for product with ID '%s': %w", productIDStr, err)
	}

	if err := s.stripe.DeactivateProduct(ctx, stripeProductID); err != nil {
		return errs.NewExternalServiceErr(err, "stripe product deletion failed")
	}

	if err := s.priceStripe.DeactivatePrices(ctx, stripePriceIDs); err != nil {
		return errs.NewExternalServiceErr(err, "stripe prices deletion failed")
	}

	if err := s.sharedRepo.DeactivatePrices(ctx, stripePriceIDs); err != nil {
		return errs.NewUnexpectedError(fmt.Errorf("database price deactivation failed: %w", err))
	}

	if err := s.repo.DeleteProduct(ctx, productID); err != nil {
		log.Printf("Service: Failed to delete product %s from DB, attempting Stripe rollback: %v", productID, err)

		if rollbackErr := s.stripe.ReactivateProduct(ctx, stripeProductID); rollbackErr != nil {
			log.Printf("Service: Failed to reactivate Stripe product %s. Data inconsistency detected! Rollback error: %v", stripeProductID, rollbackErr)
		}
		if rollbackErr := s.priceStripe.ReactivatePrices(ctx, stripePriceIDs); rollbackErr != nil {
			log.Printf("Service: Failed to reactivate Stripe prices for product %s. Data inconsistency detected! Rollback error: %v", stripeProductID, rollbackErr)
		}
		if rollbackErr := s.sharedRepo.ReactivatePrices(ctx, stripePriceIDs); rollbackErr != nil {
			log.Printf("Service: Failed to reactivate database prices for product %s. Data inconsistency detected! Rollback error: %v", stripeProductID, rollbackErr)
		}

		return fmt.Errorf("delete product with ID %s: %w", productID, err)
	}

	return nil
}
