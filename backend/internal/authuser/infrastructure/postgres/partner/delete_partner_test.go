package partnerRepository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestDeletePartner TEST_PATH=internal/authuser/infrastructure/postgres/partner/delete_partner_test.go

func TestDeletePartner(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully delete existing partner", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)
		// Create a user first since partner has foreign key constraint
		userID := td.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID

		// Insert the partner
		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Verify partner exists before deletion
		exists, err := td.CheckPartnerExistsByUserID(t, ctx, userID, testPool)
		require.NoError(t, err)
		require.True(t, exists)

		// Act
		err = repo.DeletePartner(ctx, userID)

		// Assert
		assert.NoError(t, err)

		// Verify partner no longer exists
		existsAfter, err := td.CheckPartnerExistsByUserID(t, ctx, userID, testPool)
		assert.NoError(t, err)
		assert.False(t, existsAfter)
	})

	t.Run("should return not found error when deleting non-existent partner", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)
		nonExistentUserID := uuid.New()

		// Act
		err := repo.DeletePartner(ctx, nonExistentUserID)

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrRepositoryNotFound))
	})

	t.Run("should successfully delete partner and allow creating partner with same user", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)
		userID := td.CreateTestUserForPartner(t, ctx, testPool)

		// Create and insert first partner
		firstPartnerEncx := td.NewTestPartnerEncx(t)
		firstPartnerEncx.UserID = userID

		err := td.InsertPartnerEncx(t, ctx, firstPartnerEncx, testPool)
		require.NoError(t, err)

		// Delete first partner
		err = repo.DeletePartner(ctx, userID)
		assert.NoError(t, err)

		// Act - create second partner with same user ID
		secondPartnerEncx := td.NewTestPartnerEncx(t)
		secondPartnerEncx.UserID = userID
		secondPartnerEncx.StripeAccountStatus = domain.StripeAccountStatusActive // Different status to differentiate

		err = td.InsertPartnerEncx(t, ctx, secondPartnerEncx, testPool)
		require.NoError(t, err)

		// Assert
		exists, err := td.CheckPartnerExistsByUserID(t, ctx, userID, testPool)
		assert.NoError(t, err)
		assert.True(t, exists)

		// Verify the second partner exists and is different from the first
		retrievedPartner, err := td.GetPartnerEncxByUserID(t, ctx, userID, testPool)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedPartner)
		assert.Equal(t, secondPartnerEncx.StripeAccountStatus, retrievedPartner.StripeAccountStatus)
		assert.NotEqual(t, firstPartnerEncx.StripeAccountStatus, retrievedPartner.StripeAccountStatus)
	})

	t.Run("should handle deletion of partners with different stripe statuses", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		createPartner := func(userID uuid.UUID, stripeStatus domain.StripeAccountStatus) *domain.PartnerEncx {
			partner := td.NewTestPartnerEncx(t)
			partner.UserID = userID
			partner.StripeAccountStatus = stripeStatus
			return partner
		}

		// Create partners with different stripe statuses
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "status_pending")
		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "status_active")
		userID3 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "status_restricted")

		pendingPartner := createPartner(userID1, domain.StripeAccountStatusPending)
		activePartner := createPartner(userID2, domain.StripeAccountStatusActive)
		restrictedPartner := createPartner(userID3, domain.StripeAccountStatusRestricted)

		// Insert partners
		err := td.InsertPartnerEncx(t, ctx, pendingPartner, testPool)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, activePartner, testPool)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, restrictedPartner, testPool)
		require.NoError(t, err)

		// Act & Assert - delete each partner
		for _, partner := range []*domain.PartnerEncx{pendingPartner, activePartner, restrictedPartner} {
			err := repo.DeletePartner(ctx, partner.UserID)
			assert.NoError(t, err)

			// Verify deletion
			exists, err := td.CheckPartnerExistsByUserID(t, ctx, partner.UserID, testPool)
			assert.NoError(t, err)
			assert.False(t, exists)
		}
	})

	t.Run("should handle deletion of partners with different onboarding states", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		createPartner := func(userID uuid.UUID, onboardingComplete bool) *domain.PartnerEncx {
			partner := td.NewTestPartnerEncx(t)
			partner.UserID = userID
			partner.StripeOnboardingComplete = onboardingComplete
			return partner
		}

		// Create partners with different onboarding states
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "onboarding_false")
		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "onboarding_true")

		incompletePartner := createPartner(userID1, false)
		completePartner := createPartner(userID2, true)

		// Insert partners
		err := td.InsertPartnerEncx(t, ctx, incompletePartner, testPool)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, completePartner, testPool)
		require.NoError(t, err)

		// Act & Assert - delete each partner
		for _, partner := range []*domain.PartnerEncx{incompletePartner, completePartner} {
			err := repo.DeletePartner(ctx, partner.UserID)
			assert.NoError(t, err)

			// Verify deletion
			exists, err := td.CheckPartnerExistsByUserID(t, ctx, partner.UserID, testPool)
			assert.NoError(t, err)
			assert.False(t, exists)
		}
	})

	t.Run("should handle database connection errors", func(t *testing.T) {
		// This test would typically involve mocking the database connection
		// For now, we'll skip it since we're using real testcontainers
		t.Skip("Database connection error testing requires mocking")
	})

	t.Run("should handle cascade delete behavior", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)
		userID := td.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID

		// Insert the partner
		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Verify partner exists
		exists, err := td.CheckPartnerExistsByUserID(t, ctx, userID, testPool)
		require.NoError(t, err)
		require.True(t, exists)

		// Act - Delete the user (should cascade delete the partner due to ON DELETE CASCADE)
		err = td.DeleteUserEncx(t, ctx, userID, testPool)
		assert.NoError(t, err)

		// Assert - Partner should also be deleted due to cascade
		existsAfter, err := td.CheckPartnerExistsByUserID(t, ctx, userID, testPool)
		assert.NoError(t, err)
		assert.False(t, existsAfter)
	})
}
