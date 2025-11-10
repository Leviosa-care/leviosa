package coupon_test

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

// make test-func TEST_NAME=TestGetCouponByID TEST_PATH=test/integration/catalog/coupon/get_coupon_test.go

func TestGetCouponByID(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get coupon by ID with valid admin token", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		testCoupon := td.NewValidPercentOffCoupon("Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		req := th.NewGetCouponByIDRequest(t, ctx, testServerURL, testCoupon.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response := th.ParseCouponResponse(t, resp)

		assert.Equal(t, testCoupon.ID.String(), response.ID)
		assert.Equal(t, "Test Coupon", response.Name)
		assert.Equal(t, testCoupon.StripeCouponID, response.StripeCouponID)
		assert.NotNil(t, response.PercentOff)
		assert.Equal(t, 25.0, *response.PercentOff)
		assert.Nil(t, response.AmountOff)
		assert.True(t, response.Valid)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		req := th.NewGetCouponByIDRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", "")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		req := th.NewGetCouponByIDRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		req := th.NewGetCouponByIDRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		req := th.NewGetCouponByIDRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", "invalid-token-12345")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 404 for non-existent coupon", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewGetCouponByIDRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewGetCouponByIDRequest(t, ctx, testServerURL, "invalid-uuid", accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestGetCouponByStripeID(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get coupon by Stripe ID with valid admin token", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		testCoupon := td.NewValidAmountOffCoupon("USD Test Coupon", "USD")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		req := th.NewGetCouponByStripeIDRequest(t, ctx, testServerURL, testCoupon.StripeCouponID, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response := th.ParseCouponResponse(t, resp)

		assert.Equal(t, testCoupon.ID.String(), response.ID)
		assert.Equal(t, "USD Test Coupon", response.Name)
		assert.Equal(t, testCoupon.StripeCouponID, response.StripeCouponID)
		assert.Nil(t, response.PercentOff)
		assert.NotNil(t, response.AmountOff)
		assert.Equal(t, 500, *response.AmountOff)
		assert.NotNil(t, response.Currency)
		assert.Equal(t, "USD", *response.Currency)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		req := th.NewGetCouponByStripeIDRequest(t, ctx, testServerURL, "coupon_nonexistent", "")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		req := th.NewGetCouponByStripeIDRequest(t, ctx, testServerURL, "coupon_nonexistent", accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		req := th.NewGetCouponByStripeIDRequest(t, ctx, testServerURL, "coupon_nonexistent", accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		req := th.NewGetCouponByStripeIDRequest(t, ctx, testServerURL, "coupon_nonexistent", "invalid-token-12345")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 404 for non-existent Stripe ID", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewGetCouponByStripeIDRequest(t, ctx, testServerURL, "coupon_nonexistent", accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestGetAllCoupons(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get all coupons with valid admin token", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		testCoupon1 := td.NewValidPercentOffCoupon("Coupon 1")
		testCoupon2 := td.NewValidAmountOffCoupon("Coupon 2", "USD")
		testCoupon3 := td.NewInvalidCoupon("Invalid Coupon")

		td.InsertCoupon(t, ctx, testPool, testCoupon1)
		td.InsertCoupon(t, ctx, testPool, testCoupon2)
		td.InsertCoupon(t, ctx, testPool, testCoupon3)

		req := th.NewGetAllCouponsRequest(t, ctx, testServerURL, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response := th.ParseCouponsResponse(t, resp)

		assert.Len(t, response, 3)

		couponNames := make([]string, len(response))
		for i, coupon := range response {
			couponNames[i] = coupon.Name
		}
		assert.Contains(t, couponNames, "Coupon 1")
		assert.Contains(t, couponNames, "Coupon 2")
		assert.Contains(t, couponNames, "Invalid Coupon")
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		req := th.NewGetAllCouponsRequest(t, ctx, testServerURL, "")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		req := th.NewGetAllCouponsRequest(t, ctx, testServerURL, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		req := th.NewGetAllCouponsRequest(t, ctx, testServerURL, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		req := th.NewGetAllCouponsRequest(t, ctx, testServerURL, "invalid-token-12345")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return empty array when no coupons exist", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewGetAllCouponsRequest(t, ctx, testServerURL, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response := th.ParseCouponsResponse(t, resp)

		assert.Len(t, response, 0)
	})
}

func TestGetValidCoupons(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get only valid coupons (public endpoint)", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)

		testCoupon1 := td.NewValidPercentOffCoupon("Valid Coupon 1")
		testCoupon2 := td.NewValidAmountOffCoupon("Valid Coupon 2", "EUR")
		testCoupon3 := td.NewInvalidCoupon("Invalid Coupon")
		testCoupon4 := td.NewExpiredCoupon("Expired Coupon")

		td.InsertCoupon(t, ctx, testPool, testCoupon1)
		td.InsertCoupon(t, ctx, testPool, testCoupon2)
		td.InsertCoupon(t, ctx, testPool, testCoupon3)
		td.InsertCoupon(t, ctx, testPool, testCoupon4)

		req := th.NewGetValidCouponsRequest(t, ctx, testServerURL)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response := th.ParseCouponsResponse(t, resp)

		assert.Len(t, response, 2)

		for _, coupon := range response {
			assert.True(t, coupon.Valid)
			assert.Contains(t, []string{"Valid Coupon 1", "Valid Coupon 2"}, coupon.Name)
		}
	})

	t.Run("should return empty array when no valid coupons exist", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)

		testCoupon1 := td.NewInvalidCoupon("Invalid 1")
		testCoupon2 := td.NewExpiredCoupon("Expired 1")

		td.InsertCoupon(t, ctx, testPool, testCoupon1)
		td.InsertCoupon(t, ctx, testPool, testCoupon2)

		req := th.NewGetValidCouponsRequest(t, ctx, testServerURL)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response := th.ParseCouponsResponse(t, resp)

		assert.Len(t, response, 0)
	})
}

func TestCouponResponseFormat(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return complete coupon response structure", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		testCoupon := td.NewValidRepeatingCoupon("Complete Coupon", 6)
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		req := th.NewGetCouponByIDRequest(t, ctx, testServerURL, testCoupon.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response := th.ParseCouponResponse(t, resp)

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
