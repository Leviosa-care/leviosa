package promotionCodeRepository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeletePromotionCode(t *testing.T) {
	ctx := context.Background()

	t.Run("successful deletion of existing code", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Delete Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("DELETE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Verify code exists before deletion
		existingCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, promotionCode.Code, existingCode.Code)

		// Delete promotion code
		err = repo.DeletePromotionCode(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify deletion - code should no longer exist
		deletedCode := helpers.GetPromotionCodeByIDOrNil(t, ctx, promotionCode.ID, testPool)
		assert.Nil(t, deletedCode)
	})

	t.Run("deletion of active code", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and active promotion code
		coupon := helpers.NewValidPercentOffCoupon("Delete Active Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("DELETEACTIVE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Delete active promotion code
		err = repo.DeletePromotionCode(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify deletion
		deletedCode := helpers.GetPromotionCodeByIDOrNil(t, ctx, promotionCode.ID, testPool)
		assert.Nil(t, deletedCode)
	})

	t.Run("deletion of inactive code", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and inactive promotion code
		coupon := helpers.NewValidPercentOffCoupon("Delete Inactive Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewInactivePromotionCode("DELETEINACTIVE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Delete inactive promotion code
		err = repo.DeletePromotionCode(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify deletion
		deletedCode := helpers.GetPromotionCodeByIDOrNil(t, ctx, promotionCode.ID, testPool)
		assert.Nil(t, deletedCode)
	})

	t.Run("deletion of expired code", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and expired promotion code
		coupon := helpers.NewValidPercentOffCoupon("Delete Expired Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewExpiredPromotionCode("DELETEEXPIRED", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Delete expired promotion code
		err = repo.DeletePromotionCode(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify deletion
		deletedCode := helpers.GetPromotionCodeByIDOrNil(t, ctx, promotionCode.ID, testPool)
		assert.Nil(t, deletedCode)
	})

	t.Run("deletion of code with restrictions", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code with restrictions
		coupon := helpers.NewValidPercentOffCoupon("Delete Restricted Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCodeWithRestrictions("DELETERESTRICTED", coupon.ID, []string{"premium_user", "vip_user"})
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Delete promotion code with restrictions
		err = repo.DeletePromotionCode(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify deletion
		deletedCode := helpers.GetPromotionCodeByIDOrNil(t, ctx, promotionCode.ID, testPool)
		assert.Nil(t, deletedCode)
	})

	t.Run("deletion of code with redemption limits", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code with limits
		coupon := helpers.NewValidPercentOffCoupon("Delete Limited Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCodeWithRedemptionLimits("DELETELIMITED", coupon.ID, 50)
		promotionCode.TimesRedeemed = 25
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Delete promotion code with limits
		err = repo.DeletePromotionCode(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify deletion
		deletedCode := helpers.GetPromotionCodeByIDOrNil(t, ctx, promotionCode.ID, testPool)
		assert.Nil(t, deletedCode)
	})

	t.Run("deletion of code with metadata", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code with metadata
		coupon := helpers.NewValidPercentOffCoupon("Delete Metadata Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("DELETEMETA", coupon.ID)
		promotionCode.Metadata = map[string]string{
			"campaign":   "delete_test_campaign",
			"channel":    "email",
			"created_by": "admin_user",
		}
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Delete promotion code with metadata
		err = repo.DeletePromotionCode(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify deletion
		deletedCode := helpers.GetPromotionCodeByIDOrNil(t, ctx, promotionCode.ID, testPool)
		assert.Nil(t, deletedCode)
	})

	t.Run("delete non-existent promotion code should return not found error", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)

		nonExistentID := uuid.New()

		err = repo.DeletePromotionCode(ctx, nonExistentID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("delete with nil UUID should return error", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)

		err = repo.DeletePromotionCode(ctx, uuid.Nil)

		assert.Error(t, err)
	})

	t.Run("deletion preserves other codes", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and multiple promotion codes
		coupon := helpers.NewValidPercentOffCoupon("Preserve Others Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		code1 := helpers.NewValidPromotionCode("PRESERVE1", coupon.ID)
		code2 := helpers.NewValidPromotionCode("PRESERVE2", coupon.ID)
		code3 := helpers.NewValidPromotionCode("PRESERVE3", coupon.ID)

		helpers.InsertPromotionCode(t, ctx, testPool, code1)
		helpers.InsertPromotionCode(t, ctx, testPool, code2)
		helpers.InsertPromotionCode(t, ctx, testPool, code3)

		// Delete only the middle code
		err = repo.DeletePromotionCode(ctx, code2.ID)

		require.NoError(t, err)

		// Verify only code2 was deleted
		remainingCode1, err := helpers.GetPromotionCodeByID(t, ctx, code1.ID, testPool)
		require.NoError(t, err)
		deletedCode2 := helpers.GetPromotionCodeByIDOrNil(t, ctx, code2.ID, testPool)
		remainingCode3, err := helpers.GetPromotionCodeByID(t, ctx, code3.ID, testPool)
		require.NoError(t, err)

		assert.Equal(t, "PRESERVE1", remainingCode1.Code)
		assert.Nil(t, deletedCode2)
		assert.Equal(t, "PRESERVE3", remainingCode3.Code)
	})

	t.Run("deletion does not affect coupon", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Coupon Preservation Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("DELETECOUPONTEST", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Delete promotion code
		err = repo.DeletePromotionCode(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify coupon still exists and unchanged
		remainingCoupon, err := helpers.GetCouponByID(t, ctx, coupon.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, coupon.Name, remainingCoupon.Name)
		assert.Equal(t, coupon.StripeCouponID, remainingCoupon.StripeCouponID)
	})

	t.Run("multiple deletions", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and multiple promotion codes
		coupon := helpers.NewValidPercentOffCoupon("Multiple Delete Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		codes := make([]*domain.PromotionCode, 5)
		for i := 0; i < 5; i++ {
			codes[i] = helpers.NewValidPromotionCode(fmt.Sprintf("MULTI%d", i+1), coupon.ID)
			helpers.InsertPromotionCode(t, ctx, testPool, codes[i])
		}

		// Delete all codes
		for i := 0; i < 5; i++ {
			err = repo.DeletePromotionCode(ctx, codes[i].ID)
			require.NoError(t, err)
		}

		// Verify all codes are deleted
		for i := 0; i < 5; i++ {
			deletedCode := helpers.GetPromotionCodeByIDOrNil(t, ctx, codes[i].ID, testPool)
			assert.Nil(t, deletedCode)
		}
	})

	t.Run("deletion of code with complex data", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and complex promotion code
		coupon := helpers.NewValidPercentOffCoupon("Complex Delete Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create promotion code with all possible fields set
		expiryDate := time.Now().Add(30 * 24 * time.Hour)
		promotionCode := helpers.NewValidPromotionCodeWithExpiry("COMPLEX", coupon.ID, expiryDate, 100, []string{"premium", "vip"})
		promotionCode.Metadata = map[string]string{
			"json_data":       `{"key": "value", "nested": {"inner": true}}`,
			"special_chars":   "Special: @#$%^&*()_+-=[]{}|;':\",./<>?",
			"unicode_content": "Unicode: 🎉 🚀 💡 ✨",
			"long_string":     "This is a very long string that contains multiple words and should test the deletion of complex metadata content",
		}
		promotionCode.TimesRedeemed = 42
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Delete complex promotion code
		err = repo.DeletePromotionCode(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify complete deletion
		deletedCode := helpers.GetPromotionCodeByIDOrNil(t, ctx, promotionCode.ID, testPool)
		assert.Nil(t, deletedCode)
	})
}

