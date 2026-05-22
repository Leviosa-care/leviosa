package booking

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/google/uuid"
)

// MockPaymentService is a simple mock implementation of PaymentService for testing
type MockPaymentService struct {
	paymentIntents map[string]*ports.PaymentIntentInfo
}

// NewMockPaymentService creates a new mock payment service
func NewMockPaymentService() *MockPaymentService {
	return &MockPaymentService{
		paymentIntents: make(map[string]*ports.PaymentIntentInfo),
	}
}

// CreatePaymentIntent mocks creating a Stripe payment intent
func (m *MockPaymentService) CreatePaymentIntent(ctx context.Context, amount int, currency, description string, metadata map[string]string) (string, string, error) {
	paymentIntentID := "pi_test_" + uuid.New().String()[:8]
	clientSecret := "pi_test_secret_" + uuid.New().String()[:8]

	m.paymentIntents[paymentIntentID] = &ports.PaymentIntentInfo{
		ID:           paymentIntentID,
		Status:       ports.PaymentIntentStatusRequiresPaymentMethod,
		Amount:       amount,
		Currency:     currency,
		ClientSecret: clientSecret,
		Description:  description,
		Metadata:     metadata,
	}

	return paymentIntentID, clientSecret, nil
}

// ConfirmPaymentIntent mocks confirming a payment intent
func (m *MockPaymentService) ConfirmPaymentIntent(ctx context.Context, paymentIntentID string) error {
	intent, exists := m.paymentIntents[paymentIntentID]
	if !exists {
		return fmt.Errorf("payment intent not found: %s", paymentIntentID)
	}

	intent.Status = ports.PaymentIntentStatusSucceeded
	return nil
}

// RetrievePaymentIntent mocks retrieving a payment intent
func (m *MockPaymentService) RetrievePaymentIntent(ctx context.Context, paymentIntentID string) (*ports.PaymentIntentInfo, error) {
	intent, exists := m.paymentIntents[paymentIntentID]
	if !exists {
		return nil, fmt.Errorf("payment intent not found: %s", paymentIntentID)
	}

	return intent, nil
}

// RefundPayment mocks creating a refund
func (m *MockPaymentService) RefundPayment(ctx context.Context, paymentIntentID string, amount int, reason string) (string, error) {
	intent, exists := m.paymentIntents[paymentIntentID]
	if !exists {
		return "", fmt.Errorf("payment intent not found: %s", paymentIntentID)
	}

	if intent.Status != ports.PaymentIntentStatusSucceeded {
		return "", fmt.Errorf("payment intent must be succeeded to refund")
	}

	refundID := "re_test_" + uuid.New().String()[:8]
	return refundID, nil
}

// CancelPaymentIntent mocks canceling a payment intent
func (m *MockPaymentService) CancelPaymentIntent(ctx context.Context, paymentIntentID string) error {
	intent, exists := m.paymentIntents[paymentIntentID]
	if !exists {
		return fmt.Errorf("payment intent not found: %s", paymentIntentID)
	}

	intent.Status = ports.PaymentIntentStatusCanceled
	return nil
}

// VerifyWebhookSignature mocks verifying a webhook signature
func (m *MockPaymentService) VerifyWebhookSignature(payload []byte, signature string) (*ports.WebhookEvent, error) {
	// For testing purposes, always return a valid webhook event
	return &ports.WebhookEvent{
		ID:              "evt_test_" + uuid.New().String(),
		Type:            "payment_intent.succeeded",
		PaymentIntentID: "pi_test_" + uuid.New().String(),
		Status:          "succeeded",
		Amount:          1000,
		Currency:        "usd",
		Metadata:        make(map[string]string),
	}, nil
}
