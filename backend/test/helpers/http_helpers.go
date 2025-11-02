package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	authDomain "github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	authEndpoints "github.com/Leviosa-care/leviosa/backend/internal/authuser/interface/auth"
	userEndpoints "github.com/Leviosa-care/leviosa/backend/internal/authuser/interface/user"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// NewCheckEmailSendOTPRequest creates an HTTP request for email verification with OTP
func NewCheckEmailSendOTPRequest(t *testing.T, ctx context.Context, baseURL string, request authDomain.CheckEmailAvailabilityRequest) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+authEndpoints.CheckEmailSendOTPEndpoint,
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")
	return req
}

// NewSignInRequest creates an HTTP request for user sign-in
func NewSignInRequest(t *testing.T, ctx context.Context, baseURL string, request authDomain.SignInRequest) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal sign-in request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+authEndpoints.SignInEndpoint,
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to create sign-in HTTP request")

	req.Header.Set("Content-Type", "application/json")
	return req
}

// NewInvalidJSONRequest creates an HTTP request with invalid JSON body
func NewInvalidJSONRequest(t *testing.T, ctx context.Context, baseURL, method, path string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		fmt.Sprintf("%s%s", baseURL, path),
		bytes.NewBufferString("{invalid json"),
	)
	require.NoError(t, err, "Failed to create invalid JSON HTTP request")

	req.Header.Set("Content-Type", "application/json")
	return req
}

// NewSignOutRequest creates an HTTP request for user sign-out with authentication
func NewSignOutRequest(t *testing.T, ctx context.Context, baseURL string, accessToken string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+authEndpoints.SignOutEndpoint,
		nil,
	)
	require.NoError(t, err, "Failed to create sign-out HTTP request")

	// Add session cookie if provided
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}
	return req
}

// NewSignOutRequestWithoutAuth creates an HTTP request for user sign-out without authentication
func NewSignOutRequestWithoutAuth(t *testing.T, ctx context.Context, baseURL string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+authEndpoints.SignOutEndpoint,
		nil,
	)
	require.NoError(t, err, "Failed to create sign-out HTTP request without auth")

	return req
}

// ParseJSONResponse parses a JSON HTTP response into the provided struct
func ParseJSONResponse(t *testing.T, resp *http.Response, v interface{}) {
	t.Helper()

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(v)
	require.NoError(t, err, "Failed to decode JSON response")
}

// ParseCheckEmailAvailabilityResponse parses the HTTP response for email availability
func ParseCheckEmailAvailabilityResponse(t *testing.T, resp *http.Response) (bool, string) {
	t.Helper()

	var response struct {
		Available bool   `json:"available"`
		Message   string `json:"message"`
	}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode availability response")

	return response.Available, response.Message
}

// ParseCheckEmailSendOTPResponse parses the HTTP response for email verification request
func ParseCheckEmailSendOTPResponse(t *testing.T, resp *http.Response) (string, string) {
	t.Helper()

	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode OTP response")

	return response.Message, response.Status
}

// ParseErrorResponse parses error responses from the API
func ParseErrorResponse(t *testing.T, resp *http.Response) (string, int) {
	t.Helper()

	var response struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode error response")

	errorMsg := response.Error
	if errorMsg == "" {
		errorMsg = response.Message
	}

	return errorMsg, resp.StatusCode
}

// NewValidateOTPCreatePendingUserRequest creates an HTTP request for OTP validation with user creation
func NewValidateOTPCreatePendingUserRequest(t *testing.T, ctx context.Context, baseURL string, request authDomain.ValidateOTPRequest) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal OTP validation request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+authEndpoints.ValidateOTPCreatePendingEndpoint,
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")
	return req
}

// ParseValidateOTPCreatePendingUserResponse parses the HTTP response for OTP validation with user creation
func ParseValidateOTPCreatePendingUserResponse(t *testing.T, resp *http.Response) (string, string) {
	t.Helper()

	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode validate OTP create pending user response")

	return response.Message, response.Status
}

// NewCompleteUserRequest creates an HTTP request for completing user registration
func NewCompleteUserRequest(t *testing.T, ctx context.Context, baseURL string, request authDomain.CompleteUserRequest, accessToken string) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal complete user request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+authEndpoints.CompleteUserEndpoint,
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")

	// Add session cookie if provided
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// ParseCompleteUserResponse parses the HTTP response for complete user request
func ParseCompleteUserResponse(t *testing.T, resp *http.Response) string {
	t.Helper()

	var response struct {
		Message string `json:"message"`
	}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode complete user response")

	return response.Message
}

// NewGetPendingUsersRequest creates an HTTP request for getting pending users
func NewGetPendingUsersRequest(t *testing.T, ctx context.Context, baseURL string, accessToken string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		baseURL+userEndpoints.GetPendingUsersEndpoint,
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

	// Add session cookie if provided
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}
	return req
}

// ParseGetPendingUsersResponse parses the HTTP response for get pending users request
func ParseGetPendingUsersResponse(t *testing.T, resp *http.Response) []*authDomain.UserResponse {
	t.Helper()

	var users []*authDomain.UserResponse
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&users)
	require.NoError(t, err, "Failed to decode get pending users response")

	return users
}

// NewGetAllUsersRequest creates an HTTP request for getting all users
func NewGetAllUsersRequest(t *testing.T, ctx context.Context, baseURL string, accessToken string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		baseURL+userEndpoints.GetAllUsersEndpoint,
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

	// Add session cookie if provided
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}
	return req
}

// NewGetAllUsersRequestWithoutAuth creates an HTTP request for getting all users without authentication
func NewGetAllUsersRequestWithoutAuth(t *testing.T, ctx context.Context, baseURL string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		baseURL+userEndpoints.GetAllUsersEndpoint,
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

	// Explicitly do not add any authorization headers
	return req
}

// ParseGetAllUsersResponse parses the HTTP response for get all users request
func ParseGetAllUsersResponse(t *testing.T, resp *http.Response) []*authDomain.UserResponse {
	t.Helper()

	var users []*authDomain.UserResponse
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&users)
	require.NoError(t, err, "Failed to decode get all users response")

	return users
}

// NewGetUserRequest creates an HTTP request for getting the current user's profile
func NewGetUserRequest(t *testing.T, ctx context.Context, baseURL string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		baseURL+userEndpoints.GetUserEndpoint,
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

	return req
}

// NewGetUserRequestWithAuth creates an HTTP request for getting user profile with authentication cookie
func NewGetUserRequestWithAuth(t *testing.T, ctx context.Context, baseURL string, accessToken string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		baseURL+userEndpoints.GetUserEndpoint,
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

	// Add authentication cookie
	cookie := &http.Cookie{
		Name:  ck.AccessTokenCookieName,
		Value: accessToken,
	}
	req.AddCookie(cookie)

	return req
}

// NewGetUserRequestWithMockAuth creates an HTTP request for getting user profile with mock auth data
// Note: This is for testing purposes when auth middleware is not present in the test setup
func NewGetUserRequestWithMockAuth(t *testing.T, ctx context.Context, baseURL string, userID uuid.UUID) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		baseURL+userEndpoints.GetUserEndpoint,
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

	// Note: In a real scenario with auth middleware, we would add the session cookie here
	// For now, this documents the intended behavior
	return req
}

// ParseGetUserResponse parses the HTTP response for get user request
func ParseGetUserResponse(t *testing.T, resp *http.Response) *authDomain.UserResponse {
	t.Helper()

	var user *authDomain.UserResponse
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&user)
	require.NoError(t, err, "Failed to decode get user response")

	return user
}

// NewApproveUserRequest creates an HTTP request for approving a user
func NewApproveUserRequest(t *testing.T, ctx context.Context, baseURL string, request authDomain.ApproveUserRequest, accessToken string) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		baseURL+userEndpoints.ApproveUserEndpoint,
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")

	// Add session cookie if provided
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}
	return req
}

// NewMalformedApproveUserRequest creates an HTTP request with malformed JSON
func NewMalformedApproveUserRequest(t *testing.T, ctx context.Context, baseURL string, accessToken string) *http.Request {
	t.Helper()

	malformedJSON := `{"user_id": "not-a-uuid", "role": }`

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		baseURL+userEndpoints.ApproveUserEndpoint,
		bytes.NewBuffer([]byte(malformedJSON)),
	)
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")

	// Add session cookie if provided
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}
	return req
}

// ParseApproveUserResponse parses the HTTP response for approve user request
func ParseApproveUserResponse(t *testing.T, resp *http.Response) map[string]string {
	t.Helper()

	var response map[string]string
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode approve user response")

	return response
}

// NewUpdateUserRoleRequest creates an HTTP request for updating a user's role
func NewUpdateUserRoleRequest(t *testing.T, ctx context.Context, baseURL string, userID uuid.UUID, role string, accessToken string) *http.Request {
	t.Helper()

	requestBody := struct {
		Role string `json:"role"`
	}{
		Role: role,
	}

	jsonData, err := json.Marshal(requestBody)
	require.NoError(t, err, "Failed to marshal request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/admin/users/%s/role", baseURL, userID.String()),
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to create HTTP request")

	// Add authentication cookie if access token is provided
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewMalformedUpdateUserRoleRequest creates an HTTP request with malformed JSON for role update
func NewMalformedUpdateUserRoleRequest(t *testing.T, ctx context.Context, baseURL string, userID uuid.UUID, accessToken string) *http.Request {
	t.Helper()

	malformedJSON := `{"role": }`

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/admin/users/%s/role", baseURL, userID.String()),
		bytes.NewBuffer([]byte(malformedJSON)),
	)
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")

	// Add session cookie if provided
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}
	return req
}

// ParseUpdateUserRoleResponse parses the HTTP response for update user role request
func ParseUpdateUserRoleResponse(t *testing.T, resp *http.Response) map[string]string {
	t.Helper()

	var response map[string]string
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode update user role response")

	return response
}

// NewRefreshSessionRequest creates an HTTP request for refreshing a session with refresh token cookie
func NewRefreshSessionRequest(t *testing.T, ctx context.Context, baseURL string, refreshToken string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+ck.RefreshEndpoint,
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

	// Add refresh token cookie
	cookie := &http.Cookie{
		Name:  ck.RefreshTokenCookieName,
		Value: refreshToken,
	}
	req.AddCookie(cookie)

	return req
}

// NewRefreshSessionRequestWithoutToken creates an HTTP request for refreshing a session without token
func NewRefreshSessionRequestWithoutToken(t *testing.T, ctx context.Context, baseURL string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+ck.RefreshEndpoint,
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

	// Explicitly do not add any cookies
	return req
}

// ParseRefreshSessionResponse parses the HTTP response for session refresh
func ParseRefreshSessionResponse(t *testing.T, resp *http.Response) (string, string) {
	t.Helper()

	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode refresh session response")

	return response.Message, response.Status
}

// NewDeleteUserByAdminRequest creates an HTTP request for admin deleting a user
func NewDeleteUserByAdminRequest(t *testing.T, ctx context.Context, baseURL string, userID uuid.UUID, accessToken string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/admin/auth/users/%s", baseURL, userID.String()),
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

	// Add session cookie if provided
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}
	return req
}

// NewDeleteUserByAdminRequestWithoutAuth creates an HTTP request for admin deleting a user without authentication
func NewDeleteUserByAdminRequestWithoutAuth(t *testing.T, ctx context.Context, baseURL string, userID uuid.UUID) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/admin/auth/users/%s", baseURL, userID.String()),
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

	// Explicitly do not add any authorization headers
	return req
}

// NewDeleteOwnAccountRequest creates an HTTP request for user deleting their own account
func NewDeleteOwnAccountRequest(t *testing.T, ctx context.Context, baseURL string, accessToken string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		baseURL+authEndpoints.DeleteOwnAccountEndpoint,
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

	// Add session cookie if provided
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}
	return req
}

// NewDeleteOwnAccountRequestWithoutAuth creates an HTTP request for user deleting their own account without authentication
func NewDeleteOwnAccountRequestWithoutAuth(t *testing.T, ctx context.Context, baseURL string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		baseURL+authEndpoints.DeleteOwnAccountEndpoint,
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

	// Explicitly do not add any authorization headers
	return req
}

// ParseDeleteUserResponse parses the HTTP response for delete user request
func ParseDeleteUserResponse(t *testing.T, resp *http.Response) string {
	t.Helper()

	var response struct {
		Message string `json:"message"`
	}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode delete user response")

	return response.Message
}

// NewGetUserByIDRequest creates an HTTP request for getting a specific user by ID (admin endpoint)
func NewGetUserByIDRequest(t *testing.T, ctx context.Context, baseURL string, userID uuid.UUID, accessToken string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/admin/users/%s", baseURL, userID.String()),
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

	// Add session cookie if provided
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}
	return req
}

// NewGetUserByIDRequestWithoutAuth creates an HTTP request for getting a specific user by ID without authentication
func NewGetUserByIDRequestWithoutAuth(t *testing.T, ctx context.Context, baseURL string, userID uuid.UUID) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/admin/users/%s", baseURL, userID.String()),
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

	// Explicitly do not add any authorization headers
	return req
}

// ParseGetUserByIDResponse parses the HTTP response for get user by ID request
func ParseGetUserByIDResponse(t *testing.T, resp *http.Response) *authDomain.UserResponse {
	t.Helper()

	var user *authDomain.UserResponse
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&user)
	require.NoError(t, err, "Failed to decode get user by ID response")

	return user
}

// NewUpdateUserRequest creates an HTTP request for updating user profile without authentication
func NewUpdateUserRequest(t *testing.T, ctx context.Context, baseURL string, request authDomain.UpdateUserRequest) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal update user request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		baseURL+userEndpoints.GetUserEndpoint,
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")
	return req
}

// NewUpdateUserRequestWithAuth creates an HTTP request for updating user profile with authentication
func NewUpdateUserRequestWithAuth(t *testing.T, ctx context.Context, baseURL string, request authDomain.UpdateUserRequest, accessToken string) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal update user request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		baseURL+userEndpoints.GetUserEndpoint,
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")

	// Add authentication cookie
	cookie := &http.Cookie{
		Name:  ck.AccessTokenCookieName,
		Value: accessToken,
	}
	req.AddCookie(cookie)

	return req
}

// ParseUpdateUserResponse parses the HTTP response for update user request
func ParseUpdateUserResponse(t *testing.T, resp *http.Response) *authDomain.UserResponse {
	t.Helper()

	var user *authDomain.UserResponse
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&user)
	require.NoError(t, err, "Failed to decode update user response")

	return user
}

// NewChangePasswordRequest creates an HTTP request for changing password
func NewChangePasswordRequest(t *testing.T, ctx context.Context, baseURL string, request authDomain.ChangePasswordRequest, accessToken string) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal change password request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		baseURL+userEndpoints.ChangePasswordEndpoint,
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")

	// Add authentication cookie
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewChangePasswordRequestWithoutAuth creates an HTTP request for changing password without authentication
func NewChangePasswordRequestWithoutAuth(t *testing.T, ctx context.Context, baseURL string, request authDomain.ChangePasswordRequest) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal change password request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		baseURL+userEndpoints.ChangePasswordEndpoint,
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")
	return req
}

// ParseChangePasswordResponse parses the HTTP response for change password request
func ParseChangePasswordResponse(t *testing.T, resp *http.Response) string {
	t.Helper()

	var response struct {
		Message string `json:"message"`
	}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode change password response")

	return response.Message
}

// NewRequestPasswordResetRequest creates an HTTP request for password reset
func NewRequestPasswordResetRequest(t *testing.T, ctx context.Context, baseURL string, request authDomain.RequestPasswordResetRequest) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+authEndpoints.RequestPasswordResetEndpoint,
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")
	return req
}

// ParseRequestPasswordResetResponse parses the HTTP response for password reset request
func ParseRequestPasswordResetResponse(t *testing.T, resp *http.Response) (string, string) {
	t.Helper()

	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode password reset response")

	return response.Message, response.Status
}

// AddAuthCookie adds authentication cookie to an existing request
func AddAuthCookie(req *http.Request, accessToken string) {
	cookie := &http.Cookie{
		Name:  ck.AccessTokenCookieName,
		Value: accessToken,
	}
	req.AddCookie(cookie)
}

// NewValidatePasswordResetOTPRequest creates an HTTP request for password reset OTP validation
func NewValidatePasswordResetOTPRequest(t *testing.T, ctx context.Context, baseURL string, request authDomain.ValidatePasswordResetOTPRequest) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+authEndpoints.ValidatePasswordResetOTPEndpoint,
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")
	return req
}

// ParseValidatePasswordResetOTPResponse parses the success response from password reset OTP validation
func ParseValidatePasswordResetOTPResponse(t *testing.T, resp *http.Response) (message, status string, expiresAt string) {
	t.Helper()

	var response struct {
		Message   string `json:"message"`
		Status    string `json:"status"`
		ExpiresAt string `json:"expires_at"`
	}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode password reset OTP validation response")

	return response.Message, response.Status, response.ExpiresAt
}

// GetPasswordResetTokenCookie extracts the password reset token from response cookies
func GetPasswordResetTokenCookie(t *testing.T, resp *http.Response) *http.Cookie {
	t.Helper()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "leviosa_password_reset_token" {
			return cookie
		}
	}
	require.Fail(t, "Password reset token cookie not found in response")
	return nil
}

// NewConfirmPasswordResetRequest creates an HTTP request for password reset confirmation
func NewConfirmPasswordResetRequest(t *testing.T, ctx context.Context, baseURL string, request authDomain.ConfirmPasswordResetRequest) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+authEndpoints.ConfirmPasswordResetEndpoint,
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")
	return req
}

// NewConfirmPasswordResetRequestWithCookie creates an HTTP request with password reset token cookie
func NewConfirmPasswordResetRequestWithCookie(t *testing.T, ctx context.Context, baseURL string, request authDomain.ConfirmPasswordResetRequest, resetTokenCookie *http.Cookie) *http.Request {
	t.Helper()

	req := NewConfirmPasswordResetRequest(t, ctx, baseURL, request)

	// Add the reset token cookie
	if resetTokenCookie != nil {
		req.AddCookie(resetTokenCookie)
	}

	return req
}

// ParseConfirmPasswordResetResponse parses the success response from password reset confirmation
func ParseConfirmPasswordResetResponse(t *testing.T, resp *http.Response) (message, status string) {
	t.Helper()

	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode password reset confirmation response")

	return response.Message, response.Status
}

// NewCompletePartnerRequest creates an HTTP request for completing partner registration
func NewCompletePartnerRequest(t *testing.T, ctx context.Context, baseURL string, request authDomain.CompletePartnerRequest, accessToken string) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal complete partner request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+authEndpoints.CompletePartnerEndpoint,
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")

	// Add session cookie if provided
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// ParseCompletePartnerResponse parses the HTTP response for complete partner request
func ParseCompletePartnerResponse(t *testing.T, resp *http.Response) string {
	t.Helper()

	var response struct {
		Message string `json:"message"`
	}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode complete partner response")

	return response.Message
}
