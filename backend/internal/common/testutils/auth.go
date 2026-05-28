// Package testutils provides authentication testing utilities for all microservices.
//
// This package offers comprehensive helpers for setting up users with different roles,
// creating sessions, and testing authentication/authorization scenarios across the
// Leviosa platform microservices (catalog, settings, notification, etc.).
//
// Key Features:
// - Role-based user setup (Visitor, Standard, Premium, Guest, Partner, Administrator)
// - Session management with ENCX encryption
// - HTTP request helpers for middleware testing
// - Granular cleanup utilities
// - Support for expired/invalid session testing
//
// Basic Usage:
//
//	authCtx := &testutils.AuthTestContext{
//		Pool:   testPool,    // Your test database pool
//		Redis:  testClient,  // Your test Redis client
//		Crypto: crypto,      // Your ENCX crypto service
//	}
//
//	// Setup a user with specific role
//	accessToken := testutils.SetupAdminUser(t, ctx, authCtx)
//
//	// Create authenticated HTTP request
//	req := testutils.CreateAuthenticatedRequest("GET", "/api/settings", accessToken)
//
//	// Test your middleware/endpoint
//	resp := httptest.NewRecorder()
//	handler.ServeHTTP(resp, req)
//
//	// Cleanup after test
//	testutils.ClearAuthData(t, ctx, authCtx)
package testutils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
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
	State      string    `json:"-"`
	Email      string    `json:"-" encx:"encrypt,hash_basic"`
	FirstName  string    `json:"-" encx:"encrypt"`
	LastName   string    `json:"-" encx:"encrypt"`
	Password   string    `json:"-" encx:"hash_secure"`
	Telephone  string    `json:"-" encx:"encrypt,hash_basic"`
	Role       string    `json:"-" encx:"encrypt"`
	CreatedAt  time.Time `json:"-" encx:"encrypt"`
	LoggedInAt time.Time `json:"-" encx:"encrypt"`
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
	userEncx, err := ProcessUserEncx(ctx, authCtx.Crypto, user)
	require.NoError(t, err, "Failed to encrypt user struct")

	// Insert user into database
	insertUser(t, ctx, userEncx, authCtx.Pool)

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
	sessionEncx, err := session.ProcessSessionEncx(ctx, authCtx.Crypto, sess)
	require.NoError(t, err, "Failed to encrypt session struct")

	// Store session in Redis
	insertSession(t, ctx, sessionEncx, authCtx.Redis, 24*time.Hour)

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
	userEncx, err := ProcessUserEncx(ctx, authCtx.Crypto, user)
	require.NoError(t, err, "Failed to encrypt user struct")

	// Insert user into database
	insertUser(t, ctx, userEncx, authCtx.Pool)

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
	sessionEncx, err := session.ProcessSessionEncx(ctx, authCtx.Crypto, sess)
	require.NoError(t, err, "Failed to encrypt session struct")

	// Store session in Redis
	insertSession(t, ctx, sessionEncx, authCtx.Redis, 30*time.Minute)

	return accessToken
}

// SetupExpiredUserWithRole creates a user with expired session for the specified role
// Useful for testing session expiration and timeout scenarios
func SetupExpiredUserWithRole(t *testing.T, ctx context.Context, role identity.Role, authCtx *AuthTestContext) string {
	t.Helper()

	now := time.Now()
	userID := uuid.New()

	// Create test user
	user := &User{
		ID:         userID,
		State:      "active",
		Email:      fmt.Sprintf("expired_%s@leviosa.care", role.String()),
		FirstName:  role.String(),
		LastName:   role.String(),
		Password:   "bMPSrxQK#?rPO.[<",
		Telephone:  "0612345678",
		Role:       role.String(),
		CreatedAt:  now,
		LoggedInAt: now,
	}

	// Encrypt user data
	userEncx, err := ProcessUserEncx(ctx, authCtx.Crypto, user)
	require.NoError(t, err, "Failed to encrypt user struct")

	// Insert user into database
	insertUser(t, ctx, userEncx, authCtx.Pool)

	// Create expired session
	sessionID := uuid.New()
	accessToken, err := session.GenerateToken()
	require.NoError(t, err, "Failed to generate access token")

	refreshToken, err := session.GenerateToken()
	require.NoError(t, err, "Failed to generate refresh token")

	// Session expired 1 hour ago
	sess := &session.Session{
		ID:           sessionID,
		UserID:       userID,
		Role:         role,
		State:        session.SessionActive,
		CreatedAt:    now.Add(-2 * time.Hour),
		ExpiresAt:    now.Add(-1 * time.Hour), // Expired 1 hour ago
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// Encrypt session data
	sessionEncx, err := session.ProcessSessionEncx(ctx, authCtx.Crypto, sess)
	require.NoError(t, err, "Failed to encrypt session struct")

	// Store session in Redis with no TTL (it's already expired)
	insertSession(t, ctx, sessionEncx, authCtx.Redis, 0)

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

// SetupUserWithCustomData creates a user with custom data and active session
// Useful for testing specific user scenarios (custom emails, names, etc.)
func SetupUserWithCustomData(t *testing.T, ctx context.Context, role identity.Role, email, firstName, lastName, telephone string, authCtx *AuthTestContext) string {
	t.Helper()

	now := time.Now()
	userID := uuid.New()

	// Use provided data or defaults
	if email == "" {
		email = fmt.Sprintf("%s@leviosa.care", role.String())
	}
	if firstName == "" {
		firstName = role.String()
	}
	if lastName == "" {
		lastName = role.String()
	}
	if telephone == "" {
		telephone = "0612345678"
	}

	// Create test user with custom data
	user := &User{
		ID:         userID,
		State:      "active",
		Email:      email,
		FirstName:  firstName,
		LastName:   lastName,
		Password:   "bMPSrxQK#?rPO.[<",
		Telephone:  telephone,
		Role:       role.String(),
		CreatedAt:  now,
		LoggedInAt: now,
	}

	// Encrypt user data
	userEncx, err := ProcessUserEncx(ctx, authCtx.Crypto, user)
	require.NoError(t, err, "Failed to encrypt user struct")

	// Insert user into database
	insertUser(t, ctx, userEncx, authCtx.Pool)

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
	sessionEncx, err := session.ProcessSessionEncx(ctx, authCtx.Crypto, sess)
	require.NoError(t, err, "Failed to encrypt session struct")

	// Store session in Redis
	insertSession(t, ctx, sessionEncx, authCtx.Redis, 24*time.Hour)

	return accessToken
}

// CreateAuthHeader creates an Authorization header with Bearer token for HTTP requests
// Returns a header map that can be used with httptest.NewRequest
func CreateAuthHeader(accessToken string) http.Header {
	header := make(http.Header)
	header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	return header
}

// CreateAuthCookie creates an HTTP cookie with the access token for middleware testing
// Returns a cookie that can be added to HTTP requests using req.AddCookie()
func CreateAuthCookie(accessToken string) *http.Cookie {
	return &http.Cookie{
		Name:     cookies.AccessTokenCookieName,
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

// CreateRefreshCookie creates an HTTP cookie with the refresh token for refresh endpoint testing
// Returns a cookie that can be added to HTTP requests using req.AddCookie()
func CreateRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     cookies.RefreshTokenCookieName,
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

// CreateAuthenticatedRequest creates an HTTP request with both Authorization header and cookie
// Useful for testing different auth middleware approaches
func CreateAuthenticatedRequest(method, url string, accessToken string) *http.Request {
	req := httptest.NewRequest(method, url, nil)

	// Add Authorization header
	req.Header = CreateAuthHeader(accessToken)

	// Add cookie for cookie-based middleware
	req.AddCookie(CreateAuthCookie(accessToken))

	return req
}

// CreateAuthenticatedRequestWithBody creates an HTTP request with auth and request body
// Useful for testing POST/PUT endpoints with authentication
func CreateAuthenticatedRequestWithBody(method, url string, accessToken string, body *strings.Reader) *http.Request {
	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/json")

	// Add Authorization header
	req.Header = CreateAuthHeader(accessToken)

	// Add cookie for cookie-based middleware
	req.AddCookie(CreateAuthCookie(accessToken))

	return req
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

// ClearSessionsOnly clears only session-related test data, keeping users
// Useful when you want to test multiple session scenarios with the same users
func ClearSessionsOnly(t *testing.T, ctx context.Context, authCtx *AuthTestContext) {
	t.Helper()

	// Clear all session-related Redis keys only
	clearSessionsRedis(t, ctx, authCtx.Redis)
}

// ClearUsersOnly clears only user-related test data, keeping sessions
// Useful when you want to test session cleanup scenarios
func ClearUsersOnly(t *testing.T, ctx context.Context, authCtx *AuthTestContext) {
	t.Helper()

	// Clear users table only
	_, err := authCtx.Pool.Exec(ctx, "TRUNCATE TABLE auth.users RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Failed to clear users table")
}

// CountAuthUsers returns the number of users in the auth.users table for test verification
func CountAuthUsers(t *testing.T, ctx context.Context, pool *pgxpool.Pool) int {
	t.Helper()
	var count int
	query := `SELECT COUNT(*) FROM auth.users`
	err := pool.QueryRow(ctx, query).Scan(&count)
	require.NoError(t, err, "Failed to count auth users")
	return count
}

// CountActiveSessions returns the number of active sessions in Redis for test verification
func CountActiveSessions(t *testing.T, ctx context.Context, client *redis.Client) int {
	t.Helper()

	// Count session keys
	sessionKeys, err := client.Keys(ctx, fmt.Sprintf("%s*", session.SessionKeyPrefix)).Result()
	require.NoError(t, err, "Failed to get session keys")
	return len(sessionKeys)
}

// UserExists checks if a user exists by email hash for test verification
func UserExists(t *testing.T, ctx context.Context, email string, pool *pgxpool.Pool, crypto encx.CryptoService) bool {
	t.Helper()

	// Hash the email for lookup
	emailBytes, err := encx.SerializeValue(email)
	require.NoError(t, err, "Failed to serialize email")
	emailHash := crypto.HashBasic(ctx, emailBytes)

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM auth.users WHERE email_hash = $1)`
	err = pool.QueryRow(ctx, query, emailHash).Scan(&exists)
	require.NoError(t, err, "Failed to check user existence")
	return exists
}

// insertUser inserts a user into the auth.users table
func insertUser(t *testing.T, ctx context.Context, user *UserEncx, pool *pgxpool.Pool) {
	t.Helper()

	query := `
		INSERT INTO auth.users (
			id, state, email_hash, email_encrypted, password_hash_secure,
			first_name_encrypted, last_name_encrypted, telephone_hash, telephone_encrypted,
			role_encrypted, created_at_encrypted, logged_in_at_encrypted,
			dek_encrypted, key_version, profile_incomplete
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := pool.Exec(ctx, query,
		user.ID, user.State, user.EmailHash, user.EmailEncrypted, user.PasswordHashSecure,
		user.FirstNameEncrypted, user.LastNameEncrypted, user.TelephoneHash, user.TelephoneEncrypted,
		user.RoleEncrypted, user.CreatedAtEncrypted, user.LoggedInAtEncrypted,
		user.DEKEncrypted, user.KeyVersion, user.ProfileIncomplete)
	require.NoError(t, err, "Failed to insert test user")
}

// insertSession inserts a session into Redis with proper key formatting
func insertSession(t *testing.T, ctx context.Context, sess *session.SessionEncx, client *redis.Client, ttl time.Duration) {
	t.Helper()

	sessionKey := session.FormatSessionKey(sess.ID.String())
	accessTokenKey := session.FormatAccessTokenKey(sess.AccessTokenHash)
	refreshTokenKey := session.FormatRefreshTokenKey(sess.RefreshTokenHash)
	userSessionIndexKey := session.FormatUserSessionIndexKey(sess.UserIDHash)

	// Encode session to JSON
	sessionData, err := json.Marshal(sess)
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
