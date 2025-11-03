package coupon_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
)

// make test-func TEST_NAME=TestGetCouponByID TEST_PATH=test/integration/catalog/coupon/get_coupon_test.go

func TestGetCouponByID(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get coupon by ID", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon
		testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		req := newGetCouponByIDRequest(t, ctx, testCoupon.ID.String())

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.CouponResponse
		decodeJSONResponse(t, resp, &response)

		assert.Equal(t, testCoupon.ID.String(), response.ID)
		assert.Equal(t, "Test Coupon", response.Name)
		assert.Equal(t, testCoupon.StripeCouponID, response.StripeCouponID)
		assert.NotNil(t, response.PercentOff)
		assert.Equal(t, 25.0, *response.PercentOff)
		assert.Nil(t, response.AmountOff)
		assert.True(t, response.Valid)
	})

	t.Run("should return 404 for non-existent coupon", func(t *testing.T) {
		req := newGetCouponByIDRequest(t, ctx, "00000000-0000-0000-0000-000000000000")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		req := newGetCouponByIDRequest(t, ctx, "invalid-uuid")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestGetCouponByStripeID(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get coupon by Stripe ID", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon
		testCoupon := td.NewValidAmountOffCoupon("USD Test Coupon", "USD")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		req := newGetCouponByStripeIDRequest(t, ctx, testCoupon.StripeCouponID)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.CouponResponse
		decodeJSONResponse(t, resp, &response)

		assert.Equal(t, testCoupon.ID.String(), response.ID)
		assert.Equal(t, "USD Test Coupon", response.Name)
		assert.Equal(t, testCoupon.StripeCouponID, response.StripeCouponID)
		assert.Nil(t, response.PercentOff)
		assert.NotNil(t, response.AmountOff)
		assert.Equal(t, 500, *response.AmountOff)
		assert.NotNil(t, response.Currency)
		assert.Equal(t, "USD", *response.Currency)
	})

	t.Run("should return 404 for non-existent Stripe ID", func(t *testing.T) {
		req := newGetCouponByStripeIDRequest(t, ctx, "coupon_nonexistent")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestGetAllCoupons(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get all coupons", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create multiple test coupons
		testCoupon1 := td.NewValidPercentOffCoupon("Coupon 1")
		testCoupon2 := td.NewValidAmountOffCoupon("Coupon 2", "USD")
		testCoupon3 := td.NewInvalidCoupon("Invalid Coupon")

		td.InsertCoupon(t, ctx, testPool, testCoupon1)
		td.InsertCoupon(t, ctx, testPool, testCoupon2)
		td.InsertCoupon(t, ctx, testPool, testCoupon3)

		req := newGetAllCouponsRequest(t, ctx)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []*domain.CouponResponse
		decodeJSONResponse(t, resp, &response)

		assert.Len(t, response, 3)

		// Check that we get all coupons (valid and invalid)
		couponNames := make([]string, len(response))
		for i, coupon := range response {
			couponNames[i] = coupon.Name
		}
		assert.Contains(t, couponNames, "Coupon 1")
		assert.Contains(t, couponNames, "Coupon 2")
		assert.Contains(t, couponNames, "Invalid Coupon")
	})

	t.Run("should return empty array when no coupons exist", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		req := newGetAllCouponsRequest(t, ctx)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []*domain.CouponResponse
		decodeJSONResponse(t, resp, &response)

		assert.Len(t, response, 0)
	})
}

func TestGetValidCoupons(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get only valid coupons (admin endpoint)", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create multiple test coupons
		testCoupon1 := td.NewValidPercentOffCoupon("Valid Coupon 1")
		testCoupon2 := td.NewValidAmountOffCoupon("Valid Coupon 2", "EUR")
		testCoupon3 := td.NewInvalidCoupon("Invalid Coupon")
		testCoupon4 := td.NewExpiredCoupon("Expired Coupon")

		td.InsertCoupon(t, ctx, testPool, testCoupon1)
		td.InsertCoupon(t, ctx, testPool, testCoupon2)
		td.InsertCoupon(t, ctx, testPool, testCoupon3)
		td.InsertCoupon(t, ctx, testPool, testCoupon4)

		req := newGetValidCouponsAdminRequest(t, ctx)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []*domain.CouponResponse
		decodeJSONResponse(t, resp, &response)

		// Should only return valid coupons that haven't expired and are not at redemption limit
		assert.Len(t, response, 2)

		for _, coupon := range response {
			assert.True(t, coupon.Valid)
			assert.Contains(t, []string{"Valid Coupon 1", "Valid Coupon 2"}, coupon.Name)
		}
	})

	t.Run("should successfully get only valid coupons (public endpoint)", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create multiple test coupons
		testCoupon1 := td.NewValidPercentOffCoupon("Public Valid 1")
		testCoupon2 := td.NewInvalidCoupon("Public Invalid")

		td.InsertCoupon(t, ctx, testPool, testCoupon1)
		td.InsertCoupon(t, ctx, testPool, testCoupon2)

		req := newGetValidCouponsPublicRequest(t, ctx)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []*domain.CouponResponse
		decodeJSONResponse(t, resp, &response)

		assert.Len(t, response, 1)
		assert.True(t, response[0].Valid)
		assert.Equal(t, "Public Valid 1", response[0].Name)
	})

	t.Run("should return empty array when no valid coupons exist", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create only invalid coupons
		testCoupon1 := td.NewInvalidCoupon("Invalid 1")
		testCoupon2 := td.NewExpiredCoupon("Expired 1")

		td.InsertCoupon(t, ctx, testPool, testCoupon1)
		td.InsertCoupon(t, ctx, testPool, testCoupon2)

		req := newGetValidCouponsAdminRequest(t, ctx)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []*domain.CouponResponse
		decodeJSONResponse(t, resp, &response)

		assert.Len(t, response, 0)
	})
}

func TestCouponResponseFormat(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return complete coupon response structure", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create a comprehensive test coupon
		testCoupon := td.NewValidRepeatingCoupon("Complete Coupon", 6)
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		req := newGetCouponByIDRequest(t, ctx, testCoupon.ID.String())

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.CouponResponse
		decodeJSONResponse(t, resp, &response)

		// Verify all expected fields are present
		assert.Equal(t, testCoupon.ID.String(), response.ID)
		assert.Equal(t, testCoupon.StripeCouponID, response.StripeCouponID)
		assert.Equal(t, "Complete Coupon", response.Name)
		assert.NotNil(t, response.PercentOff)
		assert.Equal(t, "repeating", response.Duration)
		assert.NotNil(t, response.DurationInMonths)
		assert.Equal(t, 6, *response.DurationInMonths)
		assert.Equal(t, 0, response.TimesRedeemed)
		assert.True(t, response.Valid)
		assert.NotNil(t, response.CreatedAt)
		assert.NotNil(t, response.UpdatedAt)
		assert.NotNil(t, response.Metadata)
		assert.Equal(t, "true", response.Metadata["test"])
	})
}
