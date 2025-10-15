package helpers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/core/contracts/settings"
	"github.com/Leviosa-care/core/errs"
	tu "github.com/Leviosa-care/core/testutils"
	httpEndpoints "github.com/Leviosa-care/settings/internal/adapters/http"
	"github.com/Leviosa-care/settings/internal/domain"
	th "github.com/Leviosa-care/settings/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestGetCompanyPhone make test-integration-test

func TestGetCompanyPhone(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 404 when company phone not set", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewGetCompanyPhoneRequest(t, ctx, testServerURL)
		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "company telephone")
	})

	t.Run("should successfully retrieve company phone (admin endpoint)", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup: Insert company phone setting directly into database using generated ENCX functions
		phoneSetting := th.NewCompanyPhone(t, ctx)
		phoneSettingEncx, err := domain.ProcessSettingEncryptedEncx(ctx, crypto, phoneSetting)
		require.NoError(t, err)
		th.InsertCompanyPhoneEncrypted(t, ctx, phoneSettingEncx, testPool)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Test: Get the company phone (admin endpoint)
		req := th.NewGetCompanyPhoneRequest(t, ctx, testServerURL)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.GetCompanyTelephoneResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Equal(t, phoneSetting.Value, respBody.Telephone)
	})
}

// TEST=TestSetCompanyPhone make test-integration-test

func TestSetCompanyPhone(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully set company phone", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		phoneValue := "0145678910"
		request := domain.SetCompanyTelephoneRequest{Telephone: phoneValue}
		req := th.NewSetCompanyPhoneRequest(t, ctx, testServerURL, request)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.SetCompanyTelephoneResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.True(t, respBody.Success)

		// Verify data was persisted directly in database
		phoneSettingEncx := th.GetEncryptedSettingFromDB(t, ctx, settings.CompanyPhone, testPool)
		phoneSetting, err := domain.DecryptSettingEncryptedEncx(ctx, crypto, phoneSettingEncx)
		require.NoError(t, err)
		assert.Equal(t, phoneValue, phoneSetting.Value)

		// Verify RabbitMQ message was published
		th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyPhone, "0145678910")
	})

	t.Run("should return 400 for empty telephone", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetCompanyTelephoneRequest{Telephone: ""}
		req := th.NewSetCompanyPhoneRequest(t, ctx, testServerURL, request)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrInvalidValue.Error())
	})

	t.Run("should return 400 for telephone shorter than 10 characters", func(t *testing.T) {
		shortPhones := []string{
			"123",
			"12345678",
			"12345",
		}

		for _, phone := range shortPhones {
			t.Run("short phone: "+phone, func(t *testing.T) {
				th.ClearAllTestData(t, ctx, testPool)
				defer tu.ClearAuthData(t, ctx, authCtx)

				// Setup admin user and create authenticated request
				accessToken := tu.SetupAdminUser(t, ctx, authCtx)

				request := domain.SetCompanyTelephoneRequest{Telephone: phone}
				req := th.NewSetCompanyPhoneRequest(t, ctx, testServerURL, request)

				// Add authentication to the request
				authHeader := tu.CreateAuthHeader(accessToken)
				for key, values := range authHeader {
					for _, value := range values {
						req.Header.Add(key, value)
					}
				}
				req.AddCookie(tu.CreateAuthCookie(accessToken))

				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				var respBody struct {
					Error string `json:"error"`
				}
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.Contains(t, respBody.Error, errs.ErrInvalidValue.Error())
			})
		}
	})

	t.Run("should return 400 for telephone longer than 20 characters", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		longPhone := "01234567890123456789012345" // 25 characters

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetCompanyTelephoneRequest{Telephone: longPhone}
		req := th.NewSetCompanyPhoneRequest(t, ctx, testServerURL, request)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrInvalidValue.Error())
	})

	t.Run("should successfully accept various valid phone formats", func(t *testing.T) {
		validPhones := []string{
			"0123456789", // French mobile format
			"0156789012", // French Paris landline
			"0234567890", // French landline format
			"0345678901", // French landline format
			"0456789012", // French landline format
			"0612345678", // French mobile format
			"0723456789", // French mobile format
		}

		for _, phone := range validPhones {
			t.Run("valid phone: "+phone, func(t *testing.T) {
				th.ClearAllTestData(t, ctx, testPool)
				defer tu.ClearAuthData(t, ctx, authCtx)

				// Setup admin user and create authenticated request
				accessToken := tu.SetupAdminUser(t, ctx, authCtx)

				request := domain.SetCompanyTelephoneRequest{Telephone: phone}
				req := th.NewSetCompanyPhoneRequest(t, ctx, testServerURL, request)

				// Add authentication to the request
				authHeader := tu.CreateAuthHeader(accessToken)
				for key, values := range authHeader {
					for _, value := range values {
						req.Header.Add(key, value)
					}
				}
				req.AddCookie(tu.CreateAuthCookie(accessToken))

				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var respBody domain.SetCompanyTelephoneResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.True(t, respBody.Success)

				// Verify the phone was stored directly in database
				phoneSettingEncx := th.GetEncryptedSettingFromDB(t, ctx, settings.CompanyPhone, testPool)
				phoneSetting, err := domain.DecryptSettingEncryptedEncx(ctx, crypto, phoneSettingEncx)
				require.NoError(t, err)
				assert.Equal(t, phone, phoneSetting.Value)
			})
		}
	})

	t.Run("should handle whitespace trimming correctly", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Phone with leading/trailing whitespace
		phoneValue := "0123456789"
		phoneValueWithSpaces := "  " + phoneValue + "  "

		request := domain.SetCompanyTelephoneRequest{Telephone: phoneValueWithSpaces}
		req := th.NewSetCompanyPhoneRequest(t, ctx, testServerURL, request)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify the phone was trimmed and stored correctly directly in database
		phoneSettingEncx := th.GetEncryptedSettingFromDB(t, ctx, settings.CompanyPhone, testPool)
		phoneSetting, err := domain.DecryptSettingEncryptedEncx(ctx, crypto, phoneSettingEncx)
		require.NoError(t, err)
		// Should be able to retrieve the original value
		assert.Equal(t, phoneValue, phoneSetting.Value)

		th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyPhone, phoneValue)
	})

	t.Run("should return 415 for incorrect content type", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetCompanyTelephoneRequest{Telephone: "0123456789"}
		req := th.NewSetCompanyPhoneRequest(t, ctx, testServerURL, request)
		req.Header.Set("Content-Type", "text/plain")

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "unsupported media type")
	})

	t.Run("should return 400 for unknown JSON fields", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+httpEndpoints.SetCompanyPhoneEndpoint,
			strings.NewReader(`{"telephone": "0123456789", "unknown_field": "value"}`))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "invalid request body")
	})

	t.Run("should successfully update existing company phone", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		// Set initial phone
		initialPhoneSetting := th.NewCompanyPhone(t, ctx)
		initialPhoneSettingEncx, err := domain.ProcessSettingEncryptedEncx(ctx, crypto, initialPhoneSetting)
		require.NoError(t, err)
		th.InsertCompanyPhoneEncrypted(t, ctx, initialPhoneSettingEncx, testPool)

		// Update to new phone
		newPhoneValue := "0222222222"

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetCompanyTelephoneRequest{Telephone: newPhoneValue}
		req := th.NewSetCompanyPhoneRequest(t, ctx, testServerURL, request)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		phoneSettingEncx := th.GetEncryptedSettingFromDB(t, ctx, settings.CompanyPhone, testPool)
		phoneSetting, err := domain.DecryptSettingEncryptedEncx(ctx, crypto, phoneSettingEncx)
		require.NoError(t, err)
		assert.Equal(t, newPhoneValue, phoneSetting.Value)

		// Verify RabbitMQ message was published
		th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyPhone, newPhoneValue)
	})

	// TODO: Add test for encryption verification once crypto service is properly configured
	// This would involve checking that the phone number is actually encrypted in the database
}
