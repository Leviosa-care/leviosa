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

		req := th.NewGetCompanyPhoneRequest(t, ctx, testServerURL)
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

		// Setup: Insert company phoneSetting directly into database
		phoneSetting := th.NewCompanyPhone(t, ctx)
		err := crypto.ProcessStruct(ctx, phoneSetting)
		require.NoError(t, err)
		th.InsertCompanyPhoneEncrypted(t, ctx, phoneSetting, testPool)

		// Test: Get the company phone (admin endpoint)
		req := th.NewGetCompanyPhoneRequest(t, ctx, testServerURL)
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

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		phoneValue := "0145678910"
		request := domain.SetCompanyTelephoneRequest{Telephone: phoneValue}
		req := th.NewSetCompanyPhoneRequest(t, ctx, testServerURL, request)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.SetCompanyTelephoneResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.True(t, respBody.Success)

		// Verify data was persisted directly in database
		phoneSetting := th.GetEncryptedSettingFromDB(t, ctx, settings.CompanyPhone, testPool)
		err = crypto.DecryptStruct(ctx, phoneSetting)
		require.NoError(t, err)
		assert.Equal(t, phoneValue, phoneSetting.Value)

		// Verify RabbitMQ message was published
		th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyPhone, "0145678910")
	})

	t.Run("should return 400 for empty telephone", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		request := domain.SetCompanyTelephoneRequest{Telephone: ""}
		req := th.NewSetCompanyPhoneRequest(t, ctx, testServerURL, request)

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
		th.ClearAllTestData(t, ctx, testPool)

		shortPhones := []string{
			"123",
			"12345678",
			"12345",
		}

		for _, phone := range shortPhones {
			t.Run("short phone: "+phone, func(t *testing.T) {
				request := domain.SetCompanyTelephoneRequest{Telephone: phone}
				req := th.NewSetCompanyPhoneRequest(t, ctx, testServerURL, request)

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

		longPhone := "01234567890123456789012345" // 25 characters
		request := domain.SetCompanyTelephoneRequest{Telephone: longPhone}
		req := th.NewSetCompanyPhoneRequest(t, ctx, testServerURL, request)

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
		th.ClearAllTestData(t, ctx, testPool)

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

				request := domain.SetCompanyTelephoneRequest{Telephone: phone}
				req := th.NewSetCompanyPhoneRequest(t, ctx, testServerURL, request)

				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var respBody domain.SetCompanyTelephoneResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.True(t, respBody.Success)

				// Verify the phone was stored directly in database
				phoneSetting := th.GetEncryptedSettingFromDB(t, ctx, settings.CompanyPhone, testPool)
				err = crypto.DecryptStruct(ctx, phoneSetting)
				require.NoError(t, err)
				assert.Equal(t, phone, phoneSetting.Value)
			})
		}
	})

	t.Run("should handle whitespace trimming correctly", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		// Phone with leading/trailing whitespace
		phoneValue := "0123456789"
		phoneValueWithSpaces := "  " + phoneValue + "  "
		request := domain.SetCompanyTelephoneRequest{Telephone: phoneValueWithSpaces}
		req := th.NewSetCompanyPhoneRequest(t, ctx, testServerURL, request)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify the phone was trimmed and stored correctly directly in database
		phoneSetting := th.GetEncryptedSettingFromDB(t, ctx, settings.CompanyPhone, testPool)
		err = crypto.DecryptStruct(ctx, phoneSetting)
		require.NoError(t, err)
		// Should be able to retrieve the original value
		assert.Equal(t, phoneValue, phoneSetting.Value)

		th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyPhone, phoneValue)
	})

	t.Run("should return 415 for incorrect content type", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		request := domain.SetCompanyTelephoneRequest{Telephone: "0123456789"}
		req := th.NewSetCompanyPhoneRequest(t, ctx, testServerURL, request)
		req.Header.Set("Content-Type", "text/plain")

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

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/admin/settings/phone",
			strings.NewReader(`{"telephone": "0123456789", "unknown_field": "value"}`))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

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

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		// Set initial phone
		initialPhoneSetting := th.NewCompanyPhone(t, ctx)
		err := crypto.ProcessStruct(ctx, initialPhoneSetting)
		require.NoError(t, err)
		th.InsertCompanyPhoneEncrypted(t, ctx, initialPhoneSetting, testPool)

		// Update to new phone
		newPhoneValue := "0222222222"
		request2 := domain.SetCompanyTelephoneRequest{Telephone: newPhoneValue}
		req2 := th.NewSetCompanyPhoneRequest(t, ctx, testServerURL, request2)
		resp2, err := client.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()
		require.Equal(t, http.StatusOK, resp2.StatusCode)

		phoneSetting := th.GetEncryptedSettingFromDB(t, ctx, settings.CompanyPhone, testPool)
		err = crypto.DecryptStruct(ctx, phoneSetting)
		require.NoError(t, err)
		assert.Equal(t, newPhoneValue, phoneSetting.Value)

		// Verify RabbitMQ message was published
		th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyPhone, newPhoneValue)
	})

	// TODO: Add test for encryption verification once crypto service is properly configured
	// This would involve checking that the phone number is actually encrypted in the database
}
