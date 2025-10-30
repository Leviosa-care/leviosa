package partnerRepository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAllPartners TEST_PATH=internal/authuser/infrastructure/postgres/partner/get_all_partners_test.go

func TestGetAllPartners(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve all partners ordered by created_at DESC", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		createPartner := func(userID uuid.UUID, stripeStatus domain.StripeAccountStatus, onboardingComplete bool) *domain.PartnerEncx {
			partner := td.NewTestPartnerEncx(t)
			partner.UserID = userID
			partner.StripeAccountStatus = stripeStatus
			partner.StripeOnboardingComplete = onboardingComplete

			err := td.InsertPartnerEncx(t, ctx, partner, testPool)
			require.NoError(t, err)
			return partner
		}

		// Create test users first (foreign key constraint)
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "user1")
		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "user2")
		userID3 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "user3")

		// Create test partners with different Stripe statuses
		partner1 := createPartner(userID1, domain.StripeAccountStatusPending, false)
		partner2 := createPartner(userID2, domain.StripeAccountStatusActive, true)
		partner3 := createPartner(userID3, domain.StripeAccountStatusRestricted, false)

		// Act
		allPartners, err := repo.GetAllPartners(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, allPartners, 3, "Should return all partners regardless of Stripe status")

		// Verify partners are ordered by created_at DESC (newest first)
		// Since partner3 was inserted last, it should come first
		assert.Equal(t, partner3.UserID, allPartners[0].UserID, "First partner should be the most recently created")
		assert.Equal(t, partner2.UserID, allPartners[1].UserID, "Second partner should be the second created")
		assert.Equal(t, partner1.UserID, allPartners[2].UserID, "Third partner should be the first created")

		// Verify all partner Stripe statuses are preserved
		partnerStatusMap := make(map[uuid.UUID]domain.StripeAccountStatus)
		for _, partner := range allPartners {
			partnerStatusMap[partner.UserID] = partner.StripeAccountStatus
		}
		assert.Equal(t, domain.StripeAccountStatusPending, partnerStatusMap[partner1.UserID])
		assert.Equal(t, domain.StripeAccountStatusActive, partnerStatusMap[partner2.UserID])
		assert.Equal(t, domain.StripeAccountStatusRestricted, partnerStatusMap[partner3.UserID])

		// Verify onboarding states are preserved
		onboardingMap := make(map[uuid.UUID]bool)
		for _, partner := range allPartners {
			onboardingMap[partner.UserID] = partner.StripeOnboardingComplete
		}
		assert.Equal(t, false, onboardingMap[partner1.UserID])
		assert.Equal(t, true, onboardingMap[partner2.UserID])
		assert.Equal(t, false, onboardingMap[partner3.UserID])

		// Verify encrypted fields are populated (not decrypted at repository layer)
		for _, partner := range allPartners {
			assert.NotEmpty(t, partner.BioEncrypted)
			assert.NotEmpty(t, partner.ExperienceEncrypted)
			assert.NotEmpty(t, partner.CertificationsEncrypted)
			assert.NotEmpty(t, partner.CategoryIDsEncrypted)
			assert.NotEmpty(t, partner.ProductIDsEncrypted)
			assert.NotEmpty(t, partner.StripeConnectedAccountIDEncrypted)
			assert.NotEmpty(t, partner.DEKEncrypted)
			assert.Greater(t, partner.KeyVersion, 0)
		}
	})

	t.Run("should return empty slice when no partners exist", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		// Act
		allPartners, err := repo.GetAllPartners(ctx)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, allPartners, "Should return non-nil slice")
		assert.Empty(t, allPartners, "Should return empty slice when no partners exist")
	})

	t.Run("should handle partners with different Stripe account statuses", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		stripeStatuses := []domain.StripeAccountStatus{
			domain.StripeAccountStatusPending,
			domain.StripeAccountStatusActive,
			domain.StripeAccountStatusRestricted,
			domain.StripeAccountStatusDisabled,
		}

		expectedPartners := make(map[domain.StripeAccountStatus]*domain.PartnerEncx)

		for i, status := range stripeStatuses {
			userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, fmt.Sprintf("status_%d", i))
			partner := td.NewTestPartnerEncx(t)
			partner.UserID = userID
			partner.StripeAccountStatus = status

			err := td.InsertPartnerEncx(t, ctx, partner, testPool)
			require.NoError(t, err)

			expectedPartners[status] = partner
		}

		// Act
		allPartners, err := repo.GetAllPartners(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, allPartners, len(stripeStatuses), "Should return partners with all Stripe statuses")

		// Verify all Stripe statuses are represented
		foundStatuses := make(map[domain.StripeAccountStatus]bool)
		for _, partner := range allPartners {
			foundStatuses[partner.StripeAccountStatus] = true
		}

		for _, expectedStatus := range stripeStatuses {
			assert.True(t, foundStatuses[expectedStatus], "Should find partner with Stripe status: %v", expectedStatus)
		}
	})

	t.Run("should handle partners with different onboarding completion states", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		onboardingStates := []bool{true, false}
		expectedPartners := make(map[bool]*domain.PartnerEncx)

		for i, isComplete := range onboardingStates {
			userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, fmt.Sprintf("onboarding_%d", i))
			partner := td.NewTestPartnerEncx(t)
			partner.UserID = userID
			partner.StripeOnboardingComplete = isComplete

			err := td.InsertPartnerEncx(t, ctx, partner, testPool)
			require.NoError(t, err)

			expectedPartners[isComplete] = partner
		}

		// Act
		allPartners, err := repo.GetAllPartners(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, allPartners, len(onboardingStates), "Should return partners with both onboarding states")

		// Verify both onboarding states are represented
		foundStates := make(map[bool]bool)
		for _, partner := range allPartners {
			foundStates[partner.StripeOnboardingComplete] = true
		}

		for _, expectedState := range onboardingStates {
			assert.True(t, foundStates[expectedState], "Should find partner with onboarding state: %v", expectedState)
		}
	})

	t.Run("should handle partners with minimal and maximal encrypted data", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		// Partner with minimal data (empty encrypted fields)
		minimalUserID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "minimal")
		minimalPartner := td.NewTestPartnerEncx(t)
		minimalPartner.UserID = minimalUserID
		minimalPartner.BioEncrypted = []byte("")
		minimalPartner.ExperienceEncrypted = []byte("")
		minimalPartner.CertificationsEncrypted = []byte("")
		minimalPartner.CategoryIDsEncrypted = []byte("")
		minimalPartner.ProductIDsEncrypted = []byte("")
		minimalPartner.StripeConnectedAccountIDEncrypted = []byte("")

		err := td.InsertPartnerEncx(t, ctx, minimalPartner, testPool)
		require.NoError(t, err)

		// Partner with maximal data (large encrypted fields)
		maximalUserID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "maximal")
		maximalPartner := td.NewTestPartnerEncx(t)
		maximalPartner.UserID = maximalUserID

		longBio := string(make([]byte, 1000))
		for i := range longBio {
			longBio = longBio[:i] + "a" + longBio[i+1:]
		}
		maximalPartner.BioEncrypted = []byte(longBio)

		longExperience := string(make([]byte, 2000))
		for i := range longExperience {
			longExperience = longExperience[:i] + "b" + longExperience[i+1:]
		}
		maximalPartner.ExperienceEncrypted = []byte(longExperience)

		err = td.InsertPartnerEncx(t, ctx, maximalPartner, testPool)
		require.NoError(t, err)

		// Act
		allPartners, err := repo.GetAllPartners(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, allPartners, 2, "Should return both partners")

		// Find partners in results
		var minimalFound, maximalFound *domain.PartnerEncx
		for _, partner := range allPartners {
			if partner.UserID == minimalUserID {
				minimalFound = partner
			} else if partner.UserID == maximalUserID {
				maximalFound = partner
			}
		}

		assert.NotNil(t, minimalFound, "Minimal partner should be found")
		assert.NotNil(t, maximalFound, "Maximal partner should be found")

		// Verify encrypted field sizes
		assert.Equal(t, 0, len(minimalFound.BioEncrypted), "Minimal partner should have empty bio")
		assert.Equal(t, 0, len(minimalFound.ExperienceEncrypted), "Minimal partner should have empty experience")
		assert.Greater(t, len(maximalFound.BioEncrypted), 500, "Maximal partner should have large bio")
		assert.Greater(t, len(maximalFound.ExperienceEncrypted), 1000, "Maximal partner should have large experience")
	})

	t.Run("should handle large number of partners", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		const numPartners = 50
		expectedPartners := make([]*domain.PartnerEncx, numPartners)
		stripeStatuses := []domain.StripeAccountStatus{domain.StripeAccountStatusPending, domain.StripeAccountStatusActive, domain.StripeAccountStatusRestricted, domain.StripeAccountStatusDisabled}

		// Create many partners with various configurations
		for i := 0; i < numPartners; i++ {
			userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, fmt.Sprintf("partner%d", i))
			partner := td.NewTestPartnerEncx(t)
			partner.UserID = userID
			partner.StripeAccountStatus = stripeStatuses[i%len(stripeStatuses)]
			partner.StripeOnboardingComplete = i%2 == 0 // Alternate between true and false

			err := td.InsertPartnerEncx(t, ctx, partner, testPool)
			require.NoError(t, err)

			expectedPartners[i] = partner
		}

		// Act
		allPartners, err := repo.GetAllPartners(ctx)

		// Assert
		require.NoError(t, err)
		assert.Len(t, allPartners, numPartners, "Should return all partners")

		// Verify order (newest first - reverse order of insertion)
		assert.Equal(t, expectedPartners[numPartners-1].UserID, allPartners[0].UserID, "First should be last inserted")
		assert.Equal(t, expectedPartners[0].UserID, allPartners[numPartners-1].UserID, "Last should be first inserted")

		// Verify Stripe status distribution
		statusCount := make(map[domain.StripeAccountStatus]int)
		for _, partner := range allPartners {
			statusCount[partner.StripeAccountStatus]++
		}

		// Each status should appear roughly numPartners/len(statuses) times
		expectedCount := numPartners / len(stripeStatuses)
		for _, status := range stripeStatuses {
			assert.GreaterOrEqual(t, statusCount[status], expectedCount-1, "Status %v should appear at least %d times", status, expectedCount-1)
			assert.LessOrEqual(t, statusCount[status], expectedCount+1, "Status %v should appear at most %d times", status, expectedCount+1)
		}

		// Verify onboarding completion distribution
		onboardingCount := make(map[bool]int)
		for _, partner := range allPartners {
			onboardingCount[partner.StripeOnboardingComplete]++
		}

		// Should have roughly equal distribution
		expectedOnboardingCount := numPartners / 2
		assert.GreaterOrEqual(t, onboardingCount[true], expectedOnboardingCount-1, "Complete onboarding should appear at least %d times", expectedOnboardingCount-1)
		assert.GreaterOrEqual(t, onboardingCount[false], expectedOnboardingCount-1, "Incomplete onboarding should appear at least %d times", expectedOnboardingCount-1)
	})

	t.Run("should preserve all partner fields correctly", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "preserve")

		originalPartner := td.NewTestPartnerEncx(t)
		originalPartner.UserID = userID
		originalPartner.StripeAccountStatus = domain.StripeAccountStatusActive
		originalPartner.StripeOnboardingComplete = true
		originalPartner.StripeConnectedAccountIDEncrypted = []byte("acct_test123456789")

		err := td.InsertPartnerEncx(t, ctx, originalPartner, testPool)
		require.NoError(t, err)

		// Act
		allPartners, err := repo.GetAllPartners(ctx)

		// Assert
		require.NoError(t, err)
		assert.Len(t, allPartners, 1, "Should return exactly one partner")

		retrievedPartner := allPartners[0]

		// Verify all fields are preserved
		assert.Equal(t, originalPartner.UserID, retrievedPartner.UserID)
		assert.Equal(t, originalPartner.StripeAccountStatus, retrievedPartner.StripeAccountStatus)
		assert.Equal(t, originalPartner.StripeOnboardingComplete, retrievedPartner.StripeOnboardingComplete)
		assert.Equal(t, originalPartner.StripeConnectedAccountIDEncrypted, retrievedPartner.StripeConnectedAccountIDEncrypted)
		assert.Equal(t, originalPartner.BioEncrypted, retrievedPartner.BioEncrypted)
		assert.Equal(t, originalPartner.ExperienceEncrypted, retrievedPartner.ExperienceEncrypted)
		assert.Equal(t, originalPartner.CertificationsEncrypted, retrievedPartner.CertificationsEncrypted)
		assert.Equal(t, originalPartner.CategoryIDsEncrypted, retrievedPartner.CategoryIDsEncrypted)
		assert.Equal(t, originalPartner.ProductIDsEncrypted, retrievedPartner.ProductIDsEncrypted)
		assert.Equal(t, originalPartner.DEKEncrypted, retrievedPartner.DEKEncrypted)
		assert.Equal(t, originalPartner.KeyVersion, retrievedPartner.KeyVersion)

		// Verify timestamps are preserved (ignoring microsecond differences)
		assert.WithinDuration(t, originalPartner.CreatedAt, retrievedPartner.CreatedAt, time.Second)
		assert.WithinDuration(t, originalPartner.UpdatedAt, retrievedPartner.UpdatedAt, time.Second)
	})

	t.Run("should handle database connection errors gracefully", func(t *testing.T) {
		// This test would typically involve mocking the database connection
		// For comprehensive testing, we'd need to simulate connection failures
		t.Skip("Database connection error testing requires mocking or network disruption")
	})
}
