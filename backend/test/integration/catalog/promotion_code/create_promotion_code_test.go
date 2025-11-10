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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreatePromotionCode TEST_PATH=test/integration/catalog/promotion_code/create_promotion_code_test.go

func TestCreatePromotionCode(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully create promotion code with valid admin token", func(t *testing.T) {
		// Clean the database to ensure isolation for this test
		clearTables(t, ctx)

		// Setup admin user authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create a test coupon first
		testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		requestBody := domain.CreatePromotionCodeRequest{
			CouponID:             testCoupon.ID.String(),
			Code:                 "TESTCODE20",
			MaxRedemptions:       &[]int{100}[0],
			FirstTimeTransaction: false,
			Metadata:             map[string]string{"test": "true"},
		}

		req := th.NewCreatePromotionCodeRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response struct {
			ID      string `json:"id"`
			Message string `json:"message"`
		}

		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Validate response
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, "Promotion code created successfully!", response.Message)

		// Verify the promotion code was actually created in the database
		promotionCodeID, err := uuid.Parse(response.ID)
		assert.NoError(t, err)

		promotionCode, err := td.GetPromotionCodeByID(t, ctx, promotionCodeID, testPool)
		assert.NoError(t, err)
		assert.Equal(t, "TESTCODE20", promotionCode.Code)
		assert.Equal(t, testCoupon.ID, promotionCode.CouponID)
		assert.True(t, promotionCode.Active)
		assert.Equal(t, 100, *promotionCode.MaxRedemptions)
		assert.Equal(t, 0, promotionCode.TimesRedeemed)
		assert.False(t, promotionCode.FirstTimeTransaction)
	})

	t.Run("should create promotion code with expiry and restrictions", func(t *testing.T) {
		// Clean the database
		clearTables(t, ctx)

		// Setup admin user authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create a test coupon first
		testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		expiryTime := time.Now().Add(30 * 24 * time.Hour) // 30 days from now
		restrictions := domain.PromotionCodeRestrictionsRequest{
			CurrencyOptions: []string{"USD", "EUR"},
		}

		const code = "EXPIRY20"

		requestBody := domain.CreatePromotionCodeRequest{
			CouponID:              testCoupon.ID.String(),
			Code:                  code,
			ExpiresAt:             &expiryTime,
			FirstTimeTransaction:  true,
			MinimumAmount:         &[]int{1000}[0], // $10.00
			MinimumAmountCurrency: &[]string{"USD"}[0],
			Restrictions:          &restrictions,
			Metadata:              map[string]string{"test": "true", "type": "expiring"},
		}

		req := th.NewCreatePromotionCodeRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response struct {
			ID      string `json:"id"`
			Message string `json:"message"`
		}

		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify the promotion code with restrictions was created
		promotionCodeID, err := uuid.Parse(response.ID)
		assert.NoError(t, err)

		promotionCode, err := td.GetPromotionCodeByID(t, ctx, promotionCodeID, testPool)
		assert.NoError(t, err)
		assert.Equal(t, code, promotionCode.Code)
		assert.True(t, promotionCode.FirstTimeTransaction)
		assert.Equal(t, 1000, *promotionCode.MinimumAmount)
		assert.Equal(t, "USD", *promotionCode.MinimumAmountCurrency)
		assert.NotNil(t, promotionCode.Restrictions)
		assert.Equal(t, []string{"USD", "EUR"}, promotionCode.Restrictions.CurrencyOptions)
	})

	t.Run("should return 400 for invalid coupon ID", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		requestBody := domain.CreatePromotionCodeRequest{
			CouponID: "invalid-uuid",
			Code:     "INVALIDCOUPON",
		}

		req := th.NewCreatePromotionCodeRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 for non-existent coupon", func(t *testing.T) {
		// Clean the database
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		nonExistentCouponID := uuid.New()
		requestBody := domain.CreatePromotionCodeRequest{
			CouponID: nonExistentCouponID.String(),
			Code:     "NONEXISTENT",
		}

		req := th.NewCreatePromotionCodeRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for missing required fields", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		requestBody := domain.CreatePromotionCodeRequest{
			// Missing required fields
		}

		req := th.NewCreatePromotionCodeRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 409 for duplicate promotion code", func(t *testing.T) {
		// Clean the database
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create a test coupon first
		testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		// Create the first promotion code
		existingPromoCode := td.NewValidPromotionCode("DUPLICATE", testCoupon.ID)
		td.InsertPromotionCode(t, ctx, testPool, existingPromoCode)

		// Try to create another promotion code with the same code
		requestBody := domain.CreatePromotionCodeRequest{
			CouponID: testCoupon.ID.String(),
			Code:     "DUPLICATE",
		}

		req := th.NewCreatePromotionCodeRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("should return 415 for non-JSON content type", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewCreatePromotionCodeRequest(t, ctx, testServerURL, map[string]string{"test": "data"}, accessToken)
		req.Header.Set("Content-Type", "text/plain")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		clearTables(t, ctx)

		requestBody := domain.CreatePromotionCodeRequest{
			CouponID: "00000000-0000-0000-0000-000000000000",
			Code:     "TESTCODE",
		}

		req := th.NewCreatePromotionCodeRequest(t, ctx, testServerURL, requestBody, "")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		clearTables(t, ctx)

		// Create expired admin session
		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		requestBody := domain.CreatePromotionCodeRequest{
			CouponID: "00000000-0000-0000-0000-000000000000",
			Code:     "TESTCODE",
		}

		req := th.NewCreatePromotionCodeRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		clearTables(t, ctx)

		// Create standard user (not admin)
		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		requestBody := domain.CreatePromotionCodeRequest{
			CouponID: "00000000-0000-0000-0000-000000000000",
			Code:     "TESTCODE",
		}

		req := th.NewCreatePromotionCodeRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		clearTables(t, ctx)

		requestBody := domain.CreatePromotionCodeRequest{
			CouponID: "00000000-0000-0000-0000-000000000000",
			Code:     "TESTCODE",
		}

		req := th.NewCreatePromotionCodeRequest(t, ctx, testServerURL, requestBody, "invalid-token-12345")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
