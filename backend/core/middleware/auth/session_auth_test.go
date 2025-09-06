package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Leviosa-care/core/auth/cookies"
	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewSessionAuthMiddleware(t *testing.T) {
	mockRepo := &session.MockSessionRepository{}
	mockCrypto := &encx.CryptoServiceMock{}

	middleware := NewSessionAuthMiddleware(mockRepo, mockCrypto)

	assert.NotNil(t, middleware)
	assert.IsType(t, &SessionAuthMiddleware{}, middleware)
}

func TestCookieExtraction(t *testing.T) {
	// Simple test to verify cookie extraction works
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  cookies.AccessTokenCookieName,
		Value: "test_token",
	})

	cookie, err := req.Cookie(cookies.AccessTokenCookieName)
	require.NoError(t, err)
	assert.Equal(t, "test_token", cookie.Value)
}

func TestRequireAccessToken(t *testing.T) {
	tests := []struct {
		name           string
		setupCookie    func(*http.Request)
		repoResponse   []byte
		repoError      error
		expectedStatus int
		expectedInCtx  bool
	}{
		{
			name:           "missing access token cookie",
			setupCookie:    func(req *http.Request) {},
			expectedStatus: http.StatusUnauthorized,
			expectedInCtx:  false,
		},
		{
			name: "empty access token cookie",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  cookies.AccessTokenCookieName,
					Value: "",
				})
			},
			expectedStatus: http.StatusUnauthorized,
			expectedInCtx:  false,
		},
		{
			name: "valid access token - session not found",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  cookies.AccessTokenCookieName,
					Value: "valid_token",
				})
			},
			repoError:      errs.ErrRepositoryNotFound,
			expectedStatus: http.StatusUnauthorized,
			expectedInCtx:  false,
		},
		{
			name: "valid access token - database error",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  cookies.AccessTokenCookieName,
					Value: "valid_token",
				})
			},
			repoError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedInCtx:  false,
		},
		{
			name: "valid access token - invalid session JSON",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  cookies.AccessTokenCookieName,
					Value: "valid_token",
				})
			},
			repoResponse:   []byte("invalid json"),
			expectedStatus: http.StatusInternalServerError,
			expectedInCtx:  false,
		},
		{
			name: "valid access token - pending session (should work)",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  cookies.AccessTokenCookieName,
					Value: "valid_token",
				})
			},
			repoResponse: func() []byte {
				session := &session.Session{
					ID:     uuid.New(),
					UserID: uuid.New(),
					Role:   identity.Standard,
					State:  session.SessionPending,
				}
				return createValidSessionJSON(t, session)
			}(),
			expectedStatus: http.StatusOK,
			expectedInCtx:  true,
		},
		// {
		// 	name: "valid access token - active session",
		// 	setupCookie: func(req *http.Request) {
		// 		req.AddCookie(&http.Cookie{
		// 			Name:  cookies.AccessTokenCookieName,
		// 			Value: "valid_token",
		// 		})
		// 	},
		// 	repoResponse: func() []byte {
		// 		session := &session.Session{
		// 			ID:     uuid.New(),
		// 			UserID: uuid.New(),
		// 			Role:   identity.Standard,
		// 			State:  session.SessionActive,
		// 		}
		// 		return createValidSessionJSON(t, session)
		// 	}(),
		// 	expectedStatus: http.StatusOK,
		// 	expectedInCtx:  true,
		// },
		{
			name: "valid access token - expired session (invalid state)",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  cookies.AccessTokenCookieName,
					Value: "valid_token",
				})
			},
			repoResponse: func() []byte {
				session := &session.Session{
					ID:     uuid.New(),
					UserID: uuid.New(),
					Role:   identity.Standard,
					State:  "expired", // Invalid state
				}
				return createValidSessionJSON(t, session)
			}(),
			expectedStatus: http.StatusUnauthorized,
			expectedInCtx:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := &session.MockSessionRepository{}
			mockCrypto := &encx.CryptoServiceMock{}
			middleware := NewSessionAuthMiddleware(mockRepo, mockCrypto)

			if tt.repoResponse != nil || tt.repoError != nil {
				mockRepo.On("FindSessionByAccessTokenHash", mock.Anything, mock.AnythingOfType("string")).Return("test-session-id", tt.repoResponse, tt.repoError)
			}

			// Create test handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				session, found := session.SessionFromContext(r.Context())
				if tt.expectedInCtx {
					assert.True(t, found, "Expected session in context")
					assert.NotNil(t, session, "Expected non-nil session")
				} else {
					assert.False(t, found, "Expected no session in context")
				}
				w.WriteHeader(http.StatusOK)
			})

			// Create request
			req := httptest.NewRequest("GET", "/test", nil)
			tt.setupCookie(req)
			rr := httptest.NewRecorder()

			// Execute
			handler := middleware.RequireAccessToken(testHandler)
			// handler.ServeHTTP(rr, req)
			handler(rr, req)

			// Verify
			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRequireRefreshToken(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		setupCookie    func(*http.Request)
		repoResponse   []byte
		repoError      error
		expectedStatus int
		expectedInCtx  bool
	}{
		{
			name:           "non-refresh endpoint should be forbidden",
			path:           "/api/users",
			setupCookie:    func(req *http.Request) {},
			expectedStatus: http.StatusForbidden,
			expectedInCtx:  false,
		},
		{
			name:           "refresh endpoint - missing refresh token cookie",
			path:           cookies.RefreshEndpoint,
			setupCookie:    func(req *http.Request) {},
			expectedStatus: http.StatusUnauthorized,
			expectedInCtx:  false,
		},
		{
			name: "refresh endpoint - empty refresh token cookie",
			path: cookies.RefreshEndpoint,
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  cookies.RefreshTokenCookieName,
					Value: "",
				})
			},
			expectedStatus: http.StatusUnauthorized,
			expectedInCtx:  false,
		},
		{
			name: "refresh endpoint - valid refresh token - session not found",
			path: cookies.RefreshEndpoint,
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  cookies.RefreshTokenCookieName,
					Value: "valid_refresh_token",
				})
			},
			repoError:      errs.ErrRepositoryNotFound,
			expectedStatus: http.StatusUnauthorized,
			expectedInCtx:  false,
		},
		{
			name: "refresh endpoint - valid refresh token - database error",
			path: cookies.RefreshEndpoint,
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  cookies.RefreshTokenCookieName,
					Value: "valid_refresh_token",
				})
			},
			repoError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedInCtx:  false,
		},
		{
			name: "refresh endpoint - valid refresh token - active session",
			path: cookies.RefreshEndpoint,
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  cookies.RefreshTokenCookieName,
					Value: "valid_refresh_token",
				})
			},
			repoResponse: func() []byte {
				session := &session.Session{
					ID:     uuid.New(),
					UserID: uuid.New(),
					Role:   identity.Standard,
					State:  session.SessionActive,
				}
				return createValidSessionJSON(t, session)
			}(),
			expectedStatus: http.StatusOK,
			expectedInCtx:  true,
		},
		{
			name: "refresh endpoint - valid refresh token - pending session",
			path: cookies.RefreshEndpoint,
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  cookies.RefreshTokenCookieName,
					Value: "valid_refresh_token",
				})
			},
			repoResponse: func() []byte {
				session := &session.Session{
					ID:     uuid.New(),
					UserID: uuid.New(),
					Role:   identity.Standard,
					State:  session.SessionPending,
				}
				return createValidSessionJSON(t, session)
			}(),
			expectedStatus: http.StatusOK,
			expectedInCtx:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := &session.MockSessionRepository{}
			mockCrypto := &encx.CryptoServiceMock{}
			middleware := NewSessionAuthMiddleware(mockRepo, mockCrypto)

			if tt.repoResponse != nil || tt.repoError != nil {
				mockRepo.On("FindSessionByRefreshTokenHash", mock.Anything, mock.AnythingOfType("string")).Return("test-session-id", tt.repoResponse, tt.repoError)
			}

			// Create test handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				session, found := session.SessionFromContext(r.Context())
				if tt.expectedInCtx {
					assert.True(t, found, "Expected session in context")
					assert.NotNil(t, session, "Expected non-nil session")
				} else {
					assert.False(t, found, "Expected no session in context")
				}
				w.WriteHeader(http.StatusOK)
			})

			// Create request
			req := httptest.NewRequest("GET", tt.path, nil)
			tt.setupCookie(req)
			rr := httptest.NewRecorder()

			// Execute
			handler := middleware.RequireRefreshToken(testHandler)
			// handler.ServeHTTP(rr, req)
			handler(rr, req)

			// Verify
			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestSession is a test-only version of Session that can marshal plaintext fields
type TestSession struct {
	ID                 uuid.UUID            `json:"id"`
	UserID             uuid.UUID            `json:"user_id"`
	UserIDEncrypted    []byte               `json:"user_id_encrypted"`
	Role               identity.Role        `json:"role"`
	RoleEncrypted      []byte               `json:"role_encrypted"`
	State              session.SessionState `json:"state"`
	StateEncrypted     []byte               `json:"state_encrypted"`
	CreatedAt          time.Time            `json:"created_at"`
	CreatedAtEncrypted []byte               `json:"created_at_encrypted"`
	ExpiresAt          time.Time            `json:"expires_at"`
	ExpiresAtEncrypted []byte               `json:"expires_at_encrypted"`
	AccessTokenHash    string               `json:"access_token_hash"`
	RefreshTokenHash   string               `json:"refresh_token_hash"`
	DEKEncrypted       []byte               `json:"dek_encrypted"`
	KeyVersion         int                  `json:"key_version"`
}

// Helper function to create valid JSON session data for testing
func createValidSessionJSON(t *testing.T, sessionStruct *session.Session) []byte {
	t.Helper()

	// Create test session with plaintext fields that can be marshaled
	testSession := &TestSession{
		ID:                 session.ID,
		UserID:             session.UserID,
		UserIDEncrypted:    []byte("encrypted_user_id"),
		Role:               session.Role,
		RoleEncrypted:      []byte("encrypted_role"),
		State:              session.State,
		StateEncrypted:     []byte("encrypted_state"),
		CreatedAt:          time.Now(),
		CreatedAtEncrypted: []byte("encrypted_created_at"),
		ExpiresAt:          time.Now().Add(time.Hour),
		ExpiresAtEncrypted: []byte("encrypted_expires_at"),
		AccessTokenHash:    "access_token_hash",
		RefreshTokenHash:   "refresh_token_hash",
		DEKEncrypted:       []byte("encrypted_dek"),
		KeyVersion:         1,
	}

	// Marshal to JSON
	data, err := json.Marshal(testSession)
	require.NoError(t, err)

	return data
}
