package helpers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
)

// make test-func TEST_NAME=TestGetCompanyAddress TEST_PATH=test/integration/settings/company_address_test.go

func TestGetCompanyAddress(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 404 when company address not set", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)

		req := th.NewGetCompanyAddressRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, "legal address")
	})

	t.Run("should successfully retrieve company address", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)

		// Setup: Insert company address directly into database
		address := "123 Main St, New York, NY 10001"
		th.InsertCompanyAddress(t, ctx, address, testPool)

		// Test: Get the company address
		req := th.NewGetCompanyAddressRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.GetCompanyLegalAddressResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Equal(t, address, respBody.Address)
	})
}

// make test-func TEST_NAME=TestSetCompanyAddress TEST_PATH=test/integration/settings/company_address_test.go

func TestSetCompanyAddress(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully set company address", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		addressValue := "456 Business Ave, San Francisco, CA 94105"
		request := domain.SetCompanyLegalAddressRequest{Address: addressValue}
		req := th.NewSetCompanyAddressRequest(t, ctx, testServerURL, request)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.SetCompanyLegalAddressResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.True(t, respBody.Success)

		// Verify data was persisted directly in database
		retrievedAddress, err := th.GetCompanyAddress(t, ctx, testPool)
		assert.NoError(t, err)
		assert.Equal(t, addressValue, retrievedAddress)

		// Verify RabbitMQ message was published
		th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyLegalAddress, addressValue)
	})

	t.Run("should return 400 for empty address", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetCompanyLegalAddressRequest{Address: ""}
		req := th.NewSetCompanyAddressRequest(t, ctx, testServerURL, request)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, "address_required")
	})

	t.Run("should return 400 for whitespace-only address", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetCompanyLegalAddressRequest{Address: "   \n\t   "}
		req := th.NewSetCompanyAddressRequest(t, ctx, testServerURL, request)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, "address_empty")
	})

	t.Run("should return 400 for address exceeding 500 characters", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		longAddress := strings.Repeat("A", 501)
		request := domain.SetCompanyLegalAddressRequest{Address: longAddress}
		req := th.NewSetCompanyAddressRequest(t, ctx, testServerURL, request)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, "address_length")
	})

	t.Run("should successfully accept multiline addresses", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		multilineAddress := `Company Headquarters

	123 Main Street
	Suite 456
	New York, NY 10001
	United States`

		request := domain.SetCompanyLegalAddressRequest{Address: multilineAddress}
		req := th.NewSetCompanyAddressRequest(t, ctx, testServerURL, request)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.SetCompanyLegalAddressResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.True(t, respBody.Success)

		// Verify the multiline address was stored correctly directly in database
		address, err := th.GetCompanyAddress(t, ctx, testPool)
		assert.NoError(t, err)
		assert.Equal(t, multilineAddress, address)
		assert.Contains(t, address, "\n")

		// Verify RabbitMQ message was published
		th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyLegalAddress, multilineAddress)
	})

	t.Run("should successfully accept international addresses", func(t *testing.T) {
		// th.ClearAllTestData(t, ctx, testPool, nil)
		// defer tu.ClearAuthData(t, ctx, authCtx)

		internationalAddresses := []string{
			"1-2-3 Shibuya, Tokyo 150-0002, Japan",
			"Unter den Linden 1, 10117 Berlin, Germany",
			"10 Downing Street, London SW1A 2AA, United Kingdom",
			"Champs-Élysées 75008 Paris, France",
		}

		for _, address := range internationalAddresses {
			t.Run("international address", func(t *testing.T) {
				th.ClearSettingTestData(t, ctx, testPool, s3Client)
				defer tu.ClearAuthData(t, ctx, authCtx)

				// Create a test channel for RabbitMQ verification
				testCh := th.GetRabbitMQChannel(t, testMQConn)
				defer testCh.Close()

				// Purge queues to ensure clean state
				th.PurgeSettingsQueues(t, testCh)

				// Setup admin user and create authenticated request
				accessToken := tu.SetupAdminUser(t, ctx, authCtx)

				request := domain.SetCompanyLegalAddressRequest{Address: address}
				req := th.NewSetCompanyAddressRequest(t, ctx, testServerURL, request)

				// Add authentication to the request
				authHeader := tu.CreateAuthHeader(accessToken)
				for key, values := range authHeader {
					for _, value := range values {
						req.Header.Add(key, value)
					}
				}
				req.AddCookie(tu.CreateAuthCookie(accessToken))

				resp, err := client.Do(req)
				assert.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var respBody domain.SetCompanyLegalAddressResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				assert.NoError(t, err)
				assert.True(t, respBody.Success)

				// Verify RabbitMQ message was published
				th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyLegalAddress, address)
			})
		}
	})

	t.Run("should return 415 for incorrect content type", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetCompanyLegalAddressRequest{Address: "123 Test St"}
		req := th.NewSetCompanyAddressRequest(t, ctx, testServerURL, request)
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
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, "unsupported media type")
	})

	t.Run("should successfully update existing company address", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Set initial address
		oldAddress := "Old Address, Old City, OC 12345"
		th.InsertCompanyAddress(t, ctx, oldAddress, testPool)

		// Update to new address
		newAddress := "New Address, New City, NC 67890"
		request := domain.SetCompanyLegalAddressRequest{Address: newAddress}
		req := th.NewSetCompanyAddressRequest(t, ctx, testServerURL, request)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp2, err := client.Do(req)
		assert.NoError(t, err)
		defer resp2.Body.Close()
		assert.Equal(t, http.StatusOK, resp2.StatusCode)

		// Verify updated address directly in database
		address, err := th.GetCompanyAddress(t, ctx, testPool)
		assert.NoError(t, err)
		assert.Equal(t, newAddress, address)

		// Verify RabbitMQ message was published
		th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyLegalAddress, newAddress)
	})
}
