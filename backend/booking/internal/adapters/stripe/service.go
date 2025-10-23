package stripe

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/booking/internal/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/refund"
)

// Service handles Stripe payment operations for bookings
type Service struct {
	client *stripe.Client
}

// Compile-time check to ensure Service implements ports.PaymentService
var _ ports.PaymentService = (*Service)(nil)

// NewService creates a new Stripe payment service
func NewService(apiKey, baseURL string) *Service {
	var sc *stripe.Client

	if baseURL != "" {
		backends := stripe.NewBackendsWithConfig(&stripe.BackendConfig{
			URL: &baseURL,
		})
		sc = stripe.NewClient(apiKey, stripe.WithBackends(backends))
	} else {
		// Set default backend for production if not using a custom base URL
		sc = stripe.NewClient(apiKey)
	}

	return &Service{sc}
}

// CreatePaymentIntent creates a Stripe payment intent for a booking
func (s *Service) CreatePaymentIntent(ctx context.Context, amount int, currency, description string, metadata map[string]string) (string, string, error) {
	params := &stripe.PaymentIntentParams{
		Amount:      stripe.Int64(int64(amount)),
		Currency:    stripe.String(currency),
		Description: stripe.String(description),
		Metadata:    metadata,
		// Enable automatic payment methods
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}
	params.SetContext(ctx)

	pi, err := paymentintent.New(params)
	if err != nil {
		return "", "", s.classifyStripeError("create payment intent", err)
	}

	return pi.ID, pi.ClientSecret, nil
}

// ConfirmPaymentIntent confirms a payment intent
func (s *Service) ConfirmPaymentIntent(ctx context.Context, paymentIntentID string) error {
	params := &stripe.PaymentIntentConfirmParams{}
	params.SetContext(ctx)

	_, err := paymentintent.Confirm(paymentIntentID, params)
	if err != nil {
		return s.classifyStripeError("confirm payment intent", err)
	}

	return nil
}

// RetrievePaymentIntent retrieves a payment intent to check its status
func (s *Service) RetrievePaymentIntent(ctx context.Context, paymentIntentID string) (*ports.PaymentIntentInfo, error) {
	params := &stripe.PaymentIntentParams{}
	params.SetContext(ctx)

	pi, err := paymentintent.Get(paymentIntentID, params)
	if err != nil {
		return nil, s.classifyStripeError("retrieve payment intent", err)
	}

	info := &ports.PaymentIntentInfo{
		ID:           pi.ID,
		Status:       string(pi.Status),
		Amount:       int(pi.Amount),
		Currency:     string(pi.Currency),
		ClientSecret: pi.ClientSecret,
		Description:  pi.Description,
		Metadata:     pi.Metadata,
	}

	// Include error information if present
	if pi.LastPaymentError != nil {
		info.LastError = &ports.PaymentIntentError{
			Code:        string(pi.LastPaymentError.Code),
			DeclineCode: string(pi.LastPaymentError.DeclineCode),
			Message:     pi.LastPaymentError.Message,
			Type:        string(pi.LastPaymentError.Type),
		}
	}

	return info, nil
}

// RefundPayment creates a refund for a payment
func (s *Service) RefundPayment(ctx context.Context, paymentIntentID string, amount int, reason string) (string, error) {
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(paymentIntentID),
		Reason:        stripe.String(reason),
	}

	// If amount is specified, create a partial refund
	if amount > 0 {
		params.Amount = stripe.Int64(int64(amount))
	}

	params.SetContext(ctx)

	refund, err := refund.New(params)
	if err != nil {
		return "", s.classifyStripeError("create refund", err)
	}

	return refund.ID, nil
}

// CancelPaymentIntent cancels a payment intent
func (s *Service) CancelPaymentIntent(ctx context.Context, paymentIntentID string) error {
	params := &stripe.PaymentIntentCancelParams{}
	params.SetContext(ctx)

	_, err := paymentintent.Cancel(paymentIntentID, params)
	if err != nil {
		return s.classifyStripeError("cancel payment intent", err)
	}

	return nil
}

// classifyStripeError converts Stripe errors to application-specific errors
func (s *Service) classifyStripeError(operation string, err error) error {
	if stripeErr, ok := err.(*stripe.Error); ok {
		switch stripeErr.Code {
		case stripe.ErrorCodeCardDeclined:
			return fmt.Errorf("%s: card declined - %s", operation, stripeErr.Msg)
		case stripe.ErrorCodeInsufficientFunds:
			return fmt.Errorf("%s: insufficient funds - %s", operation, stripeErr.Msg)
		case stripe.ErrorCodeIncorrectCVC:
			return fmt.Errorf("%s: incorrect CVC - %s", operation, stripeErr.Msg)
		case stripe.ErrorCodeExpiredCard:
			return fmt.Errorf("%s: expired card - %s", operation, stripeErr.Msg)
		case stripe.ErrorCodeProcessingError:
			return errs.NewConnectionFailure(fmt.Sprintf("%s: payment processing error", operation), err)
		case stripe.ErrorCodeRateLimitError:
			return errs.NewResourceExhausted(fmt.Sprintf("%s: rate limit exceeded", operation), err)
		case stripe.ErrorCodeResourceMissing:
			return errs.NewNotFoundErr("payment_intent", "")
		case stripe.ErrorCodeAPIKeyExpired, stripe.ErrorCodeInvalidAPIKey:
			return errs.NewPermissionDenied(fmt.Sprintf("%s: invalid API key", operation))
		default:
			return fmt.Errorf("%s: payment error - %s", operation, stripeErr.Msg)
		}
	}

	return fmt.Errorf("%s: %w", operation, err)
}