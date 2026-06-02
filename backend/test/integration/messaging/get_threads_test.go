package messaging_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newSessionInfo creates a SessionInfo for the given role with a random ID and user ID.
func newSessionInfo(role identity.Role) *session.SessionInfo {
	return &session.SessionInfo{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Role:   role,
		State:  session.SessionActive,
	}
}

// TestListThreads_StandardUser_Returns200 tests that an authenticated Standard (client)
// user can successfully call GET /threads and receive HTTP 200 with an empty thread array.
func TestListThreads_StandardUser_Returns200(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Clean state
	td.ClearSessionsRedis(t, ctx, redisClient)

	// Seed a session for a Standard user
	accessToken := td.CreateSessionWithEncryption(t, ctx, newSessionInfo(identity.Standard), redisClient, crypto)

	// Build request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/threads", nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  ck.AccessTokenCookieName,
		Value: accessToken,
	})

	// Act
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode, "GET /threads should return 200 for Standard user")

	var threads []interface{}
	err = json.NewDecoder(resp.Body).Decode(&threads)
	require.NoError(t, err, "Response body should be valid JSON array")
	assert.Empty(t, threads, "Standard user with no threads should get an empty array")
}

// TestListThreads_PartnerUser_Returns200 verifies that Partner users still get 200.
func TestListThreads_PartnerUser_Returns200(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	td.ClearSessionsRedis(t, ctx, redisClient)

	accessToken := td.CreateSessionWithEncryption(t, ctx, newSessionInfo(identity.Partner), redisClient, crypto)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/threads", nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  ck.AccessTokenCookieName,
		Value: accessToken,
	})

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "GET /threads should return 200 for Partner user")

	var threads []interface{}
	err = json.NewDecoder(resp.Body).Decode(&threads)
	require.NoError(t, err, "Response body should be valid JSON array")
	assert.Empty(t, threads, "Partner user with no threads should get an empty array")
}

// TestListThreads_AdminUser_Returns200 verifies that Administrator users still get 200.
func TestListThreads_AdminUser_Returns200(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	td.ClearSessionsRedis(t, ctx, redisClient)

	accessToken := td.CreateSessionWithEncryption(t, ctx, newSessionInfo(identity.Administrator), redisClient, crypto)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/threads", nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  ck.AccessTokenCookieName,
		Value: accessToken,
	})

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "GET /threads should return 200 for Administrator user")

	var threads []interface{}
	err = json.NewDecoder(resp.Body).Decode(&threads)
	require.NoError(t, err, "Response body should be valid JSON array")
	assert.Empty(t, threads, "Administrator user with no threads should get an empty array")
}

// TestCreateThread_StandardUser_Returns403 verifies that POST /threads (thread initiation)
// still returns 403 for Standard users — this permission is NOT lowered.
func TestCreateThread_StandardUser_Returns403(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	td.ClearSessionsRedis(t, ctx, redisClient)

	accessToken := td.CreateSessionWithEncryption(t, ctx, newSessionInfo(identity.Standard), redisClient, crypto)

	body := map[string]string{
		"participant_id": "00000000-0000-0000-0000-000000000001",
	}
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/threads", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  ck.AccessTokenCookieName,
		Value: accessToken,
	})

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusForbidden, resp.StatusCode,
		"POST /threads should still return 403 for Standard users (thread initiation is partner-only)")
}

// TestListThreads_Unauthenticated_Returns401 verifies that requests without
// authentication are rejected.
func TestListThreads_Unauthenticated_Returns401(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/threads", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
		"GET /threads should return 401 without authentication")
}
