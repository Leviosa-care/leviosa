package promotionCode_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetPromotionCode TEST_PATH=test/integration/catalog/promotion_code/get_promotion_code_test.go

func TestGetPromotionCode(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("GetPromotionCodeByID", func(t *testing.T) {
		t.Run("should successfully get promotion code by ID with valid admin token", func(t *testing.T) {
			// Clean the database
			clearTables(t, ctx)

			// Setup admin user authentication
			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			// Create test coupon and promotion code
			testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
			td.InsertCoupon(t, ctx, testPool, testCoupon)

			testPromoCode := td.NewValidPromotionCode("GETBYID20", testCoupon.ID)
			td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

			req := th.NewGetPromotionCodeByIDRequest(t, ctx, testServerURL, testPromoCode.ID.String(), accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			response := th.ParseSimplePromotionCodeResponse(t, resp)

			assert.Equal(t, testPromoCode.ID.String(), response.ID)
			assert.Equal(t, "GETBYID20", response.Code)
			assert.Equal(t, testCoupon.ID.String(), response.CouponID)
			assert.True(t, response.Active)
		})

		t.Run("should return 404 for non-existent promotion code", func(t *testing.T) {
			clearTables(t, ctx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			req := th.NewGetPromotionCodeByIDRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("should return 400 for invalid UUID", func(t *testing.T) {
			clearTables(t, ctx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			req := th.NewGetPromotionCodeByIDRequest(t, ctx, testServerURL, "invalid-uuid", accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return 401 when access token is missing", func(t *testing.T) {
			clearTables(t, ctx)

			req := th.NewGetPromotionCodeByIDRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", "")

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 401 when session is expired", func(t *testing.T) {
			clearTables(t, ctx)

			// Create expired admin session
			accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

			req := th.NewGetPromotionCodeByIDRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
			clearTables(t, ctx)

			// Create standard user (not admin)
			accessToken := tu.SetupStandardUser(t, ctx, authCtx)

			req := th.NewGetPromotionCodeByIDRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("should return 401 when token is invalid", func(t *testing.T) {
			clearTables(t, ctx)

			req := th.NewGetPromotionCodeByIDRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", "invalid-token-12345")

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})

	t.Run("GetPromotionCodeByCode", func(t *testing.T) {
		t.Run("should successfully get promotion code by code with valid admin token", func(t *testing.T) {
			// Clean the database
			clearTables(t, ctx)

			// Setup admin user authentication
			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			// Create test coupon and promotion code
			testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
			td.InsertCoupon(t, ctx, testPool, testCoupon)

			testPromoCode := td.NewValidPromotionCode("GETBYCODE20", testCoupon.ID)
			td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

			req := th.NewGetPromotionCodeByCodeRequest(t, ctx, testServerURL, "GETBYCODE20", accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			response := th.ParseSimplePromotionCodeResponse(t, resp)

			assert.Equal(t, testPromoCode.ID.String(), response.ID)
			assert.Equal(t, "GETBYCODE20", response.Code)
			assert.Equal(t, testCoupon.ID.String(), response.CouponID)
		})

		t.Run("should return 404 for non-existent code", func(t *testing.T) {
			clearTables(t, ctx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			req := th.NewGetPromotionCodeByCodeRequest(t, ctx, testServerURL, "NONEXISTENT", accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("should return 401 when access token is missing", func(t *testing.T) {
			clearTables(t, ctx)

			req := th.NewGetPromotionCodeByCodeRequest(t, ctx, testServerURL, "NONEXISTENT", "")

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 401 when session is expired", func(t *testing.T) {
			clearTables(t, ctx)

			// Create expired admin session
			accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

			req := th.NewGetPromotionCodeByCodeRequest(t, ctx, testServerURL, "NONEXISTENT", accessToken)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
			clearTables(t, ctx)

			// Create standard user (not admin)
			accessToken := tu.SetupStandardUser(t, ctx, authCtx)

			req := th.NewGetPromotionCodeByCodeRequest(t, ctx, testServerURL, "NONEXISTENT", accessToken)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("should return 401 when token is invalid", func(t *testing.T) {
			clearTables(t, ctx)

			req := th.NewGetPromotionCodeByCodeRequest(t, ctx, testServerURL, "NONEXISTENT", "invalid-token-12345")

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})

	t.Run("GetPromotionCodeWithCoupon", func(t *testing.T) {
		t.Run("should successfully get promotion code with coupon details", func(t *testing.T) {
			// Clean the database
			clearTables(t, ctx)

			// Create test coupon and promotion code
			testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
			td.InsertCoupon(t, ctx, testPool, testCoupon)

			testPromoCode := td.NewValidPromotionCode("WITHCOUPON20", testCoupon.ID)
			td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

			req := th.NewGetPromotionCodeWithCouponRequest(t, ctx, testServerURL, "WITHCOUPON20")

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			response := th.ParsePromotionCodeResponse(t, resp)

			// Check promotion code details
			assert.Equal(t, testPromoCode.ID.String(), response.PromotionCode.ID)
			assert.Equal(t, "WITHCOUPON20", response.PromotionCode.Code)

			// Check coupon details
			assert.Equal(t, testCoupon.ID.String(), response.Coupon.ID)
			assert.Equal(t, "Test Coupon", response.Coupon.Name)
			assert.NotNil(t, response.Coupon.PercentOff)
			assert.Equal(t, 25.0, *response.Coupon.PercentOff)
		})

		t.Run("should return 404 for non-existent code", func(t *testing.T) {
			req := th.NewGetPromotionCodeWithCouponRequest(t, ctx, testServerURL, "NONEXISTENT")

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})

	t.Run("GetAllPromotionCodes", func(t *testing.T) {
		t.Run("should successfully get all promotion codes with valid admin token", func(t *testing.T) {
			// Clean the database
			clearTables(t, ctx)

			// Setup admin user authentication
			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			// Create test coupon and multiple promotion codes
			testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
			td.InsertCoupon(t, ctx, testPool, testCoupon)

			testPromoCode1 := td.NewValidPromotionCode("ALL1", testCoupon.ID)
			testPromoCode2 := td.NewValidPromotionCode("ALL2", testCoupon.ID)
			testPromoCode3 := td.NewInactivePromotionCode("ALL3", testCoupon.ID)

			td.InsertPromotionCode(t, ctx, testPool, testPromoCode1)
			td.InsertPromotionCode(t, ctx, testPool, testPromoCode2)
			td.InsertPromotionCode(t, ctx, testPool, testPromoCode3)

			req := th.NewGetAllPromotionCodesRequest(t, ctx, testServerURL, accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			response := th.ParseSimplePromotionCodesResponse(t, resp)

			assert.Len(t, response, 3)

			// Check that we get all promotion codes (active and inactive)
			codes := make([]string, len(response))
			for i, promoCode := range response {
				codes[i] = promoCode.Code
			}
			assert.Contains(t, codes, "ALL1")
			assert.Contains(t, codes, "ALL2")
			assert.Contains(t, codes, "ALL3")
		})

		t.Run("should return empty array when no promotion codes exist", func(t *testing.T) {
			// Clean the database
			clearTables(t, ctx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			req := th.NewGetAllPromotionCodesRequest(t, ctx, testServerURL, accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			response := th.ParseSimplePromotionCodesResponse(t, resp)

			assert.Len(t, response, 0)
		})

		t.Run("should return 401 when access token is missing", func(t *testing.T) {
			clearTables(t, ctx)

			req := th.NewGetAllPromotionCodesRequest(t, ctx, testServerURL, "")

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 401 when session is expired", func(t *testing.T) {
			clearTables(t, ctx)

			// Create expired admin session
			accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

			req := th.NewGetAllPromotionCodesRequest(t, ctx, testServerURL, accessToken)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
			clearTables(t, ctx)

			// Create standard user (not admin)
			accessToken := tu.SetupStandardUser(t, ctx, authCtx)

			req := th.NewGetAllPromotionCodesRequest(t, ctx, testServerURL, accessToken)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("should return 401 when token is invalid", func(t *testing.T) {
			clearTables(t, ctx)

			req := th.NewGetAllPromotionCodesRequest(t, ctx, testServerURL, "invalid-token-12345")

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})

	t.Run("GetActivePromotionCodes", func(t *testing.T) {
		t.Run("should successfully get only active promotion codes with valid admin token", func(t *testing.T) {
			// Clean the database
			clearTables(t, ctx)

			// Setup admin user authentication
			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			// Create test coupon and multiple promotion codes
			testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
			td.InsertCoupon(t, ctx, testPool, testCoupon)

			testPromoCode1 := td.NewValidPromotionCode("ACTIVE1", testCoupon.ID)
			testPromoCode2 := td.NewValidPromotionCode("ACTIVE2", testCoupon.ID)
			testPromoCode3 := td.NewInactivePromotionCode("INACTIVE1", testCoupon.ID)

			td.InsertPromotionCode(t, ctx, testPool, testPromoCode1)
			td.InsertPromotionCode(t, ctx, testPool, testPromoCode2)
			td.InsertPromotionCode(t, ctx, testPool, testPromoCode3)

			req := th.NewGetActivePromotionCodesRequest(t, ctx, testServerURL, accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			response := th.ParseSimplePromotionCodesResponse(t, resp)

			assert.Len(t, response, 2)

			// Check that we only get active promotion codes
			for _, promoCode := range response {
				assert.True(t, promoCode.Active)
				assert.Contains(t, []string{"ACTIVE1", "ACTIVE2"}, promoCode.Code)
			}
		})

		t.Run("should return empty array when no active promotion codes exist", func(t *testing.T) {
			// Clean the database
			clearTables(t, ctx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			// Create test coupon and only inactive promotion codes
			testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
			td.InsertCoupon(t, ctx, testPool, testCoupon)

			testPromoCode := td.NewInactivePromotionCode("INACTIVE1", testCoupon.ID)
			td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

			req := th.NewGetActivePromotionCodesRequest(t, ctx, testServerURL, accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			response := th.ParseSimplePromotionCodesResponse(t, resp)

			assert.Len(t, response, 0)
		})

		t.Run("should return 401 when access token is missing", func(t *testing.T) {
			clearTables(t, ctx)

			req := th.NewGetActivePromotionCodesRequest(t, ctx, testServerURL, "")

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 401 when session is expired", func(t *testing.T) {
			clearTables(t, ctx)

			// Create expired admin session
			accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

			req := th.NewGetActivePromotionCodesRequest(t, ctx, testServerURL, accessToken)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
			clearTables(t, ctx)

			// Create standard user (not admin)
			accessToken := tu.SetupStandardUser(t, ctx, authCtx)

			req := th.NewGetActivePromotionCodesRequest(t, ctx, testServerURL, accessToken)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("should return 401 when token is invalid", func(t *testing.T) {
			clearTables(t, ctx)

			req := th.NewGetActivePromotionCodesRequest(t, ctx, testServerURL, "invalid-token-12345")

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})
}

// Helper function to decode JSON response body.
