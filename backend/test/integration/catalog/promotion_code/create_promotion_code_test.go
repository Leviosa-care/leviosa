package promotionCode_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePromotionCode(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully create a new promotion code", func(t *testing.T) {
		// Clean the database to ensure isolation for this test
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

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
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreatePromotionCodeRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response struct {
			ID      string `json:"id"`
			Message string `json:"message"`
		}
		decodeJSONResponse(t, resp, &response)

		// Validate response
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, "Promotion code created successfully!", response.Message)

		// Verify the promotion code was actually created in the database
		promotionCodeID, err := uuid.Parse(response.ID)
		require.NoError(t, err)

		promotionCode, err := td.GetPromotionCodeByID(t, ctx, promotionCodeID, testPool)
		require.NoError(t, err)
		assert.Equal(t, "TESTCODE20", promotionCode.Code)
		assert.Equal(t, testCoupon.ID, promotionCode.CouponID)
		assert.True(t, promotionCode.Active)
		assert.Equal(t, 100, *promotionCode.MaxRedemptions)
		assert.Equal(t, 0, promotionCode.TimesRedeemed)
		assert.False(t, promotionCode.FirstTimeTransaction)
	})

	t.Run("should create promotion code with expiry and restrictions", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

		// Create a test coupon first
		testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		expiryTime := time.Now().Add(30 * 24 * time.Hour) // 30 days from now
		restrictions := domain.PromotionCodeRestrictionsRequest{
			CurrencyOptions: []string{"USD", "EUR"},
		}

		requestBody := domain.CreatePromotionCodeRequest{
			CouponID:              testCoupon.ID.String(),
			Code:                  "EXPIRY20",
			ExpiresAt:             &expiryTime,
			FirstTimeTransaction:  true,
			MinimumAmount:         &[]int{1000}[0], // $10.00
			MinimumAmountCurrency: &[]string{"USD"}[0],
			Restrictions:          &restrictions,
			Metadata:              map[string]string{"test": "true", "type": "expiring"},
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreatePromotionCodeRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response struct {
			ID      string `json:"id"`
			Message string `json:"message"`
		}
		decodeJSONResponse(t, resp, &response)

		// Verify the promotion code with restrictions was created
		promotionCodeID, err := uuid.Parse(response.ID)
		require.NoError(t, err)

		promotionCode, err := td.GetPromotionCodeByID(t, ctx, promotionCodeID, testPool)
		require.NoError(t, err)
		assert.Equal(t, "EXPIRY20", promotionCode.Code)
		assert.True(t, promotionCode.FirstTimeTransaction)
		assert.Equal(t, 1000, *promotionCode.MinimumAmount)
		assert.Equal(t, "USD", *promotionCode.MinimumAmountCurrency)
		assert.NotNil(t, promotionCode.Restrictions)
		assert.Equal(t, []string{"USD", "EUR"}, promotionCode.Restrictions.CurrencyOptions)
	})

	t.Run("should return 400 for invalid coupon ID", func(t *testing.T) {
		requestBody := domain.CreatePromotionCodeRequest{
			CouponID: "invalid-uuid",
			Code:     "INVALIDCOUPON",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreatePromotionCodeRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 for non-existent coupon", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

		nonExistentCouponID := uuid.New()
		requestBody := domain.CreatePromotionCodeRequest{
			CouponID: nonExistentCouponID.String(),
			Code:     "NONEXISTENT",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreatePromotionCodeRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for missing required fields", func(t *testing.T) {
		requestBody := domain.CreatePromotionCodeRequest{
			// Missing required fields
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreatePromotionCodeRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 409 for duplicate promotion code", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

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
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreatePromotionCodeRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("should return 415 for non-JSON content type", func(t *testing.T) {
		req := newCreatePromotionCodeRequest(t, ctx, []byte("not json"))
		req.Header.Set("Content-Type", "text/plain")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
	})
}
