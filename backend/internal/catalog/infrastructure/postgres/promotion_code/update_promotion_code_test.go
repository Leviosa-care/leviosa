package promotionCodeRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdatePromotionCode(t *testing.T) {
	ctx := context.Background()

	t.Run("successful metadata update", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Update Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("UPDATE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Update promotion code metadata
		newMetadata := map[string]string{
			"updated":  "true",
			"version":  "2.0",
			"campaign": "spring2024",
		}
		updateReq := &domain.UpdatePromotionCodeRequest{
			Metadata: newMetadata,
		}

		err := repo.UpdatePromotionCode(ctx, promotionCode.ID, updateReq)

		require.NoError(t, err)

		// Verify update
		updatedCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, newMetadata, updatedCode.Metadata)
		assert.Equal(t, promotionCode.Code, updatedCode.Code)     // Unchanged
		assert.Equal(t, promotionCode.Active, updatedCode.Active) // Unchanged
	})

	t.Run("successful complete update", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Complete Update Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("COMPLETE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Update both fields
		newMetadata := map[string]string{
			"status":  "updated",
			"version": "3.0",
		}
		updateReq := &domain.UpdatePromotionCodeRequest{
			Metadata: newMetadata,
		}

		err := repo.UpdatePromotionCode(ctx, promotionCode.ID, updateReq)

		require.NoError(t, err)

		// Verify update
		updatedCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, newMetadata, updatedCode.Metadata)
		assert.Equal(t, promotionCode.Code, updatedCode.Code) // Unchanged
	})

	t.Run("update with empty metadata", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code with metadata
		coupon := helpers.NewValidPercentOffCoupon("Clear Metadata Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("CLEARMETA", coupon.ID)
		promotionCode.Metadata = map[string]string{"original": "data"}
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Clear metadata
		emptyMetadata := map[string]string{}
		updateReq := &domain.UpdatePromotionCodeRequest{
			Metadata: emptyMetadata,
		}

		err := repo.UpdatePromotionCode(ctx, promotionCode.ID, updateReq)

		require.NoError(t, err)

		// Verify metadata was cleared
		updatedCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.Empty(t, updatedCode.Metadata)
	})

	t.Run("update with empty request should do nothing", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Empty Update Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("EMPTY", coupon.ID)
		originalMetadata := promotionCode.Metadata
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Empty update request
		updateReq := &domain.UpdatePromotionCodeRequest{}

		err := repo.UpdatePromotionCode(ctx, promotionCode.ID, updateReq)

		require.NoError(t, err)

		// Verify nothing changed
		updatedCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, promotionCode.Code, updatedCode.Code)
		assert.Equal(t, originalMetadata, updatedCode.Metadata)
		assert.Equal(t, promotionCode.Active, updatedCode.Active)
	})

	t.Run("update non-existent promotion code should return not found error", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)

		nonExistentID := uuid.New()
		newMetadata := map[string]string{"test": "data"}
		updateReq := &domain.UpdatePromotionCodeRequest{
			Metadata: newMetadata,
		}

		err := repo.UpdatePromotionCode(ctx, nonExistentID, updateReq)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("update with nil metadata pointer", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Nil Metadata Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("NILMETA", coupon.ID)
		originalMetadata := promotionCode.Metadata
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Update with nil metadata (should be ignored)
		updateReq := &domain.UpdatePromotionCodeRequest{
			Metadata: nil,
		}

		err := repo.UpdatePromotionCode(ctx, promotionCode.ID, updateReq)

		require.NoError(t, err)

		// Verify metadata unchanged
		updatedCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, promotionCode.Code, updatedCode.Code)
		assert.Equal(t, originalMetadata, updatedCode.Metadata) // Unchanged
	})

	t.Run("update with complex metadata", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Complex Metadata Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("COMPLEX", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Update with complex metadata
		complexMetadata := map[string]string{
			"campaign_id":      "camp_12345",
			"user_segment":     "premium_users",
			"acquisition_cost": "15.50",
			"channel":          "email",
			"ab_test_variant":  "variant_b",
		}
		updateReq := &domain.UpdatePromotionCodeRequest{
			Metadata: complexMetadata,
		}

		err := repo.UpdatePromotionCode(ctx, promotionCode.ID, updateReq)

		require.NoError(t, err)

		// Verify complex metadata
		updatedCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, "camp_12345", updatedCode.Metadata["campaign_id"])
		assert.Equal(t, "premium_users", updatedCode.Metadata["user_segment"])
		assert.Equal(t, "15.50", updatedCode.Metadata["acquisition_cost"])
		assert.Equal(t, "email", updatedCode.Metadata["channel"])
		assert.Equal(t, "variant_b", updatedCode.Metadata["ab_test_variant"])
	})

	t.Run("update preserves other fields", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code with all fields
		coupon := helpers.NewValidPercentOffCoupon("Preserve Fields Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCodeWithRestrictions("PRESERVE", coupon.ID, []string{"test_customer"})
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Update only metadata
		newMetadata := map[string]string{"updated": "metadata_only"}
		updateReq := &domain.UpdatePromotionCodeRequest{
			Metadata: newMetadata,
		}

		err := repo.UpdatePromotionCode(ctx, promotionCode.ID, updateReq)

		require.NoError(t, err)

		// Verify other fields preserved
		updatedCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, promotionCode.Code, updatedCode.Code)                                                 // Unchanged
		assert.Equal(t, promotionCode.CouponID, updatedCode.CouponID)                                         // Unchanged
		assert.Equal(t, promotionCode.StripePromotionID, updatedCode.StripePromotionID)                       // Unchanged
		assert.Equal(t, promotionCode.Active, updatedCode.Active)                                             // Unchanged
		assert.Equal(t, promotionCode.Restrictions.CurrencyOptions, updatedCode.Restrictions.CurrencyOptions) // Unchanged
		assert.Equal(t, promotionCode.TimesRedeemed, updatedCode.TimesRedeemed)                               // Unchanged
		assert.Equal(t, newMetadata, updatedCode.Metadata)                                                    // Changed
	})

	t.Run("update with special characters in metadata", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Special Chars Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("SPECIAL", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Update with special characters in metadata values
		specialMetadata := map[string]string{
			"json_data":     `{"key": "value", "number": 123}`,
			"special_chars": "Special chars: @#$%^&*()_+-=[]{}|;':\",./<>?",
			"unicode":       "Unicode: 🎉 🚀 💡",
		}
		updateReq := &domain.UpdatePromotionCodeRequest{
			Metadata: specialMetadata,
		}

		err := repo.UpdatePromotionCode(ctx, promotionCode.ID, updateReq)

		require.NoError(t, err)

		// Verify special characters preserved
		updatedCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, `{"key": "value", "number": 123}`, updatedCode.Metadata["json_data"])
		assert.Equal(t, "Special chars: @#$%^&*()_+-=[]{}|;':\",./<>?", updatedCode.Metadata["special_chars"])
		assert.Equal(t, "Unicode: 🎉 🚀 💡", updatedCode.Metadata["unicode"])
	})
}

