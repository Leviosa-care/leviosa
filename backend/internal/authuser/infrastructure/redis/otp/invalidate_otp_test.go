package otpRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestInvalidateOTP TEST_PATH=internal/authuser/infrastructure/redis/otp/invalidate_otp_test.go

func TestInvalidateOTP(t *testing.T) {
	ctx := context.Background()

	t.Run("should invalidate existing OTP successfully", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create and insert test OTP
		otpEncx := td.NewTestOTPEncx(t)

		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, 10*time.Minute)
		require.NoError(t, err)

		// Verify OTP exists before invalidation
		exists := td.CheckOTPExists(t, ctx, otpEncx.EmailHash, testClient)
		require.True(t, exists, "OTP should exist before invalidation")

		// Invalidate OTP
		err = repo.InvalidateOTP(ctx, otpEncx.EmailHash)
		assert.NoError(t, err)

		// Verify OTP no longer exists
		exists = td.CheckOTPExists(t, ctx, otpEncx.EmailHash, testClient)
		assert.False(t, exists, "OTP should not exist after invalidation")
	})

	t.Run("should handle invalidating non-existent OTP", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)
		email := "hash_nonexistent@example.com"

		// Try to invalidate non-existent OTP
		err := repo.InvalidateOTP(ctx, email)
		// NOTE: that thing returns an error brother, why ?
		assert.Error(t, err, "Invalidating non-existent OTP should not return error")

		// Verify no side effects
		exists := td.CheckOTPExists(t, ctx, email, testClient)
		assert.False(t, exists, "Non-existent OTP should remain non-existent")
	})

	t.Run("should not affect other OTP entries", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create and insert multiple OTPs
		otpEncx1 := td.NewTestOTPEncx(t)
		otpEncx1.EmailHash = "user1@example.com"
		otpEncx2 := td.NewTestOTPEncx(t)
		otpEncx2.EmailHash = "user2@example.com"
		otpEncx3 := td.NewTestOTPEncx(t)
		otpEncx3.EmailHash = "user3@example.com"

		err := td.InsertOTPEncx(t, ctx, otpEncx1, testClient, 10*time.Minute)
		require.NoError(t, err)
		err = td.InsertOTPEncx(t, ctx, otpEncx2, testClient, 10*time.Minute)
		require.NoError(t, err)
		err = td.InsertOTPEncx(t, ctx, otpEncx3, testClient, 10*time.Minute)
		require.NoError(t, err)

		// Verify all OTPs exist
		assert.True(t, td.CheckOTPExists(t, ctx, otpEncx1.EmailHash, testClient))
		assert.True(t, td.CheckOTPExists(t, ctx, otpEncx2.EmailHash, testClient))
		assert.True(t, td.CheckOTPExists(t, ctx, otpEncx3.EmailHash, testClient))

		// Invalidate only the middle OTP
		err = repo.InvalidateOTP(ctx, otpEncx2.EmailHash)
		assert.NoError(t, err)

		// Verify only the targeted OTP was removed
		assert.True(t, td.CheckOTPExists(t, ctx, otpEncx1.EmailHash, testClient), "OTP1 should still exist")
		assert.False(t, td.CheckOTPExists(t, ctx, otpEncx2.EmailHash, testClient), "OTP2 should be invalidated")
		assert.True(t, td.CheckOTPExists(t, ctx, otpEncx3.EmailHash, testClient), "OTP3 should still exist")
	})

	t.Run("should invalidate expired OTP", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create and insert OTP with very short TTL
		otpEncx := td.NewTestOTPEncx(t)

		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, 100*time.Millisecond)
		require.NoError(t, err)

		// Wait for expiration
		time.Sleep(200 * time.Millisecond)

		// Try to invalidate expired OTP (should not error even if already expired)
		err = repo.InvalidateOTP(ctx, otpEncx.EmailHash)
		assert.Error(t, err, "Invalidating expired OTP should not return error")
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "Invalidating expired OTP should not return error")

		// Verify OTP doesn't exist
		exists := td.CheckOTPExists(t, ctx, otpEncx.EmailHash, testClient)
		assert.False(t, exists, "Expired OTP should not exist")
	})

	t.Run("should handle empty email hash", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create OTP with empty email hash (edge case)
		otpEncx := td.NewTestOTPEncx(t)
		otpEncx.EmailHash = ""

		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, 5*time.Minute)
		require.NoError(t, err)

		// Verify it exists
		exists := td.CheckOTPExists(t, ctx, "", testClient)
		assert.True(t, exists, "OTP with empty hash should exist")

		// Invalidate with empty email hash
		err = repo.InvalidateOTP(ctx, "")
		assert.NoError(t, err, "Should handle empty email hash")

		// Verify it's gone
		exists = td.CheckOTPExists(t, ctx, "", testClient)
		assert.False(t, exists, "OTP with empty hash should be invalidated")
	})

	t.Run("should handle multiple invalidations of same OTP", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create and insert test OTP
		otpEncx := td.NewTestOTPEncx(t)

		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, 10*time.Minute)
		require.NoError(t, err)

		// First invalidation
		err = repo.InvalidateOTP(ctx, otpEncx.EmailHash)
		assert.NoError(t, err)

		// Second invalidation of same OTP (should not error)
		err = repo.InvalidateOTP(ctx, otpEncx.EmailHash)
		assert.Error(t, err, "Multiple invalidations should not error")
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "Multiple invalidations should not error")

		// Verify OTP doesn't exist
		exists := td.CheckOTPExists(t, ctx, otpEncx.EmailHash, testClient)
		assert.False(t, exists, "OTP should not exist after multiple invalidations")
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create and insert OTP
		otpEncx := td.NewTestOTPEncx(t)

		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, 10*time.Minute)
		require.NoError(t, err)

		// Create cancelled context
		cancelCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		// Attempt to invalidate with cancelled context
		err = repo.InvalidateOTP(cancelCtx, otpEncx.EmailHash)
		assert.Error(t, err, "Should return error for cancelled context")
		assert.ErrorIs(t, err, errs.ErrContext, "Should return context error")

		// Verify OTP still exists (invalidation didn't complete)
		exists := td.CheckOTPExists(t, ctx, otpEncx.EmailHash, testClient)
		assert.True(t, exists, "OTP should still exist after cancelled invalidation")
	})

	t.Run("should invalidate OTP with large data", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create OTP with large data
		otpEncx := td.NewTestOTPEncx(t)
		otpEncx.DEKEncrypted = make([]byte, 1024) // 1KB of data
		for i := range otpEncx.DEKEncrypted {
			otpEncx.DEKEncrypted[i] = byte(i % 256)
		}

		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, 10*time.Minute)
		require.NoError(t, err)

		// Verify large OTP exists
		exists := td.CheckOTPExists(t, ctx, otpEncx.EmailHash, testClient)
		assert.True(t, exists, "Large OTP should exist")

		// Invalidate large OTP
		err = repo.InvalidateOTP(ctx, otpEncx.EmailHash)
		assert.NoError(t, err)

		// Verify large OTP is gone
		exists = td.CheckOTPExists(t, ctx, otpEncx.EmailHash, testClient)
		assert.False(t, exists, "Large OTP should be invalidated")
	})

	// Clean up after all tests
	t.Cleanup(func() {
		td.ClearOTPKeys(t, ctx, testClient)
	})
}
