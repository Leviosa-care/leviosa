package availabilityHelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"

	"github.com/stretchr/testify/require"
)

// NewCreateAvailabilityRequest creates a request to create a new availability (partner endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewCreateAvailabilityRequest(t *testing.T, ctx context.Context, serverURL string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/availabilities", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewCreateRecurringAvailabilityRequest creates a request to create a recurring availability (partner endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewCreateRecurringAvailabilityRequest(t *testing.T, ctx context.Context, serverURL string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/availabilities/recurring", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewGetAvailabilityRequest creates a request to get an availability by ID (standard user endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetAvailabilityRequest(t *testing.T, ctx context.Context, serverURL string, availabilityID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+"/availabilities/"+availabilityID, nil)
	require.NoError(t, err)

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewGetPartnerAvailabilitiesRequest creates a request to get availabilities for a partner (standard user endpoint)
// queryParams is optional - map of query parameters (e.g., {"status": "available", "limit": "10"})
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetPartnerAvailabilitiesRequest(t *testing.T, ctx context.Context, serverURL string, partnerID string, queryParams map[string]string, accessToken string) *http.Request {
	urlStr := serverURL + "/partners/" + partnerID + "/availabilities"

	// Add query parameters if provided
	if len(queryParams) > 0 {
		params := url.Values{}
		for key, value := range queryParams {
			params.Add(key, value)
		}
		urlStr += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	require.NoError(t, err)

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewGetAvailableSlotsRequest creates a request to get available slots (standard user endpoint)
// queryParams is optional - map of query parameters (e.g., {"start_time": "2024-01-01T00:00:00Z", "limit": "10"})
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetAvailableSlotsRequest(t *testing.T, ctx context.Context, serverURL string, queryParams map[string]string, accessToken string) *http.Request {
	urlStr := serverURL + "/availabilities"

	// Add query parameters if provided
	if len(queryParams) > 0 {
		params := url.Values{}
		for key, value := range queryParams {
			params.Add(key, value)
		}
		urlStr += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	require.NoError(t, err)

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewUpdateAvailabilityRequest creates a request to update an availability (partner endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewUpdateAvailabilityRequest(t *testing.T, ctx context.Context, serverURL string, availabilityID string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, serverURL+"/availabilities/"+availabilityID, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewCancelAvailabilityRequest creates a request to cancel an availability (partner endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewCancelAvailabilityRequest(t *testing.T, ctx context.Context, serverURL string, availabilityID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/availabilities/"+availabilityID+"/cancel", nil)
	require.NoError(t, err)

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewBlockAvailabilityRequest creates a request to block an availability (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewBlockAvailabilityRequest(t *testing.T, ctx context.Context, serverURL string, availabilityID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/availabilities/"+availabilityID+"/block", nil)
	require.NoError(t, err)

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewCheckAvailabilityConflictRequest creates a request to check for availability conflicts (partner endpoint)
// queryParams should contain start_time, end_time, and optionally exclude_id
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewCheckAvailabilityConflictRequest(t *testing.T, ctx context.Context, serverURL string, partnerID string, queryParams map[string]string, accessToken string) *http.Request {
	urlStr := serverURL + "/partners/" + partnerID + "/availabilities/conflict"

	// Add query parameters if provided
	if len(queryParams) > 0 {
		params := url.Values{}
		for key, value := range queryParams {
			params.Add(key, value)
		}
		urlStr += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	require.NoError(t, err)

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewGetRoomGapsRequest creates a request to get time gaps in a room's schedule (partner endpoint)
// date should be in YYYY-MM-DD format
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetRoomGapsRequest(t *testing.T, ctx context.Context, serverURL string, roomID string, date string, accessToken string) *http.Request {
	urlStr := serverURL + "/availabilities/rooms/" + roomID + "/gaps?date=" + date

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	require.NoError(t, err)

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewSuggestBlocksRequest creates a request to get availability block suggestions (partner endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewSuggestBlocksRequest(t *testing.T, ctx context.Context, serverURL string, partnerID string, roomID string, accessToken string) *http.Request {
	urlStr := serverURL + "/partners/" + partnerID + "/rooms/" + roomID + "/suggest-blocks"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	require.NoError(t, err)

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}
