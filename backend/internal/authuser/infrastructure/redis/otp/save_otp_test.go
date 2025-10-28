package otpRepository_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestSaveOTP TEST_PATH=internal/authuser/infrastructure/redis/otp/save_otp_test.go

func TestSaveOTP(t *testing.T) {
	ctx := context.Background()

	t.Run("should save OTP successfully with TTL", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create test OTP
		otpEncx := td.NewTestOTPEncx(t)

		otpData, err := json.Marshal(otpEncx)
		require.NoError(t, err)

		ttl := 5 * time.Minute

		// Save OTP
		err = repo.SaveOTP(ctx, otpEncx.EmailHash, otpData, ttl)
		assert.NoError(t, err)

		// Verify OTP was saved
		exists := td.CheckOTPExists(t, ctx, otpEncx.EmailHash, testClient)
		assert.True(t, exists, "OTP should exist in Redis")

		// Verify TTL is set correctly (allow some variance for test execution time)
		actualTTL := td.GetOTPTTL(t, ctx, otpEncx.EmailHash, testClient)
		assert.Greater(t, actualTTL, 4*time.Minute, "TTL should be greater than 4 minutes")
		assert.LessOrEqual(t, actualTTL, ttl, "TTL should not exceed set value")
	})

	t.Run("should save OTP with zero TTL (no expiration)", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create test OTP
		otpEncx := td.NewTestOTPEncx(t)

		otpData, err := json.Marshal(otpEncx)
		require.NoError(t, err)

		// Save OTP with no TTL
		err = repo.SaveOTP(ctx, otpEncx.EmailHash, otpData, 0)
		assert.NoError(t, err)

		// Verify OTP was saved
		exists := td.CheckOTPExists(t, ctx, otpEncx.EmailHash, testClient)
		assert.True(t, exists, "OTP should exist in Redis")

		// Verify no TTL is set (returns -1 for keys without expiration)
		actualTTL := td.GetOTPTTL(t, ctx, otpEncx.EmailHash, testClient)
		assert.Equal(t, time.Duration(-1), actualTTL, "Key should have no expiration")
	})

	t.Run("should overwrite existing OTP", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create and save first OTP
		otp1Encx := td.NewTestOTPEncx(t)

		otpData1, err := json.Marshal(otp1Encx)
		require.NoError(t, err)

		err = repo.SaveOTP(ctx, otp1Encx.EmailHash, otpData1, 10*time.Minute)
		require.NoError(t, err)

		// Create and save second OTP with same email hash
		otp2Encx := td.NewTestOTPEncx(t)

		otpData2, err := json.Marshal(otp2Encx)
		require.NoError(t, err)

		err = repo.SaveOTP(ctx, otp2Encx.EmailHash, otpData2, 5*time.Minute)
		assert.NoError(t, err)

		// Verify only the second OTP exists and has correct encrypted data
		retrievedOTP, err := td.GetOTPEncxByEmailHash(t, ctx, otp2Encx.EmailHash, testClient)
		require.NoError(t, err)
		assert.Equal(t, otp2Encx.EmailHash, retrievedOTP.EmailHash, "Should retrieve the correct OTP")
		assert.Equal(t, otp2Encx.CodeEncrypted, retrievedOTP.CodeEncrypted, "Should retrieve the overwritten encrypted code")

		// Verify TTL was updated
		actualTTL := td.GetOTPTTL(t, ctx, otp2Encx.EmailHash, testClient)
		assert.Greater(t, actualTTL, 4*time.Minute, "TTL should be updated")
		assert.LessOrEqual(t, actualTTL, 5*time.Minute, "TTL should not exceed new value")
	})

	t.Run("should handle empty email hash", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		otpEncx := td.NewTestOTPEncx(t)

		otpData, err := json.Marshal(otpEncx)
		require.NoError(t, err)

		// Save OTP with empty email hash
		err = repo.SaveOTP(ctx, "", otpData, 5*time.Minute)
		// This should still work (Redis allows empty keys), but key would be just the prefix
		assert.NoError(t, err)
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create cancelled context
		cancelCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		otpEncx := td.NewTestOTPEncx(t)

		otpData, err := json.Marshal(otpEncx)
		require.NoError(t, err)

		// Attempt to save with cancelled context
		err = repo.SaveOTP(cancelCtx, otpEncx.EmailHash, otpData, 5*time.Minute)
		assert.Error(t, err, "Should return error for cancelled context")
		assert.ErrorIs(t, err, errs.ErrContext, "Should return context error")
	})

	// Clean up after all tests
	t.Cleanup(func() {
		td.ClearOTPKeys(t, ctx, testClient)
	})
}
