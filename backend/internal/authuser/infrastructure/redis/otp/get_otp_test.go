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

// make test-func TEST_NAME=TestGetOTP TEST_PATH=internal/authuser/infrastructure/redis/otp/get_otp_test.go

func TestGetOTP(t *testing.T) {
	ctx := context.Background()

	t.Run("should retrieve existing OTP successfully", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create and insert test OTP
		otpEncx := td.NewTestOTPEncx(t)

		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, 10*time.Minute)
		require.NoError(t, err)

		// Retrieve OTP
		otpData, err := repo.GetOTP(ctx, otpEncx.EmailHash)
		assert.NoError(t, err)

		// Verify data integrity - check fields that are actually stored
		var retrievedOTP map[string]any
		err = json.Unmarshal(otpData, &retrievedOTP)
		assert.NoError(t, err)

		// Only check fields that are stored (have JSON tags, not json:"-")
		assert.Equal(t, otpEncx.EmailHash, retrievedOTP["email_hash"])
		assert.Equal(t, float64(otpEncx.Attempts), retrievedOTP["attempts"]) // JSON numbers are float64
		assert.Equal(t, float64(otpEncx.KeyVersion), retrievedOTP["key_version"])

		// Email and Code fields should NOT be present as they have json:"-"
		assert.NotContains(t, retrievedOTP, "email", "Email should not be stored")
		assert.NotContains(t, retrievedOTP, "code", "Code should not be stored")
	})

	t.Run("should return error for non-existent OTP", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Try to retrieve non-existent OTP
		_, err := repo.GetOTP(ctx, "hash_nonexistent@example.com")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "Should return not found error")
	})

	t.Run("should retrieve OTP near expiration", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create and insert OTP with short TTL
		otpEncx := td.NewTestOTPEncx(t)

		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, 2*time.Second)
		require.NoError(t, err)

		// Wait a moment but retrieve before expiration
		time.Sleep(500 * time.Millisecond)

		// Should still be retrievable
		otpData, err := repo.GetOTP(ctx, otpEncx.EmailHash)
		assert.NoError(t, err)
		assert.NotEmpty(t, otpData)

		// Wait for expiration
		time.Sleep(2 * time.Second)

		// Should now be expired and not retrievable
		_, err = repo.GetOTP(ctx, otpEncx.EmailHash)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "Should return not found error for expired OTP")
	})

	t.Run("should handle empty email hash", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create OTP with empty hash (edge case)
		otpEncx := td.NewTestOTPEncx(t)
		otpEncx.EmailHash = ""

		err := td.InsertOTPEncx(t, ctx, otpEncx, testClient, 5*time.Minute)
		require.NoError(t, err)

		// Try to retrieve with empty hash
		otpData, err := repo.GetOTP(ctx, "")
		assert.NoError(t, err)
		assert.NotEmpty(t, otpData)
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

		// Attempt to retrieve with cancelled context
		_, err = repo.GetOTP(cancelCtx, otpEncx.EmailHash)
		assert.Error(t, err, "Should return error for cancelled context")
		assert.ErrorIs(t, err, errs.ErrContext, "Should return context error")
	})

	t.Run("should return raw bytes without modification", func(t *testing.T) {
		// Clean up before test
		td.ClearOTPKeys(t, ctx, testClient)

		// Create test data with specific byte sequence
		// otpEncx := td.NewTestOTP("bytes@example.com")
		otpEncx := td.NewTestOTPEncx(t)

		originalData, err := json.Marshal(otpEncx)
		require.NoError(t, err)

		td.InsertOTPEncx(t, ctx, otpEncx, testClient, 10*time.Minute)

		// Retrieve OTP
		retrievedData, err := repo.GetOTP(ctx, otpEncx.EmailHash)
		assert.NoError(t, err)

		// Compare byte-for-byte
		assert.Equal(t, originalData, retrievedData, "Retrieved data should match original bytes exactly")
	})

	// Clean up after all tests
	t.Cleanup(func() {
		td.ClearOTPKeys(t, ctx, testClient)
	})
}
