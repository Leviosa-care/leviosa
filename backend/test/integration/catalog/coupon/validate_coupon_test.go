package coupon_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestValidateCoupon TEST_PATH=test/integration/catalog/coupon/validate_coupon_test.go

func TestValidateCoupon(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully validate a valid coupon", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create a valid test coupon
		testCoupon := td.NewValidPercentOffCoupon("Valid Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		requestBody := struct {
			StripeCouponID string `json:"stripeCouponId"`
		}{
			StripeCouponID: testCoupon.StripeCouponID,
		}

		req := th.NewValidateCouponRequest(t, ctx, testServerURL, requestBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Valid  bool        `json:"valid"`
			Coupon interface{} `json:"coupon"`
		}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.Valid)
		assert.NotNil(t, response.Coupon)
	})

	t.Run("should return invalid for non-existent coupon", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		requestBody := struct {
			StripeCouponID string `json:"stripeCouponId"`
		}{
			StripeCouponID: "coupon_nonexistent",
		}

		req := th.NewValidateCouponRequest(t, ctx, testServerURL, requestBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Valid  bool   `json:"valid"`
			Reason string `json:"reason"`
		}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.Valid)
		assert.Equal(t, "coupon not found", response.Reason)
	})

	t.Run("should return invalid for invalid coupon", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create an invalid test coupon
		testCoupon := td.NewInvalidCoupon("Invalid Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		requestBody := struct {
			StripeCouponID string `json:"stripeCouponId"`
		}{
			StripeCouponID: testCoupon.StripeCouponID,
		}

		req := th.NewValidateCouponRequest(t, ctx, testServerURL, requestBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Valid  bool   `json:"valid"`
			Reason string `json:"reason"`
		}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.Valid)
		assert.Contains(t, response.Reason, "coupon is not valid")
	})

	t.Run("should return invalid for expired coupon", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create an expired test coupon
		testCoupon := td.NewExpiredCoupon("Expired Test Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		requestBody := struct {
			StripeCouponID string `json:"stripeCouponId"`
		}{
			StripeCouponID: testCoupon.StripeCouponID,
		}

		req := th.NewValidateCouponRequest(t, ctx, testServerURL, requestBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Valid  bool   `json:"valid"`
			Reason string `json:"reason"`
		}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.Valid)
		assert.Contains(t, response.Reason, "coupon has expired")
	})

	t.Run("should return invalid for coupon that reached redemption limit", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon that has reached its redemption limit
		testCoupon := td.NewValidPercentOffCouponWithRedemptionLimits("Limited Coupon", 5)
		// Set times redeemed to the limit
		testCoupon.TimesRedeemed = 5
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		requestBody := struct {
			StripeCouponID string `json:"stripeCouponId"`
		}{
			StripeCouponID: testCoupon.StripeCouponID,
		}

		req := th.NewValidateCouponRequest(t, ctx, testServerURL, requestBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Valid  bool   `json:"valid"`
			Reason string `json:"reason"`
		}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.Valid)
		assert.Contains(t, response.Reason, "coupon has reached its redemption limit")
	})

	t.Run("should return 400 for empty stripe coupon ID", func(t *testing.T) {
		requestBody := struct {
			StripeCouponID string `json:"stripeCouponId"`
		}{
			StripeCouponID: "",
		}

		req := th.NewValidateCouponRequest(t, ctx, testServerURL, requestBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Valid  bool   `json:"valid"`
			Reason string `json:"reason"`
		}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.Valid)
		assert.Contains(t, response.Reason, "stripe coupon ID cannot be empty")
	})

	t.Run("should return 400 for invalid request body", func(t *testing.T) {
		req := th.NewValidateCouponRequest(t, ctx, testServerURL, ("{invalid json"))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 415 for non-JSON content type", func(t *testing.T) {
		req := th.NewValidateCouponRequest(t, ctx, testServerURL, ("not json"))
		req.Header.Set("Content-Type", "text/plain")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
	})
}
