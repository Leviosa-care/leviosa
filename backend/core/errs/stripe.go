package errs

import (
	"context"
	"errors"
	"fmt"

	"github.com/stripe/stripe-go/v82"
)

// ClassifyStripeError converts Stripe errors to appropriate domain errors
func ClassifyStripeError(operation string, err error) error {
	var stripeErr *stripe.Error
	if errors.As(err, &stripeErr) {
		switch stripeErr.Type {
		case "card_error":
			return NewInvalidValueErr(fmt.Sprintf("%s failed: %s", operation, stripeErr.Msg))
		case "invalid_request_error":
			return NewInvalidValueErr(fmt.Sprintf("%s failed: %s", operation, stripeErr.Msg))
		case "authentication_error":
			return NewPermissionErr(fmt.Sprintf("%s failed: %s", operation, stripeErr.Msg))
		case "api_connection_error", "api_error":
			return NewExternalServiceErr(err, fmt.Sprintf("%s connection failed", operation))
		case "rate_limit_error":
			return NewRateLimitErr(err, operation)
		default:
			return NewInternalErr(fmt.Errorf("%s failed: %w", operation, err))
		}
	}

	// Handle context errors
	if errors.Is(err, context.Canceled) {
		return ErrQueryCancelled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrQueryCancelled
	}

	// Default to internal error
	return NewInternalErr(fmt.Errorf("%s failed: %w", operation, err))
}

