package admin

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/google/uuid"
)

type MockPaymentService struct {
	paymentIntents map[string]*ports.PaymentIntentInfo
}

func NewMockPaymentService() *MockPaymentService {
	return &MockPaymentService{
		paymentIntents: make(map[string]*ports.PaymentIntentInfo),
	}
}

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

func (m *MockPaymentService) ConfirmPaymentIntent(ctx context.Context, paymentIntentID string) error {
	intent, exists := m.paymentIntents[paymentIntentID]
	if !exists {
		return fmt.Errorf("payment intent not found: %s", paymentIntentID)
	}

	intent.Status = ports.PaymentIntentStatusSucceeded
	return nil
}

func (m *MockPaymentService) RetrievePaymentIntent(ctx context.Context, paymentIntentID string) (*ports.PaymentIntentInfo, error) {
	intent, exists := m.paymentIntents[paymentIntentID]
	if !exists {
		return nil, fmt.Errorf("payment intent not found: %s", paymentIntentID)
	}

	return intent, nil
}

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

func (m *MockPaymentService) CancelPaymentIntent(ctx context.Context, paymentIntentID string) error {
	intent, exists := m.paymentIntents[paymentIntentID]
	if !exists {
		return fmt.Errorf("payment intent not found: %s", paymentIntentID)
	}

	intent.Status = ports.PaymentIntentStatusCanceled
	return nil
}

func (m *MockPaymentService) VerifyWebhookSignature(payload []byte, signature string) (*ports.WebhookEvent, error) {
	return nil, nil
}
