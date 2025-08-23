package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/core/contracts/settings"
	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPhone(t *testing.T) {
	ctx := context.Background()

	t.Run("successful retrieval of phone setting", func(t *testing.T) {
		// Arrange
		expectedValueEncrypted := []byte("encrypted_phone_value_bytes")
		expectedDEKEncrypted := []byte("encrypted_dek_bytes")
		expectedKeyVersion := 1
		now := time.Now()

		// Insert test data directly into the encrypted table
		insertQuery := `
			INSERT INTO settings.encrypted (key, value_encrypted, created_at, updated_at, dek_encrypted, key_version) 
			VALUES ($1, $2, $3, $4, $5, $6)`

		_, err := testPool.Exec(ctx, insertQuery,
			settings.CompanyPhone, // Assuming this is the value of domain.CompanyPhoneKey
			expectedValueEncrypted,
			now,
			now,
			expectedDEKEncrypted,
			expectedKeyVersion,
		)
		require.NoError(t, err)

		// Clean up after test
		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", settings.CompanyPhone)
		}()

		// Act
		result, err := repo.GetPhone(ctx)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, settings.CompanyPhone, result.Key) // Assuming domain.CompanyPhoneKey value
		assert.Equal(t, expectedValueEncrypted, result.ValueEncrypted)
		assert.Equal(t, expectedDEKEncrypted, result.DEKEncrypted)
		assert.WithinDuration(t, now, result.CreatedAt, time.Second)
		assert.WithinDuration(t, now, result.UpdatedAt, time.Second)
	})

	t.Run("phone setting not found returns not found error", func(t *testing.T) {
		// Ensure no phone setting exists
		_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", settings.CompanyPhone)

		// Act
		result, err := repo.GetPhone(ctx)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		// Verify it's a not found error (adjust based on your errs package)
		assert.Contains(t, err.Error(), "res") // Based on your error message
	})

	t.Run("context timeout returns context error", func(t *testing.T) {
		// Arrange - insert test data first
		insertQuery := `
			INSERT INTO settings.encrypted (key, value_encrypted, created_at, updated_at, dek_encrypted, key_version) 
			VALUES ($1, $2, $3, $4, $5, $6)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery,
			settings.CompanyPhone,
			[]byte("test_value"),
			now,
			now,
			[]byte("test_dek"),
			1,
		)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", settings.CompanyPhone)
		}()

		// Create context with very short timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)
		defer cancel()

		// Wait a bit to ensure context times out
		time.Sleep(1 * time.Millisecond)

		// Act
		result, err := repo.GetPhone(timeoutCtx)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		// Should be a context error (adjust based on your errs package)
	})

	t.Run("context cancellation returns context error", func(t *testing.T) {
		// Arrange - insert test data first
		insertQuery := `
			INSERT INTO settings.encrypted (key, value_encrypted, created_at, updated_at, dek_encrypted, key_version) 
			VALUES ($1, $2, $3, $4, $5, $6)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery,
			settings.CompanyPhone,
			[]byte("test_value"),
			now,
			now,
			[]byte("test_dek"),
			1,
		)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", settings.CompanyPhone)
		}()

		// Create cancelled context
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel()

		// Act
		result, err := repo.GetPhone(cancelledCtx)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		// Should be a context error (adjust based on your errs package)
	})

	t.Run("phone setting with empty encrypted values", func(t *testing.T) {
		// Arrange
		now := time.Now()

		insertQuery := `
			INSERT INTO settings.encrypted (key, value_encrypted, created_at, updated_at, dek_encrypted, key_version) 
			VALUES ($1, $2, $3, $4, $5, $6)`

		_, err := testPool.Exec(ctx, insertQuery,
			settings.CompanyPhone,
			[]byte{}, // Empty byte slice
			now,
			now,
			[]byte{}, // Empty byte slice
			0,
		)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", settings.CompanyPhone)
		}()

		// Act
		result, err := repo.GetPhone(ctx)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, []byte{}, result.ValueEncrypted)
		assert.Equal(t, []byte{}, result.DEKEncrypted)
	})

	t.Run("phone setting with large encrypted values", func(t *testing.T) {
		// Arrange
		// Create large byte arrays to simulate real encrypted data
		largeValueEncrypted := make([]byte, 1000)
		largeDEKEncrypted := make([]byte, 500)

		// Fill with some test data
		for i := range largeValueEncrypted {
			largeValueEncrypted[i] = byte(i % 256)
		}
		for i := range largeDEKEncrypted {
			largeDEKEncrypted[i] = byte((i + 100) % 256)
		}

		now := time.Now()

		insertQuery := `
			INSERT INTO settings.encrypted (key, value_encrypted, created_at, updated_at, dek_encrypted, key_version) 
			VALUES ($1, $2, $3, $4, $5, $6)`

		_, err := testPool.Exec(ctx, insertQuery,
			settings.CompanyPhone,
			largeValueEncrypted,
			now,
			now,
			largeDEKEncrypted,
			1,
		)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", settings.CompanyPhone)
		}()

		// Act
		result, err := repo.GetPhone(ctx)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, largeValueEncrypted, result.ValueEncrypted)
		assert.Equal(t, largeDEKEncrypted, result.DEKEncrypted)
		assert.Len(t, result.ValueEncrypted, 1000)
		assert.Len(t, result.DEKEncrypted, 500)
	})
}

// Benchmark tests
func BenchmarkRepository_GetPhone(b *testing.B) {
	ctx := context.Background()

	// Setup test data
	setting := &domain.SettingEncrypted[string]{
		Key:            settings.CompanyPhone,
		Value:          "0100000000",
		ValueEncrypted: []byte("benchmark_encrypted_phone_bytes"),
		DEK:            []byte("benchmark_dek"),
		DEKEncrypted:   []byte("benchmark_encrypted_dek_bytes"),
		KeyVersion:     1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := repo.SetPhone(ctx, setting)
	require.NoError(b, err)

	defer func() {
		_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", settings.CompanyPhone)
	}()

	b.ResetTimer()

	for b.Loop() {
		_, err := repo.GetPhone(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}
