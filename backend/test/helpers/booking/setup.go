// Package booking provides test setup helpers for booking integration tests.
//
// This file contains helper functions for setting up authenticated users with
// proper allocations and sessions for testing booking endpoints.
package booking

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	userDomain "github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	th "github.com/Leviosa-care/leviosa/backend/test/helpers"
	talloc "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/allocation"
)

// SetupPartnerUser creates a partner user with partner record in the database.
// Returns the user ID for creating allocations.
func SetupPartnerUser(t *testing.T, ctx context.Context, email string, pool *pgxpool.Pool, crypto encx.CryptoService) uuid.UUID {
	t.Helper()

	user := th.NewTestUser(t, email, "John", "DOE")
	user.Role = identity.Partner.String()
	userEncx, err := userDomain.ProcessUserEncx(ctx, crypto, user)
	require.NoError(t, err)
	err = th.InsertUserEncx(t, ctx, userEncx, pool)
	require.NoError(t, err)

	partner := th.NewTestPartner(t, user.ID)
	partner.StripeAccountStatus = userDomain.StripeAccountStatusActive
	partner.StripeOnboardingComplete = true
	partnerEncx, err := userDomain.ProcessPartnerEncx(ctx, crypto, partner)
	require.NoError(t, err)
	err = th.InsertPartnerEncx(t, ctx, partnerEncx, pool)
	require.NoError(t, err)

	return user.ID
}

// CreateSessionForUser creates an active session for a user with the specified role.
// Returns the access token for HTTP authentication.
func CreateSessionForUser(t *testing.T, ctx context.Context, userID uuid.UUID, role identity.Role, redisClient *redis.Client, crypto encx.CryptoService) string {
	t.Helper()

	now := time.Now()
	sessionID := uuid.New()

	// Generate session tokens
	accessToken, err := session.GenerateToken()
	require.NoError(t, err, "Failed to generate access token")

	refreshToken, err := session.GenerateToken()
	require.NoError(t, err, "Failed to generate refresh token")

	// Create session object
	sess := &session.Session{
		ID:           sessionID,
		UserID:       userID,
		Role:         role,
		State:        session.SessionActive,
		CreatedAt:    now,
		ExpiresAt:    now.Add(24 * time.Hour),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// Encrypt session data
	sessionEncx, err := session.ProcessSessionEncx(ctx, crypto, sess)
	require.NoError(t, err, "Failed to encrypt session")

	// Store session in Redis
	sessionKey := session.FormatSessionKey(sessionEncx.ID.String())
	accessTokenKey := session.FormatAccessTokenKey(sessionEncx.AccessTokenHash)
	refreshTokenKey := session.FormatRefreshTokenKey(sessionEncx.RefreshTokenHash)
	userSessionIndexKey := session.FormatUserSessionIndexKey(sessionEncx.UserIDHash)

	// Encode session to JSON
	sessionData, err := json.Marshal(sessionEncx)
	require.NoError(t, err, "Failed to marshal session")

	ttl := 24 * time.Hour

	// Store session data
	err = redisClient.Set(ctx, sessionKey, sessionData, ttl).Err()
	require.NoError(t, err, "Failed to store session")

	// Store access token mapping
	err = redisClient.Set(ctx, accessTokenKey, sessionEncx.ID.String(), ttl).Err()
	require.NoError(t, err, "Failed to store access token mapping")

	// Store refresh token mapping
	err = redisClient.Set(ctx, refreshTokenKey, sessionEncx.ID.String(), ttl).Err()
	require.NoError(t, err, "Failed to store refresh token mapping")

	// Add session ID to user session index
	err = redisClient.SAdd(ctx, userSessionIndexKey, sessionEncx.ID.String()).Err()
	require.NoError(t, err, "Failed to add to user session index")

	return accessToken
}

// SetupAuthenticatedPartnerWithAllocation creates a complete partner setup:
// - Partner user with partner record
// - Room allocation for the specified room
// - Active session with access token
//
// Returns both the access token and user ID.
func SetupAuthenticatedPartnerWithAllocation(
	t *testing.T,
	ctx context.Context,
	email string,
	roomID uuid.UUID,
	pool *pgxpool.Pool,
	redisClient *redis.Client,
	crypto encx.CryptoService,
) (accessToken string, userID uuid.UUID) {
	t.Helper()

	// Create partner user
	userID = SetupPartnerUser(t, ctx, email, pool, crypto)

	// Create allocation for the room
	allocation := talloc.NewTestSharedAllocation(t, roomID, userID)
	talloc.InsertAllocation(t, ctx, allocation, pool, crypto)

	// Create session for authentication
	accessToken = CreateSessionForUser(t, ctx, userID, identity.Partner, redisClient, crypto)

	return accessToken, userID
}

// SetupAdminWithAllocation creates an admin user with room allocation.
// Returns both the access token and user ID.
func SetupAdminWithAllocation(
	t *testing.T,
	ctx context.Context,
	roomID uuid.UUID,
	pool *pgxpool.Pool,
	redisClient *redis.Client,
	crypto encx.CryptoService,
) (accessToken string, userID uuid.UUID) {
	t.Helper()

	// Create admin user directly in database
	userID = uuid.New()
	user := th.NewTestUser(t, "admin@test.com", "Admin", "User")
	user.ID = userID
	user.Role = identity.Administrator.String()
	userEncx, err := userDomain.ProcessUserEncx(ctx, crypto, user)
	require.NoError(t, err)
	err = th.InsertUserEncx(t, ctx, userEncx, pool)
	require.NoError(t, err)

	// Create allocation for the room
	allocation := talloc.NewTestSharedAllocation(t, roomID, userID)
	talloc.InsertAllocation(t, ctx, allocation, pool, crypto)

	// Create session for admin
	accessToken = CreateSessionForUser(t, ctx, userID, identity.Administrator, redisClient, crypto)

	return accessToken, userID
}
