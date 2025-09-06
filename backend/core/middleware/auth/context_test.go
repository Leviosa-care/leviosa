package auth

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSessionFromContext(t *testing.T) {
	tests := []struct {
		name           string
		contextSetup   func() context.Context
		expectedExists bool
		expectedNil    bool
	}{
		{
			name: "valid session in context",
			contextSetup: func() context.Context {
				sessionStruct := &session.Session{
					ID:               uuid.New(),
					UserID:           uuid.New(),
					Role:             identity.Partner,
					State:            session.SessionActive,
					CreatedAt:        time.Now(),
					ExpiresAt:        time.Now().Add(time.Hour),
					AccessTokenHash:  "test_access_hash",
					RefreshTokenHash: "test_refresh_hash",
				}
				return context.WithValue(context.Background(), session.GetSessionContextKey(), sessionStruct)
			},
			expectedExists: true,
			expectedNil:    false,
		},
		{
			name: "no session in context",
			contextSetup: func() context.Context {
				return context.Background()
			},
			expectedExists: false,
			expectedNil:    true,
		},
		{
			name: "wrong type in context",
			contextSetup: func() context.Context {
				return context.WithValue(context.Background(), session.GetSessionContextKey(), "not a session")
			},
			expectedExists: false,
			expectedNil:    true,
		},
		{
			name: "nil session in context",
			contextSetup: func() context.Context {
				return context.WithValue(context.Background(), session.GetSessionContextKey(), (*session.Session)(nil))
			},
			expectedExists: false,
			expectedNil:    true,
		},
		{
			name: "different context key",
			contextSetup: func() context.Context {
				sessionStruct := &session.Session{
					ID:               uuid.New(),
					UserID:           uuid.New(),
					Role:             identity.Visitor,
					State:            session.SessionActive,
					CreatedAt:        time.Now(),
					ExpiresAt:        time.Now().Add(time.Hour),
					AccessTokenHash:  "test_access_hash",
					RefreshTokenHash: "test_refresh_hash",
				}
				// Use wrong key type
				type wrongKey struct{}
				return context.WithValue(context.Background(), wrongKey{}, sessionStruct)
			},
			expectedExists: false,
			expectedNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.contextSetup()

			session, exists := session.SessionFromContext(ctx)

			assert.Equal(t, tt.expectedExists, exists, "unexpected existence result")

			if tt.expectedNil {
				assert.Nil(t, session, "session should be nil")
			} else {
				assert.NotNil(t, session, "session should not be nil")
			}
		})
	}
}

func TestSessionFromContext_ValidSession(t *testing.T) {
	// Test that we get back the exact same session we put in
	originalSession := &session.Session{
		ID:               uuid.New(),
		UserID:           uuid.New(),
		Role:             identity.Administrator,
		State:            session.SessionActive,
		CreatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
		AccessTokenHash:  "specific_access_hash",
		RefreshTokenHash: "specific_refresh_hash",
	}

	ctx := context.WithValue(context.Background(), session.GetSessionContextKey(), originalSession)

	retrievedSession, exists := session.SessionFromContext(ctx)

	assert.True(t, exists, "session should exist in context")
	assert.NotNil(t, retrievedSession, "retrieved session should not be nil")

	// Verify it's the exact same session object
	assert.Same(t, originalSession, retrievedSession, "should return the exact same session object")

	// Verify session contents
	assert.Equal(t, originalSession.ID, retrievedSession.ID)
	assert.Equal(t, originalSession.UserID, retrievedSession.UserID)
	assert.Equal(t, originalSession.Role, retrievedSession.Role)
	assert.Equal(t, originalSession.State, retrievedSession.State)
	assert.Equal(t, originalSession.AccessTokenHash, retrievedSession.AccessTokenHash)
	assert.Equal(t, originalSession.RefreshTokenHash, retrievedSession.RefreshTokenHash)
}

func TestSessionInfoFromContext(t *testing.T) {
	// Test that SessionInfoFromContext works correctly
	sessionInfo := &session.SessionInfo{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Role:   identity.Partner,
		State:  session.SessionActive,
	}

	ctx := context.WithValue(context.Background(), session.GetSessionContextKey(), sessionInfo)

	// Test SessionInfoFromContext
	retrievedInfo, ok := session.SessionInfoFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, sessionInfo.ID, retrievedInfo.ID)
	assert.Equal(t, sessionInfo.UserID, retrievedInfo.UserID)
	assert.Equal(t, sessionInfo.Role, retrievedInfo.Role)
	assert.Equal(t, sessionInfo.State, retrievedInfo.State)

	// Test backward compatibility with SessionFromContext
	retrievedSession, ok := session.SessionFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, sessionInfo.ID, retrievedSession.ID)
	assert.Equal(t, sessionInfo.UserID, retrievedSession.UserID)
	assert.Equal(t, sessionInfo.Role, retrievedSession.Role)
	assert.Equal(t, sessionInfo.State, retrievedSession.State)
}

func TestSessionContextKey_Uniqueness(t *testing.T) {
	// Test that sessionContextKey is unique and doesn't conflict with other keys
	sessionStruct := &session.Session{
		ID:               uuid.New(),
		UserID:           uuid.New(),
		Role:             identity.Partner,
		State:            session.SessionActive,
		CreatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(time.Hour),
		AccessTokenHash:  "test_access_hash",
		RefreshTokenHash: "test_refresh_hash",
	}

	// Create context with multiple values using different key types
	type otherKey struct{}

	ctx := context.Background()
	ctx = context.WithValue(ctx, session.GetSessionContextKey(), sessionStruct)
	ctx = context.WithValue(ctx, otherKey{}, "other value")
	ctx = context.WithValue(ctx, "string_key", "string value")

	// SessionFromContext should only retrieve the session, not other values
	retrievedSession, exists := session.SessionFromContext(ctx)

	assert.True(t, exists, "session should exist")
	assert.Same(t, sessionStruct, retrievedSession, "should retrieve correct session")

	// Verify other values are still there but don't interfere
	otherValue := ctx.Value(otherKey{})
	assert.Equal(t, "other value", otherValue)

	stringValue := ctx.Value("string_key")
	assert.Equal(t, "string value", stringValue)
}

func TestSessionContextKey_ZeroValue(t *testing.T) {
	// Test that zero value of sessionContextKey works correctly
	key1 := session.GetSessionContextKey()
	key2 := session.GetSessionContextKey()

	// Both should be equal (same zero value)
	assert.Equal(t, key1, key2, "sessionContextKey zero values should be equal")

	sessionStruct := &session.Session{
		ID:               uuid.New(),
		UserID:           uuid.New(),
		Role:             identity.Visitor,
		State:            session.SessionActive,
		CreatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(time.Hour),
		AccessTokenHash:  "test_access_hash",
		RefreshTokenHash: "test_refresh_hash",
	}

	// Should work with both key instances
	ctx1 := context.WithValue(context.Background(), key1, sessionStruct)
	ctx2 := context.WithValue(context.Background(), key2, sessionStruct)

	session1, exists1 := session.SessionFromContext(ctx1)
	session2, exists2 := session.SessionFromContext(ctx2)

	assert.True(t, exists1, "session should exist in ctx1")
	assert.True(t, exists2, "session should exist in ctx2")
	assert.Same(t, sessionStruct, session1, "should retrieve correct session from ctx1")
	assert.Same(t, sessionStruct, session2, "should retrieve correct session from ctx2")
}

func TestSessionFromContext_ConcurrentAccess(t *testing.T) {
	// Test that SessionFromContext is safe for concurrent access
	sessionStruct := &session.Session{
		ID:               uuid.New(),
		UserID:           uuid.New(),
		Role:             identity.Administrator,
		State:            session.SessionActive,
		CreatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(time.Hour),
		AccessTokenHash:  "concurrent_access_hash",
		RefreshTokenHash: "concurrent_refresh_hash",
	}

	ctx := context.WithValue(context.Background(), session.GetSessionContextKey(), sessionStruct)

	// Run multiple goroutines accessing the same context
	results := make(chan bool, 10)

	for range 10 {
		go func() {
			retrievedSession, exists := session.SessionFromContext(ctx)
			results <- exists && retrievedSession != nil && retrievedSession.AccessTokenHash == "concurrent_access_hash"
		}()
	}

	// Verify all goroutines got the correct result
	for range 10 {
		result := <-results
		assert.True(t, result, "concurrent access should work correctly")
	}
}
