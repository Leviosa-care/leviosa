package promotionCode_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
)

func TestValidatePromotionCode(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully validate active promotion code", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon and promotion code
		testCoupon := td.NewValidPercentOffCoupon("Valid Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		testPromoCode := td.NewValidPromotionCode("VALID20", testCoupon.ID)
		td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

		requestBody := domain.ValidatePromotionCodeRequest{
			Code:          "VALID20",
			OrderAmount:   &[]int{2000}[0], // $20.00
			OrderCurrency: &[]string{"USD"}[0],
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newValidatePromotionCodeRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.ValidatePromotionCodeResponse
		decodeJSONResponse(t, resp, &response)

		assert.True(t, response.Valid)
		assert.NotNil(t, response.PromotionCode)
		assert.Equal(t, "VALID20", response.PromotionCode.PromotionCode.Code)
		assert.Equal(t, testCoupon.Name, response.PromotionCode.Coupon.Name)
	})

	t.Run("should return invalid for non-existent promotion code", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)

		requestBody := domain.ValidatePromotionCodeRequest{
			Code: "NONEXISTENT",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newValidatePromotionCodeRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.ValidatePromotionCodeResponse
		decodeJSONResponse(t, resp, &response)

		assert.False(t, response.Valid)
		assert.Equal(t, "promotion code not found", response.Reason)
	})

	t.Run("should return invalid for inactive promotion code", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon and inactive promotion code
		testCoupon := td.NewValidPercentOffCoupon("Valid Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		testPromoCode := td.NewInactivePromotionCode("INACTIVE20", testCoupon.ID)
		td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

		requestBody := domain.ValidatePromotionCodeRequest{
			Code: "INACTIVE20",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newValidatePromotionCodeRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.ValidatePromotionCodeResponse
		decodeJSONResponse(t, resp, &response)

		assert.False(t, response.Valid)
		assert.Equal(t, "promotion code is not active", response.Reason)
	})

	t.Run("should return invalid for expired promotion code", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon and expired promotion code
		testCoupon := td.NewValidPercentOffCoupon("Valid Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		testPromoCode := td.NewExpiredPromotionCode("EXPIRED20", testCoupon.ID)
		td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

		requestBody := domain.ValidatePromotionCodeRequest{
			Code: "EXPIRED20",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newValidatePromotionCodeRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.ValidatePromotionCodeResponse
		decodeJSONResponse(t, resp, &response)

		assert.False(t, response.Valid)
		assert.Equal(t, "promotion code has expired", response.Reason)
	})

	t.Run("should return invalid for promotion code that reached redemption limit", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon and promotion code with limited redemptions
		testCoupon := td.NewValidPercentOffCoupon("Valid Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		testPromoCode := td.NewValidPromotionCodeWithRedemptionLimits("LIMITED20", testCoupon.ID, 5)
		// Set times redeemed to the limit
		testPromoCode.TimesRedeemed = 5
		td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

		requestBody := domain.ValidatePromotionCodeRequest{
			Code: "LIMITED20",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newValidatePromotionCodeRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.ValidatePromotionCodeResponse
		decodeJSONResponse(t, resp, &response)

		assert.False(t, response.Valid)
		assert.Equal(t, "promotion code has reached its redemption limit", response.Reason)
	})

	t.Run("should return invalid for order amount below minimum", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon and promotion code with minimum amount
		testCoupon := td.NewValidPercentOffCoupon("Valid Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		testPromoCode := td.NewValidPromotionCodeWithMinAmount("MINAMT20", testCoupon.ID, 5000, "USD") // $50.00 minimum
		td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

		requestBody := domain.ValidatePromotionCodeRequest{
			Code:          "MINAMT20",
			OrderAmount:   &[]int{3000}[0], // $30.00 - below minimum
			OrderCurrency: &[]string{"USD"}[0],
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newValidatePromotionCodeRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.ValidatePromotionCodeResponse
		decodeJSONResponse(t, resp, &response)

		assert.False(t, response.Valid)
		assert.Contains(t, response.Reason, "order amount must be at least 5000 USD")
	})

	t.Run("should return invalid for restricted currency", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon and promotion code with currency restrictions
		testCoupon := td.NewValidPercentOffCoupon("Valid Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		testPromoCode := td.NewValidPromotionCodeWithRestrictions("RESTRICTED20", testCoupon.ID, []string{"USD", "EUR"})
		td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

		requestBody := domain.ValidatePromotionCodeRequest{
			Code:          "RESTRICTED20",
			OrderAmount:   &[]int{2000}[0],
			OrderCurrency: &[]string{"GBP"}[0], // Not in allowed currencies
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newValidatePromotionCodeRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.ValidatePromotionCodeResponse
		decodeJSONResponse(t, resp, &response)

		assert.False(t, response.Valid)
		assert.Equal(t, "promotion code is not valid for currency GBP", response.Reason)
	})

	t.Run("should return 400 for invalid request body", func(t *testing.T) {
		req := newValidatePromotionCodeRequest(t, ctx, []byte("{invalid json"))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 415 for non-JSON content type", func(t *testing.T) {
		req := newValidatePromotionCodeRequest(t, ctx, []byte("not json"))
		req.Header.Set("Content-Type", "text/plain")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
	})
}
