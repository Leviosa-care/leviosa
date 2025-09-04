package authPayment

import (
	"context"

	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
)

// FindCustomerByUserID finds a Stripe customer by user ID metadata
// Uses the Search API to find customers by metadata since List API doesn't support metadata filtering
func (s *service) FindCustomerByUserID(ctx context.Context, userID uuid.UUID) (*stripe.Customer, error) {
	// Try using Stripe Search API first (more efficient if available)
	searchParams := &stripe.CustomerSearchParams{
		SearchParams: stripe.SearchParams{
			Query: "metadata['user_id']:'" + userID.String() + "'",
		},
	}

	// Iterate through search results using the Seq2 iterator
	for customer, err := range s.client.V1Customers.Search(ctx, searchParams) {
		if err != nil {
			// If search is not available (e.g., in India), fall back to list all and filter
			return s.findCustomerByUserIDFallback(ctx, userID)
		}
		// Return the first matching customer
		return customer, nil
	}

	// No customers found in search results
	return nil, errs.ErrRepositoryNotFound
}

// findCustomerByUserIDFallback lists all customers and filters by metadata in memory
// This is less efficient but works when Search API is unavailable
func (s *service) findCustomerByUserIDFallback(ctx context.Context, userID uuid.UUID) (*stripe.Customer, error) {
	params := &stripe.CustomerListParams{}
	userIDStr := userID.String()

	for customer, err := range s.client.V1Customers.List(ctx, params) {
		if err != nil {
			return nil, errs.ClassifyStripeError("find customer by user ID fallback", err)
		}

		// Check if this customer has the matching user_id in metadata
		if customer.Metadata != nil {
			if metaUserID, exists := customer.Metadata["user_id"]; exists && metaUserID == userIDStr {
				return customer, nil
			}
		}
	}

	return nil, errs.ErrRepositoryNotFound
}
