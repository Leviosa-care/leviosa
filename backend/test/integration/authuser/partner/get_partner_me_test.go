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

// make test-func TEST_NAME=TestGetPartnerMe TEST_PATH=test/integration/authuser/partner/get_partner_me_test.go

func TestGetPartnerMe(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get own partner profile", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create user
		user := td.NewTestUser(t, "partner@example.com", "John", "Partner")
		user.State = domain.Active
		user.Role = identity.PartnerStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Create partner
		partner := td.NewTestPartner(t, user.ID)
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Create session for user partner
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

		// Act - get own partner profile
		req := td.NewGetPartnerMeRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert HTTP response
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response body
		responsePartner := td.ParsePartnerResponse(t, resp)
		require.NotNil(t, responsePartner, "Response partner should not be nil")

		// Get encrypted partner from database
		dbPartnerEncx, err := td.GetPartnerEncxByUserID(t, ctx, user.ID, testPool)
		require.NoError(t, err, "Failed to get partner from database")

		// Decrypt partner
		dbPartner, err := domain.DecryptPartnerEncx(ctx, crypto, dbPartnerEncx)
		require.NoError(t, err, "Failed to decrypt partner")

		// Verify all fields match database
		assert.Equal(t, dbPartner.ID, responsePartner.ID, "UserID mismatch")
		assert.Equal(t, dbPartner.Bio, responsePartner.Bio, "Bio mismatch")
		assert.Equal(t, dbPartner.Experience, responsePartner.Experience, "Experience mismatch")
		assert.Equal(t, dbPartner.CategoryIDs, responsePartner.CategoryIDs, "CategoryIDs mismatch")
		assert.Equal(t, dbPartner.ProductIDs, responsePartner.ProductIDs, "ProductIDs mismatch")
		assert.WithinDuration(t, dbPartner.CreatedAt, responsePartner.CreatedAt, time.Second, "CreatedAt mismatch")
		assert.WithinDuration(t, dbPartner.UpdatedAt, responsePartner.UpdatedAt, time.Second, "UpdatedAt mismatch")
	})

	t.Run("should return 404 when authenticated user has no partner profile", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create user without partner profile
		userWithoutPartner := td.NewTestUser(t, "user@example.com", "Regular", "User")
		userWithoutPartner.State = domain.Active
		userWithoutPartner.Role = identity.PartnerStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, userWithoutPartner)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Create session for user without partner profile
		now := time.Now()
		sessionID := uuid.New()

		// Generate valid base64 tokens for testing
		accessToken, err := session.GenerateToken()
		require.NoError(t, err)

		refreshToken, err := session.GenerateToken()
		require.NoError(t, err)

		standardSession := &session.Session{
			ID:           sessionID,
			UserID:       userWithoutPartner.ID,
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

		// Act - try to get own partner profile (doesn't exist)
		req := td.NewGetPartnerMeRequest(t, ctx, testServerURL, accessToken)

		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Act - try to access /me without session cookie
		req := td.NewGetPartnerMeRequest(t, ctx, testServerURL, "")

		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 500 when partner DEK is corrupted", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

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

		// Create session
		now := time.Now()
		sessionID := uuid.New()
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

		// Corrupt the DEK to simulate decryption failure
		td.CorruptPartnerDEK(t, ctx, partner.ID, testPool)

		// Act
		req := td.NewGetPartnerMeRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert - should return 500 due to decryption failure
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should return 500 when key version is invalid", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

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

		// Create session
		now := time.Now()
		sessionID := uuid.New()
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

		// Set an invalid key version to simulate decryption failure
		td.SetInvalidKeyVersion(t, ctx, partner.ID, testPool, 99999)

		// Act
		req := td.NewGetPartnerMeRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert - should return 500 due to decryption failure
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
