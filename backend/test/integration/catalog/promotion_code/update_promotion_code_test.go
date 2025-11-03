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
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdatePromotionCode TEST_PATH=test/integration/catalog/promotion_code/update_promotion_code_test.go

func TestUpdatePromotionCode(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully update promotion code metadata", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

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
		jsonBody, _ := json.Marshal(updateRequest)

		req := newUpdatePromotionCodeRequest(t, ctx, testPromoCode.ID.String(), jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Message string `json:"message"`
		}
		decodeJSONResponse(t, resp, &response)
		assert.Equal(t, "Promotion code updated successfully!", response.Message)

		// Verify the update in database
		updatedPromoCode, err := td.GetPromotionCodeByID(t, ctx, testPromoCode.ID, testPool)
		assert.NoError(t, err)
		assert.Equal(t, "true", updatedPromoCode.Metadata["updated"])
		assert.Equal(t, "false", updatedPromoCode.Metadata["test"])
		assert.Equal(t, "Updated promotion code", updatedPromoCode.Metadata["description"])
	})

	t.Run("should return 404 for non-existent promotion code", func(t *testing.T) {
		updateRequest := domain.UpdatePromotionCodeRequest{
			Metadata: map[string]string{"test": "value"},
		}
		jsonBody, _ := json.Marshal(updateRequest)

		req := newUpdatePromotionCodeRequest(t, ctx, "00000000-0000-0000-0000-000000000000", jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		updateRequest := domain.UpdatePromotionCodeRequest{
			Metadata: map[string]string{"test": "value"},
		}
		jsonBody, _ := json.Marshal(updateRequest)

		req := newUpdatePromotionCodeRequest(t, ctx, "invalid-uuid", jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 415 for non-JSON content type", func(t *testing.T) {
		req := newUpdatePromotionCodeRequest(t, ctx, "00000000-0000-0000-0000-000000000000", []byte("not json"))
		req.Header.Set("Content-Type", "text/plain")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
	})
}

// make test-func TEST_NAME=TestDeactivatePromotionCode TEST_PATH=test/integration/catalog/promotion_code/update_promotion_code_test.go

func TestDeactivatePromotionCode(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully deactivate promotion code", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon and promotion code
		testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		testPromoCode := td.NewValidPromotionCode("DEACTIVATE20", testCoupon.ID)
		td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

		// Verify it's initially active
		assert.True(t, td.GetPromotionCodeActiveStatus(t, ctx, testPromoCode.ID, testPool))

		req := newDeactivatePromotionCodeRequest(t, ctx, testPromoCode.ID.String())

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Message string `json:"message"`
		}
		decodeJSONResponse(t, resp, &response)
		assert.Equal(t, "Promotion code deactivated successfully!", response.Message)

		// Verify the promotion code is deactivated in database
		assert.False(t, td.GetPromotionCodeActiveStatus(t, ctx, testPromoCode.ID, testPool))
	})

	t.Run("should return 404 for non-existent promotion code", func(t *testing.T) {
		req := newDeactivatePromotionCodeRequest(t, ctx, "00000000-0000-0000-0000-000000000000")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		req := newDeactivatePromotionCodeRequest(t, ctx, "invalid-uuid")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// make test-func TEST_NAME=TestDeletePromotionCode TEST_PATH=test/integration/catalog/promotion_code/update_promotion_code_test.go

func TestDeletePromotionCode(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully delete promotion code", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon and promotion code
		testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		testPromoCode := td.NewValidPromotionCode("DELETE20", testCoupon.ID)
		td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

		// Verify it exists
		foundPromoCode := td.GetPromotionCodeByIDOrNil(t, ctx, testPromoCode.ID, testPool)
		assert.NotNil(t, foundPromoCode)

		req := newDeletePromotionCodeRequest(t, ctx, testPromoCode.ID.String())

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Message string `json:"message"`
		}
		decodeJSONResponse(t, resp, &response)
		assert.Equal(t, "Promotion code deleted successfully!", response.Message)

		// Verify the promotion code is deleted from database
		deletedPromoCode := td.GetPromotionCodeByIDOrNil(t, ctx, testPromoCode.ID, testPool)
		assert.Nil(t, deletedPromoCode)
	})

	t.Run("should return 404 for non-existent promotion code", func(t *testing.T) {
		req := newDeletePromotionCodeRequest(t, ctx, "00000000-0000-0000-0000-000000000000")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		req := newDeletePromotionCodeRequest(t, ctx, "invalid-uuid")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// make test-func TEST_NAME=TestPromotionCodeBusinessLogic TEST_PATH=test/integration/catalog/promotion_code/update_promotion_code_test.go

func TestPromotionCodeBusinessLogic(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should validate promotion code meets minimum amount requirement", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

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
	})

	t.Run("should handle promotion code with currency restrictions correctly", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

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
	})

	t.Run("should handle promotion code expiry correctly", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

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
		assert.Equal(t, "FUTURE20", response.PromotionCode.PromotionCode.Code)
	})
}
