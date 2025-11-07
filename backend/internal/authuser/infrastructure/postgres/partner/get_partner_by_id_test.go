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

// make test-func TEST_NAME=TestGetPartnerByID TEST_PATH=internal/authuser/infrastructure/postgres/partner/get_partner_by_id_test.go

func TestGetPartnerByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve partner by ID", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		// Create a user first since partner has foreign key constraint
		userID := td.CreateTestUserForPartner(t, ctx, testPool)

		// Create a partner with the user ID
		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID
		partnerEncx.CategoryIDs = []uuid.UUID{uuid.New()}
		partnerEncx.ProductIDs = []uuid.UUID{uuid.New()}
		partnerEncx.StripeAccountStatus = domain.StripeAccountStatusActive
		partnerEncx.StripeOnboardingComplete = true

		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act
		retrievedPartnerEncx, err := repo.GetPartnerByID(ctx, partnerEncx.ID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, retrievedPartnerEncx)
		assert.Equal(t, partnerEncx.ID, retrievedPartnerEncx.ID)
		assert.Equal(t, partnerEncx.UserID, retrievedPartnerEncx.UserID)
		assert.Equal(t, domain.StripeAccountStatusActive, retrievedPartnerEncx.StripeAccountStatus)
		assert.True(t, retrievedPartnerEncx.StripeOnboardingComplete)

		// Verify encrypted fields are populated
		assert.NotEmpty(t, retrievedPartnerEncx.Bio)
		assert.NotEmpty(t, retrievedPartnerEncx.Experience)
		assert.NotEmpty(t, retrievedPartnerEncx.CategoryIDs)
		assert.NotEmpty(t, retrievedPartnerEncx.ProductIDs)
		assert.NotEmpty(t, retrievedPartnerEncx.StripeConnectedAccountIDEncrypted)
		assert.NotEmpty(t, retrievedPartnerEncx.DEKEncrypted)
		assert.Greater(t, retrievedPartnerEncx.KeyVersion, 0)
	})

	t.Run("should return not found error when partner does not exist", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)
		nonExistentPartnerID := uuid.New()

		// Act
		partner, err := repo.GetPartnerByID(ctx, nonExistentPartnerID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, partner)
		assert.True(t, errors.Is(err, errs.ErrRepositoryNotFound))
	})

	t.Run("should handle partner with different Stripe statuses", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		stripeStatuses := []domain.StripeAccountStatus{domain.StripeAccountStatusPending, domain.StripeAccountStatusActive, domain.StripeAccountStatusRestricted, domain.StripeAccountStatusDisabled}

		for _, status := range stripeStatuses {
			userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, string(status))

			partnerEncx := td.NewTestPartnerEncx(t)
			partnerEncx.UserID = userID
			partnerEncx.StripeAccountStatus = status

			err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
			require.NoError(t, err)

			// Act
			retrievedPartnerEncx, err := repo.GetPartnerByID(ctx, partnerEncx.ID)

			// Assert
			assert.NoError(t, err, "Should successfully retrieve partner with status %s", status)
			assert.NotNil(t, retrievedPartnerEncx, "Partner should not be nil for status %s", status)
			assert.Equal(t, partnerEncx.ID, retrievedPartnerEncx.ID, "Partner ID should match for status %s", status)
			assert.Equal(t, status, retrievedPartnerEncx.StripeAccountStatus, "Stripe status should match for %s", status)
		}
	})

	t.Run("should handle partner with different onboarding completion states", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		onboardingStates := []bool{true, false}

		for i, isComplete := range onboardingStates {
			userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, fmt.Sprintf("onboarding_%d", i))

			partnerEncx := td.NewTestPartnerEncx(t)
			partnerEncx.UserID = userID
			partnerEncx.StripeOnboardingComplete = isComplete

			err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
			require.NoError(t, err)

			// Act
			retrievedPartnerEncx, err := repo.GetPartnerByID(ctx, partnerEncx.ID)

			// Assert
			assert.NoError(t, err, "Should successfully retrieve partner with onboarding %t", isComplete)
			assert.NotNil(t, retrievedPartnerEncx, "Partner should not be nil for onboarding %t", isComplete)
			assert.Equal(t, partnerEncx.ID, retrievedPartnerEncx.ID, "Partner ID should match for onboarding %t", isComplete)
			assert.Equal(t, isComplete, retrievedPartnerEncx.StripeOnboardingComplete, "Onboarding state should match for %t", isComplete)
		}
	})

	t.Run("should handle partner with minimal encrypted data", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "minimal")

		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID
		// Set all encrypted fields to empty bytes
		partnerEncx.Bio = ""
		partnerEncx.Experience = ""
		partnerEncx.CategoryIDs = []uuid.UUID{}
		partnerEncx.ProductIDs = []uuid.UUID{}
		partnerEncx.StripeConnectedAccountIDEncrypted = []byte("")

		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act
		retrievedPartnerEncx, err := repo.GetPartnerByID(ctx, partnerEncx.ID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, retrievedPartnerEncx)
		assert.Equal(t, partnerEncx.ID, retrievedPartnerEncx.ID)
		assert.Equal(t, userID, retrievedPartnerEncx.UserID)

		// Verify encrypted fields are empty (or nil for database representation)
		assert.Equal(t, 0, len(retrievedPartnerEncx.Bio))
		assert.Equal(t, 0, len(retrievedPartnerEncx.Experience))
		assert.Equal(t, 0, len(retrievedPartnerEncx.CategoryIDs))
		assert.Equal(t, 0, len(retrievedPartnerEncx.ProductIDs))
		assert.Equal(t, 0, len(retrievedPartnerEncx.StripeConnectedAccountIDEncrypted))

		// Non-encrypted fields should still be populated
		assert.NotEmpty(t, retrievedPartnerEncx.DEKEncrypted)
		assert.Greater(t, retrievedPartnerEncx.KeyVersion, 0)
	})

	t.Run("should handle partner with maximal encrypted data", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "maximal")

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
		retrievedPartnerEncx, err := repo.GetPartnerByID(ctx, partnerEncx.ID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, retrievedPartnerEncx)
		assert.Equal(t, partnerEncx.ID, retrievedPartnerEncx.ID)
		assert.Equal(t, userID, retrievedPartnerEncx.UserID)

		// Verify large encrypted fields are preserved
		assert.Greater(t, len(retrievedPartnerEncx.Bio), 500, "Bio should be large")
		assert.Greater(t, len(retrievedPartnerEncx.Experience), 1000, "Experience should be large")
		assert.Equal(t, partnerEncx.StripeConnectedAccountIDEncrypted, retrievedPartnerEncx.StripeConnectedAccountIDEncrypted)
	})

	t.Run("should retrieve correct partner when multiple partners exist", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		// Create multiple partners
		var partners []*domain.PartnerEncx
		for i := 0; i < 5; i++ {
			userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, fmt.Sprintf("multi_%d", i))

			partnerEncx := td.NewTestPartnerEncx(t)
			partnerEncx.UserID = userID
			partnerEncx.Bio = fmt.Sprintf("Bio for partner %d", i)

			err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
			require.NoError(t, err)

			partners = append(partners, partnerEncx)
		}

		// Act - retrieve the third partner
		targetPartner := partners[2]
		retrievedPartnerEncx, err := repo.GetPartnerByID(ctx, targetPartner.ID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, retrievedPartnerEncx)
		assert.Equal(t, targetPartner.ID, retrievedPartnerEncx.ID)
		assert.Equal(t, targetPartner.UserID, retrievedPartnerEncx.UserID)
		assert.Equal(t, targetPartner.Bio, retrievedPartnerEncx.Bio)
	})

	t.Run("should preserve timestamp fields correctly", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		userID := td.CreateTestUserForPartner(t, ctx, testPool)

		partnerEncx := td.NewTestPartnerEncx(t)
		partnerEncx.UserID = userID

		err := td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act
		retrievedPartnerEncx, err := repo.GetPartnerByID(ctx, partnerEncx.ID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, retrievedPartnerEncx)
		assert.False(t, retrievedPartnerEncx.CreatedAt.IsZero(), "CreatedAt should not be zero")
		assert.False(t, retrievedPartnerEncx.UpdatedAt.IsZero(), "UpdatedAt should not be zero")
		assert.WithinDuration(t, partnerEncx.CreatedAt, retrievedPartnerEncx.CreatedAt, 1e9, "CreatedAt should match within 1 second")
		assert.WithinDuration(t, partnerEncx.UpdatedAt, retrievedPartnerEncx.UpdatedAt, 1e9, "UpdatedAt should match within 1 second")
	})

	t.Run("should handle database connection errors", func(t *testing.T) {
		// This test would typically involve mocking the database connection
		// For now, we'll skip it since we're using real testcontainers
		t.Skip("Database connection error testing requires mocking")
	})
}
