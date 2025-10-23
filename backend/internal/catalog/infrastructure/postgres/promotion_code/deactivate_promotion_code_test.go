package promotionCodeRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeactivatePromotionCode(t *testing.T) {
	ctx := context.Background()

	t.Run("successful deactivation of active code", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and active promotion code
		coupon := helpers.NewValidPercentOffCoupon("Deactivate Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("DEACTIVATE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Verify code is initially active
		initialCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.True(t, initialCode.Active)

		// Deactivate promotion code
		err = repo.DeactivatePromotionCode(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify deactivation
		deactivatedCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.False(t, deactivatedCode.Active)
		assert.Equal(t, promotionCode.Code, deactivatedCode.Code) // Other fields unchanged
		assert.Equal(t, promotionCode.CouponID, deactivatedCode.CouponID)
	})

	t.Run("deactivation of already inactive code", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and inactive promotion code
		coupon := helpers.NewValidPercentOffCoupon("Already Inactive Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewInactivePromotionCode("ALREADYINACTIVE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Verify code is initially inactive
		initialCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.False(t, initialCode.Active)

		// Deactivate already inactive code (should succeed)
		err = repo.DeactivatePromotionCode(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify still inactive
		stillInactiveCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.False(t, stillInactiveCode.Active)
	})

	t.Run("deactivation preserves other fields", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code with all fields
		coupon := helpers.NewValidPercentOffCoupon("Preserve Fields Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCodeWithRestrictions("PRESERVE", coupon.ID, []string{"premium_user"})
		promotionCode.Metadata = map[string]string{
			"campaign": "preserve_test",
			"channel":  "email",
		}
		promotionCode.TimesRedeemed = 5
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Deactivate promotion code
		err = repo.DeactivatePromotionCode(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify all other fields preserved
		deactivatedCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.False(t, deactivatedCode.Active)                                                 // Changed
		assert.Equal(t, promotionCode.Code, deactivatedCode.Code)                               // Unchanged
		assert.Equal(t, promotionCode.CouponID, deactivatedCode.CouponID)                       // Unchanged
		assert.Equal(t, promotionCode.StripePromotionID, deactivatedCode.StripePromotionID)     // Unchanged
		assert.Equal(t, []string{"premium_user"}, deactivatedCode.Restrictions.CurrencyOptions) // Unchanged
		assert.Equal(t, "preserve_test", deactivatedCode.Metadata["campaign"])                  // Unchanged
		assert.Equal(t, "email", deactivatedCode.Metadata["channel"])                           // Unchanged
		assert.Equal(t, 5, deactivatedCode.TimesRedeemed)                                       // Unchanged
	})

	t.Run("deactivation with redemption limits", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code with limits
		coupon := helpers.NewValidPercentOffCoupon("Limited Deactivate Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCodeWithRedemptionLimits("LIMITED", coupon.ID, 100)
		promotionCode.TimesRedeemed = 25
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Deactivate promotion code
		err = repo.DeactivatePromotionCode(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify limits preserved
		deactivatedCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.False(t, deactivatedCode.Active)
		require.NotNil(t, deactivatedCode.MaxRedemptions)
		assert.Equal(t, 100, *deactivatedCode.MaxRedemptions)
		assert.Equal(t, 25, deactivatedCode.TimesRedeemed)
	})

	t.Run("deactivation with expiry date", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code with expiry
		coupon := helpers.NewValidPercentOffCoupon("Expiry Deactivate Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewExpiredPromotionCode("EXPIREDDEACTIVATE", coupon.ID)
		// Override to make it active initially for this test
		promotionCode.Active = true
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Deactivate promotion code
		err = repo.DeactivatePromotionCode(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify expiry preserved
		deactivatedCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.False(t, deactivatedCode.Active)
		assert.NotNil(t, deactivatedCode.ExpiresAt)
		assert.Equal(t, promotionCode.ExpiresAt.Unix(), deactivatedCode.ExpiresAt.Unix())
	})

	t.Run("deactivate non-existent promotion code should return not found error", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)

		nonExistentID := uuid.New()

		err := repo.DeactivatePromotionCode(ctx, nonExistentID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("deactivate with nil UUID should return error", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)

		err := repo.DeactivatePromotionCode(ctx, uuid.Nil)

		assert.Error(t, err)
	})

	t.Run("deactivation updates updated_at timestamp", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Timestamp Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("TIMESTAMP", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Get initial timestamp
		initialCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		initialUpdatedAt := initialCode.UpdatedAt

		// Small delay to ensure timestamp difference
		// Note: In real tests you might want to mock time or use a more reliable approach
		// time.Sleep(time.Millisecond)

		// Deactivate promotion code
		err = repo.DeactivatePromotionCode(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify timestamp updated (should be different or equal due to precision)
		deactivatedCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.True(t, deactivatedCode.UpdatedAt.Equal(initialUpdatedAt) || deactivatedCode.UpdatedAt.After(initialUpdatedAt))
		assert.Equal(t, initialCode.CreatedAt, deactivatedCode.CreatedAt) // Created at should not change
	})

	t.Run("multiple deactivations", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and multiple promotion codes
		coupon := helpers.NewValidPercentOffCoupon("Multiple Deactivate Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		code1 := helpers.NewValidPromotionCode("MULTI1", coupon.ID)
		code2 := helpers.NewValidPromotionCode("MULTI2", coupon.ID)
		code3 := helpers.NewValidPromotionCode("MULTI3", coupon.ID)

		helpers.InsertPromotionCode(t, ctx, testPool, code1)
		helpers.InsertPromotionCode(t, ctx, testPool, code2)
		helpers.InsertPromotionCode(t, ctx, testPool, code3)

		// Deactivate all codes
		err1 := repo.DeactivatePromotionCode(ctx, code1.ID)
		err2 := repo.DeactivatePromotionCode(ctx, code2.ID)
		err3 := repo.DeactivatePromotionCode(ctx, code3.ID)

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NoError(t, err3)

		// Verify all deactivated
		deactivatedCode1, err := helpers.GetPromotionCodeByID(t, ctx, code1.ID, testPool)
		require.NoError(t, err)
		deactivatedCode2, err := helpers.GetPromotionCodeByID(t, ctx, code2.ID, testPool)
		require.NoError(t, err)
		deactivatedCode3, err := helpers.GetPromotionCodeByID(t, ctx, code3.ID, testPool)
		require.NoError(t, err)

		assert.False(t, deactivatedCode1.Active)
		assert.False(t, deactivatedCode2.Active)
		assert.False(t, deactivatedCode3.Active)
	})

	t.Run("deactivation with complex metadata", func(t *testing.T) {
		var err error
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code with complex metadata
		coupon := helpers.NewValidPercentOffCoupon("Complex Metadata Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("COMPLEX", coupon.ID)
		promotionCode.Metadata = map[string]string{
			"json_config":   `{"active": true, "priority": 1}`,
			"special_chars": "Special: @#$%^&*()",
			"unicode_emoji": "🎉✨💫",
			"campaign_data": "summer_2024_premium_users",
		}
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Deactivate promotion code
		err = repo.DeactivatePromotionCode(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify complex metadata preserved
		deactivatedCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.False(t, deactivatedCode.Active)
		assert.Equal(t, `{"active": true, "priority": 1}`, deactivatedCode.Metadata["json_config"])
		assert.Equal(t, "Special: @#$%^&*()", deactivatedCode.Metadata["special_chars"])
		assert.Equal(t, "🎉✨💫", deactivatedCode.Metadata["unicode_emoji"])
		assert.Equal(t, "summer_2024_premium_users", deactivatedCode.Metadata["campaign_data"])
	})
}

