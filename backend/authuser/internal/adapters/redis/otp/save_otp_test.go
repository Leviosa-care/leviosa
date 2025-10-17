package otpRepository_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	th "github.com/Leviosa-care/authuser/test/helpers"
	"github.com/hengadev/encx"

	"github.com/Leviosa-care/core/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveOTP(t *testing.T) {
	ctx := context.Background()

	crypto, err := encx.NewTestCrypto(t)
	require.NoError(t, err)

	t.Run("should save OTP successfully with TTL", func(t *testing.T) {
		// Clean up before test
		th.ClearOTPKeys(t, ctx, testClient)

		// Create test OTP
		otp := th.NewTestOTP("test@example.com")

		otpEncx, err := domain.ProcessOTPEncx(ctx, crypto, otp)
		require.NoError(t, err)

		otpData, err := json.Marshal(otp)
		require.NoError(t, err)

		ttl := 5 * time.Minute

		// Save OTP
		err = repo.SaveOTP(ctx, otpEncx.EmailHash, otpData, ttl)
		require.NoError(t, err)

		// Verify OTP was saved
		exists := th.CheckOTPExists(t, ctx, otpEncx.EmailHash, testClient)
		assert.True(t, exists, "OTP should exist in Redis")

		// Verify TTL is set correctly (allow some variance for test execution time)
		actualTTL := th.GetOTPTTL(t, ctx, otpEncx.EmailHash, testClient)
		assert.Greater(t, actualTTL, 4*time.Minute, "TTL should be greater than 4 minutes")
		assert.LessOrEqual(t, actualTTL, ttl, "TTL should not exceed set value")
	})

	t.Run("should save OTP with zero TTL (no expiration)", func(t *testing.T) {
		// Clean up before test
		th.ClearOTPKeys(t, ctx, testClient)

		// Create test OTP
		otp := th.NewTestOTP("noexpiry@example.com")

		otpEncx, err := domain.ProcessOTPEncx(ctx, crypto, otp)
		require.NoError(t, err)

		otpData, err := json.Marshal(otp)
		require.NoError(t, err)

		// Save OTP with no TTL
		err = repo.SaveOTP(ctx, otpEncx.EmailHash, otpData, 0)
		require.NoError(t, err)

		// Verify OTP was saved
		exists := th.CheckOTPExists(t, ctx, otpEncx.EmailHash, testClient)
		assert.True(t, exists, "OTP should exist in Redis")

		// Verify no TTL is set (returns -1 for keys without expiration)
		actualTTL := th.GetOTPTTL(t, ctx, otpEncx.EmailHash, testClient)
		assert.Equal(t, time.Duration(-1), actualTTL, "Key should have no expiration")
	})

	t.Run("should overwrite existing OTP", func(t *testing.T) {
		// Clean up before test
		th.ClearOTPKeys(t, ctx, testClient)

		// Create and save first OTP
		otp1 := th.NewTestOTP("overwrite@example.com")

		// otp1.CodeEncrypted = []byte("encrypted-code-111111")

		otp1Encx, err := domain.ProcessOTPEncx(ctx, crypto, otp)
		require.NoError(t, err)

		otpData1, err := json.Marshal(otp1)
		require.NoError(t, err)

		err = repo.SaveOTP(ctx, otp1Encx.EmailHash, otpData1, 10*time.Minute)
		require.NoError(t, err)

		// Create and save second OTP with same email hash
		otp2 := th.NewTestOTP("overwrite@example.com")
		// otp2.CodeEncrypted = []byte("encrypted-code-222222")

		otp2Encx, err := domain.ProcessOTPEncx(ctx, crypto, otp)
		require.NoError(t, err)

		otpData2, err := json.Marshal(otp2)
		require.NoError(t, err)

		err = repo.SaveOTP(ctx, otp2Encx.EmailHash, otpData2, 5*time.Minute)
		require.NoError(t, err)

		// Verify only the second OTP exists and has correct encrypted data
		retrievedOTP, err := th.GetOTPEncxByEmailHash(t, ctx, otp2.EmailHash, testClient)
		require.NoError(t, err)
		assert.Equal(t, otp2.EmailHash, retrievedOTP.EmailHash, "Should retrieve the correct OTP")
		assert.Equal(t, otp2.CodeEncrypted, retrievedOTP.CodeEncrypted, "Should retrieve the overwritten encrypted code")

		// Code field should be empty since it's not stored (json:"-")
		assert.Empty(t, retrievedOTP.Code, "Code should not be stored/retrieved")

		// Verify TTL was updated
		actualTTL := th.GetOTPTTL(t, ctx, otp2.EmailHash, testClient)
		assert.Greater(t, actualTTL, 4*time.Minute, "TTL should be updated")
		assert.LessOrEqual(t, actualTTL, 5*time.Minute, "TTL should not exceed new value")
	})

	t.Run("should handle empty email hash", func(t *testing.T) {
		// Clean up before test
		th.ClearOTPKeys(t, ctx, testClient)

		otp := th.NewTestOTP("empty@example.com")
		otpData, err := json.Marshal(otp)
		require.NoError(t, err)

		// Save OTP with empty email hash
		err = repo.SaveOTP(ctx, "", otpData, 5*time.Minute)
		// This should still work (Redis allows empty keys), but key would be just the prefix
		require.NoError(t, err)
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		// Clean up before test
		th.ClearOTPKeys(t, ctx, testClient)

		// Create cancelled context
		cancelCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		otp := th.NewTestOTP("cancelled@example.com")
		otpData, err := json.Marshal(otp)
		require.NoError(t, err)

		// Attempt to save with cancelled context
		err = repo.SaveOTP(cancelCtx, otp.EmailHash, otpData, 5*time.Minute)
		assert.Error(t, err, "Should return error for cancelled context")
		assert.ErrorIs(t, err, errs.ErrContext, "Should return context error")
	})

	// Clean up after all tests
	t.Cleanup(func() {
		th.ClearOTPKeys(t, ctx, testClient)
	})
}
