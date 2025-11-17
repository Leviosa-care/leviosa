package partner_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME='^TestGetAllPartners$$' TEST_PATH=test/integration/authuser/partner/get_all_partners_test.go

func TestGetAllPartners(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get all partners", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create multiple test partners with timestamps to verify ordering
		createdPartners := make([]*domain.Partner, 0, 3)
		for i := 0; i < 3; i++ {
			user := td.NewTestUser(t,
				"partner"+string(rune(i))+"@example.com",
				"Partner",
				string(rune('A'+i)))
			user.State = domain.Active
			user.Role = identity.PartnerStr
			testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
			require.NoError(t, err)
			err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
			require.NoError(t, err)

			partner := td.NewTestPartner(t, user.ID)
			partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
			require.NoError(t, err)
			td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
			require.NoError(t, err)
			createdPartners = append(createdPartners, partner)

			// Create associated sessions
			now := time.Now()
			sessionID := uuid.New()

			// Generate valid base64 tokens for testing
			accessToken, err := session.GenerateToken()
			require.NoError(t, err)

			refreshToken, err := session.GenerateToken()
			require.NoError(t, err)

			standardSession := &session.Session{
				ID:           sessionID,
				UserID:       user.ID,
				Role:         identity.Partner,
				State:        session.SessionActive,
				CreatedAt:    now,
				ExpiresAt:    now.Add(24 * time.Hour),
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			}

			standardSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, standardSession)
			require.NoError(t, err)

			td.InsertSessionEncx(t, ctx, redisClient, standardSessionEncx, time.Hour)

			// Small delay to ensure different timestamps
			time.Sleep(10 * time.Millisecond)
		}

		// Act
		req := td.NewGetAllPartnersRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)

		// Assert HTTP response
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response body
		partners := td.ParsePartnersListResponse(t, resp)
		require.Len(t, partners, 3, "Expected 3 partners in response")

		// Verify partners are ordered by created_at DESC (newest first)
		for i := 0; i < len(partners)-1; i++ {
			assert.True(t, partners[i].CreatedAt.After(partners[i+1].CreatedAt) || partners[i].CreatedAt.Equal(partners[i+1].CreatedAt),
				"Partners should be ordered by created_at DESC")
		}

		// Verify each partner against database
		for _, responsePartner := range partners {
			// Get encrypted partner from database by ID
			partnerEncx, err := td.GetPartnerEncxByID(t, ctx, responsePartner.ID, testPool)
			require.NoError(t, err, "Failed to get partner from database")

			// Decrypt partner
			dbPartner, err := domain.DecryptPartnerEncx(ctx, crypto, partnerEncx)
			require.NoError(t, err, "Failed to decrypt partner")

			// Compare fields
			assert.Equal(t, dbPartner.ID, responsePartner.ID, "ID mismatch")
			assert.Equal(t, dbPartner.Bio, responsePartner.Bio, "Bio mismatch")
			assert.Equal(t, dbPartner.Experience, responsePartner.Experience, "Experience mismatch")
			assert.Equal(t, dbPartner.CategoryIDs, responsePartner.CategoryIDs, "CategoryIDs mismatch")
			assert.Equal(t, dbPartner.ProductIDs, responsePartner.ProductIDs, "ProductIDs mismatch")
		}
	})

	t.Run("should return empty array when no partners exist", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Act
		req := td.NewGetAllPartnersRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)

		// Assert HTTP response
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse and verify empty array structure
		partners := td.ParsePartnersListResponse(t, resp)
		assert.Empty(t, partners, "Expected empty partners array")
		assert.NotNil(t, partners, "Partners array should not be nil, should be empty array")
	})

	t.Run("should return 500 when partner DEK is corrupted", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create test user and partner
		user := td.NewTestUser(t, "partner@example.com", "John", "Partner")
		user.State = domain.Active
		user.Role = identity.PartnerStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		partner := td.NewTestPartner(t, user.ID)
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Corrupt the DEK to simulate decryption failure
		td.CorruptPartnerDEK(t, ctx, partner.ID, testPool)

		// Act
		req := td.NewGetAllPartnersRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)

		// Assert - should return 500 due to decryption failure
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should return 500 when key version is invalid", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create test user and partner
		user := td.NewTestUser(t, "partner2@example.com", "Jane", "Partner")
		user.State = domain.Active
		user.Role = identity.PartnerStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		partner := td.NewTestPartner(t, user.ID)
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Set an invalid key version to simulate decryption failure
		td.SetInvalidKeyVersion(t, ctx, partner.ID, testPool, 99999)

		// Act
		req := td.NewGetAllPartnersRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)

		// Assert - should return 500 due to decryption failure
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
