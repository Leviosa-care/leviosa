package buildingHelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"

	"github.com/stretchr/testify/require"
)

// NewCreateBuildingRequest creates a request to create a new building (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewCreateBuildingRequest(t *testing.T, ctx context.Context, serverURL string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/buildings", bytes.NewReader(jsonBody))
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

// NewGetBuildingByIDRequest creates a request to get a building by ID (public endpoint)
// accessToken is optional - if empty, no auth cookie is added (tests public access)
func NewGetBuildingByIDRequest(t *testing.T, ctx context.Context, serverURL string, buildingID interface{}, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+"/buildings/"+buildingID.(string), nil)
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

// NewGetAllBuildingsRequest creates a request to get all buildings (public endpoint)
// queryParams is optional - map of query parameters (e.g., {"is_active": "true", "limit": "10"})
// accessToken is optional - if empty, no auth cookie is added (tests public access)
func NewGetAllBuildingsRequest(t *testing.T, ctx context.Context, serverURL string, queryParams map[string]string, accessToken string) *http.Request {
	url := serverURL + "/buildings"

	// Add query parameters if provided
	if len(queryParams) > 0 {
		url += "?"
		first := true
		for key, value := range queryParams {
			if !first {
				url += "&"
			}
			url += key + "=" + value
			first = false
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

// NewUpdateBuildingRequest creates a request to update a building (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewUpdateBuildingRequest(t *testing.T, ctx context.Context, serverURL string, buildingID string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, serverURL+"/buildings/"+buildingID, bytes.NewReader(jsonBody))
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
