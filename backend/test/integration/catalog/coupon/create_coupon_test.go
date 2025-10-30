package coupon_test

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

func TestCreateCoupon(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully create a new percent-off coupon", func(t *testing.T) {
		// Clean the database to ensure isolation for this test
		td.ClearCouponsTable(t, ctx, testPool)

		percentOff := 25.0
		maxRedemptions := 100
		requestBody := domain.CreateCouponRequest{
			Name:           "25% Off Coupon",
			PercentOff:     &percentOff,
			Duration:       "once",
			MaxRedemptions: &maxRedemptions,
			Metadata:       map[string]string{"test": "true", "type": "percent"},
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreateCouponRequest(t, ctx, jsonBody)
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
		assert.Equal(t, "Coupon created successfully!", response.Message)

		// Verify the coupon was actually created in the database
		couponID, err := uuid.Parse(response.ID)
		require.NoError(t, err)

		coupon, err := td.GetCouponByID(t, ctx, couponID, testPool)
		require.NoError(t, err)
		assert.Equal(t, "25% Off Coupon", coupon.Name)
		assert.Equal(t, 25.0, *coupon.PercentOff)
		assert.Nil(t, coupon.AmountOff)
		assert.Nil(t, coupon.Currency)
		assert.Equal(t, domain.CouponDurationOnce, coupon.Duration)
		assert.Equal(t, 100, *coupon.MaxRedemptions)
		assert.Equal(t, 0, coupon.TimesRedeemed)
		assert.True(t, coupon.IsValid)
	})

	t.Run("should successfully create a new amount-off coupon", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		amountOff := 500 // $5.00
		currency := "USD"
		requestBody := domain.CreateCouponRequest{
			Name:      "$5 Off Coupon",
			AmountOff: &amountOff,
			Currency:  &currency,
			Duration:  "once",
			Metadata:  map[string]string{"test": "true", "type": "amount"},
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreateCouponRequest(t, ctx, jsonBody)
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

		// Verify the coupon was created in the database
		couponID, err := uuid.Parse(response.ID)
		require.NoError(t, err)

		coupon, err := td.GetCouponByID(t, ctx, couponID, testPool)
		require.NoError(t, err)
		assert.Equal(t, "$5 Off Coupon", coupon.Name)
		assert.Nil(t, coupon.PercentOff)
		assert.Equal(t, 500, *coupon.AmountOff)
		assert.Equal(t, "USD", *coupon.Currency)
		assert.Equal(t, domain.CouponDurationOnce, coupon.Duration)
	})

	t.Run("should successfully create a repeating coupon with duration in months", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		percentOff := 15.0
		durationInMonths := 6
		redeemBy := time.Now().Add(365 * 24 * time.Hour) // 1 year from now

		requestBody := domain.CreateCouponRequest{
			Name:             "15% Off for 6 Months",
			PercentOff:       &percentOff,
			Duration:         "repeating",
			DurationInMonths: &durationInMonths,
			RedeemBy:         &redeemBy,
			Metadata:         map[string]string{"test": "true", "type": "repeating"},
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreateCouponRequest(t, ctx, jsonBody)
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

		// Verify the coupon was created in the database
		couponID, err := uuid.Parse(response.ID)
		require.NoError(t, err)

		coupon, err := td.GetCouponByID(t, ctx, couponID, testPool)
		require.NoError(t, err)
		assert.Equal(t, "15% Off for 6 Months", coupon.Name)
		assert.Equal(t, 15.0, *coupon.PercentOff)
		assert.Equal(t, domain.CouponDurationRepeating, coupon.Duration)
		assert.Equal(t, 6, *coupon.DurationInMonths)
		assert.NotNil(t, coupon.RedeemBy)
	})

	t.Run("should return 400 for missing required fields", func(t *testing.T) {
		requestBody := domain.CreateCouponRequest{
			Name: "Invalid Coupon",
			// Missing duration and discount amount
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreateCouponRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid duration", func(t *testing.T) {
		percentOff := 25.0
		requestBody := domain.CreateCouponRequest{
			Name:       "Invalid Duration",
			PercentOff: &percentOff,
			Duration:   "invalid_duration", // Invalid duration
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreateCouponRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for both percentOff and amountOff provided", func(t *testing.T) {
		percentOff := 25.0
		amountOff := 500
		currency := "USD"
		requestBody := domain.CreateCouponRequest{
			Name:       "Invalid Discount",
			PercentOff: &percentOff,
			AmountOff:  &amountOff,
			Currency:   &currency,
			Duration:   "once",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreateCouponRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for neither percentOff nor amountOff provided", func(t *testing.T) {
		requestBody := domain.CreateCouponRequest{
			Name:     "No Discount",
			Duration: "once",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreateCouponRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for amountOff without currency", func(t *testing.T) {
		amountOff := 500
		requestBody := domain.CreateCouponRequest{
			Name:      "No Currency",
			AmountOff: &amountOff,
			Duration:  "once",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreateCouponRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for repeating duration without durationInMonths", func(t *testing.T) {
		percentOff := 25.0
		requestBody := domain.CreateCouponRequest{
			Name:       "Repeating No Months",
			PercentOff: &percentOff,
			Duration:   "repeating",
			// Missing durationInMonths
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreateCouponRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid percentOff value", func(t *testing.T) {
		percentOff := 150.0 // Over 100%
		requestBody := domain.CreateCouponRequest{
			Name:       "Invalid Percent",
			PercentOff: &percentOff,
			Duration:   "once",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreateCouponRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 415 for non-JSON content type", func(t *testing.T) {
		req := newCreateCouponRequest(t, ctx, []byte("not json"))
		req.Header.Set("Content-Type", "text/plain")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
	})
}

