package promotionCodeRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePromotionCode(t *testing.T) {
	ctx := context.Background()

	t.Run("successful promotion code creation", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first (required for foreign key constraint)
		coupon := helpers.NewValidPercentOffCoupon("Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create promotion code
		promotionCode := helpers.NewValidPromotionCode("SAVE25", coupon.ID)
		codeID, err := repo.CreatePromotionCode(ctx, promotionCode)

		require.NoError(t, err)
		assert.NotEmpty(t, codeID)

		// Verify creation
		createdUUID, err := uuid.Parse(codeID)
		require.NoError(t, err)
		savedCode, err := helpers.GetPromotionCodeByID(t, ctx, createdUUID, testPool)
		require.NoError(t, err)
		assert.Equal(t, promotionCode.Code, savedCode.Code)
		assert.Equal(t, promotionCode.CouponID, savedCode.CouponID)
		assert.Equal(t, promotionCode.StripePromotionID, savedCode.StripePromotionID)
		assert.Equal(t, promotionCode.Active, savedCode.Active)
		assert.Equal(t, promotionCode.Metadata, savedCode.Metadata)
	})

	t.Run("successful creation with restrictions", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Restricted Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create promotion code with restrictions
		promotionCode := helpers.NewValidPromotionCodeWithRestrictions("RESTRICTED", coupon.ID, []string{"customer_123", "customer_456"})
		codeID, err := repo.CreatePromotionCode(ctx, promotionCode)

		require.NoError(t, err)

		// Verify restrictions
		createdUUID, err := uuid.Parse(codeID)
		require.NoError(t, err)
		savedCode, err := helpers.GetPromotionCodeByID(t, ctx, createdUUID, testPool)
		require.NoError(t, err)
		assert.Equal(t, []string{"customer_123", "customer_456"}, savedCode.Restrictions.CurrencyOptions)
	})

	t.Run("successful creation with expiry date", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Expiring Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create promotion code with expiry
		expiryDate := time.Now().Add(7 * 24 * time.Hour).Truncate(time.Microsecond)
		promotionCode := helpers.NewValidPromotionCodeWithExpiry("EXPIRES", coupon.ID, expiryDate)
		codeID, err := repo.CreatePromotionCode(ctx, promotionCode)

		require.NoError(t, err)

		// Verify expiry
		createdUUID, err := uuid.Parse(codeID)
		require.NoError(t, err)
		savedCode, err := helpers.GetPromotionCodeByID(t, ctx, createdUUID, testPool)
		require.NoError(t, err)
		assert.NotNil(t, savedCode.ExpiresAt)
		assert.True(t, savedCode.ExpiresAt.Equal(expiryDate))
	})

	t.Run("successful creation with redemption limits", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Limited Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create promotion code with limits
		promotionCode := helpers.NewValidPromotionCodeWithRedemptionLimits("LIMITED", coupon.ID, 50)
		codeID, err := repo.CreatePromotionCode(ctx, promotionCode)

		require.NoError(t, err)

		// Verify limits
		createdUUID, err := uuid.Parse(codeID)
		require.NoError(t, err)
		savedCode, err := helpers.GetPromotionCodeByID(t, ctx, createdUUID, testPool)
		require.NoError(t, err)
		require.NotNil(t, savedCode.MaxRedemptions)
		assert.Equal(t, 50, *savedCode.MaxRedemptions)
		assert.Equal(t, 0, savedCode.TimesRedeemed)
	})

	t.Run("creation with existing code should fail", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Duplicate Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create first promotion code
		promotionCode1 := helpers.NewValidPromotionCode("DUPLICATE", coupon.ID)
		_, err := repo.CreatePromotionCode(ctx, promotionCode1)
		require.NoError(t, err)

		// Try to create another with same code
		promotionCode2 := helpers.NewValidPromotionCode("DUPLICATE", coupon.ID)
		_, err = repo.CreatePromotionCode(ctx, promotionCode2)

		assert.Error(t, err)
	})

	t.Run("creation with non-existent coupon ID should fail", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		nonExistentCouponID := uuid.New()
		promotionCode := helpers.NewValidPromotionCode("ORPHAN", nonExistentCouponID)

		_, err := repo.CreatePromotionCode(ctx, promotionCode)

		assert.Error(t, err)
	})

	t.Run("creation with empty code should fail", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Empty Code Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Try to create promotion code with empty code
		promotionCode := helpers.NewValidPromotionCode("", coupon.ID)

		_, err := repo.CreatePromotionCode(ctx, promotionCode)

		assert.Error(t, err)
	})

	t.Run("successful creation with metadata", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Metadata Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create promotion code with metadata
		promotionCode := helpers.NewValidPromotionCode("METADATA", coupon.ID)
		promotionCode.Metadata = map[string]string{
			"campaign": "summer2024",
			"channel":  "email",
		}
		codeID, err := repo.CreatePromotionCode(ctx, promotionCode)

		require.NoError(t, err)

		// Verify metadata
		createdUUID, err := uuid.Parse(codeID)
		require.NoError(t, err)
		savedCode, err := helpers.GetPromotionCodeByID(t, ctx, createdUUID, testPool)
		require.NoError(t, err)
		assert.Equal(t, "summer2024", savedCode.Metadata["campaign"])
		assert.Equal(t, "email", savedCode.Metadata["channel"])
	})

	t.Run("successful creation with all optional fields", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Complete Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Create promotion code with all fields
		expiryDate := time.Now().Add(30 * 24 * time.Hour).Truncate(time.Microsecond)
		promotionCode := helpers.NewValidPromotionCodeWithExpiry("COMPLETE", coupon.ID, expiryDate, 100, []string{"vip_customer"})
		codeID, err := repo.CreatePromotionCode(ctx, promotionCode)

		require.NoError(t, err)

		// Verify all fields
		createdUUID, err := uuid.Parse(codeID)
		require.NoError(t, err)
		savedCode, err := helpers.GetPromotionCodeByID(t, ctx, createdUUID, testPool)
		require.NoError(t, err)
		assert.NotNil(t, savedCode.ExpiresAt)
		assert.NotNil(t, savedCode.MaxRedemptions)
		assert.Equal(t, 100, *savedCode.MaxRedemptions)
		assert.Equal(t, []string{"vip_customer"}, savedCode.Restrictions.CurrencyOptions)
		assert.True(t, savedCode.Active)
	})
}

