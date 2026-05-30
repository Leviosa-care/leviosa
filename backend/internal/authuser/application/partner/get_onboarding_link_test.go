package partner_test

import (
	"context"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// serializeString encodes a string in encx compact binary format: [4-byte LE length][UTF-8 data]
func serializeString(s string) []byte {
	data := []byte(s)
	result := make([]byte, 4+len(data))
	binary.LittleEndian.PutUint32(result[:4], uint32(len(data)))
	copy(result[4:], data)
	return result
}

func TestGetOnboardingLink(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	returnURL := "https://example.com/staff/profile?stripe=return"
	refreshURL := "https://example.com/staff/profile?stripe=refresh"

	t.Run("should return onboarding link when StripeConnectedAccountID is present", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		partnerID := uuid.New()
		stripeAccountID := "acct_existing_123"
		expectedURL := "https://connect.stripe.com/setup/e/acct_existing_123"

		// Build encrypted partner with existing Stripe account ID
		partnerEncx := buildPartnerEncx(partnerID, userID, stripeAccountID, domain.StripeAccountStatusPending)

		mockPartnerRepo.On("GetPartnerByUserID", ctx, userID).Return(partnerEncx, nil)

		// Crypto stubs for DecryptPartnerEncx — return properly serialized string bytes
		mockCrypto.On("DecryptDEKWithVersion", ctx, mock.Anything, mock.Anything).Return([]byte("test-dek-32-bytes-long-enough!!!"), nil)
		mockCrypto.On("DecryptData", ctx, mock.Anything, mock.Anything).Return(serializeString(stripeAccountID), nil)

		// Stripe should NOT create a new account — only create an Account Link
		mockStripe.On("CreateAccountLink", ctx, stripeAccountID, "account_onboarding", returnURL, refreshURL).Return(expectedURL, nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		url, err := svc.GetOnboardingLink(ctx, userID, returnURL, refreshURL)

		require.NoError(t, err)
		assert.Equal(t, expectedURL, url)

		// Verify CreateConnectedAccount was NOT called
		mockStripe.AssertNotCalled(t, "CreateConnectedAccount", ctx, userID)
		mockStripe.AssertExpectations(t)
		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("should create Stripe account first when StripeConnectedAccountID is empty", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		partnerID := uuid.New()
		newAccountID := "acct_newly_created_456"
		expectedURL := "https://connect.stripe.com/setup/e/acct_newly_created_456"

		// Build encrypted partner with empty Stripe account ID
		partnerEncx := buildPartnerEncx(partnerID, userID, "", domain.StripeAccountStatusPending)

		mockPartnerRepo.On("GetPartnerByUserID", ctx, userID).Return(partnerEncx, nil)

		// Crypto stubs for DecryptPartnerEncx (StripeConnectedAccountID empty — no encrypted bytes to decrypt)
		mockCrypto.On("DecryptDEKWithVersion", ctx, mock.Anything, mock.Anything).Return([]byte("test-dek-32-bytes-long-enough!!!"), nil)

		// Stripe creates a new connected account first
		mockStripe.On("CreateConnectedAccount", ctx, userID).Return(newAccountID, nil)

		// Crypto stubs for ProcessPartnerEncx (re-encryption with new account ID)
		mockCrypto.On("GenerateDEK").Return([]byte("test-dek-32-bytes-long-enough!!!"), nil)
		mockCrypto.On("EncryptData", ctx, mock.Anything, mock.Anything).Return([]byte("encrypted"), nil)
		mockCrypto.On("EncryptDEK", ctx, mock.Anything).Return([]byte("encrypted-dek"), nil)
		mockCrypto.On("GetAlias").Return("test-alias")
		mockCrypto.On("GetCurrentKEKVersion", ctx, "test-alias").Return(1, nil)

		// Partner is updated with new Stripe account ID
		mockPartnerRepo.On("UpdatePartner", ctx, mock.AnythingOfType("*domain.PartnerEncx")).Return(nil)

		// Then create the Account Link
		mockStripe.On("CreateAccountLink", ctx, newAccountID, "account_onboarding", returnURL, refreshURL).Return(expectedURL, nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		url, err := svc.GetOnboardingLink(ctx, userID, returnURL, refreshURL)

		require.NoError(t, err)
		assert.Equal(t, expectedURL, url)

		// Verify CreateConnectedAccount was called
		mockStripe.AssertCalled(t, "CreateConnectedAccount", ctx, userID)
		mockStripe.AssertExpectations(t)
		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("should return error when partner not found", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		mockPartnerRepo.On("GetPartnerByUserID", ctx, userID).Return(nil, errors.New("not found"))

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		url, err := svc.GetOnboardingLink(ctx, userID, returnURL, refreshURL)

		require.Error(t, err)
		assert.Empty(t, url)

		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("should return error when CreateConnectedAccount fails", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		partnerID := uuid.New()
		partnerEncx := buildPartnerEncx(partnerID, userID, "", domain.StripeAccountStatusPending)

		mockPartnerRepo.On("GetPartnerByUserID", ctx, userID).Return(partnerEncx, nil)
		mockCrypto.On("DecryptDEKWithVersion", ctx, mock.Anything, mock.Anything).Return([]byte("test-dek-32-bytes-long-enough!!!"), nil)

		// Stripe account creation fails
		mockStripe.On("CreateConnectedAccount", ctx, userID).Return("", errors.New("stripe api error"))

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		url, err := svc.GetOnboardingLink(ctx, userID, returnURL, refreshURL)

		require.Error(t, err)
		assert.Empty(t, url)
		assert.Contains(t, err.Error(), "create connected account")

		mockStripe.AssertExpectations(t)
	})

	t.Run("should return error when CreateAccountLink fails", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		partnerID := uuid.New()
		stripeAccountID := "acct_existing_789"
		partnerEncx := buildPartnerEncx(partnerID, userID, stripeAccountID, domain.StripeAccountStatusPending)

		mockPartnerRepo.On("GetPartnerByUserID", ctx, userID).Return(partnerEncx, nil)
		mockCrypto.On("DecryptDEKWithVersion", ctx, mock.Anything, mock.Anything).Return([]byte("test-dek-32-bytes-long-enough!!!"), nil)
		mockCrypto.On("DecryptData", ctx, mock.Anything, mock.Anything).Return(serializeString(stripeAccountID), nil)

		// Account Link creation fails
		mockStripe.On("CreateAccountLink", ctx, stripeAccountID, "account_onboarding", returnURL, refreshURL).Return("", errors.New("stripe account link error"))

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		url, err := svc.GetOnboardingLink(ctx, userID, returnURL, refreshURL)

		require.Error(t, err)
		assert.Empty(t, url)
		assert.Contains(t, err.Error(), "create account link")

		mockStripe.AssertExpectations(t)
	})
}

// buildPartnerEncx creates a PartnerEncx with the given fields.
// For encrypted fields it uses placeholder bytes — the mocks handle actual crypto.
func buildPartnerEncx(partnerID, userID uuid.UUID, stripeConnectedAccountID string, status domain.StripeAccountStatus) *domain.PartnerEncx {
	encx := &domain.PartnerEncx{
		ID:                      partnerID,
		UserID:                  userID,
		Bio:                     "test bio",
		Experience:              "test experience",
		StripeAccountStatus:     status,
		StripeOnboardingComplete: false,
		DEKEncrypted:            []byte("encrypted-dek"),
		KeyVersion:              1,
	}

	// If there's a Stripe account ID, add encrypted bytes so the decrypt mock path is taken
	if stripeConnectedAccountID != "" {
		encx.StripeConnectedAccountIDEncrypted = []byte("encrypted-stripe-id")
	}

	return encx
}
