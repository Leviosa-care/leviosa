package testdata

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/core/contracts/settings"
	"github.com/Leviosa-care/settings/internal/domain"
	td "github.com/Leviosa-care/settings/test/testdata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCompanyAddress(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 404 when company address not set", func(t *testing.T) {
		td.ClearSettingsTable(t, ctx, testPool)

		req := td.NewGetCompanyAddressRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "legal address")
	})

	t.Run("should successfully retrieve company address", func(t *testing.T) {
		td.ClearSettingsTable(t, ctx, testPool)

		// Setup: Insert company address directly into database
		td.InsertCompanyAddress(t, ctx, "123 Main St, New York, NY 10001", testPool)

		// Test: Get the company address
		req := td.NewGetCompanyAddressRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.GetCompanyLegalAddressResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Equal(t, "123 Main St, New York, NY 10001", respBody.Address)
	})
}

func TestSetCompanyAddress(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully set company address", func(t *testing.T) {
		td.ClearSettingsTable(t, ctx, testPool)
		
		// Create a test channel for RabbitMQ verification
		testCh := td.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()
		
		// Purge queues to ensure clean state
		td.PurgeSettingsQueues(t, testCh)

		request := domain.SetCompanyLegalAddressRequest{Address: "456 Business Ave, San Francisco, CA 94105"}
		req := td.NewSetCompanyAddressRequest(t, ctx, testServerURL, request)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.SetCompanyLegalAddressResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.True(t, respBody.Success)

		// Verify data was persisted directly in database
		address, err := td.GetCompanyAddressFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, "456 Business Ave, San Francisco, CA 94105", address)

		// Verify RabbitMQ message was published
		td.VerifySettingsUpdateMessage(t, testCh, settings.CompanyLegalAddress, "456 Business Ave, San Francisco, CA 94105")
	})

	t.Run("should return 400 for empty address", func(t *testing.T) {
		td.ClearSettingsTable(t, ctx, testPool)

		request := domain.SetCompanyLegalAddressRequest{Address: ""}
		req := td.NewSetCompanyAddressRequest(t, ctx, testServerURL, request)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "address_required")
	})

	t.Run("should return 400 for whitespace-only address", func(t *testing.T) {
		td.ClearSettingsTable(t, ctx, testPool)

		request := domain.SetCompanyLegalAddressRequest{Address: "   \n\t   "}
		req := td.NewSetCompanyAddressRequest(t, ctx, testServerURL, request)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "address_empty")
	})

	t.Run("should return 400 for address exceeding 500 characters", func(t *testing.T) {
		td.ClearSettingsTable(t, ctx, testPool)

		longAddress := strings.Repeat("A", 501)
		request := domain.SetCompanyLegalAddressRequest{Address: longAddress}
		req := td.NewSetCompanyAddressRequest(t, ctx, testServerURL, request)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "address_length")
	})

	t.Run("should successfully accept multiline addresses", func(t *testing.T) {
		td.ClearSettingsTable(t, ctx, testPool)

		multilineAddress := `Company Headquarters
123 Main Street
Suite 456
New York, NY 10001
United States`

		request := domain.SetCompanyLegalAddressRequest{Address: multilineAddress}
		req := td.NewSetCompanyAddressRequest(t, ctx, testServerURL, request)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.SetCompanyLegalAddressResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.True(t, respBody.Success)

		// Verify the multiline address was stored correctly directly in database
		address, err := td.GetCompanyAddressFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, multilineAddress, address)
		assert.Contains(t, address, "\n")
	})

	t.Run("should successfully accept international addresses", func(t *testing.T) {
		td.ClearAllTestData(t, ctx, testPool)

		internationalAddresses := []string{
			"1-2-3 Shibuya, Tokyo 150-0002, Japan",
			"Unter den Linden 1, 10117 Berlin, Germany",
			"10 Downing Street, London SW1A 2AA, United Kingdom",
			"Champs-Élysées 75008 Paris, France",
		}

		for _, address := range internationalAddresses {
			t.Run("international address", func(t *testing.T) {
				td.ClearSettingsTable(t, ctx, testPool)

				request := domain.SetCompanyLegalAddressRequest{Address: address}
				req := td.NewSetCompanyAddressRequest(t, ctx, testServerURL, request)

				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var respBody domain.SetCompanyLegalAddressResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.True(t, respBody.Success)
			})
		}
	})

	t.Run("should return 415 for incorrect content type", func(t *testing.T) {
		td.ClearSettingsTable(t, ctx, testPool)

		request := domain.SetCompanyLegalAddressRequest{Address: "123 Test St"}
		req := td.NewSetCompanyAddressRequest(t, ctx, testServerURL, request)
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

	t.Run("should successfully update existing company address", func(t *testing.T) {
		td.ClearSettingsTable(t, ctx, testPool)

		// Set initial address
		request1 := domain.SetCompanyLegalAddressRequest{Address: "Old Address, Old City, OC 12345"}
		req1 := td.NewSetCompanyAddressRequest(t, ctx, testServerURL, request1)
		resp1, err := client.Do(req1)
		require.NoError(t, err)
		defer resp1.Body.Close()
		require.Equal(t, http.StatusOK, resp1.StatusCode)

		// Update to new address
		request2 := domain.SetCompanyLegalAddressRequest{Address: "New Address, New City, NC 67890"}
		req2 := td.NewSetCompanyAddressRequest(t, ctx, testServerURL, request2)
		resp2, err := client.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()
		require.Equal(t, http.StatusOK, resp2.StatusCode)

		// Verify updated address directly in database
		address, err := td.GetCompanyAddressFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, "New Address, New City, NC 67890", address)
	})
}
