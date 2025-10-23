package couponRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/catalog/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateCoupon(t *testing.T) {
	ctx := context.Background()

	t.Run("successful name update", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon
		coupon := helpers.NewValidPercentOffCoupon("Original Name")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Update coupon name
		newName := "Updated Name"
		updateReq := &domain.UpdateCouponRequest{
			Name: &newName,
		}

		err := repo.UpdateCoupon(ctx, coupon.ID, updateReq)

		require.NoError(t, err)

		// Verify update
		updatedCoupon, err := helpers.GetCouponByID(t, ctx, coupon.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, newName, updatedCoupon.Name)
		assert.Equal(t, coupon.PercentOff, updatedCoupon.PercentOff) // Unchanged
	})

	t.Run("successful metadata update", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon
		coupon := helpers.NewValidPercentOffCoupon("Metadata Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Update coupon metadata
		newMetadata := map[string]string{
			"updated": "true",
			"version": "2.0",
		}
		updateReq := &domain.UpdateCouponRequest{
			Metadata: newMetadata,
		}

		err := repo.UpdateCoupon(ctx, coupon.ID, updateReq)

		require.NoError(t, err)

		// Verify update
		updatedCoupon, err := helpers.GetCouponByID(t, ctx, coupon.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, newMetadata, updatedCoupon.Metadata)
		assert.Equal(t, coupon.Name, updatedCoupon.Name) // Unchanged
	})

	t.Run("successful complete update", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon
		coupon := helpers.NewValidPercentOffCoupon("Complete Update Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Update both name and metadata
		newName := "Completely Updated"
		newMetadata := map[string]string{
			"status":  "updated",
			"version": "3.0",
		}
		updateReq := &domain.UpdateCouponRequest{
			Name:     &newName,
			Metadata: newMetadata,
		}

		err := repo.UpdateCoupon(ctx, coupon.ID, updateReq)

		require.NoError(t, err)

		// Verify update
		updatedCoupon, err := helpers.GetCouponByID(t, ctx, coupon.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, newName, updatedCoupon.Name)
		assert.Equal(t, newMetadata, updatedCoupon.Metadata)
		assert.Equal(t, coupon.PercentOff, updatedCoupon.PercentOff) // Unchanged
	})

	t.Run("update with empty metadata", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon with metadata
		coupon := helpers.NewValidPercentOffCoupon("Clear Metadata Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Clear metadata
		emptyMetadata := map[string]string{}
		updateReq := &domain.UpdateCouponRequest{
			Metadata: emptyMetadata,
		}

		err := repo.UpdateCoupon(ctx, coupon.ID, updateReq)

		require.NoError(t, err)

		// Verify metadata was cleared
		updatedCoupon, err := helpers.GetCouponByID(t, ctx, coupon.ID, testPool)
		require.NoError(t, err)
		assert.Empty(t, updatedCoupon.Metadata)
	})

	t.Run("update with empty request should do nothing", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon
		coupon := helpers.NewValidPercentOffCoupon("Empty Update Test")
		originalMetadata := coupon.Metadata
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Empty update request
		updateReq := &domain.UpdateCouponRequest{}

		err := repo.UpdateCoupon(ctx, coupon.ID, updateReq)

		require.NoError(t, err)

		// Verify nothing changed
		updatedCoupon, err := helpers.GetCouponByID(t, ctx, coupon.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, coupon.Name, updatedCoupon.Name)
		assert.Equal(t, originalMetadata, updatedCoupon.Metadata)
	})

	t.Run("update non-existent coupon should return not found error", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		nonExistentID := uuid.New()
		newName := "Non-existent Update"
		updateReq := &domain.UpdateCouponRequest{
			Name: &newName,
		}

		err := repo.UpdateCoupon(ctx, nonExistentID, updateReq)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("update with nil name pointer", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon
		coupon := helpers.NewValidPercentOffCoupon("Nil Name Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Update with nil name (should be ignored)
		newMetadata := map[string]string{"updated": "metadata_only"}
		updateReq := &domain.UpdateCouponRequest{
			Name:     nil,
			Metadata: newMetadata,
		}

		err := repo.UpdateCoupon(ctx, coupon.ID, updateReq)

		require.NoError(t, err)

		// Verify name unchanged, metadata updated
		updatedCoupon, err := helpers.GetCouponByID(t, ctx, coupon.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, coupon.Name, updatedCoupon.Name) // Unchanged
		assert.Equal(t, newMetadata, updatedCoupon.Metadata)
	})

	t.Run("update with nil metadata pointer", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon
		coupon := helpers.NewValidPercentOffCoupon("Nil Metadata Test")
		originalMetadata := coupon.Metadata
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Update with nil metadata (should be ignored)
		newName := "Updated Name Only"
		updateReq := &domain.UpdateCouponRequest{
			Name:     &newName,
			Metadata: nil,
		}

		err := repo.UpdateCoupon(ctx, coupon.ID, updateReq)

		require.NoError(t, err)

		// Verify name updated, metadata unchanged
		updatedCoupon, err := helpers.GetCouponByID(t, ctx, coupon.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, newName, updatedCoupon.Name)
		assert.Equal(t, originalMetadata, updatedCoupon.Metadata) // Unchanged
	})
}

