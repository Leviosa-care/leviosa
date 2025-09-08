package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Leviosa-care/authuser/internal/domain"

	ck "github.com/Leviosa-care/core/auth/cookies"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// NewCheckEmailSendOTPRequest creates an HTTP request for email verification with OTP
func NewCheckEmailSendOTPRequest(t *testing.T, ctx context.Context, baseURL string, request domain.CheckEmailAvailabilityRequest) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/auth/email", baseURL),
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")
	return req
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
func NewValidateOTPCreatePendingUserRequest(t *testing.T, ctx context.Context, baseURL string, request domain.ValidateOTPRequest) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal OTP validation request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/auth/otp", baseURL),
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
func NewCompleteUserRequest(t *testing.T, ctx context.Context, baseURL string, request domain.CompleteUserRequest, accessToken string) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal complete user request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/auth/complete", baseURL),
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
		fmt.Sprintf("%s/admin/auth/users/pending", baseURL),
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
func ParseGetPendingUsersResponse(t *testing.T, resp *http.Response) []*domain.UserResponse {
	t.Helper()

	var users []*domain.UserResponse
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
		fmt.Sprintf("%s/admin/users", baseURL),
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
		fmt.Sprintf("%s/admin/users", baseURL),
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

	// Explicitly do not add any authorization headers
	return req
}

// ParseGetAllUsersResponse parses the HTTP response for get all users request
func ParseGetAllUsersResponse(t *testing.T, resp *http.Response) []*domain.UserResponse {
	t.Helper()

	var users []*domain.UserResponse
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
		fmt.Sprintf("%s/users/me", baseURL),
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
		fmt.Sprintf("%s/users/me", baseURL),
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
		fmt.Sprintf("%s/users/me", baseURL),
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

	// Note: In a real scenario with auth middleware, we would add the session cookie here
	// For now, this documents the intended behavior
	return req
}

// ParseGetUserResponse parses the HTTP response for get user request
func ParseGetUserResponse(t *testing.T, resp *http.Response) *domain.UserResponse {
	t.Helper()

	var user *domain.UserResponse
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&user)
	require.NoError(t, err, "Failed to decode get user response")

	return user
}

// NewApproveUserRequest creates an HTTP request for approving a user
func NewApproveUserRequest(t *testing.T, ctx context.Context, baseURL string, request domain.ApproveUserRequest, accessToken string) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/admin/users/approve", baseURL),
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
		fmt.Sprintf("%s/admin/users/approve", baseURL),
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

// NewRefreshSessionRequest creates an HTTP request for refreshing a session with refresh token cookie
func NewRefreshSessionRequest(t *testing.T, ctx context.Context, baseURL string, refreshToken string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/auth/refresh", baseURL),
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
		fmt.Sprintf("%s/auth/refresh", baseURL),
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
