package roomHelpers

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

// NewCreateRoomRequest creates a request to create a new room (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewCreateRoomRequest(t *testing.T, ctx context.Context, serverURL string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/rooms", bytes.NewReader(jsonBody))
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

// NewGetRoomByIDRequest creates a request to get a room by ID (public endpoint)
// accessToken is optional - if empty, no auth cookie is added (tests public access)
func NewGetRoomByIDRequest(t *testing.T, ctx context.Context, serverURL string, roomID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+"/rooms/"+roomID, nil)
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

// NewListRoomsRequest creates a request to list rooms (public endpoint)
// queryParams is optional - map of query parameters (e.g., {"building_id": "uuid", "is_active": "true", "limit": "10"})
// accessToken is optional - if empty, no auth cookie is added (tests public access)
func NewListRoomsRequest(t *testing.T, ctx context.Context, serverURL string, queryParams map[string]string, accessToken string) *http.Request {
	urlStr := serverURL + "/rooms"

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

// NewUpdateRoomRequest creates a request to update a room (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewUpdateRoomRequest(t *testing.T, ctx context.Context, serverURL string, roomID string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, serverURL+"/rooms/"+roomID, bytes.NewReader(jsonBody))
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

// NewDeleteRoomRequest creates a request to delete a room (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewDeleteRoomRequest(t *testing.T, ctx context.Context, serverURL string, roomID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, serverURL+"/rooms/"+roomID, nil)
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

// NewGetRoomCountRequest creates a request to get room count (public endpoint)
// queryParams is optional - map of query parameters (e.g., {"building_id": "uuid", "is_active": "true", "name": "Consultation"})
// accessToken is optional - if empty, no auth cookie is added (tests public access)
func NewGetRoomCountRequest(t *testing.T, ctx context.Context, serverURL string, queryParams map[string]string, accessToken string) *http.Request {
	urlStr := serverURL + "/rooms/count"

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

// NewGetRoomsByBuildingRequest creates a request to get rooms by building ID (public endpoint)
// queryParams is optional - map of query parameters (e.g., {"active_only": "true"})
// accessToken is optional - if empty, no auth cookie is added (tests public access)
func NewGetRoomsByBuildingRequest(t *testing.T, ctx context.Context, serverURL string, buildingID interface{}, queryParams map[string]string, accessToken string) *http.Request {
	urlStr := serverURL + "/buildings/" + buildingID.(string) + "/rooms"

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

// NewSetRoomEquipmentRequest creates a request to set room equipment (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewSetRoomEquipmentRequest(t *testing.T, ctx context.Context, serverURL string, roomID string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, serverURL+"/rooms/"+roomID+"/equipment", bytes.NewReader(jsonBody))
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

// NewSetRoomRateRequest creates a request to set room hourly rate (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewSetRoomRateRequest(t *testing.T, ctx context.Context, serverURL string, roomID string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, serverURL+"/rooms/"+roomID+"/rate", bytes.NewReader(jsonBody))
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

// NewClearRoomRateRequest creates a request to clear room hourly rate (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewClearRoomRateRequest(t *testing.T, ctx context.Context, serverURL string, roomID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, serverURL+"/rooms/"+roomID+"/rate", nil)
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
