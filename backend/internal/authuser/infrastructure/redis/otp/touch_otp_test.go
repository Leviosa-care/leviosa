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

// make test-func TEST_NAME=TestTouchOTP TEST_PATH=internal/authuser/infrastructure/redis/otp/touch_otp_test.go

func TestTouchOTP(t *testing.T) {
	ctx := context.Background()

	t.Run("should update TTL for existing OTP", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create and insert test OTP with initial TTL
		otpEncx := td.NewTestOTPEncx(t)

		initialTTL := 5 * time.Minute
		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, initialTTL)
		require.NoError(t, err)

		// Get initial TTL
		ttlBefore := td.GetOTPTTL(t, ctx, otpEncx.EmailHash, testClient)
		assert.Greater(t, ttlBefore, 4*time.Minute, "Initial TTL should be close to 5 minutes")

		// Update TTL to a longer duration
		newTTL := 15 * time.Minute
		err = repo.TouchOTP(ctx, otpEncx.EmailHash, newTTL)
		assert.NoError(t, err)

		// Verify TTL was updated
		ttlAfter := td.GetOTPTTL(t, ctx, otpEncx.EmailHash, testClient)
		assert.Greater(t, ttlAfter, 14*time.Minute, "New TTL should be close to 15 minutes")
		assert.Greater(t, ttlAfter, ttlBefore, "New TTL should be greater than initial TTL")
	})

	t.Run("should update TTL to shorter duration", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create and insert test OTP with long initial TTL
		otpEncx := td.NewTestOTPEncx(t)
		initialTTL := 20 * time.Minute

		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, initialTTL)
		require.NoError(t, err)

		// Update TTL to shorter duration
		newTTL := 2 * time.Minute
		err = repo.TouchOTP(ctx, otpEncx.EmailHash, newTTL)
		assert.NoError(t, err)

		// Verify TTL was updated to shorter duration
		ttlAfter := td.GetOTPTTL(t, ctx, otpEncx.EmailHash, testClient)
		assert.Greater(t, ttlAfter, 1*time.Minute, "New TTL should be close to 2 minutes")
		assert.Less(t, ttlAfter, 3*time.Minute, "New TTL should be less than 3 minutes")
	})

	t.Run("should handle touching non-existent OTP", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Try to touch non-existent OTP
		err := repo.TouchOTP(ctx, "hash_nonexistent@example.com", 10*time.Minute)
		assert.NoError(t, err, "Touching non-existent OTP should not return error")

		// Verify no OTP was created
		exists := td.CheckOTPExists(t, ctx, "hash_nonexistent@example.com", testClient)
		assert.False(t, exists, "Non-existent OTP should remain non-existent")
	})

	t.Run("should set TTL to zero (no expiration)", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create and insert test OTP with TTL
		otpEncx := td.NewTestOTPEncx(t)
		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, 10*time.Minute)
		require.NoError(t, err)

		// Verify initial TTL exists
		ttlBefore := td.GetOTPTTL(t, ctx, otpEncx.EmailHash, testClient)
		assert.Greater(t, ttlBefore, 9*time.Minute, "Initial TTL should exist")

		// Update TTL to zero (no expiration)
		err = repo.TouchOTP(ctx, otpEncx.EmailHash, 0)
		assert.NoError(t, err)

		// Verify TTL is now -1 (no expiration)
		// Note: Redis EXPIRE with 0 actually expires the key immediately, not removes expiration
		// So we should expect the key to be gone, not persist without expiration
		ttlAfter := td.GetOTPTTL(t, ctx, otpEncx.EmailHash, testClient)
		// Redis EXPIRE 0 expires immediately, so TTL will be -2 (key doesn't exist)
		assert.Equal(t, time.Duration(-2), ttlAfter, "Key should be expired (TTL = -2 means key doesn't exist)")

		// Verify OTP is gone (since EXPIRE 0 expires immediately)
		exists := td.CheckOTPExists(t, ctx, otpEncx.EmailHash, testClient)
		assert.False(t, exists, "OTP should be gone since EXPIRE 0 expires immediately")
	})

	t.Run("should handle very short TTL", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create and insert test OTP
		otpEncx := td.NewTestOTPEncx(t)
		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, 10*time.Minute)
		require.NoError(t, err)

		// Update TTL to very short duration (1 second minimum for Redis)
		shortTTL := 1 * time.Second
		err = repo.TouchOTP(ctx, otpEncx.EmailHash, shortTTL)
		assert.NoError(t, err)

		// Verify TTL was set
		ttlAfter := td.GetOTPTTL(t, ctx, otpEncx.EmailHash, testClient)
		assert.Greater(t, ttlAfter, time.Duration(0), "TTL should be positive")
		assert.LessOrEqual(t, ttlAfter, 1*time.Second, "TTL should be very short")

		// Wait for expiration and verify OTP is gone
		time.Sleep(1100 * time.Millisecond) // Wait slightly longer than TTL
		exists := td.CheckOTPExists(t, ctx, otpEncx.EmailHash, testClient)
		assert.False(t, exists, "OTP should expire with very short TTL")
	})

	t.Run("should handle touching expired OTP", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create and insert OTP with very short TTL (1 second minimum for Redis)
		otpEncx := td.NewTestOTPEncx(t)
		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, 1*time.Second)
		require.NoError(t, err)

		// Wait for expiration
		time.Sleep(1100 * time.Millisecond)

		// Verify OTP is expired/gone
		exists := td.CheckOTPExists(t, ctx, otpEncx.EmailHash, testClient)
		assert.False(t, exists, "OTP should be expired")

		// Try to touch expired OTP
		err = repo.TouchOTP(ctx, otpEncx.EmailHash, 10*time.Minute)
		assert.NoError(t, err, "Touching expired OTP should not error")

		// Verify OTP is still gone (touch doesn't resurrect expired keys)
		exists = td.CheckOTPExists(t, ctx, otpEncx.EmailHash, testClient)
		assert.False(t, exists, "Expired OTP should remain gone")
	})

	t.Run("should handle empty email hash", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create OTP with empty email hash
		otpEncx := td.NewTestOTPEncx(t)
		otpEncx.EmailHash = ""
		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, 5*time.Minute)
		require.NoError(t, err)

		// Touch with empty email hash
		err = repo.TouchOTP(ctx, "", 15*time.Minute)
		require.NoError(t, err, "Should handle empty email hash")

		// Verify TTL was updated
		ttlAfter := td.GetOTPTTL(t, ctx, "", testClient)
		assert.Greater(t, ttlAfter, 14*time.Minute, "TTL should be updated for empty hash key")
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create and insert OTP
		otpEncx := td.NewTestOTPEncx(t)
		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, 10*time.Minute)

		// Create cancelled context
		cancelCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		// Attempt to touch with cancelled context
		err = repo.TouchOTP(cancelCtx, otpEncx.EmailHash, 20*time.Minute)
		assert.Error(t, err, "Should return error for cancelled context")
		assert.ErrorIs(t, err, errs.ErrContext, "Should return context error")

		// Verify TTL wasn't updated (operation was cancelled)
		ttlAfter := td.GetOTPTTL(t, ctx, otpEncx.EmailHash, testClient)
		assert.Greater(t, ttlAfter, 9*time.Minute, "TTL should remain close to original")
		assert.Less(t, ttlAfter, 11*time.Minute, "TTL should not be updated to 20 minutes")
	})

	t.Run("should preserve OTP data when updating TTL", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create OTP with specific data
		otpEncx := td.NewTestOTPEncx(t)
		otpEncx.CodeEncrypted = []byte("987654")
		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, 10*time.Minute)
		require.NoError(t, err)

		// Get original data
		originalOTP, err := td.GetOTPEncxByEmailHash(t, ctx, otpEncx.EmailHash, testClient)
		require.NoError(t, err)

		// Touch to update TTL
		err = repo.TouchOTP(ctx, otpEncx.EmailHash, 30*time.Minute)
		assert.NoError(t, err)

		// Verify data is preserved
		updatedOTP, err := td.GetOTPEncxByEmailHash(t, ctx, otpEncx.EmailHash, testClient)
		assert.NoError(t, err)

		td.ValidateOTPData(t, originalOTP, updatedOTP)

		// Verify TTL was actually updated
		ttlAfter := td.GetOTPTTL(t, ctx, otpEncx.EmailHash, testClient)
		assert.Greater(t, ttlAfter, 29*time.Minute, "TTL should be updated to ~30 minutes")
	})

	// Clean up after all tests
	t.Cleanup(func() {
		td.ClearOTPKeys(t, ctx, testClient)
	})
}
