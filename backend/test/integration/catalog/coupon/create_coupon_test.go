package coupon_test

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
)

// make test-func TEST_NAME=TestCreateCoupon TEST_PATH=test/integration/catalog/coupon/create_coupon_test.go

func TestCreateCoupon(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully create coupon with valid admin token", func(t *testing.T) {
		// Clean the database to ensure isolation for this test
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		percentOff := 25.0
		maxRedemptions := 100
		requestBody := domain.CreateCouponRequest{
			Name:           "25% Off Coupon",
			PercentOff:     &percentOff,
			Duration:       "once",
			MaxRedemptions: &maxRedemptions,
			Metadata:       map[string]string{"test": "true", "type": "percent"},
		}

		req := th.NewCreateCouponRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response struct {
			ID      string `json:"id"`
			Message string `json:"message"`
		}

		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		// Validate response
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, "Coupon created successfully!", response.Message)

		// Verify the coupon was actually created in the database
		couponID, err := uuid.Parse(response.ID)
		assert.NoError(t, err)

		coupon, err := td.GetCouponByID(t, ctx, couponID, testPool)
		assert.NoError(t, err)
		assert.Equal(t, "25% Off Coupon", coupon.Name)
		assert.Equal(t, 25.0, *coupon.PercentOff)
		assert.Nil(t, coupon.AmountOff)
		assert.Nil(t, coupon.Currency)
		assert.Equal(t, domain.CouponDurationOnce, coupon.Duration)
		assert.Equal(t, 100, *coupon.MaxRedemptions)
		assert.Equal(t, 0, coupon.TimesRedeemed)
		assert.True(t, coupon.IsValid)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)

		percentOff := 25.0
		requestBody := domain.CreateCouponRequest{
			Name:       "Test Coupon",
			PercentOff: &percentOff,
			Duration:   "once",
		}

		req := th.NewCreateCouponRequest(t, ctx, testServerURL, requestBody, "")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		percentOff := 25.0
		requestBody := domain.CreateCouponRequest{
			Name:       "Test Coupon",
			PercentOff: &percentOff,
			Duration:   "once",
		}

		req := th.NewCreateCouponRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		percentOff := 25.0
		requestBody := domain.CreateCouponRequest{
			Name:       "Test Coupon",
			PercentOff: &percentOff,
			Duration:   "once",
		}

		req := th.NewCreateCouponRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)

		percentOff := 25.0
		requestBody := domain.CreateCouponRequest{
			Name:       "Test Coupon",
			PercentOff: &percentOff,
			Duration:   "once",
		}

		req := th.NewCreateCouponRequest(t, ctx, testServerURL, requestBody, "invalid-token-12345")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 400 for missing required fields", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		requestBody := domain.CreateCouponRequest{
			Name: "Invalid Coupon",
			// Missing duration and discount amount
		}

		req := th.NewCreateCouponRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid duration", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		percentOff := 25.0
		requestBody := domain.CreateCouponRequest{
			Name:       "Invalid Duration",
			PercentOff: &percentOff,
			Duration:   "invalid_duration", // Invalid duration
		}

		req := th.NewCreateCouponRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for both percentOff and amountOff provided", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

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

		req := th.NewCreateCouponRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for neither percentOff nor amountOff provided", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		requestBody := domain.CreateCouponRequest{
			Name:     "No Discount",
			Duration: "once",
		}

		req := th.NewCreateCouponRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for amountOff without currency", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		amountOff := 500
		requestBody := domain.CreateCouponRequest{
			Name:      "No Currency",
			AmountOff: &amountOff,
			Duration:  "once",
		}

		req := th.NewCreateCouponRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for repeating duration without durationInMonths", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		percentOff := 25.0
		requestBody := domain.CreateCouponRequest{
			Name:       "Repeating No Months",
			PercentOff: &percentOff,
			Duration:   "repeating",
			// Missing durationInMonths
		}

		req := th.NewCreateCouponRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid percentOff value", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		percentOff := 150.0 // Over 100%
		requestBody := domain.CreateCouponRequest{
			Name:       "Invalid Percent",
			PercentOff: &percentOff,
			Duration:   "once",
		}

		req := th.NewCreateCouponRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
