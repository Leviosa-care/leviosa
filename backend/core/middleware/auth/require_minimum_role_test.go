package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRequireMinimumRole(t *testing.T) {
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
			mockCrypto := encx.NewCryptoServiceMock()

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
			mockRepo.On("FindSessionByAccessTokenHash", mock.Anything, mock.AnythingOfType("string")).Return(session.ID.String(), sessionData, nil)

			middleware := NewSessionAuthMiddleware(mockRepo, mockCrypto)

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
