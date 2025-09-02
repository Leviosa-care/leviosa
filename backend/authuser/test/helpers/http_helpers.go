package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/middleware/auth"
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
func NewCompleteUserRequest(t *testing.T, ctx context.Context, baseURL string, request domain.CompleteUserRequest, sessionToken string) *http.Request {
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
	if sessionToken != "" {
		cookie := &http.Cookie{
			Name:  auth.AccessTokenCookieName,
			Value: sessionToken,
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
func NewGetPendingUsersRequest(t *testing.T, ctx context.Context, baseURL string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/admin/auth/users/pending", baseURL),
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

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
func NewGetAllUsersRequest(t *testing.T, ctx context.Context, baseURL string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/admin/users", baseURL),
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")

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
