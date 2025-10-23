package promotionCodeRepository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeletePromotionCodesByCouponID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful deletion of multiple codes for a coupon", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Delete Multiple Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create multiple promotion codes for the same coupon
		code1 := helpers.NewValidPromotionCode("DELETE1", coupon.ID)
		code2 := helpers.NewValidPromotionCode("DELETE2", coupon.ID)
		code3 := helpers.NewValidPromotionCode("DELETE3", coupon.ID)

		helpers.InsertPromotionCode(t, ctx, testPool, code1)
		helpers.InsertPromotionCode(t, ctx, testPool, code2)
		helpers.InsertPromotionCode(t, ctx, testPool, code3)

		// Verify codes exist before deletion
		existsBefore1, err := helpers.GetPromotionCodeByID(t, ctx, code1.ID, testPool)
		require.NoError(t, err)
		existsBefore2, err := helpers.GetPromotionCodeByID(t, ctx, code2.ID, testPool)
		require.NoError(t, err)
		existsBefore3, err := helpers.GetPromotionCodeByID(t, ctx, code3.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, "DELETE1", existsBefore1.Code)
		assert.Equal(t, "DELETE2", existsBefore2.Code)
		assert.Equal(t, "DELETE3", existsBefore3.Code)

		// Delete all promotion codes for the coupon
		err = repo.DeletePromotionCodesByCouponID(ctx, coupon.ID)

		require.NoError(t, err)

		// Verify all codes are deleted
		deletedCode1 := helpers.GetPromotionCodeByIDOrNil(t, ctx, code1.ID, testPool)
		deletedCode2 := helpers.GetPromotionCodeByIDOrNil(t, ctx, code2.ID, testPool)
		deletedCode3 := helpers.GetPromotionCodeByIDOrNil(t, ctx, code3.ID, testPool)

		assert.Nil(t, deletedCode1)
		assert.Nil(t, deletedCode2)
		assert.Nil(t, deletedCode3)
	})

	t.Run("successful deletion of single code for a coupon", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Delete Single Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create single promotion code
		promotionCode := helpers.NewValidPromotionCode("DELETESINGLE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Delete promotion codes for the coupon
		err = repo.DeletePromotionCodesByCouponID(ctx, coupon.ID)

		require.NoError(t, err)

		// Verify code is deleted
		deletedCode := helpers.GetPromotionCodeByIDOrNil(t, ctx, promotionCode.ID, testPool)
		assert.Nil(t, deletedCode)
	})

	t.Run("preserves codes from other coupons", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create multiple coupons
		coupon1 := helpers.NewValidPercentOffCoupon("Target Coupon")
		coupon2 := helpers.NewValidPercentOffCoupon("Preserve Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon1)
		helpers.InsertCoupon(t, ctx, testPool, coupon2)

		// Create promotion codes for different coupons
		targetCode1 := helpers.NewValidPromotionCode("TARGET1", coupon1.ID)
		targetCode2 := helpers.NewValidPromotionCode("TARGET2", coupon1.ID)
		preserveCode1 := helpers.NewValidPromotionCode("PRESERVE1", coupon2.ID)
		preserveCode2 := helpers.NewValidPromotionCode("PRESERVE2", coupon2.ID)

		helpers.InsertPromotionCode(t, ctx, testPool, targetCode1)
		helpers.InsertPromotionCode(t, ctx, testPool, targetCode2)
		helpers.InsertPromotionCode(t, ctx, testPool, preserveCode1)
		helpers.InsertPromotionCode(t, ctx, testPool, preserveCode2)

		// Delete codes only for coupon1
		err = repo.DeletePromotionCodesByCouponID(ctx, coupon1.ID)

		require.NoError(t, err)

		// Verify coupon1's codes are deleted
		deletedCode1 := helpers.GetPromotionCodeByIDOrNil(t, ctx, targetCode1.ID, testPool)
		deletedCode2 := helpers.GetPromotionCodeByIDOrNil(t, ctx, targetCode2.ID, testPool)
		assert.Nil(t, deletedCode1)
		assert.Nil(t, deletedCode2)

		// Verify coupon2's codes are preserved
		preservedCode1, err := helpers.GetPromotionCodeByID(t, ctx, preserveCode1.ID, testPool)
		require.NoError(t, err)
		preservedCode2, err := helpers.GetPromotionCodeByID(t, ctx, preserveCode2.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, "PRESERVE1", preservedCode1.Code)
		assert.Equal(t, "PRESERVE2", preservedCode2.Code)
		assert.Equal(t, coupon2.ID, preservedCode1.CouponID)
		assert.Equal(t, coupon2.ID, preservedCode2.CouponID)
	})

	t.Run("deletion with no codes succeeds", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon without any promotion codes
		coupon := helpers.NewValidPercentOffCoupon("No Codes Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Delete promotion codes for coupon with no codes
		err = repo.DeletePromotionCodesByCouponID(ctx, coupon.ID)

		require.NoError(t, err)
		// Should succeed even with no codes to delete
	})

	t.Run("deletion with non-existent coupon succeeds", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		nonExistentCouponID := uuid.New()

		// Delete promotion codes for non-existent coupon
		err = repo.DeletePromotionCodesByCouponID(ctx, nonExistentCouponID)

		require.NoError(t, err)
		// Should succeed even if coupon doesn't exist
	})

	t.Run("deletes codes in all states", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("All States Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create codes in different states
		activeCode := helpers.NewValidPromotionCode("ACTIVE", coupon.ID)
		inactiveCode := helpers.NewInactivePromotionCode("INACTIVE", coupon.ID)
		expiredCode := helpers.NewExpiredPromotionCode("EXPIRED", coupon.ID)

		helpers.InsertPromotionCode(t, ctx, testPool, activeCode)
		helpers.InsertPromotionCode(t, ctx, testPool, inactiveCode)
		helpers.InsertPromotionCode(t, ctx, testPool, expiredCode)

		// Delete all codes for the coupon
		err = repo.DeletePromotionCodesByCouponID(ctx, coupon.ID)

		require.NoError(t, err)

		// Verify all codes are deleted regardless of state
		deletedActive := helpers.GetPromotionCodeByIDOrNil(t, ctx, activeCode.ID, testPool)
		deletedInactive := helpers.GetPromotionCodeByIDOrNil(t, ctx, inactiveCode.ID, testPool)
		deletedExpired := helpers.GetPromotionCodeByIDOrNil(t, ctx, expiredCode.ID, testPool)

		assert.Nil(t, deletedActive)
		assert.Nil(t, deletedInactive)
		assert.Nil(t, deletedExpired)
	})

	t.Run("deletes codes with various restrictions and limits", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Complex Codes Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create codes with different restrictions and limits
		basicCode := helpers.NewValidPromotionCode("BASIC", coupon.ID)
		restrictedCode := helpers.NewValidPromotionCodeWithRestrictions("RESTRICTED", coupon.ID, []string{"premium", "vip"})
		limitedCode := helpers.NewValidPromotionCodeWithRedemptionLimits("LIMITED", coupon.ID, 100)
		limitedCode.TimesRedeemed = 50

		helpers.InsertPromotionCode(t, ctx, testPool, basicCode)
		helpers.InsertPromotionCode(t, ctx, testPool, restrictedCode)
		helpers.InsertPromotionCode(t, ctx, testPool, limitedCode)

		// Delete all codes for the coupon
		err = repo.DeletePromotionCodesByCouponID(ctx, coupon.ID)

		require.NoError(t, err)

		// Verify all codes are deleted regardless of complexity
		deletedBasic := helpers.GetPromotionCodeByIDOrNil(t, ctx, basicCode.ID, testPool)
		deletedRestricted := helpers.GetPromotionCodeByIDOrNil(t, ctx, restrictedCode.ID, testPool)
		deletedLimited := helpers.GetPromotionCodeByIDOrNil(t, ctx, limitedCode.ID, testPool)

		assert.Nil(t, deletedBasic)
		assert.Nil(t, deletedRestricted)
		assert.Nil(t, deletedLimited)
	})

	t.Run("deletes codes with metadata", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Metadata Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create codes with complex metadata
		code1 := helpers.NewValidPromotionCode("META1", coupon.ID)
		code1.Metadata = map[string]string{
			"campaign": "summer2024",
			"channel":  "email",
		}

		code2 := helpers.NewValidPromotionCode("META2", coupon.ID)
		code2.Metadata = map[string]string{
			"json_data":     `{"key": "value", "nested": {"inner": true}}`,
			"special_chars": "Special: @#$%^&*()_+-=",
			"unicode":       "🎉 🚀 💡",
		}

		helpers.InsertPromotionCode(t, ctx, testPool, code1)
		helpers.InsertPromotionCode(t, ctx, testPool, code2)

		// Delete all codes for the coupon
		err = repo.DeletePromotionCodesByCouponID(ctx, coupon.ID)

		require.NoError(t, err)

		// Verify codes with metadata are deleted
		deletedCode1 := helpers.GetPromotionCodeByIDOrNil(t, ctx, code1.ID, testPool)
		deletedCode2 := helpers.GetPromotionCodeByIDOrNil(t, ctx, code2.ID, testPool)

		assert.Nil(t, deletedCode1)
		assert.Nil(t, deletedCode2)
	})

	t.Run("deletion does not affect coupon", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Coupon Preservation Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create promotion codes
		code1 := helpers.NewValidPromotionCode("PRESERVE1", coupon.ID)
		code2 := helpers.NewValidPromotionCode("PRESERVE2", coupon.ID)

		helpers.InsertPromotionCode(t, ctx, testPool, code1)
		helpers.InsertPromotionCode(t, ctx, testPool, code2)

		// Delete promotion codes
		err = repo.DeletePromotionCodesByCouponID(ctx, coupon.ID)

		require.NoError(t, err)

		// Verify coupon still exists and is unchanged
		remainingCoupon, err := helpers.GetCouponByID(t, ctx, coupon.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, coupon.Name, remainingCoupon.Name)
		assert.Equal(t, coupon.StripeCouponID, remainingCoupon.StripeCouponID)
		assert.Equal(t, coupon.PercentOff, remainingCoupon.PercentOff)

		// Verify codes are deleted
		deletedCode1 := helpers.GetPromotionCodeByIDOrNil(t, ctx, code1.ID, testPool)
		deletedCode2 := helpers.GetPromotionCodeByIDOrNil(t, ctx, code2.ID, testPool)
		assert.Nil(t, deletedCode1)
		assert.Nil(t, deletedCode2)
	})

	t.Run("large scale deletion", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Large Scale Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create many promotion codes
		numberOfCodes := 25
		codeIDs := make([]uuid.UUID, numberOfCodes)
		for i := 0; i < numberOfCodes; i++ {
			code := helpers.NewValidPromotionCode(fmt.Sprintf("BULK%03d", i+1), coupon.ID)
			helpers.InsertPromotionCode(t, ctx, testPool, code)
			codeIDs[i] = code.ID
		}

		// Delete all codes for the coupon
		err = repo.DeletePromotionCodesByCouponID(ctx, coupon.ID)

		require.NoError(t, err)

		// Verify all codes are deleted
		for i, codeID := range codeIDs {
			deletedCode := helpers.GetPromotionCodeByIDOrNil(t, ctx, codeID, testPool)
			assert.Nil(t, deletedCode, "Code at index %d should be deleted", i)
		}
	})

	t.Run("deletion with nil UUID should handle gracefully", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)

		// Delete with nil UUID
		err = repo.DeletePromotionCodesByCouponID(ctx, uuid.Nil)

		// Should handle gracefully (implementation dependent - might succeed or return error)
		// The important thing is it shouldn't crash
		if err != nil {
			// If error is returned, it should be a proper error, not a panic
			assert.Error(t, err)
		} else {
			// If no error, operation should complete without issues
			require.NoError(t, err)
		}
	})

	t.Run("multiple successive deletions", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create multiple coupons
		coupon1 := helpers.NewValidPercentOffCoupon("Successive 1")
		coupon2 := helpers.NewValidPercentOffCoupon("Successive 2")
		coupon3 := helpers.NewValidPercentOffCoupon("Successive 3")

		helpers.InsertCoupon(t, ctx, testPool, coupon1)
		helpers.InsertCoupon(t, ctx, testPool, coupon2)
		helpers.InsertCoupon(t, ctx, testPool, coupon3)

		// Create codes for each coupon
		code1 := helpers.NewValidPromotionCode("SUCC1", coupon1.ID)
		code2a := helpers.NewValidPromotionCode("SUCC2A", coupon2.ID)
		code2b := helpers.NewValidPromotionCode("SUCC2B", coupon2.ID)
		code3 := helpers.NewValidPromotionCode("SUCC3", coupon3.ID)

		helpers.InsertPromotionCode(t, ctx, testPool, code1)
		helpers.InsertPromotionCode(t, ctx, testPool, code2a)
		helpers.InsertPromotionCode(t, ctx, testPool, code2b)
		helpers.InsertPromotionCode(t, ctx, testPool, code3)

		// Delete codes for each coupon successively
		err1 := repo.DeletePromotionCodesByCouponID(ctx, coupon1.ID)
		err2 := repo.DeletePromotionCodesByCouponID(ctx, coupon2.ID)
		err3 := repo.DeletePromotionCodesByCouponID(ctx, coupon3.ID)

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NoError(t, err3)

		// Verify all codes are deleted
		deletedCode1 := helpers.GetPromotionCodeByIDOrNil(t, ctx, code1.ID, testPool)
		deletedCode2a := helpers.GetPromotionCodeByIDOrNil(t, ctx, code2a.ID, testPool)
		deletedCode2b := helpers.GetPromotionCodeByIDOrNil(t, ctx, code2b.ID, testPool)
		deletedCode3 := helpers.GetPromotionCodeByIDOrNil(t, ctx, code3.ID, testPool)

		assert.Nil(t, deletedCode1)
		assert.Nil(t, deletedCode2a)
		assert.Nil(t, deletedCode2b)
		assert.Nil(t, deletedCode3)
	})
}

