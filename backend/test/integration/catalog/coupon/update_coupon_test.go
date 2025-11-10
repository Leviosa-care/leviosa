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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdateCoupon TEST_PATH=test/integration/catalog/coupon/update_coupon_test.go

func TestUpdateCoupon(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully update coupon name with valid admin token", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create test coupon
		testCoupon := td.NewValidPercentOffCoupon("Original Name")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		// Update request
		newName := "Updated Coupon Name"
		updateRequest := domain.UpdateCouponRequest{
			Name: &newName,
		}

		req := th.NewUpdateCouponRequest(t, ctx, testServerURL, testCoupon.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Message string `json:"message"`
		}

		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Coupon updated successfully!", response.Message)

		// Verify the update in database
		updatedCoupon, err := td.GetCouponByID(t, ctx, testCoupon.ID, testPool)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Coupon Name", updatedCoupon.Name)
		// Verify other fields remain unchanged
		assert.Equal(t, testCoupon.StripeCouponID, updatedCoupon.StripeCouponID)
		assert.Equal(t, *testCoupon.PercentOff, *updatedCoupon.PercentOff)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)

		newName := "Test Update"
		updateRequest := domain.UpdateCouponRequest{
			Name: &newName,
		}

		req := th.NewUpdateCouponRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", updateRequest, "")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		newName := "Test Update"
		updateRequest := domain.UpdateCouponRequest{
			Name: &newName,
		}

		req := th.NewUpdateCouponRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", updateRequest, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		newName := "Test Update"
		updateRequest := domain.UpdateCouponRequest{
			Name: &newName,
		}

		req := th.NewUpdateCouponRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", updateRequest, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		td.ClearCouponsTable(t, ctx, testPool)

		newName := "Test Update"
		updateRequest := domain.UpdateCouponRequest{
			Name: &newName,
		}

		req := th.NewUpdateCouponRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", updateRequest, "invalid-token-12345")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 404 for non-existent coupon", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		newName := "Non-existent Update"
		updateRequest := domain.UpdateCouponRequest{
			Name: &newName,
		}

		req := th.NewUpdateCouponRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for empty name", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create test coupon
		testCoupon := td.NewValidPercentOffCoupon("Valid Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		// Update request with empty name
		emptyName := ""
		updateRequest := domain.UpdateCouponRequest{
			Name: &emptyName,
		}

		req := th.NewUpdateCouponRequest(t, ctx, testServerURL, testCoupon.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		newName := "Invalid UUID Test"
		updateRequest := domain.UpdateCouponRequest{
			Name: &newName,
		}

		req := th.NewUpdateCouponRequest(t, ctx, testServerURL, "invalid-uuid", updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestDeactivateCoupon(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully deactivate coupon with valid admin token", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create test coupon
		testCoupon := td.NewValidPercentOffCoupon("Deactivate Test")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		// Verify it's initially valid
		assert.True(t, td.GetCouponValidStatus(t, ctx, testCoupon.ID, testPool))

		req := th.NewDeactivateCouponRequest(t, ctx, testServerURL, testCoupon.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Message string `json:"message"`
		}
		assert.Equal(t, "Coupon deactivated successfully!", response.Message)

		// Verify the coupon is deactivated in database
		assert.False(t, td.GetCouponValidStatus(t, ctx, testCoupon.ID, testPool))
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		req := th.NewDeactivateCouponRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", "")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		req := th.NewDeactivateCouponRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		req := th.NewDeactivateCouponRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		req := th.NewDeactivateCouponRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", "invalid-token-12345")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 404 for non-existent coupon", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewDeactivateCouponRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewDeactivateCouponRequest(t, ctx, testServerURL, "invalid-uuid", accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestDeleteCoupon(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully delete coupon with valid admin token", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create test coupon
		testCoupon := td.NewValidPercentOffCoupon("Delete Test")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		// Verify it exists
		foundCoupon := td.GetCouponByIDOrNil(t, ctx, testCoupon.ID, testPool)
		assert.NotNil(t, foundCoupon)

		req := th.NewDeleteCouponRequest(t, ctx, testServerURL, testCoupon.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Message string `json:"message"`
		}
		assert.Equal(t, "Coupon deleted successfully!", response.Message)

		// Verify the coupon is deleted from database
		deletedCoupon := td.GetCouponByIDOrNil(t, ctx, testCoupon.ID, testPool)
		assert.Nil(t, deletedCoupon)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		req := th.NewDeleteCouponRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", "")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		req := th.NewDeleteCouponRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		req := th.NewDeleteCouponRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		req := th.NewDeleteCouponRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", "invalid-token-12345")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 404 for non-existent coupon", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewDeleteCouponRequest(t, ctx, testServerURL, "00000000-0000-0000-0000-000000000000", accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewDeleteCouponRequest(t, ctx, testServerURL, "invalid-uuid", accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestCouponBusinessLogic(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should handle coupon validation with complex scenarios", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon with future expiry
		testCoupon := td.NewValidCouponWithRedeemBy("Future Expiry", time.Now().Add(30*24*time.Hour))
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		// Validate the coupon (public endpoint)
		requestBody := struct {
			StripeCouponID string `json:"stripeCouponId"`
		}{
			StripeCouponID: testCoupon.StripeCouponID,
		}

		req := th.NewValidateCouponRequest(t, ctx, testServerURL, requestBody)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Valid  bool `json:"valid"`
			Coupon any  `json:"coupon"`
		}

		assert.True(t, response.Valid)
		assert.NotNil(t, response.Coupon)
	})

	t.Run("should handle coupon with redemption limits correctly", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon with redemption limits but not reached
		testCoupon := td.NewValidPercentOffCouponWithRedemptionLimits("Limited Coupon", 10)
		testCoupon.TimesRedeemed = 5 // Not at limit yet
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		// Validate the coupon (public endpoint)
		requestBody := struct {
			StripeCouponID string `json:"stripeCouponId"`
		}{
			StripeCouponID: testCoupon.StripeCouponID,
		}

		req := th.NewValidateCouponRequest(t, ctx, testServerURL, requestBody)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Valid  bool `json:"valid"`
			Coupon any  `json:"coupon"`
		}

		assert.True(t, response.Valid)
		assert.NotNil(t, response.Coupon)
	})

	t.Run("should properly handle forever duration coupons", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create a forever duration coupon
		testCoupon := td.NewValidForeverCoupon("Forever Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		req := th.NewGetCouponByIDRequest(t, ctx, testServerURL, testCoupon.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response := th.ParseCouponResponse(t, resp)

		assert.Nil(t, response.RedeemBy)
		assert.True(t, response.Valid)
	})
}
