package otpRepository_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/test/helpers"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/middleware/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOTP(t *testing.T) {
	ctx := context.Background()

	t.Run("should retrieve existing OTP successfully", func(t *testing.T) {
		// Clean up before test
		helpers.ClearOTPKeys(t, ctx, testClient)

		// Create and insert test OTP
		otp := helpers.NewValidOTP("get@example.com")
		helpers.InsertOTP(t, ctx, otp, testClient, 10*time.Minute)

		// Retrieve OTP
		otpData, err := repo.GetOTP(ctx, otp.EmailHash)
		require.NoError(t, err)

		// Verify data integrity - check fields that are actually stored
		var retrievedOTP map[string]any
		err = json.Unmarshal(otpData, &retrievedOTP)
		require.NoError(t, err)

		// Only check fields that are stored (have JSON tags, not json:"-")
		assert.Equal(t, otp.EmailHash, retrievedOTP["email_hash"])
		assert.Equal(t, float64(otp.Attempts), retrievedOTP["attempts"]) // JSON numbers are float64
		assert.Equal(t, float64(otp.KeyVersion), retrievedOTP["key_version"])

		// Email and Code fields should NOT be present as they have json:"-"
		assert.NotContains(t, retrievedOTP, "email", "Email should not be stored")
		assert.NotContains(t, retrievedOTP, "code", "Code should not be stored")
	})

	t.Run("should return error for non-existent OTP", func(t *testing.T) {
		// Clean up before test
		helpers.ClearOTPKeys(t, ctx, testClient)

		// Try to retrieve non-existent OTP
		_, err := repo.GetOTP(ctx, "hash_nonexistent@example.com")
		require.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "Should return not found error")
	})

	t.Run("should retrieve OTP with all fields intact", func(t *testing.T) {
		// Clean up before test
		helpers.ClearOTPKeys(t, ctx, testClient)

		// Create OTP with all fields populated
		otp := helpers.NewOTPWithAttempts("complete@example.com", 2)
		otp.CodeEncrypted = []byte("encrypted-code-data")
		otp.DEKEncrypted = []byte("encrypted-dek-data")
		otp.KeyVersion = 1

		helpers.InsertOTP(t, ctx, otp, testClient, 15*time.Minute)

		// Retrieve OTP
		otpData, err := repo.GetOTP(ctx, otp.EmailHash)
		require.NoError(t, err)

		// Deserialize and verify all fields using auth.OTP
		var retrievedOTP auth.OTP
		err = json.Unmarshal(otpData, &retrievedOTP)
		require.NoError(t, err)

		helpers.ValidateOTPData(t, otp, &retrievedOTP)
	})

	t.Run("should retrieve OTP near expiration", func(t *testing.T) {
		// Clean up before test
		helpers.ClearOTPKeys(t, ctx, testClient)

		// Create and insert OTP with short TTL
		otp := helpers.NewValidOTP("expiring@example.com")
		helpers.InsertOTP(t, ctx, otp, testClient, 2*time.Second)

		// Wait a moment but retrieve before expiration
		time.Sleep(500 * time.Millisecond)

		// Should still be retrievable
		otpData, err := repo.GetOTP(ctx, otp.EmailHash)
		require.NoError(t, err)
		assert.NotEmpty(t, otpData)

		// Wait for expiration
		time.Sleep(2 * time.Second)

		// Should now be expired and not retrievable
		_, err = repo.GetOTP(ctx, otp.EmailHash)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "Should return not found error for expired OTP")
	})

	t.Run("should handle empty email hash", func(t *testing.T) {
		// Clean up before test
		helpers.ClearOTPKeys(t, ctx, testClient)

		// Create OTP with empty hash (edge case)
		otp := helpers.NewValidOTP("empty@example.com")
		otp.EmailHash = ""
		helpers.InsertOTP(t, ctx, otp, testClient, 5*time.Minute)

		// Try to retrieve with empty hash
		otpData, err := repo.GetOTP(ctx, "")
		require.NoError(t, err)
		assert.NotEmpty(t, otpData)
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		// Clean up before test
		helpers.ClearOTPKeys(t, ctx, testClient)

		// Create and insert OTP
		otp := helpers.NewValidOTP("cancelled@example.com")
		helpers.InsertOTP(t, ctx, otp, testClient, 10*time.Minute)

		// Create cancelled context
		cancelCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		// Attempt to retrieve with cancelled context
		_, err := repo.GetOTP(cancelCtx, otp.EmailHash)
		assert.Error(t, err, "Should return error for cancelled context")
		assert.ErrorIs(t, err, errs.ErrContext, "Should return context error")
	})

	t.Run("should return raw bytes without modification", func(t *testing.T) {
		// Clean up before test
		helpers.ClearOTPKeys(t, ctx, testClient)

		// Create test data with specific byte sequence
		otp := helpers.NewValidOTP("bytes@example.com")
		originalData, err := json.Marshal(otp)
		require.NoError(t, err)

		helpers.InsertOTP(t, ctx, otp, testClient, 10*time.Minute)

		// Retrieve OTP
		retrievedData, err := repo.GetOTP(ctx, otp.EmailHash)
		require.NoError(t, err)

		// Compare byte-for-byte
		assert.Equal(t, originalData, retrievedData, "Retrieved data should match original bytes exactly")
	})

	// Clean up after all tests
	t.Cleanup(func() {
		helpers.ClearOTPKeys(t, ctx, testClient)
	})
}
