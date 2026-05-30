package partner_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetPublicPartners TEST_PATH=internal/authuser/application/partner/get_public_partners_test.go

func TestGetPublicPartners(t *testing.T) {
	ctx := context.Background()

	t.Run("should return only active non-disabled partners with user info", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		userID := uuid.New()
		partnerID := uuid.New()

		// Build a PublicPartnerRow for an active partner
		partnerEncx := td.NewTestPartnerEncxWithUserID(t, userID)
		partnerEncx.ID = partnerID
		partnerEncx.Occupation = "Kinésithérapeute"
		partnerEncx.Quote = "Le mouvement est la vie"
		partnerEncx.Tags = []string{"sport", "rééducation"}
		partnerEncx.StripeAccountStatus = domain.StripeAccountStatusActive

		row := &domain.PublicPartnerRow{
			PartnerEncx:        partnerEncx,
			FirstNameEncrypted: []byte("encrypted-first"),
			LastNameEncrypted:  []byte("encrypted-last"),
			PictureEncrypted:   []byte("encrypted-pic"),
			UserDEKEncrypted:   []byte("user-dek"),
			UserKeyVersion:     1,
		}

		mockPartnerRepo.On("GetPublicPartners", ctx).Return([]*domain.PublicPartnerRow{row}, nil)

		// Partner decryption stubs
		mockCrypto.On("DecryptDEKWithVersion", ctx, mock.Anything, mock.Anything).Return([]byte("dek"), nil)
		mockCrypto.On("DecryptData", ctx, mock.Anything, mock.Anything).Return(serializeString("decrypted-value"), nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		partners, err := svc.GetPublicPartners(ctx)

		require.NoError(t, err)
		require.Len(t, partners, 1)
		assert.Equal(t, partnerID, partners[0].ID)
		assert.Equal(t, "decrypted-value", partners[0].FirstName)
		assert.Equal(t, "decrypted-value", partners[0].LastName)
		assert.Equal(t, "Kinésithérapeute", partners[0].Occupation)
		assert.Equal(t, "Le mouvement est la vie", partners[0].Quote)
		assert.Equal(t, []string{"sport", "rééducation"}, partners[0].Tags)

		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("should return empty array when no public partners exist", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		mockPartnerRepo.On("GetPublicPartners", ctx).Return([]*domain.PublicPartnerRow{}, nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		partners, err := svc.GetPublicPartners(ctx)

		require.NoError(t, err)
		assert.Empty(t, partners)
		assert.NotNil(t, partners)

		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("should handle partners with empty picture gracefully", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		userID := uuid.New()
		partnerID := uuid.New()

		partnerEncx := td.NewTestPartnerEncxWithUserID(t, userID)
		partnerEncx.ID = partnerID
		partnerEncx.Occupation = "Psychologue"
		partnerEncx.StripeAccountStatus = domain.StripeAccountStatusActive

		row := &domain.PublicPartnerRow{
			PartnerEncx:        partnerEncx,
			FirstNameEncrypted: []byte("encrypted-first"),
			LastNameEncrypted:  []byte("encrypted-last"),
			PictureEncrypted:   nil, // No picture
			UserDEKEncrypted:   []byte("user-dek"),
			UserKeyVersion:     1,
		}

		mockPartnerRepo.On("GetPublicPartners", ctx).Return([]*domain.PublicPartnerRow{row}, nil)

		mockCrypto.On("DecryptDEKWithVersion", ctx, mock.Anything, mock.Anything).Return([]byte("dek"), nil)
		mockCrypto.On("DecryptData", ctx, mock.Anything, mock.Anything).Return(serializeString("decrypted-value"), nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		partners, err := svc.GetPublicPartners(ctx)

		require.NoError(t, err)
		require.Len(t, partners, 1)
		assert.Equal(t, "", partners[0].Picture, "Picture should be empty when not set")

		mockPartnerRepo.AssertExpectations(t)
	})
}
