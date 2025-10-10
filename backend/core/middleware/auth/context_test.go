package auth

import (
	"context"
	"testing"

	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSessionInfoFromContext(t *testing.T) {
	tests := []struct {
		name           string
		contextSetup   func() context.Context
		expectedExists bool
		expectedNil    bool
	}{
		{
			name: "valid session info in context",
			contextSetup: func() context.Context {
				sessionInfo := &session.SessionInfo{
					ID:     uuid.New(),
					UserID: uuid.New(),
					Role:   identity.Partner,
					State:  session.SessionActive,
				}
				return context.WithValue(context.Background(), session.GetSessionContextKey(), sessionInfo)
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
				return context.WithValue(context.Background(), session.GetSessionContextKey(), (*session.SessionInfo)(nil))
			},
			expectedExists: false,
			expectedNil:    true,
		},
		{
			name: "different context key",
			contextSetup: func() context.Context {
				sessionInfo := &session.SessionInfo{
					ID:     uuid.New(),
					UserID: uuid.New(),
					Role:   identity.Visitor,
					State:  session.SessionActive,
				}
				// Use wrong key type
				type wrongKey struct{}
				return context.WithValue(context.Background(), wrongKey{}, sessionInfo)
			},
			expectedExists: false,
			expectedNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.contextSetup()

			sessionInfo, exists := session.SessionInfoFromContext(ctx)

			assert.Equal(t, tt.expectedExists, exists, "unexpected existence result")

			if tt.expectedNil {
				assert.Nil(t, sessionInfo, "session should be nil")
			} else {
				assert.NotNil(t, sessionInfo, "session should not be nil")
			}
		})
	}
}

func TestSessionInfoFromContext_ValidSession(t *testing.T) {
	// Test that we get back the exact same session info we put in
	originalSessionInfo := &session.SessionInfo{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Role:   identity.Administrator,
		State:  session.SessionActive,
	}

	ctx := context.WithValue(context.Background(), session.GetSessionContextKey(), originalSessionInfo)

	retrievedSessionInfo, exists := session.SessionInfoFromContext(ctx)

	assert.True(t, exists, "session should exist in context")
	assert.NotNil(t, retrievedSessionInfo, "retrieved session should not be nil")

	// Verify it's the exact same session object
	assert.Same(t, originalSessionInfo, retrievedSessionInfo, "should return the exact same session object")

	// Verify session contents
	assert.Equal(t, originalSessionInfo.ID, retrievedSessionInfo.ID)
	assert.Equal(t, originalSessionInfo.UserID, retrievedSessionInfo.UserID)
	assert.Equal(t, originalSessionInfo.Role, retrievedSessionInfo.Role)
	assert.Equal(t, originalSessionInfo.State, retrievedSessionInfo.State)
}

func TestSessionInfoFromContext_BackwardCompatibility(t *testing.T) {
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
}

func TestSessionContextKey_Uniqueness(t *testing.T) {
	// Test that sessionContextKey is unique and doesn't conflict with other keys
	sessionInfo := &session.SessionInfo{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Role:   identity.Partner,
		State:  session.SessionActive,
	}

	// Create context with multiple values using different key types
	type otherKey struct{}

	ctx := context.Background()
	ctx = context.WithValue(ctx, session.GetSessionContextKey(), sessionInfo)
	ctx = context.WithValue(ctx, otherKey{}, "other value")
	ctx = context.WithValue(ctx, "string_key", "string value")

	// SessionInfoFromContext should only retrieve the session info, not other values
	retrievedSessionInfo, exists := session.SessionInfoFromContext(ctx)

	assert.True(t, exists, "session should exist")
	assert.Same(t, sessionInfo, retrievedSessionInfo, "should retrieve correct session")

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

	sessionInfo := &session.SessionInfo{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Role:   identity.Visitor,
		State:  session.SessionActive,
	}

	// Should work with both key instances
	ctx1 := context.WithValue(context.Background(), key1, sessionInfo)
	ctx2 := context.WithValue(context.Background(), key2, sessionInfo)

	sessionInfo1, exists1 := session.SessionInfoFromContext(ctx1)
	sessionInfo2, exists2 := session.SessionInfoFromContext(ctx2)

	assert.True(t, exists1, "session should exist in ctx1")
	assert.True(t, exists2, "session should exist in ctx2")
	assert.Same(t, sessionInfo, sessionInfo1, "should retrieve correct session from ctx1")
	assert.Same(t, sessionInfo, sessionInfo2, "should retrieve correct session from ctx2")
}

func TestSessionInfoFromContext_ConcurrentAccess(t *testing.T) {
	// Test that SessionInfoFromContext is safe for concurrent access
	sessionInfo := &session.SessionInfo{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Role:   identity.Administrator,
		State:  session.SessionActive,
	}

	ctx := context.WithValue(context.Background(), session.GetSessionContextKey(), sessionInfo)

	// Run multiple goroutines accessing the same context
	results := make(chan bool, 10)

	for range 10 {
		go func() {
			retrievedSessionInfo, exists := session.SessionInfoFromContext(ctx)
			results <- exists && retrievedSessionInfo != nil && retrievedSessionInfo.Role == identity.Administrator
		}()
	}

	// Verify all goroutines got the correct result
	for range 10 {
		result := <-results
		assert.True(t, result, "concurrent access should work correctly")
	}
}
