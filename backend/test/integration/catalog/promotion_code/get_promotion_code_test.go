package promotionCode_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
)

// make test-func TEST_NAME=TestGetPromotionCodeByID TEST_PATH=test/integration/catalog/promotion_code/get_promotion_code_test.go

func TestGetPromotionCodeByID(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get promotion code by ID", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon and promotion code
		testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		testPromoCode := td.NewValidPromotionCode("GETBYID20", testCoupon.ID)
		td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

		req := newGetPromotionCodeByIDRequest(t, ctx, testPromoCode.ID.String())

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.PromotionCodeResponse
		decodeJSONResponse(t, resp, &response)

		assert.Equal(t, testPromoCode.ID.String(), response.ID)
		assert.Equal(t, "GETBYID20", response.Code)
		assert.Equal(t, testCoupon.ID.String(), response.CouponID)
		assert.True(t, response.Active)
	})

	t.Run("should return 404 for non-existent promotion code", func(t *testing.T) {
		req := newGetPromotionCodeByIDRequest(t, ctx, "00000000-0000-0000-0000-000000000000")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		req := newGetPromotionCodeByIDRequest(t, ctx, "invalid-uuid")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// make test-func TEST_NAME=TestGetPromotionCodeByCode TEST_PATH=test/integration/catalog/promotion_code/get_promotion_code_test.go

func TestGetPromotionCodeByCode(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get promotion code by code", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon and promotion code
		testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		testPromoCode := td.NewValidPromotionCode("GETBYCODE20", testCoupon.ID)
		td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

		req := newGetPromotionCodeByCodeRequest(t, ctx, "GETBYCODE20")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.PromotionCodeResponse
		decodeJSONResponse(t, resp, &response)

		assert.Equal(t, testPromoCode.ID.String(), response.ID)
		assert.Equal(t, "GETBYCODE20", response.Code)
		assert.Equal(t, testCoupon.ID.String(), response.CouponID)
	})

	t.Run("should return 404 for non-existent code", func(t *testing.T) {
		req := newGetPromotionCodeByCodeRequest(t, ctx, "NONEXISTENT")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

// make test-func TEST_NAME=TestGetPromotionCodeWithCoupon TEST_PATH=test/integration/catalog/promotion_code/get_promotion_code_test.go

func TestGetPromotionCodeWithCoupon(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get promotion code with coupon details", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon and promotion code
		testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		testPromoCode := td.NewValidPromotionCode("WITHCOUPON20", testCoupon.ID)
		td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

		req := newGetPromotionCodeWithCouponRequest(t, ctx, "WITHCOUPON20")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.PromotionCodeWithCouponResponse
		decodeJSONResponse(t, resp, &response)

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
		req := newGetPromotionCodeWithCouponRequest(t, ctx, "NONEXISTENT")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

// make test-func TEST_NAME=TestGetAllPromotionCodes TEST_PATH=test/integration/catalog/promotion_code/get_promotion_code_test.go

func TestGetAllPromotionCodes(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get all promotion codes", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon and multiple promotion codes
		testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		testPromoCode1 := td.NewValidPromotionCode("ALL1", testCoupon.ID)
		testPromoCode2 := td.NewValidPromotionCode("ALL2", testCoupon.ID)
		testPromoCode3 := td.NewInactivePromotionCode("ALL3", testCoupon.ID)

		td.InsertPromotionCode(t, ctx, testPool, testPromoCode1)
		td.InsertPromotionCode(t, ctx, testPool, testPromoCode2)
		td.InsertPromotionCode(t, ctx, testPool, testPromoCode3)

		req := newGetAllPromotionCodesRequest(t, ctx)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []*domain.PromotionCodeResponse
		decodeJSONResponse(t, resp, &response)

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
		td.ClearPromotionCodesTable(t, ctx, testPool)

		req := newGetAllPromotionCodesRequest(t, ctx)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []*domain.PromotionCodeResponse
		decodeJSONResponse(t, resp, &response)

		assert.Len(t, response, 0)
	})
}

// make test-func TEST_NAME=TestGetActivePromotionCodes TEST_PATH=test/integration/catalog/promotion_code/get_promotion_code_test.go

func TestGetActivePromotionCodes(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get only active promotion codes", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon and multiple promotion codes
		testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		testPromoCode1 := td.NewValidPromotionCode("ACTIVE1", testCoupon.ID)
		testPromoCode2 := td.NewValidPromotionCode("ACTIVE2", testCoupon.ID)
		testPromoCode3 := td.NewInactivePromotionCode("INACTIVE1", testCoupon.ID)

		td.InsertPromotionCode(t, ctx, testPool, testPromoCode1)
		td.InsertPromotionCode(t, ctx, testPool, testPromoCode2)
		td.InsertPromotionCode(t, ctx, testPool, testPromoCode3)

		req := newGetActivePromotionCodesRequest(t, ctx)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []*domain.PromotionCodeResponse
		decodeJSONResponse(t, resp, &response)

		assert.Len(t, response, 2)

		// Check that we only get active promotion codes
		for _, promoCode := range response {
			assert.True(t, promoCode.Active)
			assert.Contains(t, []string{"ACTIVE1", "ACTIVE2"}, promoCode.Code)
		}
	})

	t.Run("should return empty array when no active promotion codes exist", func(t *testing.T) {
		// Clean the database
		td.ClearPromotionCodesTable(t, ctx, testPool)
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon and only inactive promotion codes
		testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		testPromoCode := td.NewInactivePromotionCode("INACTIVE1", testCoupon.ID)
		td.InsertPromotionCode(t, ctx, testPool, testPromoCode)

		req := newGetActivePromotionCodesRequest(t, ctx)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []*domain.PromotionCodeResponse
		decodeJSONResponse(t, resp, &response)

		assert.Len(t, response, 0)
	})
}
