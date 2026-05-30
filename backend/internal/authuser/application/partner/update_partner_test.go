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

// make test-func TEST_NAME=TestUpdatePartnerPresentationFields TEST_PATH=internal/authuser/application/partner/update_partner_test.go

func TestUpdatePartnerPresentationFields(t *testing.T) {
	ctx := context.Background()

	t.Run("should persist occupation, quote, and tags when provided", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		partnerID := uuid.New()
		userID := uuid.New()

		// Existing partner (with no presentation fields)
		existingEncx := td.NewTestPartnerEncxWithUserID(t, userID)
		existingEncx.ID = partnerID
		existingEncx.Occupation = ""
		existingEncx.Quote = ""
		existingEncx.Tags = nil

		mockPartnerRepo.On("GetPartnerByID", ctx, partnerID).Return(existingEncx, nil)

		// Crypto stubs for re-encryption
		mockCrypto.On("GenerateDEK").Return([]byte("test-dek-32-bytes-long-enough!!!"), nil)
		mockCrypto.On("EncryptData", ctx, mock.Anything, mock.Anything).Return([]byte("encrypted"), nil)
		mockCrypto.On("EncryptDEK", ctx, mock.Anything).Return([]byte("encrypted-dek"), nil)
		mockCrypto.On("GetAlias").Return("test-alias")
		mockCrypto.On("GetCurrentKEKVersion", ctx, "test-alias").Return(1, nil)
		mockCrypto.On("DecryptDEKWithVersion", ctx, mock.Anything, mock.Anything).Return([]byte("dek"), nil)
		mockCrypto.On("DecryptData", ctx, mock.Anything, mock.Anything).Return(serializeString("acct_test"), nil)

		mockPartnerRepo.On("UpdatePartner", ctx, mock.AnythingOfType("*domain.PartnerEncx")).Return(nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		occupation := "Kinésithérapeute du sport"
		quote := "Le mouvement est la vie"
		tags := []string{"sport", "rééducation", "blessures"}
		request := &domain.UpdatePartnerRequest{
			Occupation: &occupation,
			Quote:      &quote,
			Tags:       &tags,
		}

		resp, err := svc.UpdatePartner(ctx, partnerID, request)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, occupation, resp.Occupation)
		assert.Equal(t, quote, resp.Quote)
		assert.Equal(t, tags, resp.Tags)

		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("should not clear existing values when fields are omitted", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		partnerID := uuid.New()
		userID := uuid.New()

		existingEncx := td.NewTestPartnerEncxWithUserID(t, userID)
		existingEncx.ID = partnerID
		existingEncx.Occupation = "Existing Occupation"
		existingEncx.Quote = "Existing Quote"
		existingEncx.Tags = []string{"existing", "tags"}

		mockPartnerRepo.On("GetPartnerByID", ctx, partnerID).Return(existingEncx, nil)

		mockCrypto.On("GenerateDEK").Return([]byte("test-dek-32-bytes-long-enough!!!"), nil)
		mockCrypto.On("EncryptData", ctx, mock.Anything, mock.Anything).Return([]byte("encrypted"), nil)
		mockCrypto.On("EncryptDEK", ctx, mock.Anything).Return([]byte("encrypted-dek"), nil)
		mockCrypto.On("GetAlias").Return("test-alias")
		mockCrypto.On("GetCurrentKEKVersion", ctx, "test-alias").Return(1, nil)
		mockCrypto.On("DecryptDEKWithVersion", ctx, mock.Anything, mock.Anything).Return([]byte("dek"), nil)
		mockCrypto.On("DecryptData", ctx, mock.Anything, mock.Anything).Return(serializeString("acct_test"), nil)

		mockPartnerRepo.On("UpdatePartner", ctx, mock.AnythingOfType("*domain.PartnerEncx")).Return(nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		// Only update bio, not the presentation fields
		newBio := "Updated bio only"
		request := &domain.UpdatePartnerRequest{
			Bio: &newBio,
		}

		resp, err := svc.UpdatePartner(ctx, partnerID, request)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, newBio, resp.Bio)
		assert.Equal(t, "Existing Occupation", resp.Occupation, "Occupation should remain unchanged")
		assert.Equal(t, "Existing Quote", resp.Quote, "Quote should remain unchanged")
		assert.Equal(t, []string{"existing", "tags"}, resp.Tags, "Tags should remain unchanged")

		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("should reject occupation exceeding 200 characters", func(t *testing.T) {
		svc := setupPartnerService(new(mockPartnerRepository), new(mockUserRepository), new(mockProductService), new(mockCategoryService), new(mockCryptoService), new(mockStripeService))

		longOccupation := ""
		for i := 0; i < 201; i++ {
			longOccupation += "a"
		}
		request := &domain.UpdatePartnerRequest{
			Occupation: &longOccupation,
		}

		_, err := svc.UpdatePartner(ctx, uuid.New(), request)
		assert.Error(t, err)
	})

	t.Run("should reject quote exceeding 300 characters", func(t *testing.T) {
		svc := setupPartnerService(new(mockPartnerRepository), new(mockUserRepository), new(mockProductService), new(mockCategoryService), new(mockCryptoService), new(mockStripeService))

		longQuote := ""
		for i := 0; i < 301; i++ {
			longQuote += "a"
		}
		request := &domain.UpdatePartnerRequest{
			Quote: &longQuote,
		}

		_, err := svc.UpdatePartner(ctx, uuid.New(), request)
		assert.Error(t, err)
	})

	t.Run("should reject more than 20 tags", func(t *testing.T) {
		svc := setupPartnerService(new(mockPartnerRepository), new(mockUserRepository), new(mockProductService), new(mockCategoryService), new(mockCryptoService), new(mockStripeService))

		tooManyTags := make([]string, 21)
		for i := range tooManyTags {
			tooManyTags[i] = "tag"
		}
		request := &domain.UpdatePartnerRequest{
			Tags: &tooManyTags,
		}

		_, err := svc.UpdatePartner(ctx, uuid.New(), request)
		assert.Error(t, err)
	})

	t.Run("should reject empty tag values", func(t *testing.T) {
		svc := setupPartnerService(new(mockPartnerRepository), new(mockUserRepository), new(mockProductService), new(mockCategoryService), new(mockCryptoService), new(mockStripeService))

		tagsWithEmpty := []string{"valid", ""}
		request := &domain.UpdatePartnerRequest{
			Tags: &tagsWithEmpty,
		}

		_, err := svc.UpdatePartner(ctx, uuid.New(), request)
		assert.Error(t, err)
	})
}
