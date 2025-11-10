package promotionCode_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdatePromotionCode TEST_PATH=test/integration/catalog/promotion_code/update_promotion_code_test.go

func TestUpdatePromotionCode(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("UpdatePromotionCodeMetadata", func(t *testing.T) {
		t.Run("should successfully update promotion code metadata with valid admin token", func(t *testing.T) {
			// Clean the database
			clearTables(t, ctx)

			// Setup admin user authentication
			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			// Create test coupon and promotion code
			testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
			td.InsertCoupon(t, ctx, testPool, testCoupon)

			testPromoCode := td.NewValidPromotionCode("UPDATE20", testCoupon.ID)
			td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

			// Update request
			updateRequest := domain.UpdatePromotionCodeRequest{
				Metadata: map[string]string{
					"updated":     "true",
					"test":        "false",
					"description": "Updated promotion code",
				},
			}

			req := th.NewUpdatePromotionCodeRequest(t, ctx, testServerURL, testPromoCode.ID.String(), updateRequest, accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response struct {
				Message string `json:"message"`
			}

			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Equal(t, "Promotion code updated successfully!", response.Message)

			// Verify the update in database
			updatedPromoCode, err := td.GetPromotionCodeByID(t, ctx, testPromoCode.ID, testPool)
			assert.NoError(t, err)
			assert.Equal(t, "true", updatedPromoCode.Metadata["updated"])
			assert.Equal(t, "false", updatedPromoCode.Metadata["test"])
			assert.Equal(t, "Updated promotion code", updatedPromoCode.Metadata["description"])
		})

		t.Run("should return 404 for non-existent promotion code", func(t *testing.T) {
			clearTables(t, ctx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			updateRequest := domain.UpdatePromotionCodeRequest{
				Metadata: map[string]string{"test": "value"},
			}

			req := th.NewUpdatePromotionCodeRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", updateRequest, accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("should return 400 for invalid UUID", func(t *testing.T) {
			clearTables(t, ctx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			updateRequest := domain.UpdatePromotionCodeRequest{
				Metadata: map[string]string{"test": "value"},
			}

			req := th.NewUpdatePromotionCodeRequest(t, ctx, testServerURL, "invalid-uuid", updateRequest, accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return 415 for non-JSON content type", func(t *testing.T) {
			clearTables(t, ctx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			req := th.NewUpdatePromotionCodeRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", map[string]string{"test": "data"}, accessToken)
			req.Header.Set("Content-Type", "text/plain")

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
		})

		t.Run("should return 401 when access token is missing", func(t *testing.T) {
			clearTables(t, ctx)

			updateRequest := domain.UpdatePromotionCodeRequest{
				Metadata: map[string]string{"test": "value"},
			}

			req := th.NewUpdatePromotionCodeRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", updateRequest, "")

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 401 when session is expired", func(t *testing.T) {
			clearTables(t, ctx)

			// Create expired admin session
			accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

			updateRequest := domain.UpdatePromotionCodeRequest{
				Metadata: map[string]string{"test": "value"},
			}

			req := th.NewUpdatePromotionCodeRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", updateRequest, accessToken)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
			clearTables(t, ctx)

			// Create standard user (not admin)
			accessToken := tu.SetupStandardUser(t, ctx, authCtx)

			updateRequest := domain.UpdatePromotionCodeRequest{
				Metadata: map[string]string{"test": "value"},
			}

			req := th.NewUpdatePromotionCodeRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", updateRequest, accessToken)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("should return 401 when token is invalid", func(t *testing.T) {
			clearTables(t, ctx)

			updateRequest := domain.UpdatePromotionCodeRequest{
				Metadata: map[string]string{"test": "value"},
			}

			req := th.NewUpdatePromotionCodeRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", updateRequest, "invalid-token-12345")

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})

	t.Run("DeactivatePromotionCode", func(t *testing.T) {
		t.Run("should successfully deactivate promotion code with valid admin token", func(t *testing.T) {
			// Clean the database
			clearTables(t, ctx)

			// Setup admin user authentication
			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			// Create test coupon and promotion code
			testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
			td.InsertCoupon(t, ctx, testPool, testCoupon)

			testPromoCode := td.NewValidPromotionCode("DEACTIVATE20", testCoupon.ID)
			td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

			// Verify it's initially active
			assert.True(t, td.GetPromotionCodeActiveStatus(t, ctx, testPromoCode.ID, testPool))

			req := th.NewDeactivatePromotionCodeRequest(t, ctx, testServerURL, testPromoCode.ID.String(), accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response struct {
				Message string `json:"message"`
			}

			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Equal(t, "Promotion code deactivated successfully!", response.Message)

			// Verify the promotion code is deactivated in database
			assert.False(t, td.GetPromotionCodeActiveStatus(t, ctx, testPromoCode.ID, testPool))
		})

		t.Run("should return 404 for non-existent promotion code", func(t *testing.T) {
			clearTables(t, ctx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			req := th.NewDeactivatePromotionCodeRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("should return 400 for invalid UUID", func(t *testing.T) {
			clearTables(t, ctx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			req := th.NewDeactivatePromotionCodeRequest(t, ctx, testServerURL, "invalid-uuid", accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return 401 when access token is missing", func(t *testing.T) {
			clearTables(t, ctx)

			req := th.NewDeactivatePromotionCodeRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", "")

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 401 when session is expired", func(t *testing.T) {
			clearTables(t, ctx)

			// Create expired admin session
			accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

			req := th.NewDeactivatePromotionCodeRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
			clearTables(t, ctx)

			// Create standard user (not admin)
			accessToken := tu.SetupStandardUser(t, ctx, authCtx)

			req := th.NewDeactivatePromotionCodeRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("should return 401 when token is invalid", func(t *testing.T) {
			clearTables(t, ctx)

			req := th.NewDeactivatePromotionCodeRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", "invalid-token-12345")

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})

	t.Run("DeletePromotionCode", func(t *testing.T) {
		t.Run("should successfully delete promotion code with valid admin token", func(t *testing.T) {
			// Clean the database
			clearTables(t, ctx)

			// Setup admin user authentication
			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			// Create test coupon and promotion code
			testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
			td.InsertCoupon(t, ctx, testPool, testCoupon)

			testPromoCode := td.NewValidPromotionCode("DELETE20", testCoupon.ID)
			td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

			// Verify it exists
			foundPromoCode := td.GetPromotionCodeByIDOrNil(t, ctx, testPromoCode.ID, testPool)
			assert.NotNil(t, foundPromoCode)

			req := th.NewDeletePromotionCodeRequest(t, ctx, testServerURL, testPromoCode.ID.String(), accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response struct {
				Message string `json:"message"`
			}

			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Equal(t, "Promotion code deleted successfully!", response.Message)

			// Verify the promotion code is deleted from database
			deletedPromoCode := td.GetPromotionCodeByIDOrNil(t, ctx, testPromoCode.ID, testPool)
			assert.Nil(t, deletedPromoCode)
		})

		t.Run("should return 404 for non-existent promotion code", func(t *testing.T) {
			clearTables(t, ctx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			req := th.NewDeletePromotionCodeRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("should return 400 for invalid UUID", func(t *testing.T) {
			clearTables(t, ctx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			req := th.NewDeletePromotionCodeRequest(t, ctx, testServerURL, "invalid-uuid", accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return 401 when access token is missing", func(t *testing.T) {
			clearTables(t, ctx)

			req := th.NewDeletePromotionCodeRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", "")

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 401 when session is expired", func(t *testing.T) {
			clearTables(t, ctx)

			// Create expired admin session
			accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

			req := th.NewDeletePromotionCodeRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
			clearTables(t, ctx)

			// Create standard user (not admin)
			accessToken := tu.SetupStandardUser(t, ctx, authCtx)

			req := th.NewDeletePromotionCodeRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("should return 401 when token is invalid", func(t *testing.T) {
			clearTables(t, ctx)

			req := th.NewDeletePromotionCodeRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", "invalid-token-12345")

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})

	t.Run("PromotionCodeBusinessLogic", func(t *testing.T) {
		t.Run("should validate promotion code meets minimum amount requirement", func(t *testing.T) {
			// Clean the database
			clearTables(t, ctx)

			// Create test coupon and promotion code with minimum amount
			testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
			td.InsertCoupon(t, ctx, testPool, testCoupon)

			testPromoCode := td.NewValidPromotionCodeWithMinAmount("MINAMT50", testCoupon.ID, 5000, "USD") // $50.00 minimum
			td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

			// Test with valid amount (above minimum)
			requestBody := domain.ValidatePromotionCodeRequest{
				Code:          "MINAMT50",
				OrderAmount:   &[]int{6000}[0], // $60.00 - above minimum
				OrderCurrency: &[]string{"USD"}[0],
			}

			req := th.NewValidatePromotionCodeRequest(t, ctx, testServerURL, requestBody)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response domain.ValidatePromotionCodeResponse

			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.True(t, response.Valid)
			assert.NotNil(t, response.PromotionCode)
		})

		t.Run("should handle promotion code with currency restrictions correctly", func(t *testing.T) {
			// Clean the database
			clearTables(t, ctx)

			// Create test coupon and promotion code with currency restrictions
			testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
			td.InsertCoupon(t, ctx, testPool, testCoupon)

			testPromoCode := td.NewValidPromotionCodeWithRestrictions("CURRENCY20", testCoupon.ID, []string{"USD", "EUR"})
			td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

			// Test with allowed currency
			requestBody := domain.ValidatePromotionCodeRequest{
				Code:          "CURRENCY20",
				OrderAmount:   &[]int{2000}[0],
				OrderCurrency: &[]string{"USD"}[0], // Allowed currency
			}

			req := th.NewValidatePromotionCodeRequest(t, ctx, testServerURL, requestBody)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response domain.ValidatePromotionCodeResponse

			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.True(t, response.Valid)
			assert.NotNil(t, response.PromotionCode)
		})

		t.Run("should handle promotion code expiry correctly", func(t *testing.T) {
			// Clean the database
			clearTables(t, ctx)

			// Create test coupon and promotion code with future expiry
			testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
			td.InsertCoupon(t, ctx, testPool, testCoupon)
			futureExpiry := time.Now().Add(30 * 24 * time.Hour) // 30 days from now
			testPromoCode := td.NewValidPromotionCodeWithExpiry("FUTURE20", testCoupon.ID, futureExpiry, 10)
			td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

			// Validate the promotion code
			requestBody := domain.ValidatePromotionCodeRequest{
				Code: "FUTURE20",
			}

			req := th.NewValidatePromotionCodeRequest(t, ctx, testServerURL, requestBody)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response domain.ValidatePromotionCodeResponse

			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.True(t, response.Valid)
			assert.NotNil(t, response.PromotionCode)
			assert.Equal(t, "FUTURE20", response.PromotionCode.PromotionCode.Code)
		})
	})
}
