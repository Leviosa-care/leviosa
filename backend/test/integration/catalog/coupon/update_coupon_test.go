package coupon_test

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

// make test-func TEST_NAME=TestUpdateCoupon TEST_PATH=test/integration/catalog/coupon/update_coupon_test.go

func TestUpdateCoupon(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully update coupon name", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon
		testCoupon := td.NewValidPercentOffCoupon("Original Name")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		// Update request
		newName := "Updated Coupon Name"
		updateRequest := domain.UpdateCouponRequest{
			Name: &newName,
		}
		jsonBody, _ := json.Marshal(updateRequest)

		req := newUpdateCouponRequest(t, ctx, testCoupon.ID.String(), jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Message string `json:"message"`
		}
		decodeJSONResponse(t, resp, &response)
		assert.Equal(t, "Coupon updated successfully!", response.Message)

		// Verify the update in database
		updatedCoupon, err := td.GetCouponByID(t, ctx, testCoupon.ID, testPool)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Coupon Name", updatedCoupon.Name)
		// Verify other fields remain unchanged
		assert.Equal(t, testCoupon.StripeCouponID, updatedCoupon.StripeCouponID)
		assert.Equal(t, *testCoupon.PercentOff, *updatedCoupon.PercentOff)
	})

	t.Run("should successfully update coupon metadata", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon
		testCoupon := td.NewValidPercentOffCoupon("Metadata Test")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		// Update request with new metadata
		updateRequest := domain.UpdateCouponRequest{
			Metadata: map[string]string{
				"updated":     "true",
				"test":        "false",
				"description": "Updated coupon metadata",
				"version":     "2.0",
			},
		}
		jsonBody, _ := json.Marshal(updateRequest)

		req := newUpdateCouponRequest(t, ctx, testCoupon.ID.String(), jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify the metadata update in database
		updatedCoupon, err := td.GetCouponByID(t, ctx, testCoupon.ID, testPool)
		assert.NoError(t, err)
		assert.Equal(t, "true", updatedCoupon.Metadata["updated"])
		assert.Equal(t, "false", updatedCoupon.Metadata["test"]) // Should replace original
		assert.Equal(t, "Updated coupon metadata", updatedCoupon.Metadata["description"])
		assert.Equal(t, "2.0", updatedCoupon.Metadata["version"])
	})

	t.Run("should successfully update both name and metadata", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon
		testCoupon := td.NewValidAmountOffCoupon("Combined Update Test", "USD")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		// Update request with both name and metadata
		newName := "Combined Updated Name"
		updateRequest := domain.UpdateCouponRequest{
			Name: &newName,
			Metadata: map[string]string{
				"updated": "true",
				"type":    "combined_update",
			},
		}
		jsonBody, _ := json.Marshal(updateRequest)

		req := newUpdateCouponRequest(t, ctx, testCoupon.ID.String(), jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify both updates in database
		updatedCoupon, err := td.GetCouponByID(t, ctx, testCoupon.ID, testPool)
		assert.NoError(t, err)
		assert.Equal(t, "Combined Updated Name", updatedCoupon.Name)
		assert.Equal(t, "true", updatedCoupon.Metadata["updated"])
		assert.Equal(t, "combined_update", updatedCoupon.Metadata["type"])
	})

	t.Run("should return 404 for non-existent coupon", func(t *testing.T) {
		newName := "Non-existent Update"
		updateRequest := domain.UpdateCouponRequest{
			Name: &newName,
		}
		jsonBody, _ := json.Marshal(updateRequest)

		req := newUpdateCouponRequest(t, ctx, "00000000-0000-0000-0000-000000000000", jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for empty name", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon
		testCoupon := td.NewValidPercentOffCoupon("Valid Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		// Update request with empty name
		emptyName := ""
		updateRequest := domain.UpdateCouponRequest{
			Name: &emptyName,
		}
		jsonBody, _ := json.Marshal(updateRequest)

		req := newUpdateCouponRequest(t, ctx, testCoupon.ID.String(), jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		newName := "Invalid UUID Test"
		updateRequest := domain.UpdateCouponRequest{
			Name: &newName,
		}
		jsonBody, _ := json.Marshal(updateRequest)

		req := newUpdateCouponRequest(t, ctx, "invalid-uuid", jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 415 for non-JSON content type", func(t *testing.T) {
		req := newUpdateCouponRequest(t, ctx, "00000000-0000-0000-0000-000000000000", []byte("not json"))
		req.Header.Set("Content-Type", "text/plain")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
	})
}

func TestDeactivateCoupon(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully deactivate coupon", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon
		testCoupon := td.NewValidPercentOffCoupon("Deactivate Test")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		// Verify it's initially valid
		assert.True(t, td.GetCouponValidStatus(t, ctx, testCoupon.ID, testPool))

		req := newDeactivateCouponRequest(t, ctx, testCoupon.ID.String())

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Message string `json:"message"`
		}
		decodeJSONResponse(t, resp, &response)
		assert.Equal(t, "Coupon deactivated successfully!", response.Message)

		// Verify the coupon is deactivated in database
		assert.False(t, td.GetCouponValidStatus(t, ctx, testCoupon.ID, testPool))
	})

	t.Run("should return 404 for non-existent coupon", func(t *testing.T) {
		req := newDeactivateCouponRequest(t, ctx, "00000000-0000-0000-0000-000000000000")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		req := newDeactivateCouponRequest(t, ctx, "invalid-uuid")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestDeleteCoupon(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully delete coupon", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create test coupon
		testCoupon := td.NewValidPercentOffCoupon("Delete Test")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		// Verify it exists
		foundCoupon := td.GetCouponByIDOrNil(t, ctx, testCoupon.ID, testPool)
		assert.NotNil(t, foundCoupon)

		req := newDeleteCouponRequest(t, ctx, testCoupon.ID.String())

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Message string `json:"message"`
		}
		decodeJSONResponse(t, resp, &response)
		assert.Equal(t, "Coupon deleted successfully!", response.Message)

		// Verify the coupon is deleted from database
		deletedCoupon := td.GetCouponByIDOrNil(t, ctx, testCoupon.ID, testPool)
		assert.Nil(t, deletedCoupon)
	})

	t.Run("should return 404 for non-existent coupon", func(t *testing.T) {
		req := newDeleteCouponRequest(t, ctx, "00000000-0000-0000-0000-000000000000")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		req := newDeleteCouponRequest(t, ctx, "invalid-uuid")

		resp, err := client.Do(req)
		assert.NoError(t, err)
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

		// Validate the coupon
		requestBody := struct {
			StripeCouponID string `json:"stripeCouponId"`
		}{
			StripeCouponID: testCoupon.StripeCouponID,
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newValidateCouponRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Valid  bool `json:"valid"`
			Coupon any  `json:"coupon"`
		}
		decodeJSONResponse(t, resp, &response)

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

		// Validate the coupon
		requestBody := struct {
			StripeCouponID string `json:"stripeCouponId"`
		}{
			StripeCouponID: testCoupon.StripeCouponID,
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newValidateCouponRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Valid  bool `json:"valid"`
			Coupon any  `json:"coupon"`
		}
		decodeJSONResponse(t, resp, &response)

		assert.True(t, response.Valid)
		assert.NotNil(t, response.Coupon)
	})

	t.Run("should properly handle forever duration coupons", func(t *testing.T) {
		// Clean the database
		td.ClearCouponsTable(t, ctx, testPool)

		// Create a forever duration coupon
		testCoupon := td.NewValidForeverCoupon("Forever Coupon")
		td.InsertCoupon(t, ctx, testPool, testCoupon)

		req := newGetCouponByIDRequest(t, ctx, testCoupon.ID.String())

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.CouponResponse
		decodeJSONResponse(t, resp, &response)

		assert.Nil(t, response.RedeemBy)
		assert.True(t, response.Valid)
	})
}
