package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-unit-postgres TEST=TestGetEncryptedSetting

func TestGetEncryptedSetting(t *testing.T) {
	ctx := context.Background()

	t.Run("successful retrieval of encrypted setting", func(t *testing.T) {
		// Arrange
		key := "test_encrypted_key"
		valueEncrypted := []byte("encrypted-value-123")
		dekEncrypted := []byte("dek-encrypted-456")
		keyVersion := 1
		metadata := encx.EncryptionMetadata{
			PepperVersion:    1,
			KEKAlias:         "test-alias",
			EncryptionTime:   1234567890,
			GeneratorVersion: "1.0.0",
		}

		now := time.Now()

		// Insert test data directly into encrypted table (id is auto-generated)
		insertQuery := `
			INSERT INTO settings.encrypted (
				key, value_encrypted, created_at, updated_at,
				dek_encrypted, key_version, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err := testPool.Exec(ctx, insertQuery,
			key, valueEncrypted, now, now,
			dekEncrypted, keyVersion, metadata)
		require.NoError(t, err)

		// Clean up after test
		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetEncryptedSetting(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, key, result.Key)
		assert.NotEmpty(t, result.ID) // ID should be auto-generated as string
		assert.Equal(t, valueEncrypted, result.ValueEncrypted)
		assert.Equal(t, dekEncrypted, result.DEKEncrypted)
		assert.Equal(t, keyVersion, result.KeyVersion)
		assert.Equal(t, metadata, result.Metadata)
		assert.False(t, result.CreatedAt.IsZero())
		assert.False(t, result.UpdatedAt.IsZero())
	})

	t.Run("successful retrieval with empty encrypted value", func(t *testing.T) {
		// Arrange
		key := "test_empty_encrypted"
		valueEncrypted := []byte{}
		dekEncrypted := []byte("dek-encrypted-empty")
		keyVersion := 2
		metadata := encx.EncryptionMetadata{
			PepperVersion:    1,
			KEKAlias:         "test-alias",
			EncryptionTime:   1234567890,
			GeneratorVersion: "1.0.0",
		}

		now := time.Now()

		insertQuery := `
			INSERT INTO settings.encrypted (
				key, value_encrypted, created_at, updated_at,
				dek_encrypted, key_version, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err := testPool.Exec(ctx, insertQuery,
			key, valueEncrypted, now, now,
			dekEncrypted, keyVersion, metadata)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetEncryptedSetting(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, key, result.Key)
		assert.NotEmpty(t, result.ID)
		assert.Empty(t, result.ValueEncrypted)
		assert.Equal(t, dekEncrypted, result.DEKEncrypted)
		assert.Equal(t, keyVersion, result.KeyVersion)
	})

	t.Run("unsuccessful retrieval due to not null violation thanks to nil encrypted value", func(t *testing.T) {
		// Arrange
		key := "test_nil_encrypted"
		dekEncrypted := []byte("dek-encrypted-nil")
		keyVersion := 3
		metadata := encx.EncryptionMetadata{
			PepperVersion:    1,
			KEKAlias:         "test-alias",
			EncryptionTime:   1234567890,
			GeneratorVersion: "1.0.0",
		}

		now := time.Now()

		insertQuery := `
			INSERT INTO settings.encrypted (
				key, value_encrypted, created_at, updated_at,
				dek_encrypted, key_version, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err := testPool.Exec(ctx, insertQuery,
			key, nil, now, now,
			dekEncrypted, keyVersion, metadata)
		assert.Contains(t, err.Error(), "not-null constraint")
	})

	t.Run("successful retrieval with large encrypted data", func(t *testing.T) {
		// Arrange
		key := "test_large_encrypted"
		// Simulate large encrypted data (1KB)
		valueEncrypted := make([]byte, 1024)
		for i := range valueEncrypted {
			valueEncrypted[i] = byte(i % 256)
		}
		dekEncrypted := []byte("dek-encrypted-large")
		keyVersion := 4
		metadata := encx.EncryptionMetadata{
			PepperVersion:    1,
			KEKAlias:         "test-alias",
			EncryptionTime:   1234567890,
			GeneratorVersion: "1.0.0",
		}

		now := time.Now()

		insertQuery := `
			INSERT INTO settings.encrypted (
				key, value_encrypted, created_at, updated_at,
				dek_encrypted, key_version, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err := testPool.Exec(ctx, insertQuery,
			key, valueEncrypted, now, now,
			dekEncrypted, keyVersion, metadata)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetEncryptedSetting(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, key, result.Key)
		assert.NotEmpty(t, result.ID)
		assert.Len(t, result.ValueEncrypted, 1024)
		assert.Equal(t, valueEncrypted, result.ValueEncrypted)
		assert.Equal(t, dekEncrypted, result.DEKEncrypted)
		assert.Equal(t, keyVersion, result.KeyVersion)
	})

	t.Run("successful retrieval with special characters in key", func(t *testing.T) {
		// Arrange
		key := "test.key-with_special.chars_and_numbers_123"
		valueEncrypted := []byte("encrypted-special-chars")
		dekEncrypted := []byte("dek-encrypted-special")
		keyVersion := 5
		metadata := encx.EncryptionMetadata{
			PepperVersion:    1,
			KEKAlias:         "test-alias",
			EncryptionTime:   1234567890,
			GeneratorVersion: "1.0.0",
		}

		now := time.Now()

		insertQuery := `
			INSERT INTO settings.encrypted (
				key, value_encrypted, created_at, updated_at,
				dek_encrypted, key_version, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err := testPool.Exec(ctx, insertQuery,
			key, valueEncrypted, now, now,
			dekEncrypted, keyVersion, metadata)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetEncryptedSetting(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, key, result.Key)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, valueEncrypted, result.ValueEncrypted)
		assert.Equal(t, dekEncrypted, result.DEKEncrypted)
		assert.Equal(t, keyVersion, result.KeyVersion)
	})

	t.Run("key not found returns repository not found error", func(t *testing.T) {
		// Arrange
		nonExistentKey := "non_existent_encrypted_key_12345"

		// Act
		result, err := repo.GetEncryptedSetting(ctx, nonExistentKey)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)

		// Check if it's the expected error type from errs.ClassifyPgError
		assert.Contains(t, err.Error(), "get encrypted setting for key 'non_existent_encrypted_key_12345'")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Arrange
		key := "test_context_cancel_encrypted"
		valueEncrypted := []byte("encrypted-context-cancel")
		dekEncrypted := []byte("dek-encrypted-context")
		keyVersion := 6
		metadata := encx.EncryptionMetadata{
			PepperVersion:    1,
			KEKAlias:         "test-alias",
			EncryptionTime:   1234567890,
			GeneratorVersion: "1.0.0",
		}

		now := time.Now()

		insertQuery := `
			INSERT INTO settings.encrypted (
				key, value_encrypted, created_at, updated_at,
				dek_encrypted, key_version, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err := testPool.Exec(ctx, insertQuery,
			key, valueEncrypted, now, now,
			dekEncrypted, keyVersion, metadata)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", key)
		}()

		// Create cancelled context
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel()

		// Act
		result, err := repo.GetEncryptedSetting(cancelledCtx, key)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("successful retrieval with zero key version", func(t *testing.T) {
		// Arrange
		key := "test_zero_key_version"
		valueEncrypted := []byte("encrypted-zero-version")
		dekEncrypted := []byte("dek-encrypted-zero")
		keyVersion := 0
		metadata := encx.EncryptionMetadata{
			PepperVersion:    1,
			KEKAlias:         "test-alias",
			EncryptionTime:   1234567890,
			GeneratorVersion: "1.0.0",
		}

		now := time.Now()

		insertQuery := `
			INSERT INTO settings.encrypted (
				key, value_encrypted, created_at, updated_at,
				dek_encrypted, key_version, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err := testPool.Exec(ctx, insertQuery,
			key, valueEncrypted, now, now,
			dekEncrypted, keyVersion, metadata)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetEncryptedSetting(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, key, result.Key)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, keyVersion, result.KeyVersion)
	})
}

func BenchmarkGetEncryptedSetting(b *testing.B) {
	ctx := context.Background()

	b.Run("encrypted setting retrieval", func(b *testing.B) {
		// Arrange
		key := "benchmark_encrypted_key"
		valueEncrypted := make([]byte, 512) // 512 bytes of encrypted data
		for i := range valueEncrypted {
			valueEncrypted[i] = byte(i % 256)
		}
		dekEncrypted := []byte("dek-encrypted-benchmark")
		keyVersion := 10
		metadata := encx.EncryptionMetadata{
			PepperVersion:    1,
			KEKAlias:         "test-alias",
			EncryptionTime:   1234567890,
			GeneratorVersion: "1.0.0",
		}

		now := time.Now()

		insertQuery := `
			INSERT INTO settings.encrypted (
				key, value_encrypted, created_at, updated_at,
				dek_encrypted, key_version, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err := testPool.Exec(ctx, insertQuery,
			key, valueEncrypted, now, now,
			dekEncrypted, keyVersion, metadata)
		require.NoError(b, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", key)
		}()

		// Benchmark the retrieval
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, err := repo.GetEncryptedSetting(ctx, key)
			if err != nil {
				b.Fatal(err)
			}
			if result == nil {
				b.Fatal("result should not be nil")
			}
			if result.Key != key {
				b.Fatalf("expected key %s, got %s", key, result.Key)
			}
		}
	})

	b.Run("encrypted setting retrieval with different sizes", func(b *testing.B) {
		testCases := []struct {
			name string
			size int
			key  string
		}{
			{"small", 64, "benchmark_small"},
			{"medium", 256, "benchmark_medium"},
			{"large", 1024, "benchmark_large"},
		}

		for _, tc := range testCases {
			b.Run(tc.name, func(b *testing.B) {
				// Arrange
				valueEncrypted := make([]byte, tc.size)
				for i := range valueEncrypted {
					valueEncrypted[i] = byte(i % 256)
				}
				dekEncrypted := []byte("dek-encrypted-" + tc.name)
				keyVersion := 10
				metadata := encx.EncryptionMetadata{
					PepperVersion:    1,
					KEKAlias:         "test-alias",
					EncryptionTime:   1234567890,
					GeneratorVersion: "1.0.0",
				}

				now := time.Now()

				insertQuery := `
					INSERT INTO settings.encrypted (
						key, value_encrypted, created_at, updated_at,
						dek_encrypted, key_version, metadata
					) VALUES ($1, $2, $3, $4, $5, $6, $7)`

				_, err := testPool.Exec(ctx, insertQuery,
					tc.key, valueEncrypted, now, now,
					dekEncrypted, keyVersion, metadata)
				require.NoError(b, err)

				defer func() {
					_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", tc.key)
				}()

				// Benchmark the retrieval
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					result, err := repo.GetEncryptedSetting(ctx, tc.key)
					if err != nil {
						b.Fatal(err)
					}
					if result == nil {
						b.Fatal("result should not be nil")
					}
				}
			})
		}
	})
}
