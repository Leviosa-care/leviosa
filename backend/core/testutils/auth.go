package testutils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

// User represents the minimal user structure needed for auth tests
// This avoids coupling to specific domain models in authuser module
type User struct {
	ID         uuid.UUID `json:"-"`
	State      string    `json:"-" encx:"encrypt"`
	Email      string    `json:"-" encx:"encrypt,hash_basic"`
	FirstName  string    `json:"-" encx:"encrypt"`
	LastName   string    `json:"-" encx:"encrypt"`
	Password   string    `json:"-" encx:"hash_bcrypt"`
	Telephone  string    `json:"-" encx:"encrypt,hash_basic"`
	Role       string    `json:"-" encx:"encrypt"`
	CreatedAt  time.Time `json:"-" encx:"encrypt"`
	LoggedInAt time.Time `json:"-" encx:"encrypt"`

	// Encrypted fields
	StateEncrypted     []byte `json:"state_encrypted"`
	EmailEncrypted     []byte `json:"email_encrypted"`
	EmailHash          string `json:"email_hash"`
	FirstNameEncrypted []byte `json:"first_name_encrypted"`
	LastNameEncrypted  []byte `json:"last_name_encrypted"`
	PasswordHash       string `json:"password_hash"`
	TelephoneEncrypted []byte `json:"telephone_encrypted"`
	TelephoneHash      string `json:"telephone_hash"`
	RoleEncrypted      []byte `json:"role_encrypted"`
	CreatedAtEncrypted []byte `json:"created_at_encrypted"`
	LoggedInAtEncrypted []byte `json:"logged_in_at_encrypted"`
	DEKEncrypted       []byte `json:"dek_encrypted"`
	KeyVersion         int    `json:"key_version"`
}

// AuthTestContext holds the necessary dependencies for auth test utilities
type AuthTestContext struct {
	Pool   *pgxpool.Pool
	Redis  *redis.Client
	Crypto encx.CryptoService
}

// SetupUserWithRole creates a user and active session for the specified role
// Returns the access token that can be used in Authorization headers
func SetupUserWithRole(t *testing.T, ctx context.Context, role identity.Role, authCtx *AuthTestContext) string {
	t.Helper()
	
	now := time.Now()
	userID := uuid.New()
	
	// Create test user
	user := &User{
		ID:         userID,
		State:      "active",
		Email:      fmt.Sprintf("%s@leviosa.care", role.String()),
		FirstName:  role.String(),
		LastName:   role.String(),
		Password:   "bMPSrxQK#?rPO.[<",
		Telephone:  "0612345678",
		Role:       role.String(),
		CreatedAt:  now,
		LoggedInAt: now,
	}
	
	// Encrypt user data
	err := authCtx.Crypto.ProcessStruct(ctx, user)
	require.NoError(t, err, "Failed to encrypt user struct")
	
	// Insert user into database
	insertUser(t, ctx, user, authCtx.Pool)
	
	// Create active session
	sessionID := uuid.New()
	accessToken, err := session.GenerateToken()
	require.NoError(t, err, "Failed to generate access token")
	
	refreshToken, err := session.GenerateToken()
	require.NoError(t, err, "Failed to generate refresh token")
	
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
	err = authCtx.Crypto.ProcessStruct(ctx, sess)
	require.NoError(t, err, "Failed to encrypt session struct")
	
	// Store session in Redis
	insertSession(t, ctx, sess, authCtx.Redis, 24*time.Hour)
	
	return accessToken
}

// SetupVisitorUser creates a visitor user with active session
func SetupVisitorUser(t *testing.T, ctx context.Context, authCtx *AuthTestContext) string {
	t.Helper()
	return SetupUserWithRole(t, ctx, identity.Visitor, authCtx)
}

// SetupStandardUser creates a standard user with active session
func SetupStandardUser(t *testing.T, ctx context.Context, authCtx *AuthTestContext) string {
	t.Helper()
	return SetupUserWithRole(t, ctx, identity.Standard, authCtx)
}

// SetupPremiumUser creates a premium user with active session
func SetupPremiumUser(t *testing.T, ctx context.Context, authCtx *AuthTestContext) string {
	t.Helper()
	return SetupUserWithRole(t, ctx, identity.Premium, authCtx)
}

// SetupGuestUser creates a guest user with active session
func SetupGuestUser(t *testing.T, ctx context.Context, authCtx *AuthTestContext) string {
	t.Helper()
	return SetupUserWithRole(t, ctx, identity.Guest, authCtx)
}

// SetupPartnerUser creates a partner user with active session
func SetupPartnerUser(t *testing.T, ctx context.Context, authCtx *AuthTestContext) string {
	t.Helper()
	return SetupUserWithRole(t, ctx, identity.Partner, authCtx)
}

// SetupAdminUser creates an administrator user with active session
func SetupAdminUser(t *testing.T, ctx context.Context, authCtx *AuthTestContext) string {
	t.Helper()
	return SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
}

// SetupPendingUserWithRole creates a user with pending session for the specified role
// Useful for testing registration flows
func SetupPendingUserWithRole(t *testing.T, ctx context.Context, role identity.Role, authCtx *AuthTestContext) string {
	t.Helper()
	
	now := time.Now()
	userID := uuid.New()
	
	// Create test user in unverified state
	user := &User{
		ID:         userID,
		State:      "unverified",
		Email:      fmt.Sprintf("pending_%s@leviosa.care", role.String()),
		FirstName:  role.String(),
		LastName:   role.String(),
		Password:   "bMPSrxQK#?rPO.[<",
		Telephone:  "0612345678",
		Role:       role.String(),
		CreatedAt:  now,
		LoggedInAt: now,
	}
	
	// Encrypt user data
	err := authCtx.Crypto.ProcessStruct(ctx, user)
	require.NoError(t, err, "Failed to encrypt user struct")
	
	// Insert user into database
	insertUser(t, ctx, user, authCtx.Pool)
	
	// Create pending session
	sessionID := uuid.New()
	accessToken, err := session.GenerateToken()
	require.NoError(t, err, "Failed to generate access token")
	
	refreshToken, err := session.GenerateToken()
	require.NoError(t, err, "Failed to generate refresh token")
	
	sess := &session.Session{
		ID:           sessionID,
		UserID:       userID,
		Role:         role,
		State:        session.SessionPending,
		CreatedAt:    now,
		ExpiresAt:    now.Add(30 * time.Minute), // Shorter duration for pending sessions
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	
	// Encrypt session data
	err = authCtx.Crypto.ProcessStruct(ctx, sess)
	require.NoError(t, err, "Failed to encrypt session struct")
	
	// Store session in Redis
	insertSession(t, ctx, sess, authCtx.Redis, 30*time.Minute)
	
	return accessToken
}

// SetupMultipleUsers creates users for multiple roles and returns their access tokens
// Useful for testing role-based authorization across different privilege levels
func SetupMultipleUsers(t *testing.T, ctx context.Context, roles []identity.Role, authCtx *AuthTestContext) map[identity.Role]string {
	t.Helper()
	
	tokens := make(map[identity.Role]string)
	for _, role := range roles {
		tokens[role] = SetupUserWithRole(t, ctx, role, authCtx)
	}
	
	return tokens
}

// ClearAuthData cleans up all auth-related test data (users and sessions)
func ClearAuthData(t *testing.T, ctx context.Context, authCtx *AuthTestContext) {
	t.Helper()
	
	// Clear users table
	_, err := authCtx.Pool.Exec(ctx, "TRUNCATE TABLE auth.users RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Failed to clear users table")
	
	// Clear all session-related Redis keys
	clearSessionsRedis(t, ctx, authCtx.Redis)
}

// insertUser inserts a user into the auth.users table
func insertUser(t *testing.T, ctx context.Context, user *User, pool *pgxpool.Pool) {
	t.Helper()
	
	query := `
		INSERT INTO auth.users (
			id, state, email_hash, email_encrypted, password_hash,
			first_name_encrypted, last_name_encrypted, telephone_hash, telephone_encrypted,
			role_encrypted, created_at_encrypted, logged_in_at_encrypted, 
			dek_encrypted, key_version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	
	_, err := pool.Exec(ctx, query,
		user.ID, user.StateEncrypted, user.EmailHash, user.EmailEncrypted, user.PasswordHash,
		user.FirstNameEncrypted, user.LastNameEncrypted, user.TelephoneHash, user.TelephoneEncrypted,
		user.RoleEncrypted, user.CreatedAtEncrypted, user.LoggedInAtEncrypted,
		user.DEKEncrypted, user.KeyVersion)
	require.NoError(t, err, "Failed to insert test user")
}

// insertSession inserts a session into Redis with proper key formatting
func insertSession(t *testing.T, ctx context.Context, sess *session.Session, client *redis.Client, ttl time.Duration) {
	t.Helper()
	
	sessionKey := session.FormatSessionKey(sess.ID.String())
	accessTokenKey := session.FormatAccessTokenKey(sess.AccessTokenHash)
	refreshTokenKey := session.FormatRefreshTokenKey(sess.RefreshTokenHash)
	userSessionIndexKey := session.FormatUserSessionIndexKey(sess.UserIDHash)
	
	// Encode session to JSON
	sessionData, err := sess.MarshalJSON()
	require.NoError(t, err, "Failed to marshal session")
	
	// Store session data
	err = client.Set(ctx, sessionKey, sessionData, ttl).Err()
	require.NoError(t, err, "Failed to store session")
	
	// Store access token mapping
	err = client.Set(ctx, accessTokenKey, sess.ID.String(), ttl).Err()
	require.NoError(t, err, "Failed to store access token mapping")
	
	// Store refresh token mapping
	err = client.Set(ctx, refreshTokenKey, sess.ID.String(), ttl).Err()
	require.NoError(t, err, "Failed to store refresh token mapping")
	
	// Add session ID to user session index
	err = client.SAdd(ctx, userSessionIndexKey, sess.ID.String()).Err()
	require.NoError(t, err, "Failed to add to user session index")
}

// clearSessionsRedis clears all session-related Redis keys
func clearSessionsRedis(t *testing.T, ctx context.Context, client *redis.Client) {
	t.Helper()
	
	// Clear session keys
	sessionKeys, err := client.Keys(ctx, fmt.Sprintf("%s*", session.SessionKeyPrefix)).Result()
	require.NoError(t, err, "Failed to get session keys")
	if len(sessionKeys) > 0 {
		err = client.Del(ctx, sessionKeys...).Err()
		require.NoError(t, err, "Failed to delete session keys")
	}
	
	// Clear access token keys
	accessTokenKeys, err := client.Keys(ctx, fmt.Sprintf("%s*", session.AccessTokenKeyPrefix)).Result()
	require.NoError(t, err, "Failed to get access token keys")
	if len(accessTokenKeys) > 0 {
		err = client.Del(ctx, accessTokenKeys...).Err()
		require.NoError(t, err, "Failed to delete access token keys")
	}
	
	// Clear refresh token keys
	refreshTokenKeys, err := client.Keys(ctx, fmt.Sprintf("%s*", session.RefreshTokenKeyPrefix)).Result()
	require.NoError(t, err, "Failed to get refresh token keys")
	if len(refreshTokenKeys) > 0 {
		err = client.Del(ctx, refreshTokenKeys...).Err()
		require.NoError(t, err, "Failed to delete refresh token keys")
	}
	
	// Clear user session index keys
	userSessionKeys, err := client.Keys(ctx, "authuser:user_sessions:*").Result()
	require.NoError(t, err, "Failed to get user session keys")
	if len(userSessionKeys) > 0 {
		err = client.Del(ctx, userSessionKeys...).Err()
		require.NoError(t, err, "Failed to delete user session keys")
	}
}