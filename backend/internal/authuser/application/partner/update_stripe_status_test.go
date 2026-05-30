package partner_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateStripeAccountStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("should update status to active when partner found", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		partnerID := uuid.New()
		userID := uuid.New()
		stripeAccountID := "acct_test_123"

		partnerEncx := buildPartnerEncx(partnerID, userID, stripeAccountID, domain.StripeAccountStatusPending)

		mockPartnerRepo.On("GetAllPartnersWithStripeAccount", ctx).Return([]*domain.PartnerEncx{partnerEncx}, nil)
		mockCrypto.On("DecryptDEKWithVersion", ctx, mock.Anything, mock.Anything).Return([]byte("test-dek-32-bytes-long-enough!!!"), nil)
		mockCrypto.On("DecryptData", ctx, mock.Anything, mock.Anything).Return(serializeString(stripeAccountID), nil)
		mockPartnerRepo.On("UpdatePartnerStripeStatus", ctx, partnerID, domain.StripeAccountStatusActive).Return(nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		updatedID, err := svc.UpdateStripeAccountStatus(ctx, stripeAccountID, domain.StripeAccountStatusActive)

		require.NoError(t, err)
		assert.Equal(t, partnerID, updatedID)

		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("should return not found when no partner matches stripe account ID", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		partnerID := uuid.New()
		userID := uuid.New()
		differentStripeID := "acct_other_456"

		partnerEncx := buildPartnerEncx(partnerID, userID, differentStripeID, domain.StripeAccountStatusPending)

		mockPartnerRepo.On("GetAllPartnersWithStripeAccount", ctx).Return([]*domain.PartnerEncx{partnerEncx}, nil)
		mockCrypto.On("DecryptDEKWithVersion", ctx, mock.Anything, mock.Anything).Return([]byte("test-dek-32-bytes-long-enough!!!"), nil)
		mockCrypto.On("DecryptData", ctx, mock.Anything, mock.Anything).Return(serializeString(differentStripeID), nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		updatedID, err := svc.UpdateStripeAccountStatus(ctx, "acct_not_found_789", domain.StripeAccountStatusActive)

		require.Error(t, err)
		assert.Equal(t, uuid.Nil, updatedID)

		// UpdatePartnerStripeStatus should NOT have been called
		mockPartnerRepo.AssertNotCalled(t, "UpdatePartnerStripeStatus")
	})

	t.Run("should update status to restricted", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		partnerID := uuid.New()
		userID := uuid.New()
		stripeAccountID := "acct_test_restricted"

		partnerEncx := buildPartnerEncx(partnerID, userID, stripeAccountID, domain.StripeAccountStatusPending)

		mockPartnerRepo.On("GetAllPartnersWithStripeAccount", ctx).Return([]*domain.PartnerEncx{partnerEncx}, nil)
		mockCrypto.On("DecryptDEKWithVersion", ctx, mock.Anything, mock.Anything).Return([]byte("test-dek-32-bytes-long-enough!!!"), nil)
		mockCrypto.On("DecryptData", ctx, mock.Anything, mock.Anything).Return(serializeString(stripeAccountID), nil)
		mockPartnerRepo.On("UpdatePartnerStripeStatus", ctx, partnerID, domain.StripeAccountStatusRestricted).Return(nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		updatedID, err := svc.UpdateStripeAccountStatus(ctx, stripeAccountID, domain.StripeAccountStatusRestricted)

		require.NoError(t, err)
		assert.Equal(t, partnerID, updatedID)

		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("should update status to disabled", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		partnerID := uuid.New()
		userID := uuid.New()
		stripeAccountID := "acct_test_disabled"

		partnerEncx := buildPartnerEncx(partnerID, userID, stripeAccountID, domain.StripeAccountStatusActive)

		mockPartnerRepo.On("GetAllPartnersWithStripeAccount", ctx).Return([]*domain.PartnerEncx{partnerEncx}, nil)
		mockCrypto.On("DecryptDEKWithVersion", ctx, mock.Anything, mock.Anything).Return([]byte("test-dek-32-bytes-long-enough!!!"), nil)
		mockCrypto.On("DecryptData", ctx, mock.Anything, mock.Anything).Return(serializeString(stripeAccountID), nil)
		mockPartnerRepo.On("UpdatePartnerStripeStatus", ctx, partnerID, domain.StripeAccountStatusDisabled).Return(nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		updatedID, err := svc.UpdateStripeAccountStatus(ctx, stripeAccountID, domain.StripeAccountStatusDisabled)

		require.NoError(t, err)
		assert.Equal(t, partnerID, updatedID)

		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("should skip partners that fail decryption and continue", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		partnerID1 := uuid.New()
		userID1 := uuid.New()
		partnerID2 := uuid.New()
		userID2 := uuid.New()
		targetStripeID := "acct_target"

		partnerEncx1 := buildPartnerEncx(partnerID1, userID1, "acct_corrupt", domain.StripeAccountStatusPending)
		partnerEncx2 := buildPartnerEncx(partnerID2, userID2, targetStripeID, domain.StripeAccountStatusPending)

		mockPartnerRepo.On("GetAllPartnersWithStripeAccount", ctx).Return([]*domain.PartnerEncx{partnerEncx1, partnerEncx2}, nil)

		// First partner fails to decrypt
		mockCrypto.On("DecryptDEKWithVersion", ctx, mock.Anything, mock.Anything).Return([]byte("test-dek-32-bytes-long-enough!!!"), nil).Once()
		mockCrypto.On("DecryptData", ctx, mock.Anything, mock.Anything).Return(nil, errors.New("decryption failed")).Once()

		// Second partner decrypts successfully
		mockCrypto.On("DecryptDEKWithVersion", ctx, mock.Anything, mock.Anything).Return([]byte("test-dek-32-bytes-long-enough!!!"), nil).Once()
		mockCrypto.On("DecryptData", ctx, mock.Anything, mock.Anything).Return(serializeString(targetStripeID), nil).Once()

		mockPartnerRepo.On("UpdatePartnerStripeStatus", ctx, partnerID2, domain.StripeAccountStatusActive).Return(nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		updatedID, err := svc.UpdateStripeAccountStatus(ctx, targetStripeID, domain.StripeAccountStatusActive)

		require.NoError(t, err)
		assert.Equal(t, partnerID2, updatedID)

		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("should return error when GetAllPartnersWithStripeAccount fails", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		mockPartnerRepo.On("GetAllPartnersWithStripeAccount", ctx).Return(nil, errors.New("db error"))

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		updatedID, err := svc.UpdateStripeAccountStatus(ctx, "acct_any", domain.StripeAccountStatusActive)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "get partners with stripe account")
		assert.Equal(t, uuid.Nil, updatedID)

		mockPartnerRepo.AssertExpectations(t)
	})
}
