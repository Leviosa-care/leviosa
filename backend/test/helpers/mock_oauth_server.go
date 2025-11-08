package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
)

// MockOAuthServer provides a lightweight OAuth 2.0 server for testing
type MockOAuthServer struct {
	Server       *httptest.Server
	mu           sync.Mutex
	codes        map[string]*MockOAuthCode     // code -> user info
	accessTokens map[string]*MockOAuthUserInfo // access_token -> user info
	states       map[string]bool                // state -> valid
}

// MockOAuthCode represents a mock OAuth authorization code
type MockOAuthCode struct {
	Code     string
	UserInfo *MockOAuthUserInfo
}

// MockOAuthUserInfo represents user information returned by the OAuth provider
type MockOAuthUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// TokenResponse represents the OAuth token endpoint response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// SetupMockOAuthServer creates and starts a mock OAuth server
func SetupMockOAuthServer(t *testing.T) *MockOAuthServer {
	t.Helper()

	mockServer := &MockOAuthServer{
		codes:        make(map[string]*MockOAuthCode),
		accessTokens: make(map[string]*MockOAuthUserInfo),
		states:       make(map[string]bool),
	}

	mux := http.NewServeMux()

	// OAuth authorization endpoint
	mux.HandleFunc("/oauth/authorize", mockServer.handleAuthorize)

	// OAuth token endpoint
	mux.HandleFunc("/oauth/token", mockServer.handleToken)

	// OAuth userinfo endpoint
	mux.HandleFunc("/oauth/userinfo", mockServer.handleUserInfo)

	// OpenID configuration endpoint (for discovery)
	mux.HandleFunc("/.well-known/openid-configuration", mockServer.handleOpenIDConfig)

	mockServer.Server = httptest.NewServer(mux)

	return mockServer
}

// TeardownMockOAuthServer stops the mock OAuth server
func (m *MockOAuthServer) TeardownMockOAuthServer() {
	if m.Server != nil {
		m.Server.Close()
	}
}

// RegisterUser registers a user that can be authenticated via the mock OAuth flow
func (m *MockOAuthServer) RegisterUser(userInfo *MockOAuthUserInfo) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// In a real OAuth flow, we don't pre-register users, but for testing
	// we'll create a code that can be exchanged for this user's info
	code := "mock_code_" + userInfo.ID
	m.codes[code] = &MockOAuthCode{
		Code:     code,
		UserInfo: userInfo,
	}
}

// CreateAuthorizationCode creates an authorization code for a specific user
func (m *MockOAuthServer) CreateAuthorizationCode(userInfo *MockOAuthUserInfo) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	code := "mock_code_" + generateTestString(32)
	m.codes[code] = &MockOAuthCode{
		Code:     code,
		UserInfo: userInfo,
	}

	return code
}

// ValidateState registers a state parameter as valid
func (m *MockOAuthServer) ValidateState(state string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.states[state] = true
}

// handleAuthorize handles OAuth authorization requests
func (m *MockOAuthServer) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	clientID := query.Get("client_id")
	redirectURI := query.Get("redirect_uri")
	state := query.Get("state")
	responseType := query.Get("response_type")

	// Validate required parameters
	if clientID == "" || redirectURI == "" || state == "" || responseType != "code" {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	// For testing, auto-approve and redirect with a code
	// In real OAuth, this would show a consent screen
	m.mu.Lock()
	// Find or create a default test user
	code := "mock_code_" + generateTestString(32)
	m.codes[code] = &MockOAuthCode{
		Code: code,
		UserInfo: &MockOAuthUserInfo{
			ID:            "test_user_123",
			Email:         "testuser@example.com",
			EmailVerified: true,
			Name:          "Test User",
			GivenName:     "Test",
			FamilyName:    "User",
			Picture:       "https://example.com/picture.jpg",
		},
	}
	m.mu.Unlock()

	// Redirect back to the application with the authorization code
	redirectURL, _ := url.Parse(redirectURI)
	redirectQuery := redirectURL.Query()
	redirectQuery.Set("code", code)
	redirectQuery.Set("state", state)
	redirectURL.RawQuery = redirectQuery.Encode()

	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}

// handleToken handles OAuth token exchange requests
func (m *MockOAuthServer) handleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method_not_allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	grantType := r.FormValue("grant_type")
	code := r.FormValue("code")
	redirectURI := r.FormValue("redirect_uri")

	// Validate grant type
	if grantType != "authorization_code" {
		http.Error(w, "unsupported_grant_type", http.StatusBadRequest)
		return
	}

	// Validate required parameters
	if code == "" || redirectURI == "" {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	// Exchange code for user info
	m.mu.Lock()
	oauthCode, exists := m.codes[code]
	if !exists {
		m.mu.Unlock()
		http.Error(w, "invalid_grant", http.StatusBadRequest)
		return
	}

	// Generate access token
	accessToken := "mock_access_token_" + generateTestString(32)
	m.accessTokens[accessToken] = oauthCode.UserInfo

	// Remove code (one-time use)
	delete(m.codes, code)
	m.mu.Unlock()

	// Return token response
	tokenResponse := TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: "mock_refresh_token_" + generateTestString(32),
		Scope:        "openid email profile",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResponse)
}

// handleUserInfo handles OAuth userinfo requests
func (m *MockOAuthServer) handleUserInfo(w http.ResponseWriter, r *http.Request) {
	// Extract access token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse Bearer token
	var accessToken string
	if _, err := fmt.Sscanf(authHeader, "Bearer %s", &accessToken); err != nil {
		http.Error(w, "invalid_token", http.StatusUnauthorized)
		return
	}

	// Lookup user info by access token
	m.mu.Lock()
	userInfo, exists := m.accessTokens[accessToken]
	m.mu.Unlock()

	if !exists {
		http.Error(w, "invalid_token", http.StatusUnauthorized)
		return
	}

	// Return user info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

// handleOpenIDConfig handles OpenID Connect discovery requests
func (m *MockOAuthServer) handleOpenIDConfig(w http.ResponseWriter, r *http.Request) {
	baseURL := m.Server.URL

	config := map[string]interface{}{
		"issuer":                 baseURL,
		"authorization_endpoint": baseURL + "/oauth/authorize",
		"token_endpoint":         baseURL + "/oauth/token",
		"userinfo_endpoint":      baseURL + "/oauth/userinfo",
		"response_types_supported": []string{"code"},
		"subject_types_supported":  []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// GetBaseURL returns the base URL of the mock server
func (m *MockOAuthServer) GetBaseURL() string {
	return m.Server.URL
}

// GetAuthorizationURL returns the full authorization URL
func (m *MockOAuthServer) GetAuthorizationURL(clientID, redirectURI, state string) string {
	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("response_type", "code")
	params.Set("state", state)
	params.Set("scope", "openid email profile")

	return m.Server.URL + "/oauth/authorize?" + params.Encode()
}

// SetupMockOAuthEnvironment configures environment variables for the mock OAuth server
func SetupMockOAuthEnvironment(t *testing.T, provider, baseURL, clientID, clientSecret string) {
	t.Helper()

	switch provider {
	case "google":
		t.Setenv("GOOGLE_CLIENT_ID", clientID)
		t.Setenv("GOOGLE_CLIENT_SECRET", clientSecret)
		t.Setenv("GOOGLE_AUTH_URL", baseURL+"/oauth/authorize")
		t.Setenv("GOOGLE_TOKEN_URL", baseURL+"/oauth/token")
		t.Setenv("GOOGLE_USERINFO_URL", baseURL+"/oauth/userinfo")
	case "apple":
		t.Setenv("APPLE_CLIENT_ID", clientID)
		t.Setenv("APPLE_CLIENT_SECRET", clientSecret)
		t.Setenv("APPLE_AUTH_URL", baseURL+"/oauth/authorize")
		t.Setenv("APPLE_TOKEN_URL", baseURL+"/oauth/token")
		t.Setenv("APPLE_USERINFO_URL", baseURL+"/oauth/userinfo")
	}

	t.Setenv("BASE_URL", "http://localhost:8080") // Test server base URL for OAuth callback
}

// NewMockOAuthUserInfo creates a mock OAuth user info object
func NewMockOAuthUserInfo(id, email, givenName, familyName string) *MockOAuthUserInfo {
	return &MockOAuthUserInfo{
		ID:            id,
		Email:         email,
		EmailVerified: true,
		Name:          givenName + " " + familyName,
		GivenName:     givenName,
		FamilyName:    familyName,
		Picture:       "https://example.com/picture.jpg",
	}
}

// CreateMockOAuthTestUser creates a complete test user for OAuth testing
func CreateMockOAuthTestUser(t *testing.T, mockServer *MockOAuthServer, id, email, givenName, familyName string) (string, *MockOAuthUserInfo) {
	t.Helper()

	userInfo := NewMockOAuthUserInfo(id, email, givenName, familyName)
	code := mockServer.CreateAuthorizationCode(userInfo)

	return code, userInfo
}

// NewOAuthCallbackRequest creates an HTTP request to handle OAuth callback
func NewOAuthCallbackRequest(t *testing.T, serverURL, provider, code, state string) *http.Request {
	callbackURL := serverURL + "/auth/oauth/" + provider + "/callback"

	// Add OAuth callback parameters
	params := url.Values{}
	if code != "" {
		params.Add("code", code)
	}
	if state != "" {
		params.Add("state", state)
	}

	req, err := http.NewRequest("GET", callbackURL+"?"+params.Encode(), nil)
	if err != nil {
		t.Fatalf("Failed to create OAuth callback request: %v", err)
	}
	return req
}

// GenerateTestGoogleID creates a test Google OAuth ID
func GenerateTestGoogleID() string {
	return "google_user_" + generateTestString(21) // Google IDs are typically 21 characters
}

// NewTestGoogleUser creates a test user with Google OAuth ID
func NewTestGoogleUser(t *testing.T, googleID, email, firstName, lastName string) *MockOAuthUserInfo {
	t.Helper()

	return &MockOAuthUserInfo{
		ID:            googleID,
		Email:         email,
		EmailVerified: true,
		Name:          firstName + " " + lastName,
		GivenName:     firstName,
		FamilyName:    lastName,
		Picture:       "https://example.com/picture.jpg",
	}
}

// generateTestString creates a random string of specified length for test data
func generateTestString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	return string(result)
}
