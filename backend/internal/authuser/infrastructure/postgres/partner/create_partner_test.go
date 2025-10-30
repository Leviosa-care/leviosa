package partnerRepository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreatePartner TEST_PATH=internal/authuser/infrastructure/postgres/partner/create_partner_test.go

func TestCreatePartner(t *testing.T) {
	ctx := context.Background()

	const existsQuery = `SELECT EXISTS(SELECT 1 FROM auth.partners WHERE user_id = $1)`
	const existsByUserIDQuery = `SELECT EXISTS(SELECT 1 FROM auth.partners WHERE user_id = $1)`

	t.Run("should successfully create a new partner", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)
		// Create a user first since partner has foreign key constraint
		userID := td.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID

		// Act
		err := repo.CreatePartner(ctx, partnerEncx)

		// Assert
		assert.NoError(t, err)

		// Verify partner was created
		var exists bool
		err = testPool.QueryRow(ctx, existsQuery, partnerEncx.UserID).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should fail to create partner with duplicate user_id", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)
		// Create a user for the first partner
		userID := td.CreateTestUserForPartner(t, ctx, testPool)

		// Create first partner
		firstPartnerEncx := td.NewTestPartnerEncx(t)
		firstPartnerEncx.UserID = userID
		err := repo.CreatePartner(ctx, firstPartnerEncx)
		require.NoError(t, err)

		// Create second partner with same user_id
		duplicatePartnerEncx := td.NewTestPartnerEncx(t)
		duplicatePartnerEncx.UserID = userID

		// Act
		err = repo.CreatePartner(ctx, duplicatePartnerEncx)

		// Assert
		assert.ErrorIs(t, err, errs.ErrUniqueViolation, "Should be a unique constraint violation")
	})

	t.Run("should fail to create partner without required fields", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		// Partner missing required encrypted fields (user_id, dek_encrypted, etc.)
		invalidPartner := &domain.PartnerEncx{
			// Missing required fields
		}

		// Act
		err := repo.CreatePartner(ctx, invalidPartner)

		// Assert
		assert.Error(t, err)
		// Should fail due to NOT NULL constraints
	})

	t.Run("should handle partner with minimal required fields", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)
		userID := td.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID

		// Act
		err := repo.CreatePartner(ctx, partnerEncx)

		// Assert
		assert.NoError(t, err)

		// Verify partner was created
		var exists bool
		err = testPool.QueryRow(ctx, existsQuery, partnerEncx.UserID).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should handle partner with all fields populated", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)
		userID := td.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID

		// Act
		err := repo.CreatePartner(ctx, partnerEncx)

		// Assert
		assert.NoError(t, err)

		// Verify partner was created
		exists, err := td.CheckPartnerExistsByUserID(t, ctx, partnerEncx.UserID, testPool)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should handle different stripe account statuses", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		statuses := []domain.StripeAccountStatus{
			domain.StripeAccountStatusPending,
			domain.StripeAccountStatusActive,
			domain.StripeAccountStatusRestricted,
			domain.StripeAccountStatusDisabled,
		}

		for i, status := range statuses {
			userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, fmt.Sprintf("status_%d", i))
			partnerEncx := td.NewTestPartnerEncx(t)
			partnerEncx.UserID = userID
			partnerEncx.StripeAccountStatus = status

			// Act
			err := repo.CreatePartner(ctx, partnerEncx)

			// Assert
			assert.NoError(t, err, "Should create partner with status %s", status)

			// Verify partner was created
			exists, err := td.CheckPartnerExistsByUserID(t, ctx, userID, testPool)
			assert.NoError(t, err)
			assert.True(t, exists, "Partner with status %s should exist", status)
		}
	})

	t.Run("should handle very long text fields", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		longBio := string(make([]byte, 1000)) // 1000 character bio
		for i := range longBio {
			longBio = longBio[:i] + "a" + longBio[i+1:]
		}

		longExperience := string(make([]byte, 2000)) // 2000 character experience
		for i := range longExperience {
			longExperience = longExperience[:i] + "b" + longExperience[i+1:]
		}

		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "")
		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID
		partnerEncx.BioEncrypted = []byte(longBio)
		partnerEncx.ExperienceEncrypted = []byte(longExperience)

		// Act
		err := repo.CreatePartner(ctx, partnerEncx)

		// Assert
		assert.NoError(t, err)

		// Verify partner was created
		exists, err := td.CheckPartnerExistsByUserID(t, ctx, partnerEncx.UserID, testPool)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should fail when context is cancelled", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		partnerEncx := td.NewTestPartnerEncx(t)

		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		// Act
		err := repo.CreatePartner(cancelledCtx, partnerEncx)

		// Assert
		assert.Error(t, err)
		// Should be classified as a context-related error by ClassifyPgError
	})

	t.Run("should handle concurrent partner creation with different user IDs", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		const count = 3
		userIDs := [count]uuid.UUID{}

		results := make(chan error, count)

		// Act - create partners concurrently
		for i := range count {
			go func(idx int, userIDs *[count]uuid.UUID) {
				uid := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, fmt.Sprintf("%d", i))
				partnerEncx := td.NewTestPartnerEncx(t)
				partnerEncx.UserID = uid
				userIDs[i] = uid

				err := repo.CreatePartner(ctx, partnerEncx)
				results <- err
			}(i, &userIDs)
		}

		// Assert - collect results
		successCount := 0
		for range count {
			err := <-results
			if err == nil {
				successCount++
			}
		}

		assert.Equal(t, count, successCount, "All concurrent creations should succeed")

		// Verify all partners were created
		for _, userID := range userIDs {
			var exists bool
			err := testPool.QueryRow(ctx, existsByUserIDQuery, userID).Scan(&exists)

			assert.NoError(t, err)
			assert.True(t, exists, "Partner with user ID %s should exist", userID)
		}
	})

	t.Run("should handle different onboarding completion states", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		onboardingStates := []bool{
			true,  // Onboarding complete
			false, // Onboarding not complete
		}

		for _, isComplete := range onboardingStates {
			userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, fmt.Sprintf("%t", isComplete))
			partnerEncx := td.NewTestPartnerEncx(t)
			partnerEncx.UserID = userID
			partnerEncx.StripeOnboardingComplete = isComplete

			// Act
			err := repo.CreatePartner(ctx, partnerEncx)

			// Assert
			assert.NoError(t, err, "Should create partner with onboarding complete %t", isComplete)

			// Verify partner was created
			exists, err := td.CheckPartnerExistsByUserID(t, ctx, userID, testPool)
			assert.NoError(t, err)
			assert.True(t, exists, "Partner with onboarding complete %t should exist", isComplete)
		}
	})

	t.Run("should verify partner count increases correctly", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		// Get initial count
		initialCount, err := td.CountPartners(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, 0, initialCount, "Table should be empty initially")

		// Create multiple partners
		const numPartners = 3
		for i := 0; i < numPartners; i++ {
			userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, fmt.Sprintf("%d", i))
			partnerEncx := td.NewTestPartnerEncx(t)
			partnerEncx.UserID = userID
			err := repo.CreatePartner(ctx, partnerEncx)
			assert.NoError(t, err)
		}

		// Verify final count
		finalCount, err := td.CountPartners(t, ctx, testPool)
		assert.NoError(t, err)
		assert.Equal(t, numPartners, finalCount, "Should have %d partners", numPartners)
	})
}
