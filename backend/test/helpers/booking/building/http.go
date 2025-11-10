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
