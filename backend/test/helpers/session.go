package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

// ClearSessionsRedis clears all session-related Redis keys for clean test state
func ClearSessionsRedis(t *testing.T, ctx context.Context, client *redis.Client) {
	t.Helper()

	// Clear all session keys
	sessionKeys, err := client.Keys(ctx, fmt.Sprintf("%s*", session.SessionKeyPrefix)).Result()
	if err != nil {
		require.NoError(t, err, "Failed to get session keys")
	}
	if len(sessionKeys) > 0 {
		err = client.Del(ctx, sessionKeys...).Err()
		require.NoError(t, err, "Failed to delete session keys")
	}

	// Clear all access token keys (new dual-token system)
	accessTokenKeys, err := client.Keys(ctx, fmt.Sprintf("%s*", session.AccessTokenKeyPrefix)).Result()
	if err != nil {
		require.NoError(t, err, "Failed to get access token keys")
	}
	if len(accessTokenKeys) > 0 {
		err = client.Del(ctx, accessTokenKeys...).Err()
		require.NoError(t, err, "Failed to delete access token keys")
	}

	// Clear all refresh token keys (new dual-token system)
	// refreshTokenKeys, err := client.Keys(ctx, "authuser:refresh:*").Result()
	refreshTokenKeys, err := client.Keys(ctx, fmt.Sprintf("%s*", session.RefreshTokenKeyPrefix)).Result()
	if err != nil {
		require.NoError(t, err, "Failed to get refresh token keys")
	}
	if len(refreshTokenKeys) > 0 {
		err = client.Del(ctx, refreshTokenKeys...).Err()
		require.NoError(t, err, "Failed to delete refresh token keys")
	}
}

// NewTestSession creates a test session with reasonable defaults using real encryption
// Returns original session for test usage
func NewTestSession(t *testing.T, crypto encx.CryptoService) (*session.Session, error) {
	t.Helper()

	now := time.Now()
	userID := uuid.New()
	sessionID := uuid.New()

	// Generate valid base64 tokens for testing
	accessToken, err := session.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := session.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	s := &session.Session{
		ID:           sessionID,
		UserID:       userID,
		Role:         identity.Visitor,
		State:        session.SessionActive,
		CreatedAt:    now,
		ExpiresAt:    now.Add(24 * time.Hour),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return s, nil
}

// NewTestSessionEncx creates a test session encx with mock defaults that does not use encryption.
// Returns both the original session and the processed SessionEncx for test usage
func NewTestSessionEncx(t *testing.T) *session.SessionEncx {
	t.Helper()

	s := &session.SessionEncx{
		ID:                 uuid.New(),
		UserIDEncrypted:    []byte("user_id_encrypted"),
		UserIDHash:         "user_id_hash_basic",
		RoleEncrypted:      []byte("role_encrypted"),
		StateEncrypted:     []byte("state_active_encrypted"),
		CreatedAtEncrypted: []byte("created_at_encrypted"),
		ExpiresAtEncrypted: []byte("expires_at_encrypted"),
		AccessTokenHash:    "access_token_hash_basic",
		RefreshTokenHash:   "refresh_token_hash_basic",
	}

	return s
}

// // NewTestSession creates a test session with reasonable defaults using real encryption
// // Returns both the original session and the processed SessionEncx for test usage
// func NewTestSession(crypto encx.CryptoService) (*session.Session, *session.SessionEncx, error) {
// 	now := time.Now()
// 	userID := uuid.New()
// 	sessionID := uuid.New()
//
// 	// Generate valid base64 tokens for testing
// 	accessToken, err := session.GenerateToken()
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("generate access token: %w", err)
// 	}
//
// 	refreshToken, err := session.GenerateToken()
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("generate refresh token: %w", err)
// 	}
//
// 	sess := &session.Session{
// 		ID:           sessionID,
// 		UserID:       userID,
// 		Role:         identity.Visitor,
// 		State:        session.SessionActive,
// 		CreatedAt:    now,
// 		ExpiresAt:    now.Add(24 * time.Hour),
// 		AccessToken:  accessToken,
// 		RefreshToken: refreshToken,
// 	}
//
// 	// Process the session to get encrypted data and hashes
// 	sessionEncx, err := session.ProcessSessionEncx(context.Background(), crypto, sess)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("process session for encryption: %w", err)
// 	}
//
// 	return sess, sessionEncx, nil
// }

// // NewTestSessionWithUserID creates a test session for a specific user ID
// func NewTestSessionWithUserID(userID uuid.UUID, crypto encx.CryptoService) (*session.Session, *session.SessionEncx, error) {
// 	now := time.Now()
// 	sessionID := uuid.New()
//
// 	// Generate valid base64 tokens for testing
// 	accessToken, err := session.GenerateToken()
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("generate access token: %w", err)
// 	}
//
// 	refreshToken, err := session.GenerateToken()
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("generate refresh token: %w", err)
// 	}
//
// 	sess := &session.Session{
// 		ID:           sessionID,
// 		UserID:       userID,
// 		Role:         identity.Visitor,
// 		State:        session.SessionPending,
// 		CreatedAt:    now,
// 		ExpiresAt:    now.Add(24 * time.Hour),
// 		AccessToken:  accessToken,
// 		RefreshToken: refreshToken,
// 	}
//
// 	// Process the session to get encrypted data and hashes
// 	sessionEncx, err := session.ProcessSessionEncx(context.Background(), crypto, sess)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("process session struct for user ID update: %w", err)
// 	}
//
// 	return sess, sessionEncx, nil
// }

// EncodeSession marshals a session to JSON bytes for Redis storage
func EncodeSession(t *testing.T, session *session.SessionEncx) []byte {
	t.Helper()
	data, err := json.Marshal(session)
	require.NoError(t, err, "Failed to marshal session")
	return data
}

// DecodeSessionEncx unmarshals JSON bytes back to a session
func DecodeSessionEncx(t *testing.T, data []byte) *session.SessionEncx {
	t.Helper()
	var session session.SessionEncx
	err := json.Unmarshal(data, &session)
	require.NoError(t, err, "Failed to unmarshal session")
	return &session
}

// DecodeSessionWithDecryption unmarshals JSON bytes back to a session and decrypts it using the new approach
func DecodeSessionWithDecryption(t *testing.T, data []byte, crypto encx.CryptoService) *session.Session {
	t.Helper()

	// First unmarshal as SessionEncx (since that's how we store it now)
	var sessionEncx session.SessionEncx
	err := json.Unmarshal(data, &sessionEncx)
	require.NoError(t, err, "Failed to unmarshal SessionEncx")

	// Decrypt the SessionEncx to get the original session
	session, err := session.DecryptSessionEncx(context.Background(), crypto, &sessionEncx)
	require.NoError(t, err, "Failed to decrypt session")

	return session
}

// InsertSessionEncx inserts a session directly into Redis using the new Encx approach
func InsertSessionEncx(t *testing.T, ctx context.Context, client *redis.Client, sessionEncx *session.SessionEncx, ttl time.Duration) {
	t.Helper()

	sessionKey := session.FormatSessionKey(sessionEncx.ID.String())
	accessTokenKey := session.FormatAccessTokenKey(sessionEncx.AccessTokenHash)
	refreshTokenKey := session.FormatRefreshTokenKey(sessionEncx.RefreshTokenHash)
	userSessionIndexKey := session.FormatUserSessionIndexKey(sessionEncx.UserIDHash)

	// Store session data as SessionEncx (for Redis storage)
	sessionData, err := json.Marshal(sessionEncx)
	require.NoError(t, err, "Failed to marshal SessionEncx for Redis storage")

	err = client.Set(ctx, sessionKey, sessionData, ttl).Err()
	require.NoError(t, err, "Failed to insert session directly")

	// Store access token mapping
	err = client.Set(ctx, accessTokenKey, sessionEncx.ID.String(), ttl).Err()
	require.NoError(t, err, "Failed to insert access token mapping directly")

	// Store refresh token mapping
	err = client.Set(ctx, refreshTokenKey, sessionEncx.ID.String(), ttl).Err()
	require.NoError(t, err, "Failed to insert refresh token mapping directly")

	// Add session ID to user session index
	err = client.SAdd(ctx, userSessionIndexKey, sessionEncx.ID.String()).Err()
	// Rollback all previous operations
	if err != nil {
		delErr := client.Del(ctx, sessionKey).Err()
		require.NoError(t, delErr)
		delErr = client.Del(ctx, accessTokenKey).Err()
		require.NoError(t, delErr)
		delErr = client.Del(ctx, refreshTokenKey).Err()
		require.NoError(t, delErr)
	}
	require.NoError(t, err, "Failed to add user session values")
}

// SessionTestData is a helper struct for multi-session tests
type SessionTestData struct {
	SessionID   uuid.UUID
	TokenHash   string
	Session     *session.SessionEncx
	SessionData []byte
}

// // TestSession is a helper struct that contains both the session and its encrypted/hash data for testing
// type TestSession struct {
// 	Session     *session.Session
// 	SessionEncx *session.SessionEncx
// }

// // CreateTestSessionWithCrypto is a convenience function that creates and returns a session
// func CreateTestSessionWithCrypto(t *testing.T, crypto encx.CryptoService) *session.Session {
// 	t.Helper()
// 	session, _, err := NewTestSession(crypto)
// 	require.NoError(t, err, "Failed to create test session with encryption")
// 	return session
// }

// // CreateTestPendingSessionWithCrypto is a convenience function that creates and returns a pending session
// func CreateTestPendingSessionWithCrypto(t *testing.T, crypto encx.CryptoService) *session.Session {
// 	t.Helper()
// 	session, _, err := NewTestPendingSession(crypto)
// 	require.NoError(t, err, "Failed to create test pending session with encryption")
// 	return session
// }

// // CreateTestSessionWithUserIDAndCrypto is a convenience function that creates and returns a session with specific user ID
// func CreateTestSessionWithUserIDAndCrypto(t *testing.T, userID uuid.UUID, crypto encx.CryptoService) *session.Session {
// 	t.Helper()
// 	session, _, err := NewTestSessionWithUserID(userID, crypto)
// 	require.NoError(t, err, "Failed to create test session with user ID and encryption")
// 	return session
// }

// // CreateTestPendingSessionWithUserIDAndCrypto is a convenience function that creates and returns a pending session with specific user ID
// func CreateTestPendingSessionWithUserIDAndCrypto(t *testing.T, userID uuid.UUID, crypto encx.CryptoService) *session.Session {
// 	t.Helper()
//
// 	// Create a session with the specific user ID and pending state
// 	sess, _, err := NewTestSessionWithUserID(userID, crypto)
// 	require.NoError(t, err, "Failed to create test session with user ID and encryption")
//
// 	sess.State = session.SessionPending
//
// 	// Re-process to update encrypted state and hashes for the pending state
// 	_, err = session.ProcessSessionEncx(context.Background(), crypto, sess)
// 	require.NoError(t, err, "Failed to re-process session struct for state update")
//
// 	return sess
// }

// CreateSessionWithEncryption creates and stores a session in Redis, returns access token
func CreateSessionWithEncryption(t *testing.T, ctx context.Context, sessionInfo *session.SessionInfo, client *redis.Client, crypto encx.CryptoService) string {
	t.Helper()

	// Generate tokens
	accessToken, err := session.GenerateToken()
	require.NoError(t, err)

	refreshToken, err := session.GenerateToken()
	require.NoError(t, err)

	// Create session object
	sess := &session.Session{
		ID:           sessionInfo.ID,
		UserID:       sessionInfo.UserID,
		Role:         sessionInfo.Role,
		State:        sessionInfo.State,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// Process encryption using the new generated function
	sessionEncx, err := session.ProcessSessionEncx(ctx, crypto, sess)
	require.NoError(t, err)

	// Store in Redis using the new approach
	ttl := 24 * time.Hour
	InsertSessionEncx(t, ctx, client, sessionEncx, ttl)

	return accessToken
}

// GetSessionByID retrieves a session by ID, returns nil if not found
func GetSessionByID(t *testing.T, ctx context.Context, sessionID uuid.UUID, client *redis.Client) *session.SessionEncx {
	t.Helper()

	sessionKey := session.FormatSessionKey(sessionID.String())
	data, err := client.Get(ctx, sessionKey).Bytes()
	if err != nil {
		// Return nil if session not found
		return nil
	}

	return DecodeSessionEncx(t, data)
}

// FindSessionByAccessTokenHash retrieves a session by its access token hash, returns nil if not found
func FindSessionByAccessTokenHash(t *testing.T, ctx context.Context, accessTokenHash string, client *redis.Client) ([]byte, error) {
	t.Helper()

	sessionKey := session.FormatAccessTokenKey(accessTokenHash)
	data, err := client.Get(ctx, sessionKey).Bytes()

	return data, err
}

// FindSessionByRefreshTokenHash retrieves a session by its access token hash, returns nil if not found
func FindSessionByRefreshTokenHash(t *testing.T, ctx context.Context, accessTokenHash string, client *redis.Client) ([]byte, error) {
	t.Helper()

	sessionKey := session.FormatRefreshTokenKey(accessTokenHash)
	data, err := client.Get(ctx, sessionKey).Bytes()

	return data, err
}

// ToCookieString generates a cookie string for the session access token
func ToCookieString(sess *session.Session) string {
	return fmt.Sprintf("access_token=%s", sess.AccessToken)
}
