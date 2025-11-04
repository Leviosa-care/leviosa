package partnerRepository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestVerifyPartner TEST_PATH=internal/authuser/infrastructure/postgres/partner/verify_partner_test.go

func TestVerifyPartner(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully verify partner with pending status", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		// Create a user first since partner has foreign key constraint
		userID := td.CreateTestUserForPartner(t, ctx, testPool)
		verifiedByUserID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "admin")

		// Create a partner with pending status
		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID
		partnerEncx.StripeAccountStatus = domain.StripeAccountStatusPending
		partnerEncx.StripeOnboardingComplete = false

		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Verify partner exists with pending status
		existingPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
		require.NoError(t, err)
		assert.Equal(t, domain.StripeAccountStatusPending, existingPartner.StripeAccountStatus)
		assert.False(t, existingPartner.StripeOnboardingComplete)

		// Act
		err = repo.VerifyPartner(ctx, userID, verifiedByUserID)

		// Assert
		assert.NoError(t, err)

		// Verify partner status was updated to active
		updatedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
		require.NoError(t, err)
		assert.Equal(t, domain.StripeAccountStatusActive, updatedPartner.StripeAccountStatus)
		assert.True(t, updatedPartner.StripeOnboardingComplete)

		// Verify updated_at timestamp was updated
		assert.True(t, updatedPartner.UpdatedAt.After(existingPartner.UpdatedAt))
	})

	t.Run("should successfully verify partner with restricted status", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		userID := td.CreateTestUserForPartner(t, ctx, testPool)
		verifiedByUserID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "admin")

		// Create a partner with restricted status
		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID
		partnerEncx.StripeAccountStatus = domain.StripeAccountStatusRestricted
		partnerEncx.StripeOnboardingComplete = false

		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act
		err = repo.VerifyPartner(ctx, userID, verifiedByUserID)

		// Assert
		assert.NoError(t, err)

		// Verify partner status was updated to active
		updatedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
		require.NoError(t, err)
		assert.Equal(t, domain.StripeAccountStatusActive, updatedPartner.StripeAccountStatus)
		assert.True(t, updatedPartner.StripeOnboardingComplete)
	})

	t.Run("should successfully verify partner with disabled status", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		userID := td.CreateTestUserForPartner(t, ctx, testPool)
		verifiedByUserID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "admin")

		// Create a partner with disabled status
		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID
		partnerEncx.StripeAccountStatus = domain.StripeAccountStatusDisabled
		partnerEncx.StripeOnboardingComplete = false

		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act
		err = repo.VerifyPartner(ctx, userID, verifiedByUserID)

		// Assert
		assert.NoError(t, err)

		// Verify partner status was updated to active
		updatedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
		require.NoError(t, err)
		assert.Equal(t, domain.StripeAccountStatusActive, updatedPartner.StripeAccountStatus)
		assert.True(t, updatedPartner.StripeOnboardingComplete)
	})

	t.Run("should return not found error when verifying non-existent partner", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		nonExistentUserID := uuid.New()
		verifiedByUserID := uuid.New()

		// Act
		err := repo.VerifyPartner(ctx, nonExistentUserID, verifiedByUserID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("should handle verification of partners with different initial statuses", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		verifiedByUserID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "admin")
		initialStatuses := []domain.StripeAccountStatus{
			domain.StripeAccountStatusPending,
			domain.StripeAccountStatusRestricted,
			domain.StripeAccountStatusDisabled,
		}

		for _, initialStatus := range initialStatuses {
			userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, string(initialStatus))

			partnerEncx := td.NewTestPartnerEncx(t)
			partnerEncx.UserID = userID
			partnerEncx.StripeAccountStatus = initialStatus
			partnerEncx.StripeOnboardingComplete = initialStatus == domain.StripeAccountStatusActive

			err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
			require.NoError(t, err)

			// Act
			err = repo.VerifyPartner(ctx, userID, verifiedByUserID)

			// Assert
			assert.NoError(t, err, "Should successfully verify partner with initial status %s", initialStatus)

			// Verify partner was updated to active status
			updatedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
			require.NoError(t, err)
			assert.Equal(t, domain.StripeAccountStatusActive, updatedPartner.StripeAccountStatus, "Partner should be active for initial status %s", initialStatus)
			assert.True(t, updatedPartner.StripeOnboardingComplete, "Onboarding should be complete for initial status %s", initialStatus)
		}
	})

	t.Run("should successfully verify partner that is already active", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		userID := td.CreateTestUserForPartner(t, ctx, testPool)
		verifiedByUserID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "admin")

		// Create a partner that is already active
		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID
		partnerEncx.StripeAccountStatus = domain.StripeAccountStatusActive
		partnerEncx.StripeOnboardingComplete = true

		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - verify an already active partner
		err = repo.VerifyPartner(ctx, userID, verifiedByUserID)

		// Assert
		assert.NoError(t, err)

		// Verify partner remains active
		updatedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
		require.NoError(t, err)
		assert.Equal(t, domain.StripeAccountStatusActive, updatedPartner.StripeAccountStatus)
		assert.True(t, updatedPartner.StripeOnboardingComplete)
	})

	t.Run("should handle verification with different verifiedByUserID values", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		userID := td.CreateTestUserForPartner(t, ctx, testPool)

		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID
		partnerEncx.StripeAccountStatus = domain.StripeAccountStatusPending
		partnerEncx.StripeOnboardingComplete = false

		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Test with different verifiedByUserID values
		testVerifierIDs := []uuid.UUID{
			uuid.New(),
			td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "verifier1"),
			td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "verifier2"),
		}

		for _, verifiedByUserID := range testVerifierIDs {
			err := td.DeletePartnerEncx(t, ctx, userID, testPool)
			require.NoError(t, err)

			// Reset partner to pending status
			partnerEncx.StripeAccountStatus = domain.StripeAccountStatusPending
			partnerEncx.StripeOnboardingComplete = false
			err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
			require.NoError(t, err)

			// Act
			err = repo.VerifyPartner(ctx, userID, verifiedByUserID)

			// Assert
			assert.NoError(t, err, "Should successfully verify partner with verifier ID %s", verifiedByUserID)

			// Verify partner was updated
			updatedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
			require.NoError(t, err)
			assert.Equal(t, domain.StripeAccountStatusActive, updatedPartner.StripeAccountStatus)
			assert.True(t, updatedPartner.StripeOnboardingComplete)
		}
	})

	t.Run("should handle verification of partners with minimal data", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		userID := td.CreateTestUserForPartner(t, ctx, testPool)
		verifiedByUserID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "admin")

		now := time.Now()

		// Create a partner with minimal data (all encrypted fields empty)
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
			CreatedAt:                         now,
			UpdatedAt:                         now,
			Metadata:                          encx.EncryptionMetadata{},
		}

		err := td.InsertPartnerEncx(t, ctx, minimalPartner, testPool)
		require.NoError(t, err)

		// Act
		err = repo.VerifyPartner(ctx, userID, verifiedByUserID)

		// Assert
		assert.NoError(t, err)

		// Verify partner status was updated to active
		updatedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
		require.NoError(t, err)
		assert.Equal(t, domain.StripeAccountStatusActive, updatedPartner.StripeAccountStatus)
		assert.True(t, updatedPartner.StripeOnboardingComplete)

		// Verify other fields remain unchanged (minimal)
		assert.Equal(t, "", updatedPartner.Bio)
		assert.Equal(t, "", updatedPartner.Experience)
	})

	t.Run("should handle verification of partners with maximal encrypted data", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "maximal")
		verifiedByUserID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "admin")

		// Create a partner with maximal data
		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID

		// Create large encrypted data
		longBio := string(make([]byte, 1000))
		for i := range longBio {
			longBio = longBio[:i] + "a" + longBio[i+1:]
		}
		partnerEncx.Bio = longBio

		longExperience := string(make([]byte, 2000))
		for i := range longExperience {
			longExperience = longExperience[:i] + "b" + longExperience[i+1:]
		}
		partnerEncx.Experience = longExperience

		partnerEncx.StripeConnectedAccountIDEncrypted = []byte("acct_test123456789abcdef")

		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act
		err = repo.VerifyPartner(ctx, userID, verifiedByUserID)

		// Assert
		assert.NoError(t, err)

		// Verify partner status was updated to active
		updatedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
		require.NoError(t, err)
		assert.Equal(t, domain.StripeAccountStatusActive, updatedPartner.StripeAccountStatus)
		assert.True(t, updatedPartner.StripeOnboardingComplete)

		// Verify large encrypted fields are preserved
		assert.Greater(t, len(updatedPartner.Bio), 500, "Bio should be large")
		assert.Greater(t, len(updatedPartner.Experience), 1000, "Experience should be large")
		assert.Equal(t, partnerEncx.StripeConnectedAccountIDEncrypted, updatedPartner.StripeConnectedAccountIDEncrypted)
	})

	t.Run("should handle verification of multiple partners in sequence", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		const numPartners = 5
		userIDs := make([]uuid.UUID, numPartners)
		verifiedByUserID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "admin")

		// Create multiple partners with different statuses
		for i := 0; i < numPartners; i++ {
			userIDs[i] = td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, fmt.Sprintf("partner%d", i))

			partnerEncx := td.NewTestPartnerEncx(t)
			partnerEncx.UserID = userIDs[i]
			partnerEncx.StripeAccountStatus = domain.StripeAccountStatusPending
			partnerEncx.StripeOnboardingComplete = false

			err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
			require.NoError(t, err)
		}

		// Act - verify all partners
		for _, userID := range userIDs {
			err := repo.VerifyPartner(ctx, userID, verifiedByUserID)
			assert.NoError(t, err, "Should successfully verify partner %s", userID)

			// Verify partner was updated
			updatedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
			require.NoError(t, err)
			assert.Equal(t, domain.StripeAccountStatusActive, updatedPartner.StripeAccountStatus, "Partner %s should be active", userID)
			assert.True(t, updatedPartner.StripeOnboardingComplete, "Partner %s should have complete onboarding", userID)
		}
	})

	t.Run("should handle database connection errors", func(t *testing.T) {
		// This test would typically involve mocking the database connection
		// For now, we'll skip it since we're using real testcontainers
		t.Skip("Database connection error testing requires mocking")
	})

	t.Run("should verify that only Stripe status and onboarding fields are updated", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		userID := td.CreateTestUserForPartner(t, ctx, testPool)
		verifiedByUserID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "admin")

		// Create partner with specific values
		originalBio := "Original bio content"
		originalExperience := "Original experience content"
		originalStripeAccountID := []byte("acct_original123456789")

		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID
		partnerEncx.Bio = originalBio
		partnerEncx.Experience = originalExperience
		partnerEncx.StripeConnectedAccountIDEncrypted = originalStripeAccountID
		partnerEncx.StripeAccountStatus = domain.StripeAccountStatusPending
		partnerEncx.StripeOnboardingComplete = false

		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Record original values
		originalPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
		require.NoError(t, err)
		originalUpdatedAt := originalPartner.UpdatedAt

		// Act
		err = repo.VerifyPartner(ctx, userID, verifiedByUserID)

		// Assert
		assert.NoError(t, err)

		// Verify partner was updated
		updatedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
		require.NoError(t, err)

		// Verify only status and onboarding fields changed
		assert.Equal(t, domain.StripeAccountStatusActive, updatedPartner.StripeAccountStatus)
		assert.True(t, updatedPartner.StripeOnboardingComplete)

		// Verify other fields are unchanged
		assert.Equal(t, originalBio, updatedPartner.Bio, "Bio should remain unchanged")
		assert.Equal(t, originalExperience, updatedPartner.Experience, "Experience should remain unchanged")
		assert.Equal(t, originalStripeAccountID, updatedPartner.StripeConnectedAccountIDEncrypted, "Stripe account ID should remain unchanged")
		assert.Equal(t, originalPartner.DEKEncrypted, updatedPartner.DEKEncrypted, "DEK should remain unchanged")
		assert.Equal(t, originalPartner.KeyVersion, updatedPartner.KeyVersion, "Key version should remain unchanged")

		// Verify updated_at timestamp was updated
		assert.True(t, updatedPartner.UpdatedAt.After(originalUpdatedAt), "Updated_at should be newer")
	})
}
