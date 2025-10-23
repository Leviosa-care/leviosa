package helpers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	httpEndpoints "github.com/Leviosa-care/settings/internal/adapters/http"
	"github.com/Leviosa-care/settings/internal/domain"
	th "github.com/Leviosa-care/settings/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestBulkSettingsHandler make test-integration-test

func TestBulkSettingsHandler(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 400 when keys parameter is missing", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+httpEndpoints.AdminBulkEndpoint, nil)
		require.NoError(t, err)

		// Add authentication to the request
		req.Header = tu.CreateAuthHeader(accessToken)
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 when keys parameter is empty", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+httpEndpoints.AdminBulkEndpoint+"?keys=", nil)
		require.NoError(t, err)

		// Add authentication to the request
		req.Header = tu.CreateAuthHeader(accessToken)
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should successfully retrieve single setting", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup: Insert company name
		th.InsertCompanyName(t, ctx, "Test Company", testPool)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewBulkSettingsRequest(t, ctx, testServerURL, []string{settings.CompanyName})

		// Add authentication to the request
		req.Header = tu.CreateAuthHeader(accessToken)
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []*settings.SettingDTO
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.Len(t, response, 1)
		assert.Equal(t, settings.CompanyName, response[0].Key)
		assert.Equal(t, "Test Company", response[0].Value)
	})

	t.Run("should successfully retrieve multiple settings", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup: Insert multiple settings
		th.InsertCompanyName(t, ctx, "Bulk Test Company", testPool)
		th.InsertCompanyEmail(t, ctx, "bulk@test.com", testPool)
		th.InsertCompanyAddress(t, ctx, "123 Bulk Test St", testPool)
		th.InsertOTPDuration(t, ctx, 10, testPool)
		th.InsertOTPLength(t, ctx, 8, testPool)
		th.InsertOTPMaxAttempts(t, ctx, 5, testPool)

		keys := []string{
			settings.CompanyName,
			settings.CompanyEmail,
			settings.CompanyLegalAddress,
			settings.OTPDuration,
			settings.OTPLength,
			settings.OTPMaxAttempts,
		}

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewBulkSettingsRequest(t, ctx, testServerURL, keys)

		// Add authentication to the request
		req.Header = tu.CreateAuthHeader(accessToken)
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []settings.SettingDTO
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.Len(t, response, 6)

		// Create a map for easier verification
		responseMap := make(map[string]string)
		for _, setting := range response {
			responseMap[setting.Key] = setting.Value
		}

		assert.Equal(t, "Bulk Test Company", responseMap[settings.CompanyName])
		assert.Equal(t, "bulk@test.com", responseMap[settings.CompanyEmail])
		assert.Equal(t, "123 Bulk Test St", responseMap[settings.CompanyLegalAddress])
		assert.Equal(t, "10", responseMap[settings.OTPDuration])
		assert.Equal(t, "8", responseMap[settings.OTPLength])
		assert.Equal(t, "5", responseMap[settings.OTPMaxAttempts])
	})

	t.Run("should handle partial success with some invalid keys", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup: Insert only some settings
		th.InsertCompanyName(t, ctx, "Partial Test Company", testPool)
		th.InsertOTPLength(t, ctx, 6, testPool)

		keys := []string{
			settings.CompanyName,  // exists
			"invalid_key",         // invalid
			settings.OTPLength,    // exists
			settings.CompanyEmail, // doesn't exist (not set)
		}

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewBulkSettingsRequest(t, ctx, testServerURL, keys)

		// Add authentication to the request
		req.Header = tu.CreateAuthHeader(accessToken)
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMultiStatus, resp.StatusCode)

		var response struct {
			Data   []settings.SettingDTO `json:"data"`
			Errors map[string]any        `json:"errors"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should have 2 successful settings
		require.Len(t, response.Data, 2)

		// Create a map for easier verification
		dataMap := make(map[string]string)
		for _, setting := range response.Data {
			dataMap[setting.Key] = setting.Value
		}

		assert.Equal(t, "Partial Test Company", dataMap[settings.CompanyName])
		assert.Equal(t, "6", dataMap[settings.OTPLength])

		// Should have 2 errors
		require.Len(t, response.Errors, 2)
		assert.Contains(t, response.Errors, "invalid_key")
		assert.Contains(t, response.Errors, settings.CompanyEmail)
	})

	t.Run("should handle all invalid keys", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		keys := []string{"invalid_key1", "invalid_key2"}

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewBulkSettingsRequest(t, ctx, testServerURL, keys)

		// Add authentication to the request
		req.Header = tu.CreateAuthHeader(accessToken)
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMultiStatus, resp.StatusCode)

		var response struct {
			Data   []settings.SettingDTO `json:"data"`
			Errors map[string]any        `json:"errors"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should have no successful settings
		assert.Len(t, response.Data, 0)

		// Should have 2 errors
		require.Len(t, response.Errors, 2)
		assert.Contains(t, response.Errors, "invalid_key1")
		assert.Contains(t, response.Errors, "invalid_key2")
	})

	t.Run("should handle missing settings with not found errors", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		keys := []string{
			settings.CompanyName,
			settings.CompanyEmail,
			settings.OTPDuration,
		}

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewBulkSettingsRequest(t, ctx, testServerURL, keys)

		// Add authentication to the request
		req.Header = tu.CreateAuthHeader(accessToken)
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMultiStatus, resp.StatusCode)

		var response struct {
			Data   []settings.SettingDTO `json:"data"`
			Errors map[string]any        `json:"errors"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should have no successful settings
		assert.Len(t, response.Data, 0)

		// Should have 3 not found errors
		require.Len(t, response.Errors, 3)
		assert.Contains(t, response.Errors, settings.CompanyName)
		assert.Contains(t, response.Errors, settings.CompanyEmail)
		assert.Contains(t, response.Errors, settings.OTPDuration)
	})

	t.Run("should handle duplicate keys in request", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup: Insert company name
		th.InsertCompanyName(t, ctx, "Duplicate Test Company", testPool)

		keys := []string{
			settings.CompanyName,
			settings.CompanyName, // duplicate
			settings.CompanyName, // another duplicate
		}

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewBulkSettingsRequest(t, ctx, testServerURL, keys)

		// Add authentication to the request
		req.Header = tu.CreateAuthHeader(accessToken)
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []settings.SettingDTO
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return the setting multiple times (once for each key in the request)
		require.Len(t, response, 3)
		for _, setting := range response {
			assert.Equal(t, settings.CompanyName, setting.Key)
			assert.Equal(t, "Duplicate Test Company", setting.Value)
		}
	})

	t.Run("should handle all supported setting types", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup: Insert all types of settings
		th.InsertCompanyName(t, ctx, "All Types Company", testPool)
		th.InsertCompanyEmail(t, ctx, "all@types.com", testPool)
		th.InsertCompanyAddress(t, ctx, "456 All Types Ave", testPool)
		th.InsertCompanyInstagram(t, ctx, "https://instagram.com/alltypes", testPool)
		th.InsertOTPDuration(t, ctx, 20, testPool)
		th.InsertOTPLength(t, ctx, 4, testPool)
		th.InsertOTPMaxAttempts(t, ctx, 7, testPool)
		// TODO: Add company phone just to get encrypted value test
		phoneValue := "0612345679"
		phoneSetting := domain.SettingEncrypted{
			ID:        "",
			Key:       settings.CompanyPhone,
			Value:     phoneValue,
			CreatedAt: now,
			UpdatedAt: now,
		}

		phoneEncx, err := domain.ProcessSettingEncryptedEncx(ctx, crypto, &phoneSetting)
		require.NoError(t, err)
		th.InsertCompanyPhoneEncrypted(t, ctx, phoneEncx, testPool)

		// Note: CompanyLogo is not included as it requires special setup

		keys := []string{
			settings.CompanyName,
			settings.CompanyEmail,
			settings.CompanyLegalAddress,
			settings.CompanyInstagram,
			settings.CompanyPhone,
			settings.OTPDuration,
			settings.OTPLength,
			settings.OTPMaxAttempts,
		}

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewBulkSettingsRequest(t, ctx, testServerURL, keys)

		// Add authentication to the request
		req.Header = tu.CreateAuthHeader(accessToken)
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []settings.SettingDTO
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.Len(t, response, 8)

		// Create a map for easier verification
		responseMap := make(map[string]string)
		for _, setting := range response {
			responseMap[setting.Key] = setting.Value
		}

		assert.Equal(t, "All Types Company", responseMap[settings.CompanyName])
		assert.Equal(t, "all@types.com", responseMap[settings.CompanyEmail])
		assert.Equal(t, "456 All Types Ave", responseMap[settings.CompanyLegalAddress])
		assert.Equal(t, "https://instagram.com/alltypes", responseMap[settings.CompanyInstagram])
		assert.Equal(t, "20", responseMap[settings.OTPDuration])
		assert.Equal(t, "4", responseMap[settings.OTPLength])
		assert.Equal(t, "7", responseMap[settings.OTPMaxAttempts])
		assert.Equal(t, phoneValue, responseMap[settings.CompanyPhone])
	})
}
