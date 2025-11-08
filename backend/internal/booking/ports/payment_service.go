package ports

import (
	"context"
)

// PaymentService defines the interface for payment processing
type PaymentService interface {
	// CreatePaymentIntent creates a Stripe payment intent for a booking
	CreatePaymentIntent(ctx context.Context, amount int, currency, description string, metadata map[string]string) (string, string, error) // returns paymentIntentID, clientSecret, error

	// ConfirmPaymentIntent confirms a payment intent
	ConfirmPaymentIntent(ctx context.Context, paymentIntentID string) error

	// RetrievePaymentIntent retrieves a payment intent to check its status
	RetrievePaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntentInfo, error)

	// RefundPayment creates a refund for a payment
	RefundPayment(ctx context.Context, paymentIntentID string, amount int, reason string) (string, error) // returns refundID, error

	// CancelPaymentIntent cancels a payment intent
	CancelPaymentIntent(ctx context.Context, paymentIntentID string) error
}

// PaymentIntentInfo represents payment intent information from Stripe
type PaymentIntentInfo struct {
	ID            string                 `json:"id"`
	Status        string                 `json:"status"`
	Amount        int                    `json:"amount"`
	Currency      string                 `json:"currency"`
	ClientSecret  string                 `json:"client_secret"`
	Description   string                 `json:"description"`
	Metadata      map[string]string      `json:"metadata"`
	LastError     *PaymentIntentError    `json:"last_error,omitempty"`
}

// PaymentIntentError represents payment error information
type PaymentIntentError struct {
	Code        string `json:"code"`
	DeclineCode string `json:"decline_code,omitempty"`
	Message     string `json:"message"`
	Type        string `json:"type"`
}

// Payment Intent Status Constants
// These represent the contract between the application layer and payment service implementations
const (
	// PaymentIntentStatusSucceeded indicates payment was successful
	PaymentIntentStatusSucceeded = "succeeded"

	// PaymentIntentStatusRequiresPaymentMethod indicates payment requires a payment method
	PaymentIntentStatusRequiresPaymentMethod = "requires_payment_method"

	// PaymentIntentStatusRequiresConfirmation indicates payment requires confirmation
	PaymentIntentStatusRequiresConfirmation = "requires_confirmation"

	// PaymentIntentStatusRequiresAction indicates payment requires additional action
	PaymentIntentStatusRequiresAction = "requires_action"

	// PaymentIntentStatusProcessing indicates payment is being processed
	PaymentIntentStatusProcessing = "processing"

	// PaymentIntentStatusCanceled indicates payment was canceled
	PaymentIntentStatusCanceled = "canceled"

	// PaymentIntentStatusPaymentFailed indicates payment failed
	PaymentIntentStatusPaymentFailed = "payment_failed"
)

// Refund Reason Constants
const (
	// RefundReasonDuplicate indicates duplicate charge
	RefundReasonDuplicate = "duplicate"

	// RefundReasonFraudulent indicates fraudulent charge
	RefundReasonFraudulent = "fraudulent"

	// RefundReasonRequestedByCustomer indicates customer requested refund
	RefundReasonRequestedByCustomer = "requested_by_customer"
)