package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewSessionAuthMiddleware(t *testing.T) {
	mockRepo := &MockSessionRepository{}

	middleware := NewSessionAuthMiddleware(mockRepo)

	assert.NotNil(t, middleware)
	assert.IsType(t, &SessionAuthMiddleware{}, middleware)
}

func TestSessionAuthMiddleware_RequireSession(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		repoResponse   []byte
		repoError      error
		expectedStatus int
		expectedInCtx  bool
		setupMock      func(*MockSessionRepository, string)
	}{
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedInCtx:  false,
			setupMock:      func(m *MockSessionRepository, token string) {},
		},
		{
			name:           "malformed authorization header - no bearer",
			authHeader:     "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
			expectedInCtx:  false,
			setupMock:      func(m *MockSessionRepository, token string) {},
		},
		{
			name:           "malformed authorization header - no token",
			authHeader:     "Bearer",
			expectedStatus: http.StatusUnauthorized,
			expectedInCtx:  false,
			setupMock:      func(m *MockSessionRepository, token string) {},
		},
		{
			name:           "empty token",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectedInCtx:  false,
			setupMock:      func(m *MockSessionRepository, token string) {},
		},
		{
			name:           "session not found",
			authHeader:     "Bearer valid_token",
			repoError:      errs.ErrRepositoryNotFound,
			expectedStatus: http.StatusUnauthorized,
			expectedInCtx:  false,
			setupMock: func(m *MockSessionRepository, token string) {
				m.On("FindSessionByTokenHash", mock.Anything, "valid_token").Return([]byte(nil), errs.ErrRepositoryNotFound)
			},
		},
		{
			name:           "repository error",
			authHeader:     "Bearer valid_token",
			repoError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedInCtx:  false,
			setupMock: func(m *MockSessionRepository, token string) {
				m.On("FindSessionByTokenHash", mock.Anything, "valid_token").Return([]byte(nil), errors.New("database error"))
			},
		},
		{
			name:           "invalid json session data",
			authHeader:     "Bearer valid_token",
			repoResponse:   []byte("invalid json"),
			expectedStatus: http.StatusInternalServerError,
			expectedInCtx:  false,
			setupMock: func(m *MockSessionRepository, token string) {
				m.On("FindSessionByTokenHash", mock.Anything, "valid_token").Return([]byte("invalid json"), nil)
			},
		},
		{
			name:       "pending session state",
			authHeader: "Bearer valid_token",
			repoResponse: createValidSessionJSON(t, &Session{
				ID:        uuid.New(),
				UserID:    uuid.New(),
				Role:      identity.Visitor,
				State:     SessionPending,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(time.Hour),
				TokenHash: "hash123",
			}),
			expectedStatus: http.StatusUnauthorized,
			expectedInCtx:  false,
			setupMock: func(m *MockSessionRepository, token string) {
				sessionData := createValidSessionJSON(t, &Session{
					ID:        uuid.New(),
					UserID:    uuid.New(),
					Role:      identity.Visitor,
					State:     SessionPending,
					CreatedAt: time.Now(),
					ExpiresAt: time.Now().Add(time.Hour),
					TokenHash: "hash123",
				})
				m.On("FindSessionByTokenHash", mock.Anything, "valid_token").Return(sessionData, nil)
			},
		},
		{
			name:       "valid active session",
			authHeader: "Bearer valid_token",
			repoResponse: createValidSessionJSON(t, &Session{
				ID:        uuid.New(),
				UserID:    uuid.New(),
				Role:      identity.Visitor,
				State:     SessionActive,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(time.Hour),
				TokenHash: "hash123",
			}),
			expectedStatus: http.StatusOK,
			expectedInCtx:  true,
			setupMock: func(m *MockSessionRepository, token string) {
				sessionData := createValidSessionJSON(t, &Session{
					ID:        uuid.New(),
					UserID:    uuid.New(),
					Role:      identity.Visitor,
					State:     SessionActive,
					CreatedAt: time.Now(),
					ExpiresAt: time.Now().Add(time.Hour),
					TokenHash: "hash123",
				})
				m.On("FindSessionByTokenHash", mock.Anything, "valid_token").Return(sessionData, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockSessionRepository{}
			tt.setupMock(mockRepo, "valid_token")

			middleware := NewSessionAuthMiddleware(mockRepo)

			// Create a test handler that checks for session in context
			var contextHasSession bool
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				session, ok := SessionFromContext(r.Context())
				contextHasSession = ok && session != nil
				w.WriteHeader(http.StatusOK)
			})

			// Wrap with middleware
			handler := middleware.RequireSession(testHandler)

			// Create request
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Record response
			w := httptest.NewRecorder()

			// Execute
			handler.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code, "unexpected status code")

			// Assert context session presence
			assert.Equal(t, tt.expectedInCtx, contextHasSession, "unexpected session presence in context")

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSessionAuthMiddleware_RequireSession_TokenExtraction(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		expectedOk  bool
		description string
	}{
		{
			name:        "valid bearer token",
			authHeader:  "Bearer abc123",
			expectedOk:  true,
			description: "should extract valid bearer token",
		},
		{
			name:        "bearer with extra spaces",
			authHeader:  "Bearer   abc123   ",
			expectedOk:  true,
			description: "should handle extra spaces after Bearer",
		},
		{
			name:        "bearer case sensitivity",
			authHeader:  "bearer abc123",
			expectedOk:  false,
			description: "should be case sensitive for Bearer",
		},
		{
			name:        "multiple spaces",
			authHeader:  "Bearer  token  with  spaces",
			expectedOk:  true,
			description: "should handle token with internal spaces",
		},
		{
			name:        "bearer only",
			authHeader:  "Bearer",
			expectedOk:  false,
			description: "should reject Bearer without token",
		},
		{
			name:        "empty token after bearer",
			authHeader:  "Bearer ",
			expectedOk:  false,
			description: "should reject empty token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockSessionRepository{}

			if tt.expectedOk {
				// Mock successful session retrieval for valid tokens
				sessionData := createValidSessionJSON(t, &Session{
					ID:        uuid.New(),
					UserID:    uuid.New(),
					Role:      identity.Visitor,
					State:     SessionActive,
					CreatedAt: time.Now(),
					ExpiresAt: time.Now().Add(time.Hour),
					TokenHash: "hash123",
				})
				mockRepo.On("FindSessionByTokenHash", mock.Anything, mock.AnythingOfType("string")).Return(sessionData, nil)
			}

			middleware := NewSessionAuthMiddleware(mockRepo)

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handler := middleware.RequireSession(testHandler)
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.authHeader)

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if tt.expectedOk {
				assert.Equal(t, http.StatusOK, w.Code, tt.description)
			} else {
				assert.Equal(t, http.StatusUnauthorized, w.Code, tt.description)
			}
		})
	}
}

// Helper function to create valid JSON session data for testing
func createValidSessionJSON(t *testing.T, session *Session) []byte {
	t.Helper()

	// Create a session with the required encrypted fields populated
	// In tests, we populate both plaintext and encrypted fields for DecodeSession to work
	testSession := &Session{
		ID:                 session.ID,
		UserIDEncrypted:    []byte("encrypted_user_id"),
		RoleEncrypted:      []byte("encrypted_role"),
		StateEncrypted:     []byte("encrypted_state"),
		CreatedAtEncrypted: []byte("encrypted_created_at"),
		ExpiresAtEncrypted: []byte("encrypted_expires_at"),
		TokenHash:          session.TokenHash,
		DEKEncrypted:       []byte("encrypted_dek"),
		KeyVersion:         1,
		// Set plaintext fields for test validation (DecodeSession only handles JSON unmarshaling)
		UserID:    session.UserID,
		Role:      session.Role,
		State:     session.State,
		CreatedAt: session.CreatedAt,
		ExpiresAt: session.ExpiresAt,
	}

	// Marshal the session to JSON using the actual struct
	jsonData, err := json.Marshal(testSession)
	require.NoError(t, err)

	return jsonData
}
