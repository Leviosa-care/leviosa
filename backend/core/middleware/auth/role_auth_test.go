package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/Leviosa-care/core/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSessionAuthMiddleware_RequireMinimumRole(t *testing.T) {
	tests := []struct {
		name           string
		userRole       identity.Role
		requiredRole   identity.Role
		expectedStatus int
		shouldCallNext bool
	}{
		{
			name:           "visitor cannot access staff endpoint",
			userRole:       identity.Visitor,
			requiredRole:   identity.Partner,
			expectedStatus: http.StatusForbidden,
			shouldCallNext: false,
		},
		{
			name:           "visitor can access visitor endpoint",
			userRole:       identity.Visitor,
			requiredRole:   identity.Visitor,
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "partner can access visitor endpoint",
			userRole:       identity.Partner,
			requiredRole:   identity.Visitor,
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "partner can access staff endpoint",
			userRole:       identity.Partner,
			requiredRole:   identity.Partner,
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "partner cannot access admin endpoint",
			userRole:       identity.Partner,
			requiredRole:   identity.Administrator,
			expectedStatus: http.StatusForbidden,
			shouldCallNext: false,
		},
		{
			name:           "admin can access any endpoint",
			userRole:       identity.Administrator,
			requiredRole:   identity.Partner,
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "admin can access admin endpoint",
			userRole:       identity.Administrator,
			requiredRole:   identity.Administrator,
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &MockSessionRepository{}

			// Create valid session data with the test role
			session := &Session{
				ID:               uuid.New(),
				UserID:           uuid.New(),
				Role:             tt.userRole,
				State:            SessionActive,
				CreatedAt:        time.Now(),
				ExpiresAt:        time.Now().Add(time.Hour),
				AccessTokenHash:  "test_access_hash",
				RefreshTokenHash: "test_refresh_hash",
			}
			sessionData := createValidSessionJSON(t, session)

			// Mock the repository call
			mockRepo.On("FindSessionByAccessToken", mock.Anything, mock.AnythingOfType("string")).Return(session.ID.String(), sessionData, nil)

			middleware := NewSessionAuthMiddleware(mockRepo)

			// Track if next handler was called
			nextCalled := false
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			// Apply role-based middleware
			handler := middleware.RequireMinimumRole(tt.requiredRole)(testHandler)

			// Create request with access token cookie
			req := httptest.NewRequest("GET", "/test", nil)
			req.AddCookie(&http.Cookie{
				Name:  AccessTokenCookieName,
				Value: "valid_token",
			})

			w := httptest.NewRecorder()
			// handler.ServeHTTP(w, req)
			handler(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code, "unexpected status code")
			assert.Equal(t, tt.shouldCallNext, nextCalled, "unexpected next handler call behavior")

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSessionAuthMiddleware_RequireAnyRole(t *testing.T) {
	tests := []struct {
		name           string
		userRole       identity.Role
		allowedRoles   []identity.Role
		expectedStatus int
		shouldCallNext bool
	}{
		{
			name:           "visitor matches visitor role",
			userRole:       identity.Visitor,
			allowedRoles:   []identity.Role{identity.Visitor},
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "visitor matches in multiple roles",
			userRole:       identity.Visitor,
			allowedRoles:   []identity.Role{identity.Partner, identity.Visitor},
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "partner matches staff in multiple roles",
			userRole:       identity.Partner,
			allowedRoles:   []identity.Role{identity.Partner, identity.Administrator},
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "visitor denied for partner-only endpoint",
			userRole:       identity.Visitor,
			allowedRoles:   []identity.Role{identity.Partner},
			expectedStatus: http.StatusForbidden,
			shouldCallNext: false,
		},
		{
			name:           "visitor denied for partner or admin endpoint",
			userRole:       identity.Visitor,
			allowedRoles:   []identity.Role{identity.Partner, identity.Administrator},
			expectedStatus: http.StatusForbidden,
			shouldCallNext: false,
		},
		{
			name:           "admin matches admin role",
			userRole:       identity.Administrator,
			allowedRoles:   []identity.Role{identity.Administrator},
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "empty roles list denies everyone",
			userRole:       identity.Administrator,
			allowedRoles:   []identity.Role{},
			expectedStatus: http.StatusForbidden,
			shouldCallNext: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &MockSessionRepository{}

			// Create valid session data with the test role
			session := &Session{
				ID:               uuid.New(),
				UserID:           uuid.New(),
				Role:             tt.userRole,
				State:            SessionActive,
				CreatedAt:        time.Now(),
				ExpiresAt:        time.Now().Add(time.Hour),
				AccessTokenHash:  "test_access_hash",
				RefreshTokenHash: "test_refresh_hash",
			}
			sessionData := createValidSessionJSON(t, session)

			// Mock the repository call
			mockRepo.On("FindSessionByAccessToken", mock.Anything, mock.AnythingOfType("string")).Return(session.ID.String(), sessionData, nil)

			middleware := NewSessionAuthMiddleware(mockRepo)

			// Track if next handler was called
			nextCalled := false
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			// Apply role-based middleware
			handler := middleware.RequireAnyRole(tt.allowedRoles...)(testHandler)

			// Create request with access token cookie
			req := httptest.NewRequest("GET", "/test", nil)
			req.AddCookie(&http.Cookie{
				Name:  AccessTokenCookieName,
				Value: "valid_token",
			})

			w := httptest.NewRecorder()
			// handler.ServeHTTP(w, req)
			handler(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code, "unexpected status code")
			assert.Equal(t, tt.shouldCallNext, nextCalled, "unexpected next handler call behavior")

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSessionAuthMiddleware_RequireAdmin(t *testing.T) {
	tests := []struct {
		name           string
		userRole       identity.Role
		expectedStatus int
		shouldCallNext bool
	}{
		{
			name:           "visitor denied admin access",
			userRole:       identity.Visitor,
			expectedStatus: http.StatusForbidden,
			shouldCallNext: false,
		},
		{
			name:           "partner denied admin access",
			userRole:       identity.Partner,
			expectedStatus: http.StatusForbidden,
			shouldCallNext: false,
		},
		{
			name:           "admin granted admin access",
			userRole:       identity.Administrator,
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &MockSessionRepository{}

			// Create valid session data with the test role
			session := &Session{
				ID:               uuid.New(),
				UserID:           uuid.New(),
				Role:             tt.userRole,
				State:            SessionActive,
				CreatedAt:        time.Now(),
				ExpiresAt:        time.Now().Add(time.Hour),
				AccessTokenHash:  "test_access_hash",
				RefreshTokenHash: "test_refresh_hash",
			}
			sessionData := createValidSessionJSON(t, session)

			// Mock the repository call
			mockRepo.On("FindSessionByAccessToken", mock.Anything, mock.AnythingOfType("string")).Return(session.ID.String(), sessionData, nil)

			middleware := NewSessionAuthMiddleware(mockRepo)

			// Track if next handler was called
			nextCalled := false
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			// Apply admin middleware (uses RequireMinimumRole internally)
			handler := middleware.RequireAdmin(testHandler)

			// Create request with access token cookie
			req := httptest.NewRequest("GET", "/test", nil)
			req.AddCookie(&http.Cookie{
				Name:  AccessTokenCookieName,
				Value: "valid_token",
			})

			w := httptest.NewRecorder()
			// handler.ServeHTTP(w, req)
			handler(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code, "unexpected status code")
			assert.Equal(t, tt.shouldCallNext, nextCalled, "unexpected next handler call behavior")

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRoleAuthMiddleware_NoSessionInContext(t *testing.T) {
	// Test behavior when session authentication fails but role middleware is still called
	// This shouldn't happen in normal flow, but tests edge case handling

	tests := []struct {
		name string
		// middlewareFn func(AuthMiddleware) func(http.Handler) http.Handler
		middlewareFn func(AuthMiddleware) func(middleware.Handler) middleware.Handler
	}{
		{
			name: "RequireMinimumRole with no session",
			middlewareFn: func(m AuthMiddleware) func(middleware.Handler) middleware.Handler {
				return m.RequireMinimumRole(identity.Visitor)
			},
		},
		{
			name: "RequireAnyRole with no session",
			middlewareFn: func(m AuthMiddleware) func(middleware.Handler) middleware.Handler {
				return m.RequireAnyRole(identity.Visitor, identity.Partner)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockSessionRepository{}
			middleware := NewSessionAuthMiddleware(mockRepo)

			// Create a handler that manually adds broken context (no session)
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create handler that simulates missing session in context
			brokenSessionHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Don't add session to context, simulate RequireSession failure
				// tt.middlewareFn(middleware)(testHandler).ServeHTTP(w, r)
				tt.middlewareFn(middleware)(testHandler)(w, r)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			brokenSessionHandler.ServeHTTP(w, req)

			// Should return 401 when no session in context
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

func TestRoleAuthMiddleware_Integration(t *testing.T) {
	// Test the full flow: RequireSession -> RequireMinimumRole
	mockRepo := &MockSessionRepository{}

	// Create session data for partner user
	session := &Session{
		ID:               uuid.New(),
		UserID:           uuid.New(),
		Role:             identity.Partner,
		State:            SessionActive,
		CreatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(time.Hour),
		AccessTokenHash:  "test_access_hash",
		RefreshTokenHash: "test_refresh_hash",
	}
	sessionData := createValidSessionJSON(t, session)
	mockRepo.On("FindSessionByAccessToken", mock.Anything, mock.AnythingOfType("string")).Return(session.ID.String(), sessionData, nil)

	middleware := NewSessionAuthMiddleware(mockRepo)

	// Create endpoint that requires partner role
	nextCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify session is available in context
		session, ok := SessionFromContext(r.Context())
		assert.True(t, ok, "session should be in context")
		assert.Equal(t, identity.Partner, session.Role, "session should have partner role")

		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Chain middlewares: RequireSession is called by RequireMinimumRole
	handler := middleware.RequireMinimumRole(identity.Partner)(testHandler)

	req := httptest.NewRequest("GET", "/partner-endpoint", nil)
	req.AddCookie(&http.Cookie{
		Name:  AccessTokenCookieName,
		Value: "valid_token",
	})

	w := httptest.NewRecorder()
	// handler.ServeHTTP(w, req)
	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, nextCalled, "next handler should be called")
	mockRepo.AssertExpectations(t)
}
