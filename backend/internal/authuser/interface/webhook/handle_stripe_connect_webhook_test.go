package webhookHandler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type mockPartnerService struct {
	mock.Mock
}

func (m *mockPartnerService) GetPartnerVerificationStatus(ctx context.Context, partnerID uuid.UUID) (bool, error) {
	args := m.Called(ctx, partnerID)
	return args.Bool(0), args.Error(1)
}

func (m *mockPartnerService) GetPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.PartnerResponse, error) {
	args := m.Called(ctx, partnerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PartnerResponse), args.Error(1)
}

func (m *mockPartnerService) GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*domain.PartnerResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PartnerResponse), args.Error(1)
}

func (m *mockPartnerService) GetAllPartners(ctx context.Context) ([]*domain.PartnerResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PartnerResponse), args.Error(1)
}

func (m *mockPartnerService) GetAllPartnersByCategory(ctx context.Context, categoryID string) ([]*domain.PartnerResponse, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PartnerResponse), args.Error(1)
}

func (m *mockPartnerService) GetAllPartnersByCategories(ctx context.Context, categoryIDs []string) ([]*domain.PartnerResponse, error) {
	args := m.Called(ctx, categoryIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PartnerResponse), args.Error(1)
}

func (m *mockPartnerService) GetAllPartnersByProduct(ctx context.Context, productID string) ([]*domain.PartnerResponse, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PartnerResponse), args.Error(1)
}

func (m *mockPartnerService) GetAllPartnersByProducts(ctx context.Context, productIDs []string) ([]*domain.PartnerResponse, error) {
	args := m.Called(ctx, productIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PartnerResponse), args.Error(1)
}

func (m *mockPartnerService) CreatePartner(ctx context.Context, userID uuid.UUID, bio, experience string, categoryIDs, productIDs []uuid.UUID) (*domain.Partner, error) {
	args := m.Called(ctx, userID, bio, experience, categoryIDs, productIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Partner), args.Error(1)
}

func (m *mockPartnerService) UpdatePartner(ctx context.Context, partnerID uuid.UUID, request *domain.UpdatePartnerRequest) (*domain.PartnerResponse, error) {
	args := m.Called(ctx, partnerID, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PartnerResponse), args.Error(1)
}

func (m *mockPartnerService) DeletePartner(ctx context.Context, partnerID uuid.UUID) error {
	args := m.Called(ctx, partnerID)
	return args.Error(0)
}

func (m *mockPartnerService) VerifyPartner(ctx context.Context, partnerID uuid.UUID, verifiedByUserID uuid.UUID) (*domain.PartnerResponse, error) {
	args := m.Called(ctx, partnerID, verifiedByUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PartnerResponse), args.Error(1)
}

func (m *mockPartnerService) GetOnboardingLink(ctx context.Context, userID uuid.UUID, returnURL, refreshURL string) (string, error) {
	args := m.Called(ctx, userID, returnURL, refreshURL)
	return args.String(0), args.Error(1)
}

func (m *mockPartnerService) UpdateStripeAccountStatus(ctx context.Context, stripeAccountID string, status domain.StripeAccountStatus) (uuid.UUID, error) {
	args := m.Called(ctx, stripeAccountID, status)
	if args.Get(0) == nil {
		return uuid.Nil, args.Error(1)
	}
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockPartnerService) GetPublicPartners(ctx context.Context) ([]*domain.PublicPartnerResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PublicPartnerResponse), args.Error(1)
}

// Verify mockPartnerService satisfies the interface
var _ ports.PartnerService = (*mockPartnerService)(nil)

type mockStripeService struct {
	mock.Mock
}

func (m *mockStripeService) CreateCustomer(ctx context.Context, userID uuid.UUID, email, firstName, lastName string) (*stripe.Customer, error) {
	args := m.Called(ctx, userID, email, firstName, lastName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func (m *mockStripeService) RetrieveCustomer(ctx context.Context, customerID string) (*stripe.Customer, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func (m *mockStripeService) UpdateCustomer(ctx context.Context, customerID string, params *stripe.CustomerUpdateParams) (*stripe.Customer, error) {
	args := m.Called(ctx, customerID, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func (m *mockStripeService) DeleteCustomer(ctx context.Context, customerID string) (*stripe.Customer, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func (m *mockStripeService) FindCustomerByUserID(ctx context.Context, userID uuid.UUID) (*stripe.Customer, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func (m *mockStripeService) CreateConnectedAccount(ctx context.Context, userID uuid.UUID) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *mockStripeService) CreateAccountLink(ctx context.Context, accountID, returnType, returnURL, refreshURL string) (string, error) {
	args := m.Called(ctx, accountID, returnType, returnURL, refreshURL)
	return args.String(0), args.Error(1)
}

func (m *mockStripeService) VerifyConnectWebhookSignature(payload []byte, signature string) (string, bool, bool, error) {
	args := m.Called(payload, signature)
	return args.String(0), args.Bool(1), args.Bool(2), args.Error(3)
}

// Verify mockStripeService satisfies the interface
var _ ports.StripeService = (*mockStripeService)(nil)

type mockAuthMiddleware struct{}

func (m *mockAuthMiddleware) RequireAccessToken(next mw.Handler) mw.Handler                     { return next }
func (m *mockAuthMiddleware) RequireRefreshToken(next mw.Handler) mw.Handler                    { return next }
func (m *mockAuthMiddleware) RequireMinimumRole(role identity.Role) func(mw.Handler) mw.Handler {
	return func(next mw.Handler) mw.Handler { return next }
}
func (m *mockAuthMiddleware) RequireAnyRole(roles ...identity.Role) func(mw.Handler) mw.Handler {
	return func(next mw.Handler) mw.Handler { return next }
}
func (m *mockAuthMiddleware) RequireAdmin(next mw.Handler) mw.Handler        { return next }
func (m *mockAuthMiddleware) RequireServiceAuth(next mw.Handler) mw.Handler { return next }

// Ensure mockAuthMiddleware satisfies the auth.AuthMiddleware interface
var _ auth.AuthMiddleware = (*mockAuthMiddleware)(nil)

// --- Tests ---

func TestHandleStripeConnectWebhook(t *testing.T) {
	t.Run("should return 200 and update partner status to active for valid account.updated event", func(t *testing.T) {
		mockSvc := new(mockPartnerService)
		mockStripe := new(mockStripeService)

		stripeAccountID := "acct_test_123"
		partnerID := uuid.New()

		mockStripe.On("VerifyConnectWebhookSignature", mock.Anything, "valid-signature").
			Return(stripeAccountID, true, true, nil)
		mockSvc.On("UpdateStripeAccountStatus", mock.Anything, stripeAccountID, domain.StripeAccountStatusActive).
			Return(partnerID, nil)

		handler := &handler{svc: mockSvc, stripe: mockStripe, authmw: &mockAuthMiddleware{}}

		body := bytes.NewBufferString(`{"type":"account.updated"}`)
		req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe/connect", body)
		req.Header.Set("Stripe-Signature", "valid-signature")
		w := httptest.NewRecorder()

		handler.HandleStripeConnectWebhook(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertCalled(t, "UpdateStripeAccountStatus", mock.Anything, stripeAccountID, domain.StripeAccountStatusActive)
	})

	t.Run("should return 200 with restricted status when charges enabled but payouts disabled", func(t *testing.T) {
		mockSvc := new(mockPartnerService)
		mockStripe := new(mockStripeService)

		stripeAccountID := "acct_test_456"
		partnerID := uuid.New()

		mockStripe.On("VerifyConnectWebhookSignature", mock.Anything, "valid-signature").
			Return(stripeAccountID, true, false, nil)
		mockSvc.On("UpdateStripeAccountStatus", mock.Anything, stripeAccountID, domain.StripeAccountStatusRestricted).
			Return(partnerID, nil)

		handler := &handler{svc: mockSvc, stripe: mockStripe, authmw: &mockAuthMiddleware{}}

		body := bytes.NewBufferString(`{"type":"account.updated"}`)
		req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe/connect", body)
		req.Header.Set("Stripe-Signature", "valid-signature")
		w := httptest.NewRecorder()

		handler.HandleStripeConnectWebhook(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertCalled(t, "UpdateStripeAccountStatus", mock.Anything, stripeAccountID, domain.StripeAccountStatusRestricted)
	})

	t.Run("should return 200 with disabled status when both charges and payouts disabled", func(t *testing.T) {
		mockSvc := new(mockPartnerService)
		mockStripe := new(mockStripeService)

		stripeAccountID := "acct_test_789"
		partnerID := uuid.New()

		mockStripe.On("VerifyConnectWebhookSignature", mock.Anything, "valid-signature").
			Return(stripeAccountID, false, false, nil)
		mockSvc.On("UpdateStripeAccountStatus", mock.Anything, stripeAccountID, domain.StripeAccountStatusDisabled).
			Return(partnerID, nil)

		handler := &handler{svc: mockSvc, stripe: mockStripe, authmw: &mockAuthMiddleware{}}

		body := bytes.NewBufferString(`{"type":"account.updated"}`)
		req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe/connect", body)
		req.Header.Set("Stripe-Signature", "valid-signature")
		w := httptest.NewRecorder()

		handler.HandleStripeConnectWebhook(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertCalled(t, "UpdateStripeAccountStatus", mock.Anything, stripeAccountID, domain.StripeAccountStatusDisabled)
	})

	t.Run("should return 400 when Stripe signature is missing", func(t *testing.T) {
		mockSvc := new(mockPartnerService)
		mockStripe := new(mockStripeService)

		handler := &handler{svc: mockSvc, stripe: mockStripe, authmw: &mockAuthMiddleware{}}

		body := bytes.NewBufferString(`{"type":"account.updated"}`)
		req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe/connect", body)
		// No Stripe-Signature header
		w := httptest.NewRecorder()

		handler.HandleStripeConnectWebhook(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		// No database write should occur
		mockSvc.AssertNotCalled(t, "UpdateStripeAccountStatus")
	})

	t.Run("should return 400 when Stripe signature verification fails", func(t *testing.T) {
		mockSvc := new(mockPartnerService)
		mockStripe := new(mockStripeService)

		mockStripe.On("VerifyConnectWebhookSignature", mock.Anything, "invalid-signature").
			Return("", false, false, errors.New("signature mismatch"))

		handler := &handler{svc: mockSvc, stripe: mockStripe, authmw: &mockAuthMiddleware{}}

		body := bytes.NewBufferString(`{"type":"account.updated"}`)
		req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe/connect", body)
		req.Header.Set("Stripe-Signature", "invalid-signature")
		w := httptest.NewRecorder()

		handler.HandleStripeConnectWebhook(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		// No database write should occur
		mockSvc.AssertNotCalled(t, "UpdateStripeAccountStatus")
	})

	t.Run("should return 200 and ignore unknown event types", func(t *testing.T) {
		mockSvc := new(mockPartnerService)
		mockStripe := new(mockStripeService)

		// Empty accountID signals unknown event type
		mockStripe.On("VerifyConnectWebhookSignature", mock.Anything, "valid-signature").
			Return("", false, false, nil)

		handler := &handler{svc: mockSvc, stripe: mockStripe, authmw: &mockAuthMiddleware{}}

		body := bytes.NewBufferString(`{"type":"some.other.event"}`)
		req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe/connect", body)
		req.Header.Set("Stripe-Signature", "valid-signature")
		w := httptest.NewRecorder()

		handler.HandleStripeConnectWebhook(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// No database write should occur
		mockSvc.AssertNotCalled(t, "UpdateStripeAccountStatus")
	})

	t.Run("should return 200 when partner not found for Stripe account", func(t *testing.T) {
		mockSvc := new(mockPartnerService)
		mockStripe := new(mockStripeService)

		stripeAccountID := "acct_unknown"

		mockStripe.On("VerifyConnectWebhookSignature", mock.Anything, "valid-signature").
			Return(stripeAccountID, true, true, nil)
		mockSvc.On("UpdateStripeAccountStatus", mock.Anything, stripeAccountID, domain.StripeAccountStatusActive).
			Return(uuid.Nil, errs.ErrRepositoryNotFound)

		handler := &handler{svc: mockSvc, stripe: mockStripe, authmw: &mockAuthMiddleware{}}

		body := bytes.NewBufferString(`{"type":"account.updated"}`)
		req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe/connect", body)
		req.Header.Set("Stripe-Signature", "valid-signature")
		w := httptest.NewRecorder()

		handler.HandleStripeConnectWebhook(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return 500 when service returns transient error", func(t *testing.T) {
		mockSvc := new(mockPartnerService)
		mockStripe := new(mockStripeService)

		stripeAccountID := "acct_transient_err"

		mockStripe.On("VerifyConnectWebhookSignature", mock.Anything, "valid-signature").
			Return(stripeAccountID, true, true, nil)
		mockSvc.On("UpdateStripeAccountStatus", mock.Anything, stripeAccountID, domain.StripeAccountStatusActive).
			Return(uuid.Nil, errors.New("db connection lost"))

		handler := &handler{svc: mockSvc, stripe: mockStripe, authmw: &mockAuthMiddleware{}}

		body := bytes.NewBufferString(`{"type":"account.updated"}`)
		req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe/connect", body)
		req.Header.Set("Stripe-Signature", "valid-signature")
		w := httptest.NewRecorder()

		handler.HandleStripeConnectWebhook(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

