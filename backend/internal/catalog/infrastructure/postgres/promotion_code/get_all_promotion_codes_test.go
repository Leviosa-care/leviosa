package promotionCodeRepository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllPromotionCodes(t *testing.T) {
	ctx := context.Background()

	t.Run("retrieves all promotion codes", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create multiple coupons
		coupon1 := helpers.NewValidPercentOffCoupon("All Codes Coupon 1")
		coupon2 := helpers.NewValidPercentOffCoupon("All Codes Coupon 2")
		helpers.InsertCoupon(t, ctx, testPool, coupon1)
		helpers.InsertCoupon(t, ctx, testPool, coupon2)

		// Insert multiple promotion codes
		activeCode := helpers.NewValidPromotionCode("ACTIVE", coupon1.ID)
		inactiveCode := helpers.NewInactivePromotionCode("INACTIVE", coupon1.ID)
		expiredCode := helpers.NewExpiredPromotionCode("EXPIRED", coupon2.ID)
		restrictedCode := helpers.NewValidPromotionCodeWithRestrictions("RESTRICTED", coupon2.ID, []string{"premium"})

		helpers.InsertPromotionCode(t, ctx, testPool, activeCode)
		helpers.InsertPromotionCode(t, ctx, testPool, inactiveCode)
		helpers.InsertPromotionCode(t, ctx, testPool, expiredCode)
		helpers.InsertPromotionCode(t, ctx, testPool, restrictedCode)

		// Get all promotion codes
		codes, err := repo.GetAllPromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 4)

		// Verify all codes are included
		codeMap := make(map[string]bool)
		for _, code := range codes {
			codeMap[code.Code] = true
		}

		assert.True(t, codeMap["ACTIVE"])
		assert.True(t, codeMap["INACTIVE"])
		assert.True(t, codeMap["EXPIRED"])
		assert.True(t, codeMap["RESTRICTED"])
	})

	t.Run("empty result when no codes exist", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Get all promotion codes from empty table
		codes, err := repo.GetAllPromotionCodes(ctx)

		require.NoError(t, err)
		assert.Empty(t, codes)
	})

	t.Run("includes codes with all possible states", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("All States Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert codes in different states
		activeCode := helpers.NewValidPromotionCode("ACTIVE", coupon.ID)
		inactiveCode := helpers.NewInactivePromotionCode("INACTIVE", coupon.ID)
		expiredCode := helpers.NewExpiredPromotionCode("EXPIRED", coupon.ID)
		limitedCode := helpers.NewValidPromotionCodeWithRedemptionLimits("LIMITED", coupon.ID, 10)
		limitedCode.TimesRedeemed = 10 // At limit

		helpers.InsertPromotionCode(t, ctx, testPool, activeCode)
		helpers.InsertPromotionCode(t, ctx, testPool, inactiveCode)
		helpers.InsertPromotionCode(t, ctx, testPool, expiredCode)
		helpers.InsertPromotionCode(t, ctx, testPool, limitedCode)

		// Get all codes
		codes, err := repo.GetAllPromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 4)

		// Verify all states are included
		stateMap := make(map[string]map[string]interface{})
		for _, code := range codes {
			stateMap[code.Code] = map[string]interface{}{
				"isActive":       code.Active,
				"expiresAt":      code.ExpiresAt,
				"timesRedeemed":  code.TimesRedeemed,
				"maxRedemptions": code.MaxRedemptions,
			}
		}

		assert.True(t, stateMap["ACTIVE"]["isActive"].(bool))
		assert.False(t, stateMap["INACTIVE"]["isActive"].(bool))
		assert.NotNil(t, stateMap["EXPIRED"]["expiresAt"])
		assert.Equal(t, 10, stateMap["LIMITED"]["timesRedeemed"].(int))
	})

	t.Run("includes all fields and metadata", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("All Fields Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code with all fields
		promotionCode := helpers.NewValidPromotionCodeWithRestrictions("COMPLETE", coupon.ID, []string{"premium", "vip"})
		promotionCode.Metadata = map[string]string{
			"campaign": "all_fields_test",
			"channel":  "email",
			"priority": "high",
		}
		promotionCode.TimesRedeemed = 15
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Get all codes
		codes, err := repo.GetAllPromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 1)

		code := codes[0]
		assert.Equal(t, "COMPLETE", code.Code)
		assert.Equal(t, coupon.ID, code.CouponID)
		assert.NotEmpty(t, code.StripePromotionID)
		assert.True(t, code.Active)
		assert.Equal(t, []string{"premium", "vip"}, code.Restrictions.CurrencyOptions)
		assert.Equal(t, "all_fields_test", code.Metadata["campaign"])
		assert.Equal(t, "email", code.Metadata["channel"])
		assert.Equal(t, "high", code.Metadata["priority"])
		assert.Equal(t, 15, code.TimesRedeemed)
		assert.NotZero(t, code.CreatedAt)
		assert.NotZero(t, code.UpdatedAt)
	})

	t.Run("returns codes ordered by creation date descending", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Ordered All Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert multiple codes in sequence
		code1 := helpers.NewValidPromotionCode("FIRST", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, code1)

		code2 := helpers.NewValidPromotionCode("SECOND", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, code2)

		code3 := helpers.NewValidPromotionCode("THIRD", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, code3)

		// Get all codes
		codes, err := repo.GetAllPromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 3)

		// Should be ordered by creation date descending (newest first)
		assert.Equal(t, "THIRD", codes[0].Code)
		assert.Equal(t, "SECOND", codes[1].Code)
		assert.Equal(t, "FIRST", codes[2].Code)
	})

	t.Run("includes codes from multiple coupons", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create multiple coupons
		coupon1 := helpers.NewValidPercentOffCoupon("Multi Coupon 1")
		coupon2 := helpers.NewValidPercentOffCoupon("Multi Coupon 2")
		coupon3 := helpers.NewValidPercentOffCoupon("Multi Coupon 3")
		helpers.InsertCoupon(t, ctx, testPool, coupon1)
		helpers.InsertCoupon(t, ctx, testPool, coupon2)
		helpers.InsertCoupon(t, ctx, testPool, coupon3)

		// Insert codes for different coupons
		code1 := helpers.NewValidPromotionCode("COUPON1CODE", coupon1.ID)
		code2a := helpers.NewValidPromotionCode("COUPON2CODEA", coupon2.ID)
		code2b := helpers.NewValidPromotionCode("COUPON2CODEB", coupon2.ID)
		code3 := helpers.NewValidPromotionCode("COUPON3CODE", coupon3.ID)

		helpers.InsertPromotionCode(t, ctx, testPool, code1)
		helpers.InsertPromotionCode(t, ctx, testPool, code2a)
		helpers.InsertPromotionCode(t, ctx, testPool, code2b)
		helpers.InsertPromotionCode(t, ctx, testPool, code3)

		// Get all codes
		codes, err := repo.GetAllPromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 4)

		// Verify codes from all coupons are included
		couponCodeCount := make(map[string]int)
		for _, code := range codes {
			couponID := code.CouponID.String()
			couponCodeCount[couponID]++
		}

		assert.Equal(t, 1, couponCodeCount[coupon1.ID.String()])
		assert.Equal(t, 2, couponCodeCount[coupon2.ID.String()])
		assert.Equal(t, 1, couponCodeCount[coupon3.ID.String()])
	})

	t.Run("handles large number of codes", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Large Scale Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert many promotion codes
		numberOfCodes := 50
		expectedCodes := make(map[string]bool)
		for i := 0; i < numberOfCodes; i++ {
			codeName := fmt.Sprintf("CODE%03d", i+1)
			promotionCode := helpers.NewValidPromotionCode(codeName, coupon.ID)
			helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)
			expectedCodes[codeName] = true
		}

		// Get all codes
		codes, err := repo.GetAllPromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, numberOfCodes)

		// Verify all codes are present
		actualCodes := make(map[string]bool)
		for _, code := range codes {
			actualCodes[code.Code] = true
		}

		for expectedCode := range expectedCodes {
			assert.True(t, actualCodes[expectedCode], "Expected code %s not found", expectedCode)
		}
	})

	t.Run("includes codes with complex metadata", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Complex Metadata Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code with complex metadata
		promotionCode := helpers.NewValidPromotionCode("COMPLEX", coupon.ID)
		promotionCode.Metadata = map[string]string{
			"json_config":      `{"settings": {"active": true, "priority": 5}}`,
			"special_chars":    "Test: @#$%^&*()_+-=[]{}|;':\",./<>?",
			"unicode_content":  "Emojis: 🎉 🚀 💡 ✨ 🎯",
			"long_description": "This is a very long description that contains multiple sentences and should test the handling of large text content in metadata fields.",
			"numeric_string":   "123456789",
			"boolean_string":   "true",
		}
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Get all codes
		codes, err := repo.GetAllPromotionCodes(ctx)

		require.NoError(t, err)
		assert.Len(t, codes, 1)

		code := codes[0]
		assert.Equal(t, `{"settings": {"active": true, "priority": 5}}`, code.Metadata["json_config"])
		assert.Equal(t, "Test: @#$%^&*()_+-=[]{}|;':\",./<>?", code.Metadata["special_chars"])
		assert.Equal(t, "Emojis: 🎉 🚀 💡 ✨ 🎯", code.Metadata["unicode_content"])
		assert.Contains(t, code.Metadata["long_description"], "very long description")
		assert.Equal(t, "123456789", code.Metadata["numeric_string"])
		assert.Equal(t, "true", code.Metadata["boolean_string"])
	})
}

