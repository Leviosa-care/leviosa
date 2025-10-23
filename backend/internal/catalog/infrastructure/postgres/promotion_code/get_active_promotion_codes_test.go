package promotionCodeRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/catalog/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetActivePromotionCodes(t *testing.T) {
	ctx := context.Background()

	t.Run("retrieves only active codes", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Active Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert active promotion code
		activeCode := helpers.NewValidPromotionCode("ACTIVE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, activeCode)

		// Insert inactive promotion code
		inactiveCode := helpers.NewInactivePromotionCode("INACTIVE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, inactiveCode)

		// Get active codes
		codes, err := repo.GetActivePromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 1)
		assert.Equal(t, activeCode.ID, codes[0].ID)
		assert.True(t, codes[0].Active)
	})

	t.Run("excludes expired codes", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Expiry Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert active promotion code
		activeCode := helpers.NewValidPromotionCode("ACTIVE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, activeCode)

		// Insert expired promotion code
		expiredCode := helpers.NewExpiredPromotionCode("EXPIRED", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, expiredCode)

		// Get active codes
		codes, err := repo.GetActivePromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 1)
		assert.Equal(t, activeCode.ID, codes[0].ID)
	})

	t.Run("excludes codes at redemption limit", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Limit Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert active promotion code
		activeCode := helpers.NewValidPromotionCode("ACTIVE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, activeCode)

		// Insert promotion code at redemption limit
		limitCode := helpers.NewValidPromotionCodeWithRedemptionLimits("LIMITED", coupon.ID, 10)
		limitCode.TimesRedeemed = 10
		helpers.InsertPromotionCode(t, ctx, testPool, limitCode)

		// Get active codes
		codes, err := repo.GetActivePromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 1)
		assert.Equal(t, activeCode.ID, codes[0].ID)
	})

	t.Run("includes codes under redemption limit", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Under Limit Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code under limit
		underLimitCode := helpers.NewValidPromotionCodeWithRedemptionLimits("UNDER", coupon.ID, 10)
		underLimitCode.TimesRedeemed = 5
		helpers.InsertPromotionCode(t, ctx, testPool, underLimitCode)

		// Get active codes
		codes, err := repo.GetActivePromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 1)
		assert.Equal(t, underLimitCode.ID, codes[0].ID)
		assert.Equal(t, 5, codes[0].TimesRedeemed)
	})

	t.Run("includes codes with no redemption limit", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Unlimited Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code without redemption limit
		unlimitedCode := helpers.NewValidPromotionCode("UNLIMITED", coupon.ID)
		unlimitedCode.MaxRedemptions = nil
		unlimitedCode.TimesRedeemed = 1000
		helpers.InsertPromotionCode(t, ctx, testPool, unlimitedCode)

		// Get active codes
		codes, err := repo.GetActivePromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 1)
		assert.Equal(t, unlimitedCode.ID, codes[0].ID)
	})

	t.Run("includes codes with future expiry", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Future Expiry Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code with future expiry
		futureDate := time.Now().Add(24 * time.Hour)
		futureCode := helpers.NewValidPromotionCodeWithExpiry("FUTURE", coupon.ID, futureDate)
		helpers.InsertPromotionCode(t, ctx, testPool, futureCode)

		// Get active codes
		codes, err := repo.GetActivePromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 1)
		assert.Equal(t, futureCode.ID, codes[0].ID)
	})

	t.Run("includes codes with no expiry date", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("No Expiry Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code without expiry
		noExpiryCode := helpers.NewValidPromotionCode("NOEXPIRY", coupon.ID)
		noExpiryCode.ExpiresAt = nil
		helpers.InsertPromotionCode(t, ctx, testPool, noExpiryCode)

		// Get active codes
		codes, err := repo.GetActivePromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 1)
		assert.Equal(t, noExpiryCode.ID, codes[0].ID)
	})

	t.Run("returns empty list when no active codes", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("No Active Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert only inactive/expired codes
		inactiveCode := helpers.NewInactivePromotionCode("INACTIVE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, inactiveCode)

		expiredCode := helpers.NewExpiredPromotionCode("EXPIRED", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, expiredCode)

		// Get active codes
		codes, err := repo.GetActivePromotionCodes(ctx)

		require.NoError(t, err)
		assert.Empty(t, codes)
	})

	t.Run("returns codes ordered by creation date descending", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Ordered Active Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert multiple active promotion codes
		code1 := helpers.NewValidPromotionCode("FIRST", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, code1)

		code2 := helpers.NewValidPromotionCode("SECOND", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, code2)

		code3 := helpers.NewValidPromotionCode("THIRD", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, code3)

		// Get active codes
		codes, err := repo.GetActivePromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 3)

		// Should be ordered by creation date descending (newest first)
		assert.Equal(t, "THIRD", codes[0].Code)
		assert.Equal(t, "SECOND", codes[1].Code)
		assert.Equal(t, "FIRST", codes[2].Code)
	})

	t.Run("returns codes from multiple coupons", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create multiple coupons
		coupon1 := helpers.NewValidPercentOffCoupon("Multi Coupon 1")
		coupon2 := helpers.NewValidPercentOffCoupon("Multi Coupon 2")
		helpers.InsertCoupon(t, ctx, testPool, coupon1)
		helpers.InsertCoupon(t, ctx, testPool, coupon2)

		// Insert active codes for different coupons
		code1 := helpers.NewValidPromotionCode("COUPON1", coupon1.ID)
		code2 := helpers.NewValidPromotionCode("COUPON2", coupon2.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, code1)
		helpers.InsertPromotionCode(t, ctx, testPool, code2)

		// Get active codes
		codes, err := repo.GetActivePromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 2)

		// Verify codes from different coupons are included
		codeMap := make(map[string]string)
		for _, code := range codes {
			codeMap[code.Code] = code.CouponID.String()
		}

		assert.Equal(t, coupon1.ID.String(), codeMap["COUPON1"])
		assert.Equal(t, coupon2.ID.String(), codeMap["COUPON2"])
	})

	t.Run("includes codes with restrictions", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Restrictions Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert codes with different restrictions
		basicCode := helpers.NewValidPromotionCode("BASIC", coupon.ID)
		restrictedCode := helpers.NewValidPromotionCodeWithRestrictions("VIP", coupon.ID, []string{"vip_customer"})

		helpers.InsertPromotionCode(t, ctx, testPool, basicCode)
		helpers.InsertPromotionCode(t, ctx, testPool, restrictedCode)

		// Get active codes
		codes, err := repo.GetActivePromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 2)

		// Verify both types of codes are included
		codeDetails := make(map[string]*domain.PromotionCode)
		for i, code := range codes {
			codeDetails[code.Code] = codes[i]
		}

		assert.Nil(t, codeDetails["BASIC"].Restrictions)
		require.NotNil(t, codeDetails["VIP"].Restrictions)
		assert.Equal(t, []string{"vip_customer"}, codeDetails["VIP"].Restrictions.CurrencyOptions)
	})

	t.Run("includes metadata and all fields", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Complete Fields Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code with all fields
		promotionCode := helpers.NewValidPromotionCode("COMPLETE", coupon.ID)
		promotionCode.Metadata = map[string]string{
			"campaign": "active2024",
			"channel":  "web",
		}
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Get active codes
		codes, err := repo.GetActivePromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 1)

		code := codes[0]
		assert.Equal(t, "active2024", code.Metadata["campaign"])
		assert.Equal(t, "web", code.Metadata["channel"])
		assert.NotZero(t, code.CreatedAt)
		assert.NotZero(t, code.UpdatedAt)
		assert.NotEmpty(t, code.StripePromotionID)
		assert.True(t, code.Active)
	})
}

