package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	otpRepository "github.com/Leviosa-care/authuser/internal/adapters/redis/otp"
	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/hengadev/encx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ClearOTPKeys removes all OTP-related keys from Redis
func ClearOTPKeys(t *testing.T, ctx context.Context, client *redis.Client) {
	t.Helper()

	// Get all OTP keys
	keys, err := client.Keys(ctx, "authuser:otp:*").Result()
	require.NoError(t, err)

	// Delete all OTP keys if any exist
	if len(keys) > 0 {
		err = client.Del(ctx, keys...).Err()
		require.NoError(t, err)
	}
}

// NewValidOTP creates a valid OTP domain object for testing
func NewValidOTP(email string) *domain.OTP {
	now := time.Now().UTC().Truncate(time.Microsecond)
	return &domain.OTP{
		Email:     email,
		Code:      "123456",
		Attempts:  0,
		CreatedAt: now,
		ExpiresAt: now.Add(10 * time.Minute),
	}
}

// NewExpiredOTP creates an expired OTP for testing expiration scenarios
func NewExpiredOTP(email string) *domain.OTP {
	past := time.Now().UTC().Add(-15 * time.Minute).Truncate(time.Microsecond)
	return &domain.OTP{
		Email:     email,
		Code:      "123456",
		Attempts:  0,
		CreatedAt: past,
		ExpiresAt: past.Add(10 * time.Minute), // Expired 5 minutes ago
	}
}

// NewOTPWithAttempts creates an OTP with specific number of attempts
func NewOTPWithAttempts(email string, attempts int) *domain.OTP {
	now := time.Now().UTC().Truncate(time.Microsecond)
	return &domain.OTP{
		Email:     email,
		Code:      "123456",
		Attempts:  attempts,
		CreatedAt: now,
		ExpiresAt: now.Add(10 * time.Minute),
	}
}

// InsertOTP directly inserts an OTP into Redis for test setup
func InsertOTP(t *testing.T, ctx context.Context, otp *domain.OTP, client *redis.Client, ttl time.Duration) {
	t.Helper()

	// Serialize OTP to JSON
	otpData, err := json.Marshal(otp)
	require.NoError(t, err, "Failed to marshal OTP for test insertion")

	// Format the key
	key := otpRepository.FormatOTPKey(otp.EmailHash)

	// Insert into Redis with TTL
	err = client.Set(ctx, key, otpData, ttl).Err()
	require.NoError(t, err, "Failed to insert OTP for email hash: %s", otp.EmailHash)
}

// GetOTPFromRedis retrieves an OTP directly from Redis for verification
func GetOTPFromRedis(t *testing.T, ctx context.Context, emailHash string, client *redis.Client) (*domain.OTP, error) {
	t.Helper()

	key := otpRepository.FormatOTPKey(emailHash)

	// Get data from Redis
	data, err := client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("OTP not found for hash %s", emailHash)
		}
		return nil, fmt.Errorf("failed to get OTP from Redis: %w", err)
	}

	// Deserialize OTP
	var otp domain.OTP
	err = json.Unmarshal([]byte(data), &otp)
	require.NoError(t, err, "Failed to unmarshal OTP from Redis")

	return &otp, nil
}

// CheckOTPExists verifies if an OTP exists in Redis
func CheckOTPExists(t *testing.T, ctx context.Context, emailHash string, client *redis.Client) bool {
	t.Helper()

	key := otpRepository.FormatOTPKey(emailHash)
	exists, err := client.Exists(ctx, key).Result()
	require.NoError(t, err, "Failed to check OTP existence")

	return exists == 1
}

// GetOTPTTL retrieves the TTL of an OTP key in Redis
func GetOTPTTL(t *testing.T, ctx context.Context, emailHash string, client *redis.Client) time.Duration {
	t.Helper()

	key := otpRepository.FormatOTPKey(emailHash)
	ttl, err := client.TTL(ctx, key).Result()
	require.NoError(t, err, "Failed to get OTP TTL")

	return ttl
}

// CreateEncryptedOTPData creates encrypted OTP data using real encryption for testing
func CreateEncryptedOTPData(t *testing.T, otp *domain.OTP, crypto encx.CryptoService) []byte {
	t.Helper()

	// Create a copy to avoid modifying the original
	encryptedOTP := *otp

	// Use crypto service to process the struct and populate encrypted fields
	err := crypto.ProcessStruct(context.Background(), &encryptedOTP)
	require.NoError(t, err, "Failed to encrypt OTP struct")

	data, err := json.Marshal(&encryptedOTP)
	require.NoError(t, err, "Failed to marshal encrypted OTP data")

	return data
}

// FormatOTPKey formats an OTP key for Redis (public for testing)
func FormatOTPKey(emailHash string) string {
	return otpRepository.FormatOTPKey(emailHash)
}

// ValidateOTPData compares two OTP objects for testing equality
// Only compares fields that are actually stored/retrieved from Redis
func ValidateOTPData(t *testing.T, expected, actual *domain.OTP) {
	t.Helper()

	// Compare fields that are stored in Redis (have JSON tags)
	assert.Equal(t, expected.EmailHash, actual.EmailHash, "EmailHash mismatch")
	assert.Equal(t, expected.CodeEncrypted, actual.CodeEncrypted, "CodeEncrypted mismatch")
	assert.Equal(t, expected.Attempts, actual.Attempts, "Attempts mismatch")
	assert.True(t, expected.CreatedAt.Equal(actual.CreatedAt), "CreatedAt mismatch")
	assert.True(t, expected.ExpiresAt.Equal(actual.ExpiresAt), "ExpiresAt mismatch")
	assert.Equal(t, expected.DEKEncrypted, actual.DEKEncrypted, "DEKEncrypted mismatch")
	assert.Equal(t, expected.KeyVersion, actual.KeyVersion, "KeyVersion mismatch")

	// Plaintext fields (Email, Code, DEK) should not be compared
	// as they are not stored/retrieved from Redis
}

// CreateOTP creates and stores an OTP in Redis for testing
func CreateOTP(t *testing.T, ctx context.Context, email string, client *redis.Client) {
	t.Helper()

	otp := NewValidOTP(email)
	// Create basic hash for email (simplified for test)
	otp.EmailHash = fmt.Sprintf("hash_%s", email)

	ttl := 10 * time.Minute
	InsertOTP(t, ctx, otp, client, ttl)
}

// GetOTP retrieves an OTP by email, returns nil if not found
func GetOTP(t *testing.T, ctx context.Context, email string, client *redis.Client) *domain.OTP {
	t.Helper()

	// Create basic hash for email (simplified for test)
	emailHash := fmt.Sprintf("hash_%s", email)

	otp, err := GetOTPFromRedis(t, ctx, emailHash, client)
	if err != nil {
		// Return nil if OTP not found
		return nil
	}

	return otp
}
