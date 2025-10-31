package partner_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestDeletePartner make test-integration-partner-test

func TestDeletePartner(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully delete partner but preserve user account", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create test user and partner
		testUser := td.NewTestUser(t, "partner@example.com", "John", "Partner")
		testUser.State = domain.Active
		testUser.Role = identity.PartnerStr
		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool, crypto)
		require.NoError(t, err)

		testPartner := &domain.Partner{
			ID:             uuid.New(),
			UserID:         testUser.ID,
			Bio:            "Test bio",
			Experience:     "Test experience",
			// Certifications: []string{"Cert1"},
			CategoryIDs:    []uuid.UUID{},
			ProductIDs:     []uuid.UUID{},
			IsVerified:     true,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		// Act
		req := td.NewDeletePartnerRequest(t, ctx, testServerURL, testPartner.ID)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify partner was deleted
		_, err = td.GetPartnerByUserID(t, ctx, testUser.ID, testPool)
		assert.Error(t, err, "Partner should be deleted")

		// Verify user still exists
		userEncx, err := td.GetUserEnxByID(t, ctx, testUser.ID, testPool, crypto)
		require.NoError(t, err, "User should still exist after partner deletion")
		user, err := domain.DecryptUserEncx(ctx, crypto, userEncx)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, user.ID)
	})

	t.Run("should return 404 when partner not found", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Act
		nonExistentID := uuid.New()
		req := td.NewDeletePartnerRequest(t, ctx, testServerURL, nonExistentID)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 403 for non-admin user", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create standard user
		standardUser := td.NewTestUser(t, "standard@example.com", "Standard", "User")
		standardUser.State = domain.Active
		standardUserEncx, err := domain.ProcessUserEncx(ctx, crypto, standardUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, standardUserEncx, testPool, crypto)
		require.NoError(t, err)

		accessToken := tu.SetupSessionForUser(t, ctx, authCtx, standardUser.ID, identity.Standard)

		// Create partner to delete
		partnerUser := td.NewTestUser(t, "partner@example.com", "Partner", "User")
		partnerUser.State = domain.Active
		partnerUserEncx, err := domain.ProcessUserEncx(ctx, crypto, partnerUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, partnerUserEncx, testPool, crypto)
		require.NoError(t, err)

		testPartner := &domain.Partner{
			ID:      uuid.New(),
			UserID:  partnerUser.ID,
			Bio:     "Test bio",
			IsVerified: false,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		// Act
		req := td.NewDeletePartnerRequest(t, ctx, testServerURL, testPartner.ID)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create partner to delete
		testUser := td.NewTestUser(t, "partner@example.com", "Partner", "User")
		testUser.State = domain.Active
		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool, crypto)
		require.NoError(t, err)

		testPartner := &domain.Partner{
			ID:      uuid.New(),
			UserID:  testUser.ID,
			Bio:     "Test bio",
			IsVerified: false,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		// Act
		req := td.NewDeletePartnerRequest(t, ctx, testServerURL, testPartner.ID)
		// No session cookie

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 400 for invalid UUID format", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Act
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodDelete,
			testServerURL+"/admin/partners/invalid-uuid",
			nil,
		)
		require.NoError(t, err)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
