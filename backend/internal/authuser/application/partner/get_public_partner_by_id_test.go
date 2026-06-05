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

// make test-func TEST_NAME=TestGetPublicPartnerByID TEST_PATH=internal/authuser/application/partner/get_public_partner_by_id_test.go

func TestGetPublicPartnerByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return active partner with user info", func(t *testing.T) {
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
		partnerEncx.Occupation = "Ostéopathe D.O"
		partnerEncx.Quote = "Le mouvement est la vie"
		partnerEncx.Tags = []string{"sport", "rééducation"}
		partnerEncx.Bio = "Bio détaillée du praticien"
		partnerEncx.Experience = "10 ans d'expérience"
		partnerEncx.StripeAccountStatus = domain.StripeAccountStatusActive

		row := &domain.PublicPartnerRow{
			PartnerEncx:        partnerEncx,
			FirstNameEncrypted: []byte("encrypted-first"),
			LastNameEncrypted:  []byte("encrypted-last"),
			PictureEncrypted:   []byte("encrypted-pic"),
			UserDEKEncrypted:   []byte("user-dek"),
			UserKeyVersion:     1,
		}

		mockPartnerRepo.On("GetPublicPartnerByID", ctx, partnerID).Return(row, nil)

		// Decryption stubs
		mockCrypto.On("DecryptDEKWithVersion", ctx, mock.Anything, mock.Anything).Return([]byte("dek"), nil)
		mockCrypto.On("DecryptData", ctx, mock.Anything, mock.Anything).Return(serializeString("decrypted-value"), nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		partner, err := svc.GetPublicPartnerByID(ctx, partnerID)

		require.NoError(t, err)
		require.NotNil(t, partner)
		assert.Equal(t, partnerID, partner.ID)
		assert.Equal(t, "decrypted-value", partner.FirstName)
		assert.Equal(t, "decrypted-value", partner.LastName)
		assert.Equal(t, "decrypted-value", partner.Picture)
		assert.Equal(t, "Ostéopathe D.O", partner.Occupation)
		assert.Equal(t, "Le mouvement est la vie", partner.Quote)
		assert.Equal(t, []string{"sport", "rééducation"}, partner.Tags)
		assert.Equal(t, "Bio détaillée du praticien", partner.Bio)
		assert.Equal(t, "10 ans d'expérience", partner.Experience)

		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("should return error when partner not found", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		nonExistentID := uuid.New()
		mockPartnerRepo.On("GetPublicPartnerByID", ctx, nonExistentID).Return(nil, assert.AnError)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		partner, err := svc.GetPublicPartnerByID(ctx, nonExistentID)

		require.Error(t, err)
		assert.Nil(t, partner)

		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("should handle partner with empty picture gracefully", func(t *testing.T) {
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

		mockPartnerRepo.On("GetPublicPartnerByID", ctx, partnerID).Return(row, nil)

		mockCrypto.On("DecryptDEKWithVersion", ctx, mock.Anything, mock.Anything).Return([]byte("dek"), nil)
		mockCrypto.On("DecryptData", ctx, mock.Anything, mock.Anything).Return(serializeString("decrypted-value"), nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		partner, err := svc.GetPublicPartnerByID(ctx, partnerID)

		require.NoError(t, err)
		require.NotNil(t, partner)
		assert.Equal(t, "", partner.Picture, "Picture should be empty when not set")

		mockPartnerRepo.AssertExpectations(t)
	})
}
