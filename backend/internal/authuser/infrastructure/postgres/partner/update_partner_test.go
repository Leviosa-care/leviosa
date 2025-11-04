package partnerRepository_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdatePartner TEST_PATH=internal/authuser/infrastructure/postgres/partner/update_partner_test.go

func TestUpdatePartner(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully update partner", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		// Create a user first since partner has foreign key constraint
		userID := td.CreateTestUserForPartner(t, ctx, testPool)

		// Create and insert initial partner
		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID
		partnerEncx.StripeAccountStatus = domain.StripeAccountStatusPending
		partnerEncx.StripeOnboardingComplete = false

		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Update partner data
		bio := "Updated professional bio with more details"
		experience := "Updated experience: 10+ years in specialized field"
		categoryIDsEncrypted := []uuid.UUID{uuid.New()}
		productIDsEncrypted := []uuid.UUID{uuid.New()}
		stripeConnectedAccountIDEncrypted := []byte("acct_updated123456789")
		dekEncrypted := []byte("updated_dek_encrypted")

		partnerEncx.Bio = bio
		partnerEncx.Experience = experience
		partnerEncx.CategoryIDs = categoryIDsEncrypted
		partnerEncx.ProductIDs = productIDsEncrypted
		partnerEncx.StripeConnectedAccountIDEncrypted = stripeConnectedAccountIDEncrypted
		partnerEncx.StripeAccountStatus = domain.StripeAccountStatusActive
		partnerEncx.StripeOnboardingComplete = true
		partnerEncx.DEKEncrypted = dekEncrypted
		partnerEncx.KeyVersion = 2

		// Act
		err = repo.UpdatePartner(ctx, partnerEncx)

		// Assert
		assert.NoError(t, err)

		// Verify the update by retrieving the partner
		updatedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
		require.NoError(t, err)
		require.NotNil(t, updatedPartner)

		// Verify data was updated
		assert.Equal(t, userID, updatedPartner.UserID)
		assert.Equal(t, bio, updatedPartner.Bio)
		assert.Equal(t, experience, updatedPartner.Experience)
		assert.Equal(t, categoryIDsEncrypted, updatedPartner.CategoryIDs)
		assert.Equal(t, productIDsEncrypted, updatedPartner.ProductIDs)
		assert.Equal(t, stripeConnectedAccountIDEncrypted, updatedPartner.StripeConnectedAccountIDEncrypted)
		assert.Equal(t, domain.StripeAccountStatusActive, updatedPartner.StripeAccountStatus)
		assert.True(t, updatedPartner.StripeOnboardingComplete)
		assert.Equal(t, dekEncrypted, updatedPartner.DEKEncrypted)
		assert.Equal(t, 2, updatedPartner.KeyVersion)

		// Verify updated_at timestamp was updated
		assert.True(t, updatedPartner.UpdatedAt.After(partnerEncx.UpdatedAt))
	})

	t.Run("should return not found error when updating non-existent partner", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		partner := td.NewTestPartnerEncx(t)
		partner.UserID = uuid.New()

		// Act
		err := repo.UpdatePartner(ctx, partner)

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrRepositoryNotFound))
	})

	t.Run("should successfully update partner Stripe status from pending to active", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		userID := td.CreateTestUserForPartner(t, ctx, testPool)

		// Create partner with pending status
		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID
		partnerEncx.StripeAccountStatus = domain.StripeAccountStatusPending
		partnerEncx.StripeOnboardingComplete = false

		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Update partner Stripe status to active
		partnerEncx.StripeAccountStatus = domain.StripeAccountStatusActive
		partnerEncx.StripeOnboardingComplete = true

		// Act
		err = repo.UpdatePartner(ctx, partnerEncx)

		// Assert
		assert.NoError(t, err)

		// Verify the status change
		updatedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
		require.NoError(t, err)
		assert.Equal(t, domain.StripeAccountStatusActive, updatedPartner.StripeAccountStatus)
		assert.True(t, updatedPartner.StripeOnboardingComplete)
	})

	t.Run("should successfully update partner with different Stripe statuses", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		stripeStatuses := []domain.StripeAccountStatus{domain.StripeAccountStatusPending, domain.StripeAccountStatusActive, domain.StripeAccountStatusRestricted, domain.StripeAccountStatusDisabled}

		for _, status := range stripeStatuses {
			userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, string(status))

			partnerEncx := td.NewTestPartnerEncx(t)
			partnerEncx.UserID = userID
			partnerEncx.StripeAccountStatus = domain.StripeAccountStatusPending
			partnerEncx.StripeOnboardingComplete = false

			err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
			require.NoError(t, err)

			// Update partner to new status
			partnerEncx.StripeAccountStatus = status
			partnerEncx.StripeOnboardingComplete = status == domain.StripeAccountStatusActive // Only active accounts have complete onboarding

			// Act
			err = repo.UpdatePartner(ctx, partnerEncx)

			// Assert
			assert.NoError(t, err, "Should successfully update partner with status %s", status)

			// Verify the status change
			updatedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
			require.NoError(t, err)
			assert.Equal(t, status, updatedPartner.StripeAccountStatus, "Status should match for %s", status)
			assert.Equal(t, status == domain.StripeAccountStatusActive, updatedPartner.StripeOnboardingComplete, "Onboarding should match for %s", status)
		}
	})

	t.Run("should successfully update partner with empty optional fields", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		userID := td.CreateTestUserForPartner(t, ctx, testPool)

		// Create partner with some fields populated
		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID

		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Update partner to empty optional encrypted fields
		partnerEncx.Bio = ""
		partnerEncx.Experience = ""
		partnerEncx.StripeConnectedAccountIDEncrypted = []byte("")

		// Act
		err = repo.UpdatePartner(ctx, partnerEncx)

		// Assert
		assert.NoError(t, err)

		// Verify the update
		updatedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
		require.NoError(t, err)

		assert.Empty(t, updatedPartner.Bio)
		assert.Empty(t, updatedPartner.Experience)
		assert.Empty(t, updatedPartner.StripeConnectedAccountIDEncrypted)

		// Non-encrypted fields should still be preserved
		assert.Equal(t, partnerEncx.StripeAccountStatus, updatedPartner.StripeAccountStatus)
		assert.Equal(t, partnerEncx.StripeOnboardingComplete, updatedPartner.StripeOnboardingComplete)
	})

	t.Run("should successfully update partner with large encrypted data", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "large_data")

		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID

		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Create large encrypted data
		longBio := string(make([]byte, 2000))
		for i := range longBio {
			longBio = longBio[:i] + "a" + longBio[i+1:]
		}
		partnerEncx.Bio = longBio

		longExperience := string(make([]byte, 3000))
		for i := range longExperience {
			longExperience = longExperience[:i] + "b" + longExperience[i+1:]
		}
		partnerEncx.Experience = longExperience

		// manyCertifications := []byte("Certification 1, Certification 2, Certification 3, Certification 4, Certification 5")

		// Act
		err = repo.UpdatePartner(ctx, partnerEncx)

		// Assert
		assert.NoError(t, err)

		// Verify the update
		updatedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
		require.NoError(t, err)

		assert.Greater(t, len(updatedPartner.Bio), 1500, "Bio should be large")
		assert.Greater(t, len(updatedPartner.Experience), 2000, "Experience should be large")
	})

	t.Run("should handle database connection errors", func(t *testing.T) {
		// This test would typically involve mocking the database connection
		// For now, we'll skip it since we're using real testcontainers
		t.Skip("Database connection error testing requires mocking")
	})

	t.Run("should successfully update all partner fields in complete partner flow", func(t *testing.T) {
		// Arrange - simulates the complete partner registration flow
		td.ClearPartnersTable(t, ctx, testPool)

		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "complete")

		// Start with minimal partner
		minimalPartner := &domain.PartnerEncx{
			UserID:     userID,
			Bio:        "",
			Experience: "",
			// CertificationsEncrypted:           []byte(""),
			CategoryIDs:                       []uuid.UUID{uuid.New()},
			ProductIDs:                        []uuid.UUID{uuid.New()},
			StripeConnectedAccountIDEncrypted: []byte(""),
			StripeAccountStatus:               domain.StripeAccountStatusPending,
			StripeOnboardingComplete:          false,
			DEKEncrypted:                      []byte("minimal_dek"),
			KeyVersion:                        1,
		}

		err := td.InsertPartnerEncx(t, ctx, minimalPartner, testPool)
		require.NoError(t, err)

		// Update with complete partner information
		bio := "Experienced healthcare professional with 15+ years of practice"
		experience := "Specialized in patient care and medical consultation"
		categories := []uuid.UUID{uuid.New()}
		products := []uuid.UUID{uuid.New()}
		stripeAccount := []byte("acct_complete_partner_123456789")

		minimalPartner.Bio = bio
		minimalPartner.Experience = experience
		minimalPartner.CategoryIDs = categories
		minimalPartner.ProductIDs = products
		minimalPartner.StripeConnectedAccountIDEncrypted = stripeAccount
		minimalPartner.StripeAccountStatus = domain.StripeAccountStatusActive
		minimalPartner.StripeOnboardingComplete = true
		minimalPartner.DEKEncrypted = []byte("complete_dek_encrypted")
		minimalPartner.KeyVersion = 2

		// Act
		err = repo.UpdatePartner(ctx, minimalPartner)

		// Assert
		assert.NoError(t, err)

		// Verify all fields were updated correctly
		completePartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
		require.NoError(t, err)

		assert.Equal(t, bio, completePartner.Bio)
		assert.Equal(t, experience, completePartner.Experience)
		assert.Equal(t, categories, completePartner.CategoryIDs)
		assert.Equal(t, products, completePartner.ProductIDs)
		assert.Equal(t, stripeAccount, completePartner.StripeConnectedAccountIDEncrypted)
		assert.Equal(t, domain.StripeAccountStatusActive, completePartner.StripeAccountStatus)
		assert.True(t, completePartner.StripeOnboardingComplete)
		assert.Equal(t, []byte("complete_dek_encrypted"), completePartner.DEKEncrypted)
		assert.Equal(t, 2, completePartner.KeyVersion)
	})

	t.Run("should update partner with different key versions", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "key_version")

		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID
		partnerEncx.KeyVersion = 1

		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Test updating with different key versions
		keyVersions := []int{2, 5, 10, 100}

		for _, keyVersion := range keyVersions {
			partnerEncx.KeyVersion = keyVersion
			partnerEncx.DEKEncrypted = []byte(fmt.Sprintf("dek_for_version_%d", keyVersion))

			// Act
			err = repo.UpdatePartner(ctx, partnerEncx)

			// Assert
			assert.NoError(t, err, "Should successfully update partner with key version %d", keyVersion)

			// Verify the update
			updatedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
			require.NoError(t, err)
			assert.Equal(t, keyVersion, updatedPartner.KeyVersion, "Key version should match for %d", keyVersion)
			assert.Equal(t, []byte(fmt.Sprintf("dek_for_version_%d", keyVersion)), updatedPartner.DEKEncrypted, "DEK should match for key version %d", keyVersion)
		}
	})
}
