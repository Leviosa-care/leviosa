package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	sessionRepository "github.com/Leviosa-care/authuser/internal/adapters/redis/session"
	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/Leviosa-care/core/middleware/auth"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

// ClearSessionsRedis clears all session-related Redis keys for clean test state
func ClearSessionsRedis(t *testing.T, ctx context.Context, client *redis.Client) {
	t.Helper()

	// Clear all session keys
	// sessionKeys, err := client.Keys(ctx, "authuser:session:*").Result()
	sessionKeys, err := client.Keys(ctx, fmt.Sprintf("%s*", sessionRepository.SessionKeyPrefix)).Result()
	if err != nil {
		require.NoError(t, err, "Failed to get session keys")
	}
	if len(sessionKeys) > 0 {
		err = client.Del(ctx, sessionKeys...).Err()
		require.NoError(t, err, "Failed to delete session keys")
	}

	// Clear all access token keys (new dual-token system)
	// accessTokenKeys, err := client.Keys(ctx, "authuser:access:*").Result()
	accessTokenKeys, err := client.Keys(ctx, fmt.Sprintf("%s*", sessionRepository.AccessTokenKeyPrefix)).Result()
	if err != nil {
		require.NoError(t, err, "Failed to get access token keys")
	}
	if len(accessTokenKeys) > 0 {
		err = client.Del(ctx, accessTokenKeys...).Err()
		require.NoError(t, err, "Failed to delete access token keys")
	}

	// Clear all refresh token keys (new dual-token system)
	// refreshTokenKeys, err := client.Keys(ctx, "authuser:refresh:*").Result()
	refreshTokenKeys, err := client.Keys(ctx, fmt.Sprintf("%s*", sessionRepository.RefreshTokenKeyPrefix)).Result()
	if err != nil {
		require.NoError(t, err, "Failed to get refresh token keys")
	}
	if len(refreshTokenKeys) > 0 {
		err = client.Del(ctx, refreshTokenKeys...).Err()
		require.NoError(t, err, "Failed to delete refresh token keys")
	}
}

// NewTestSession creates a test session with reasonable defaults using real encryption
func NewTestSession(crypto encx.CryptoService) (*auth.Session, error) {
	now := time.Now()
	userID := uuid.New()
	sessionID := uuid.New()

	// Generate valid base64 tokens for testing
	accessToken, err := auth.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := auth.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	session := &auth.Session{
		ID:           sessionID,
		UserID:       userID,
		Role:         identity.Visitor,
		State:        auth.SessionActive,
		CreatedAt:    now,
		ExpiresAt:    now.Add(24 * time.Hour),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// Use crypto service to process the struct and populate encrypted/hashed fields
	err = crypto.ProcessStruct(context.Background(), session)
	if err != nil {
		return nil, fmt.Errorf("process session struct for encryption: %w", err)
	}

	return session, nil
}

// NewTestPendingSession creates a test session with pending state
func NewTestPendingSession(crypto encx.CryptoService) (*auth.Session, error) {
	session, err := NewTestSession(crypto)
	if err != nil {
		return nil, err
	}
	session.State = auth.SessionPending

	// Re-process the struct to update encrypted state
	err = crypto.ProcessStruct(context.Background(), session)
	if err != nil {
		return nil, fmt.Errorf("re-process session struct for state update: %w", err)
	}

	return session, nil
}

// NewTestSessionWithUserID creates a test session for a specific user ID
func NewTestSessionWithUserID(userID uuid.UUID, crypto encx.CryptoService) (*auth.Session, error) {
	now := time.Now()
	sessionID := uuid.New()

	// Generate valid base64 tokens for testing
	accessToken, err := auth.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := auth.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	session := &auth.Session{
		ID:           sessionID,
		UserID:       userID,
		Role:         identity.Visitor,
		State:        auth.SessionPending,
		CreatedAt:    now,
		ExpiresAt:    now.Add(24 * time.Hour),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// Process the struct to update encrypted user ID
	err = crypto.ProcessStruct(context.Background(), session)
	if err != nil {
		return nil, fmt.Errorf("re-process session struct for user ID update: %w", err)
	}

	return session, nil
}

// EncodeSession marshals a session to JSON bytes for Redis storage
func EncodeSession(t *testing.T, session *auth.Session) []byte {
	t.Helper()
	data, err := json.Marshal(session)
	require.NoError(t, err, "Failed to marshal session")
	return data
}

// DecodeSession unmarshals JSON bytes back to a session
func DecodeSession(t *testing.T, data []byte) *auth.Session {
	t.Helper()
	var session auth.Session
	err := json.Unmarshal(data, &session)
	require.NoError(t, err, "Failed to unmarshal session")
	return &session
}

// DecodeSessionWithDecryption unmarshals JSON bytes back to a session and decrypts it
func DecodeSessionWithDecryption(t *testing.T, data []byte, crypto encx.CryptoService) *auth.Session {
	t.Helper()
	var session auth.Session
	err := json.Unmarshal(data, &session)
	require.NoError(t, err, "Failed to unmarshal session")

	// Decrypt the session to populate plaintext fields
	err = crypto.DecryptStruct(context.Background(), &session)
	require.NoError(t, err, "Failed to decrypt session")

	return &session
}

// InsertSessionDirectly inserts a session directly into Redis (bypasses repository)
func InsertSessionDirectly(t *testing.T, ctx context.Context, client *redis.Client, session *auth.Session, ttl time.Duration) {
	t.Helper()

	sessionKey := sessionRepository.FormatSessionKey(session.ID.String())
	accessTokenKey := sessionRepository.FormatAccessTokenKey(session.AccessTokenHash)
	refreshTokenKey := sessionRepository.FormatRefreshTokenKey(session.RefreshTokenHash)

	// Store session data
	sessionData := EncodeSession(t, session)
	err := client.Set(ctx, sessionKey, sessionData, ttl).Err()
	require.NoError(t, err, "Failed to insert session directly")

	// Store access token mapping
	err = client.Set(ctx, accessTokenKey, session.ID.String(), ttl).Err()
	require.NoError(t, err, "Failed to insert access token mapping directly")

	// Store refresh token mapping
	err = client.Set(ctx, refreshTokenKey, session.ID.String(), ttl).Err()
	require.NoError(t, err, "Failed to insert refresh token mapping directly")
}

// CheckSessionExistsInRedis checks if a session exists by session ID
// DEPRECATED: Use raw Redis queries in tests for Test Independence Principle
// Example: exists, err := testClient.Exists(ctx, "authuser:session:"+sessionID.String()).Result()
func CheckSessionExistsInRedis(t *testing.T, ctx context.Context, client *redis.Client, sessionID uuid.UUID) bool {
	t.Helper()
	sessionKey := fmt.Sprintf("authuser:session:%s", sessionID.String())
	exists, err := client.Exists(ctx, sessionKey).Result()
	require.NoError(t, err, "Failed to check session existence")
	return exists > 0
}

// CheckTokenMappingExistsInRedis checks if a token mapping exists by token hash
// DEPRECATED: Use raw Redis queries in tests for Test Independence Principle
// Example: exists, err := testClient.Exists(ctx, "authuser:token:"+tokenHash).Result()
func CheckTokenMappingExistsInRedis(t *testing.T, ctx context.Context, client *redis.Client, tokenHash string) bool {
	t.Helper()
	tokenKey := fmt.Sprintf("authuser:token:%s", tokenHash)
	exists, err := client.Exists(ctx, tokenKey).Result()
	require.NoError(t, err, "Failed to check token mapping existence")
	return exists > 0
}

// GetSessionFromRedis retrieves and decodes a session from Redis by session ID
func GetSessionFromRedis(t *testing.T, ctx context.Context, client *redis.Client, sessionID uuid.UUID) *auth.Session {
	t.Helper()
	sessionKey := fmt.Sprintf("authuser:session:%s", sessionID.String())
	data, err := client.Get(ctx, sessionKey).Bytes()
	require.NoError(t, err, "Failed to get session from Redis")
	return DecodeSession(t, data)
}

// GetSessionIDFromTokenHash retrieves session ID from token hash mapping
// DEPRECATED: Use raw Redis queries in tests for Test Independence Principle
// Example: sessionID, err := testClient.Get(ctx, "authuser:token:"+tokenHash).Result()
func GetSessionIDFromTokenHash(t *testing.T, ctx context.Context, client *redis.Client, tokenHash string) string {
	t.Helper()
	tokenKey := fmt.Sprintf("authuser:token:%s", tokenHash)
	sessionID, err := client.Get(ctx, tokenKey).Result()
	require.NoError(t, err, "Failed to get session ID from token hash")
	return sessionID
}

// CountSessionKeys returns the number of session keys in Redis
// DEPRECATED: Use raw Redis queries in tests for Test Independence Principle
// Example: keys, err := testClient.Keys(ctx, "authuser:session:*").Result(); count := len(keys)
func CountSessionKeys(t *testing.T, ctx context.Context, client *redis.Client) int {
	t.Helper()
	keys, err := client.Keys(ctx, "authuser:session:*").Result()
	require.NoError(t, err, "Failed to count session keys")
	return len(keys)
}

// CountTokenKeys returns the number of token keys in Redis
// DEPRECATED: Use raw Redis queries in tests for Test Independence Principle
// Example: keys, err := testClient.Keys(ctx, "authuser:token:*").Result(); count := len(keys)
func CountTokenKeys(t *testing.T, ctx context.Context, client *redis.Client) int {
	t.Helper()
	keys, err := client.Keys(ctx, "authuser:token:*").Result()
	require.NoError(t, err, "Failed to count token keys")
	return len(keys)
}

// SessionTestData is a helper struct for multi-session tests
type SessionTestData struct {
	SessionID   uuid.UUID
	TokenHash   string
	Session     *auth.Session
	SessionData []byte
}

// CreateTestSessionWithCrypto is a convenience function that creates and returns a session
func CreateTestSessionWithCrypto(t *testing.T, crypto encx.CryptoService) *auth.Session {
	t.Helper()
	session, err := NewTestSession(crypto)
	require.NoError(t, err, "Failed to create test session with encryption")
	return session
}

// CreateTestPendingSessionWithCrypto is a convenience function that creates and returns a pending session
func CreateTestPendingSessionWithCrypto(t *testing.T, crypto encx.CryptoService) *auth.Session {
	t.Helper()
	session, err := NewTestPendingSession(crypto)
	require.NoError(t, err, "Failed to create test pending session with encryption")
	return session
}

// CreateTestSessionWithUserIDAndCrypto is a convenience function that creates and returns a session with specific user ID
func CreateTestSessionWithUserIDAndCrypto(t *testing.T, userID uuid.UUID, crypto encx.CryptoService) *auth.Session {
	t.Helper()
	session, err := NewTestSessionWithUserID(userID, crypto)
	require.NoError(t, err, "Failed to create test session with user ID and encryption")
	return session
}

// CreateTestPendingSessionWithUserIDAndCrypto is a convenience function that creates and returns a pending session with specific user ID
func CreateTestPendingSessionWithUserIDAndCrypto(t *testing.T, userID uuid.UUID, crypto encx.CryptoService) *auth.Session {
	now := time.Now()
	sessionID := uuid.New()

	// Generate valid base64 tokens for testing
	accessToken, err := auth.GenerateToken()
	require.NoError(t, err)

	// print("THE ACCESS TOKEN BEFORE ENCRYPTION IS :", accessToken)
	refreshToken, err := auth.GenerateToken()
	require.NoError(t, err)

	session := &auth.Session{
		ID:           sessionID,
		UserID:       userID,
		Role:         identity.Visitor,
		State:        auth.SessionPending,
		CreatedAt:    now,
		ExpiresAt:    now.Add(24 * time.Hour),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// Re-process the struct to update encrypted state
	err = crypto.ProcessStruct(context.Background(), session)
	require.NoError(t, err, "Failed to re-process session struct for state update")

	return session
}
