package testdata

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/authuser/test/helpers"
	"github.com/Leviosa-care/core/testutils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// NextcloudOAuthTestHelper provides utilities for testing Nextcloud OAuth integration
type NextcloudOAuthTestHelper struct {
	Container *testutils.NextcloudContainer
	BaseURL   string
}

// SetupNextcloudOAuth initializes a Nextcloud container for OAuth testing
func SetupNextcloudOAuth(ctx context.Context, t *testing.T) *NextcloudOAuthTestHelper {
	container, err := testutils.SetupNextcloud(ctx, t)
	require.NoError(t, err)

	return &NextcloudOAuthTestHelper{
		Container: container,
		BaseURL:   container.BaseURL,
	}
}

// TeardownNextcloudOAuth cleans up the Nextcloud container
func (h *NextcloudOAuthTestHelper) TeardownNextcloudOAuth(ctx context.Context, t *testing.T) {
	if h.Container != nil {
		testutils.TeardownNextcloud(ctx, t, h.Container)
	}
}

// CreateNextcloudTestUser creates a test user in the Nextcloud instance
func (h *NextcloudOAuthTestHelper) CreateNextcloudTestUser(ctx context.Context, t *testing.T, username, email, displayName string) {
	err := h.Container.CreateTestUser(ctx, username, email, displayName)
	require.NoError(t, err)
}

// GetOAuthConfig returns OAuth configuration for testing
func (h *NextcloudOAuthTestHelper) GetOAuthConfig() (clientID, clientSecret, baseURL string) {
	return h.Container.GetOAuthConfig()
}

// NewOAuthStartRequest creates an HTTP request to start OAuth flow with Nextcloud
func NewOAuthStartRequest(t *testing.T, serverURL, provider string) *http.Request {
	req, err := http.NewRequest("GET", serverURL+"/auth/oauth/"+provider+"/start", nil)
	require.NoError(t, err)
	return req
}

// NewOAuthCallbackRequest creates an HTTP request to handle OAuth callback from Nextcloud
func NewOAuthCallbackRequest(t *testing.T, serverURL, provider, code, state string) *http.Request {
	callbackURL := serverURL + "/auth/oauth/" + provider + "/callback"

	// Add OAuth callback parameters
	params := url.Values{}
	params.Add("code", code)
	params.Add("state", state)

	req, err := http.NewRequest("GET", callbackURL+"?"+params.Encode(), nil)
	require.NoError(t, err)
	return req
}

// InsertNextcloudOAuthUser creates a user with Google OAuth ID in the database (using Google field for Nextcloud testing)
func InsertNextcloudOAuthUser(t *testing.T, ctx context.Context, pool *pgxpool.Pool, user *domain.User) {
	const query = `
		INSERT INTO auth.users (
			id, state, email_encrypted, email_hash, password_hash, picture_encrypted,
			created_at_encrypted, logged_in_at_encrypted, role_encrypted, birth_date_encrypted,
			last_name_encrypted, first_name_encrypted, gender_encrypted, telephone_encrypted,
			telephone_hash, postal_code_encrypted, city_encrypted, address1_encrypted,
			address2_encrypted, google_id_encrypted, apple_id_encrypted,
			stripe_customer_id_encrypted, dek_encrypted, key_version
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23, $24
		)
	`

	_, err := pool.Exec(ctx, query,
		user.ID, user.State, user.EmailEncrypted, user.EmailHash, user.PasswordHash,
		user.PictureEncrypted, user.CreatedAtEncrypted, user.LoggedInAtEncrypted,
		user.RoleEncrypted, user.BirthDateEncrypted, user.LastNameEncrypted,
		user.FirstNameEncrypted, user.GenderEncrypted, user.TelephoneEncrypted,
		user.TelephoneHash, user.PostalCodeEncrypted, user.CityEncrypted,
		user.Address1Encrypted, user.Address2Encrypted, user.GoogleIDEncrypted,
		user.AppleIDEncrypted, user.StripeCustomerIDEncrypted,
		user.DEKEncrypted, user.KeyVersion,
	)
	require.NoError(t, err)
}

// GetUserByNextcloudIDFromDB retrieves a user by Google ID directly from the database (using Google field for Nextcloud testing)
func GetUserByNextcloudIDFromDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, nextcloudID string) (*domain.User, error) {
	const query = `
		SELECT id, state, email_encrypted, email_hash, password_hash, picture_encrypted,
			   created_at_encrypted, logged_in_at_encrypted, role_encrypted, birth_date_encrypted,
			   last_name_encrypted, first_name_encrypted, gender_encrypted, telephone_encrypted,
			   telephone_hash, postal_code_encrypted, city_encrypted, address1_encrypted,
			   address2_encrypted, google_id_encrypted, apple_id_encrypted,
			   stripe_customer_id_encrypted, dek_encrypted, key_version
		FROM auth.users 
		WHERE google_id_encrypted IS NOT NULL
		AND google_id_encrypted = pgp_sym_encrypt($1, $2)
	`

	// Use a test encryption key - in real tests this would come from test setup
	encryptionKey := "test-encryption-key-32-bytes!!"

	var user domain.User
	err := pool.QueryRow(ctx, query, nextcloudID, encryptionKey).Scan(
		&user.ID, &user.State, &user.EmailEncrypted, &user.EmailHash, &user.PasswordHash,
		&user.PictureEncrypted, &user.CreatedAtEncrypted, &user.LoggedInAtEncrypted,
		&user.RoleEncrypted, &user.BirthDateEncrypted, &user.LastNameEncrypted,
		&user.FirstNameEncrypted, &user.GenderEncrypted, &user.TelephoneEncrypted,
		&user.TelephoneHash, &user.PostalCodeEncrypted, &user.CityEncrypted,
		&user.Address1Encrypted, &user.Address2Encrypted, &user.GoogleIDEncrypted,
		&user.AppleIDEncrypted, &user.StripeCustomerIDEncrypted,
		&user.DEKEncrypted, &user.KeyVersion,
	)

	return &user, err
}

// VerifyNextcloudOAuthUser verifies that a user with Google OAuth ID exists in the database (using Google field for Nextcloud testing)
func VerifyNextcloudOAuthUser(t *testing.T, ctx context.Context, pool *pgxpool.Pool, nextcloudID, email string) {
	user, err := GetUserByNextcloudIDFromDB(t, ctx, pool, nextcloudID)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, nextcloudID, user.GoogleID) // This would need decryption in real tests
}

// MockOAuthCode generates a mock OAuth authorization code for testing
func MockOAuthCode(t *testing.T) string {
	return "mock_oauth_code_" + generateTestString(16)
}

// MockOAuthState generates a mock OAuth state parameter for testing
func MockOAuthState(t *testing.T) string {
	return "mock_oauth_state_" + generateTestString(32)
}

// generateTestString creates a test string of specified length
func generateTestString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	return string(result)
}

// SetupNextcloudAsGoogleOAuthEnvironment configures environment variables to use Nextcloud as Google OAuth provider for testing
func SetupNextcloudAsGoogleOAuthEnvironment(t *testing.T, nextcloudURL, clientID, clientSecret string) {
	t.Setenv("GOOGLE_CLIENT_ID", clientID)
	t.Setenv("GOOGLE_CLIENT_SECRET", clientSecret)
	t.Setenv("USE_NEXTCLOUD_FOR_TESTING", "true")
	t.Setenv("NEXTCLOUD_TEST_URL", nextcloudURL)
	t.Setenv("BASE_URL", "http://localhost:8080") // Test server base URL
}

// ExtractOAuthRedirectURL parses the OAuth redirect URL from a response
func ExtractOAuthRedirectURL(t *testing.T, response *http.Response) *url.URL {
	location := response.Header.Get("Location")
	require.NotEmpty(t, location, "OAuth redirect location should be present")

	redirectURL, err := url.Parse(location)
	require.NoError(t, err)

	return redirectURL
}

// ValidateOAuthStartRedirect validates that an OAuth start request properly redirects to the provider
func ValidateOAuthStartRedirect(t *testing.T, response *http.Response, expectedProviderURL string) {
	require.Equal(t, http.StatusFound, response.StatusCode)

	redirectURL := ExtractOAuthRedirectURL(t, response)
	require.True(t, strings.HasPrefix(redirectURL.String(), expectedProviderURL),
		"Redirect URL should point to Nextcloud OAuth endpoint")

	// Validate OAuth parameters are present
	query := redirectURL.Query()
	require.NotEmpty(t, query.Get("client_id"), "client_id should be present")
	require.NotEmpty(t, query.Get("redirect_uri"), "redirect_uri should be present")
	require.NotEmpty(t, query.Get("state"), "state should be present")
	require.Equal(t, "code", query.Get("response_type"), "response_type should be 'code'")
}

// GenerateTestGoogleID creates a test Google OAuth ID
func GenerateTestGoogleID() string {
	return "google_user_" + generateTestString(21) // Google IDs are typically 21 characters
}

// NewTestGoogleUserWithEncryption creates a test user with Google OAuth ID and encrypts the data
func NewTestGoogleUserWithEncryption(t *testing.T, ctx context.Context, crypto interface{}, googleID, email, firstName, lastName string) *domain.User {
	user := helpers.NewTestUser(email, firstName, lastName)
	user.GoogleID = googleID
	user.State = domain.Active

	// Note: In a real test, we would use the actual crypto service to encrypt the user data
	// For testing purposes, we'll simulate the encryption
	if crypto != nil {
		user.EmailHash = "hashed_" + email
		user.EmailEncrypted = []byte("encrypted_" + email)
		user.FirstNameEncrypted = []byte("encrypted_" + firstName)
		user.LastNameEncrypted = []byte("encrypted_" + lastName)
		user.GoogleIDEncrypted = []byte("encrypted_" + googleID)
		user.DEKEncrypted = []byte("encrypted_dek")
		user.KeyVersion = 1
	}

	return user
}

// NewTestOAuthUserWithEncryption creates a test OAuth user with encryption (for linking tests)
func NewTestOAuthUserWithEncryption(t *testing.T, ctx context.Context, crypto interface{}, provider, oauthID, email, firstName, lastName string) *domain.User {
	user := helpers.NewTestUser(email, firstName, lastName)
	user.State = domain.Active

	// Set OAuth ID based on provider (for linking tests, we may not have an OAuth ID initially)
	if oauthID != "" {
		switch provider {
		case "google":
			user.GoogleID = oauthID
		case "apple":
			user.AppleID = oauthID
		}
	}

	// Simulate encryption
	if crypto != nil {
		user.EmailHash = "hashed_" + email
		user.EmailEncrypted = []byte("encrypted_" + email)
		user.FirstNameEncrypted = []byte("encrypted_" + firstName)
		user.LastNameEncrypted = []byte("encrypted_" + lastName)
		user.DEKEncrypted = []byte("encrypted_dek")
		user.KeyVersion = 1

		if oauthID != "" {
			switch provider {
			case "google":
				user.GoogleIDEncrypted = []byte("encrypted_" + oauthID)
			case "apple":
				user.AppleIDEncrypted = []byte("encrypted_" + oauthID)
			}
		}
	}

	return user
}
