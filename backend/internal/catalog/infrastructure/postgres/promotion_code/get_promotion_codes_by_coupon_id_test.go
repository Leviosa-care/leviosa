package promotionCodeRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/catalog/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPromotionCodesByCouponID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful retrieval with multiple codes", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Multi Code Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create multiple promotion codes for the same coupon
		code1 := helpers.NewValidPromotionCode("CODE1", coupon.ID)
		code2 := helpers.NewValidPromotionCode("CODE2", coupon.ID)
		code3 := helpers.NewValidPromotionCode("CODE3", coupon.ID)

		helpers.InsertPromotionCode(t, ctx, testPool, code1)
		helpers.InsertPromotionCode(t, ctx, testPool, code2)
		helpers.InsertPromotionCode(t, ctx, testPool, code3)

		// Retrieve promotion codes by coupon ID
		codes, err := repo.GetPromotionCodesByCouponID(ctx, coupon.ID)

		require.NoError(t, err)
		assert.Len(t, codes, 3)

		// Verify all codes belong to the same coupon
		codeMap := make(map[string]bool)
		for _, code := range codes {
			assert.Equal(t, coupon.ID, code.CouponID)
			codeMap[code.Code] = true
		}

		assert.True(t, codeMap["CODE1"])
		assert.True(t, codeMap["CODE2"])
		assert.True(t, codeMap["CODE3"])
	})

	t.Run("successful retrieval with single code", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Single Code Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create single promotion code
		promotionCode := helpers.NewValidPromotionCode("SINGLE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve promotion codes by coupon ID
		codes, err := repo.GetPromotionCodesByCouponID(ctx, coupon.ID)

		require.NoError(t, err)
		assert.Len(t, codes, 1)
		assert.Equal(t, promotionCode.ID, codes[0].ID)
		assert.Equal(t, "SINGLE", codes[0].Code)
		assert.Equal(t, coupon.ID, codes[0].CouponID)
	})

	t.Run("empty result for coupon with no codes", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon without any promotion codes
		coupon := helpers.NewValidPercentOffCoupon("No Codes Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Retrieve promotion codes by coupon ID
		codes, err := repo.GetPromotionCodesByCouponID(ctx, coupon.ID)

		require.NoError(t, err)
		assert.Empty(t, codes)
	})

	t.Run("empty result for non-existent coupon", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		nonExistentCouponID := uuid.New()

		// Retrieve promotion codes by non-existent coupon ID
		codes, err := repo.GetPromotionCodesByCouponID(ctx, nonExistentCouponID)

		require.NoError(t, err)
		assert.Empty(t, codes)
	})

	t.Run("excludes codes from other coupons", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create multiple coupons
		coupon1 := helpers.NewValidPercentOffCoupon("Coupon One")
		coupon2 := helpers.NewValidPercentOffCoupon("Coupon Two")
		helpers.InsertCoupon(t, ctx, testPool, coupon1)
		helpers.InsertCoupon(t, ctx, testPool, coupon2)

		// Create promotion codes for different coupons
		code1 := helpers.NewValidPromotionCode("COUPON1CODE", coupon1.ID)
		code2 := helpers.NewValidPromotionCode("COUPON2CODE", coupon2.ID)
		code3 := helpers.NewValidPromotionCode("COUPON1CODE2", coupon1.ID)

		helpers.InsertPromotionCode(t, ctx, testPool, code1)
		helpers.InsertPromotionCode(t, ctx, testPool, code2)
		helpers.InsertPromotionCode(t, ctx, testPool, code3)

		// Retrieve codes only for coupon1
		codes, err := repo.GetPromotionCodesByCouponID(ctx, coupon1.ID)

		require.NoError(t, err)
		assert.Len(t, codes, 2)

		// Verify only coupon1's codes are returned
		for _, code := range codes {
			assert.Equal(t, coupon1.ID, code.CouponID)
			assert.True(t, code.Code == "COUPON1CODE" || code.Code == "COUPON1CODE2")
		}
	})

	t.Run("includes all code states", func(t *testing.T) {
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

		// Retrieve all codes
		codes, err := repo.GetPromotionCodesByCouponID(ctx, coupon.ID)

		require.NoError(t, err)
		assert.Len(t, codes, 3)

		// Verify all states are included
		stateMap := make(map[string]bool)
		for _, code := range codes {
			stateMap[code.Code] = code.Active
		}

		assert.True(t, stateMap["ACTIVE"])
		assert.False(t, stateMap["INACTIVE"])
		assert.True(t, stateMap["EXPIRED"]) // expired but still active flag
	})

	t.Run("includes codes with various restrictions and limits", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Restrictions Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create codes with different restrictions
		basicCode := helpers.NewValidPromotionCode("BASIC", coupon.ID)
		restrictedCode := helpers.NewValidPromotionCodeWithRestrictions("RESTRICTED", coupon.ID, []string{"premium_user"})
		limitedCode := helpers.NewValidPromotionCodeWithRedemptionLimits("LIMITED", coupon.ID, 50)

		helpers.InsertPromotionCode(t, ctx, testPool, basicCode)
		helpers.InsertPromotionCode(t, ctx, testPool, restrictedCode)
		helpers.InsertPromotionCode(t, ctx, testPool, limitedCode)

		// Retrieve all codes
		codes, err := repo.GetPromotionCodesByCouponID(ctx, coupon.ID)

		require.NoError(t, err)
		assert.Len(t, codes, 3)

		// Verify different types are included
		codeDetails := make(map[string]*domain.PromotionCode)
		for i, code := range codes {
			codeDetails[code.Code] = codes[i]
		}

		// Basic code
		assert.Nil(t, codeDetails["BASIC"].Restrictions)
		assert.Nil(t, codeDetails["BASIC"].MaxRedemptions)

		// Restricted code
		require.NotNil(t, codeDetails["RESTRICTED"].Restrictions)
		assert.Equal(t, []string{"premium_user"}, codeDetails["RESTRICTED"].Restrictions.CurrencyOptions)

		// Limited code
		require.NotNil(t, codeDetails["LIMITED"].MaxRedemptions)
		assert.Equal(t, 50, *codeDetails["LIMITED"].MaxRedemptions)
	})

	t.Run("returns codes ordered by creation date descending", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Ordered Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create codes with different creation times (need to use time delay or set CreatedAt)
		code1 := helpers.NewValidPromotionCode("FIRST", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, code1)

		code2 := helpers.NewValidPromotionCode("SECOND", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, code2)

		code3 := helpers.NewValidPromotionCode("THIRD", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, code3)

		// Retrieve codes
		codes, err := repo.GetPromotionCodesByCouponID(ctx, coupon.ID)

		require.NoError(t, err)
		assert.Len(t, codes, 3)

		// Verify ordering (newest first - descending order)
		// Since all codes are created in sequence, the last created should be first
		assert.Equal(t, "THIRD", codes[0].Code)
		assert.Equal(t, "SECOND", codes[1].Code)
		assert.Equal(t, "FIRST", codes[2].Code)
	})

	t.Run("includes metadata and timestamps", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon
		coupon := helpers.NewValidPercentOffCoupon("Metadata Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create code with metadata
		promotionCode := helpers.NewValidPromotionCode("METADATA", coupon.ID)
		promotionCode.Metadata = map[string]string{
			"campaign": "summer2024",
			"channel":  "email",
		}
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve codes
		codes, err := repo.GetPromotionCodesByCouponID(ctx, coupon.ID)

		require.NoError(t, err)
		assert.Len(t, codes, 1)

		// Verify metadata and timestamps
		code := codes[0]
		assert.Equal(t, "summer2024", code.Metadata["campaign"])
		assert.Equal(t, "email", code.Metadata["channel"])
		assert.NotZero(t, code.CreatedAt)
		assert.NotZero(t, code.UpdatedAt)
	})
}

