package otpRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/test/helpers"
	"github.com/Leviosa-care/core/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidateOTP(t *testing.T) {
	ctx := context.Background()

	t.Run("should invalidate existing OTP successfully", func(t *testing.T) {
		// Clean up before test
		helpers.ClearOTPKeys(t, ctx, testClient)

		// Create and insert test OTP
		otp := helpers.NewTestOTP("invalidate@example.com")
		helpers.InsertOTP(t, ctx, otp, testClient, 10*time.Minute)

		// Verify OTP exists before invalidation
		exists := helpers.CheckOTPExists(t, ctx, otp.EmailHash, testClient)
		assert.True(t, exists, "OTP should exist before invalidation")

		// Invalidate OTP
		err := repo.InvalidateOTP(ctx, otp.EmailHash)
		require.NoError(t, err)

		// Verify OTP no longer exists
		exists = helpers.CheckOTPExists(t, ctx, otp.EmailHash, testClient)
		assert.False(t, exists, "OTP should not exist after invalidation")
	})

	t.Run("should handle invalidating non-existent OTP", func(t *testing.T) {
		// Clean up before test
		helpers.ClearOTPKeys(t, ctx, testClient)

		// Try to invalidate non-existent OTP
		err := repo.InvalidateOTP(ctx, "hash_nonexistent@example.com")
		require.NoError(t, err, "Invalidating non-existent OTP should not return error")

		// Verify no side effects
		exists := helpers.CheckOTPExists(t, ctx, "hash_nonexistent@example.com", testClient)
		assert.False(t, exists, "Non-existent OTP should remain non-existent")
	})

	t.Run("should not affect other OTP entries", func(t *testing.T) {
		// Clean up before test
		helpers.ClearOTPKeys(t, ctx, testClient)

		// Create and insert multiple OTPs
		otp1 := helpers.NewTestOTP("user1@example.com")
		otp2 := helpers.NewTestOTP("user2@example.com")
		otp3 := helpers.NewTestOTP("user3@example.com")

		helpers.InsertOTP(t, ctx, otp1, testClient, 10*time.Minute)
		helpers.InsertOTP(t, ctx, otp2, testClient, 10*time.Minute)
		helpers.InsertOTP(t, ctx, otp3, testClient, 10*time.Minute)

		// Verify all OTPs exist
		assert.True(t, helpers.CheckOTPExists(t, ctx, otp1.EmailHash, testClient))
		assert.True(t, helpers.CheckOTPExists(t, ctx, otp2.EmailHash, testClient))
		assert.True(t, helpers.CheckOTPExists(t, ctx, otp3.EmailHash, testClient))

		// Invalidate only the middle OTP
		err := repo.InvalidateOTP(ctx, otp2.EmailHash)
		require.NoError(t, err)

		// Verify only the targeted OTP was removed
		assert.True(t, helpers.CheckOTPExists(t, ctx, otp1.EmailHash, testClient), "OTP1 should still exist")
		assert.False(t, helpers.CheckOTPExists(t, ctx, otp2.EmailHash, testClient), "OTP2 should be invalidated")
		assert.True(t, helpers.CheckOTPExists(t, ctx, otp3.EmailHash, testClient), "OTP3 should still exist")
	})

	t.Run("should invalidate expired OTP", func(t *testing.T) {
		// Clean up before test
		helpers.ClearOTPKeys(t, ctx, testClient)

		// Create and insert OTP with very short TTL
		otp := helpers.NewTestOTP("expired@example.com")
		helpers.InsertOTP(t, ctx, otp, testClient, 100*time.Millisecond)

		// Wait for expiration
		time.Sleep(200 * time.Millisecond)

		// Try to invalidate expired OTP (should not error even if already expired)
		err := repo.InvalidateOTP(ctx, otp.EmailHash)
		require.NoError(t, err, "Invalidating expired OTP should not return error")

		// Verify OTP doesn't exist
		exists := helpers.CheckOTPExists(t, ctx, otp.EmailHash, testClient)
		assert.False(t, exists, "Expired OTP should not exist")
	})

	t.Run("should handle empty email hash", func(t *testing.T) {
		// Clean up before test
		helpers.ClearOTPKeys(t, ctx, testClient)

		// Create OTP with empty email hash (edge case)
		otp := helpers.NewTestOTP("empty@example.com")
		otp.EmailHash = ""
		helpers.InsertOTP(t, ctx, otp, testClient, 5*time.Minute)

		// Verify it exists
		exists := helpers.CheckOTPExists(t, ctx, "", testClient)
		assert.True(t, exists, "OTP with empty hash should exist")

		// Invalidate with empty email hash
		err := repo.InvalidateOTP(ctx, "")
		require.NoError(t, err, "Should handle empty email hash")

		// Verify it's gone
		exists = helpers.CheckOTPExists(t, ctx, "", testClient)
		assert.False(t, exists, "OTP with empty hash should be invalidated")
	})

	t.Run("should handle multiple invalidations of same OTP", func(t *testing.T) {
		// Clean up before test
		helpers.ClearOTPKeys(t, ctx, testClient)

		// Create and insert test OTP
		otp := helpers.NewTestOTP("double@example.com")
		helpers.InsertOTP(t, ctx, otp, testClient, 10*time.Minute)

		// First invalidation
		err := repo.InvalidateOTP(ctx, otp.EmailHash)
		require.NoError(t, err)

		// Second invalidation of same OTP (should not error)
		err = repo.InvalidateOTP(ctx, otp.EmailHash)
		require.NoError(t, err, "Multiple invalidations should not error")

		// Verify OTP doesn't exist
		exists := helpers.CheckOTPExists(t, ctx, otp.EmailHash, testClient)
		assert.False(t, exists, "OTP should not exist after multiple invalidations")
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		// Clean up before test
		helpers.ClearOTPKeys(t, ctx, testClient)

		// Create and insert OTP
		otp := helpers.NewTestOTP("cancelled@example.com")
		helpers.InsertOTP(t, ctx, otp, testClient, 10*time.Minute)

		// Create cancelled context
		cancelCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		// Attempt to invalidate with cancelled context
		err := repo.InvalidateOTP(cancelCtx, otp.EmailHash)
		assert.Error(t, err, "Should return error for cancelled context")
		assert.ErrorIs(t, err, errs.ErrContext, "Should return context error")

		// Verify OTP still exists (invalidation didn't complete)
		exists := helpers.CheckOTPExists(t, ctx, otp.EmailHash, testClient)
		assert.True(t, exists, "OTP should still exist after cancelled invalidation")
	})

	t.Run("should invalidate OTP with large data", func(t *testing.T) {
		// Clean up before test
		helpers.ClearOTPKeys(t, ctx, testClient)

		// Create OTP with large data
		otp := helpers.NewTestOTP("large@example.com")
		otp.DEK = make([]byte, 1024) // 1KB of data
		for i := range otp.DEK {
			otp.DEK[i] = byte(i % 256)
		}

		helpers.InsertOTP(t, ctx, otp, testClient, 10*time.Minute)

		// Verify large OTP exists
		exists := helpers.CheckOTPExists(t, ctx, otp.EmailHash, testClient)
		assert.True(t, exists, "Large OTP should exist")

		// Invalidate large OTP
		err := repo.InvalidateOTP(ctx, otp.EmailHash)
		require.NoError(t, err)

		// Verify large OTP is gone
		exists = helpers.CheckOTPExists(t, ctx, otp.EmailHash, testClient)
		assert.False(t, exists, "Large OTP should be invalidated")
	})

	// Clean up after all tests
	t.Cleanup(func() {
		helpers.ClearOTPKeys(t, ctx, testClient)
	})
}
