package metricsHelpers

import (
	"context"
	"net/http"
	"testing"

	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"

	"github.com/stretchr/testify/require"
)

// NewGetRoomMetricsRequest creates a request to get utilization metrics for a room (partner endpoint)
// startDate and endDate should be in YYYY-MM-DD format
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetRoomMetricsRequest(t *testing.T, ctx context.Context, serverURL string, roomID string, startDate, endDate string, accessToken string) *http.Request {
	urlStr := serverURL + "/rooms/" + roomID + "/metrics?start_date=" + startDate + "&end_date=" + endDate

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

// NewGetPartnerMetricsRequest creates a request to get aggregated metrics for a partner's rooms (partner endpoint)
// startDate and endDate should be in YYYY-MM-DD format
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetPartnerMetricsRequest(t *testing.T, ctx context.Context, serverURL string, partnerID string, startDate, endDate string, accessToken string) *http.Request {
	urlStr := serverURL + "/partners/" + partnerID + "/metrics?start_date=" + startDate + "&end_date=" + endDate

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
