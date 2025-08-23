package postgres_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetString(t *testing.T) {
	ctx := context.Background()

	t.Run("successful insertion", func(t *testing.T) {
		// Arrange
		setting := &domain.Setting[string]{
			Key:       "test_string_setting",
			Value:     "hello world",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetString(ctx, setting)

		// Assert
		require.NoError(t, err)

		// Verify the record was actually inserted
		var storedValue string
		var storedID int
		query := "SELECT id, value FROM settings.plain WHERE key = $1"
		err = testPool.QueryRow(ctx, query, setting.Key).Scan(&storedID, &storedValue)
		require.NoError(t, err)
		assert.Greater(t, storedID, 0) // ID should be auto-generated and positive
		assert.Equal(t, "hello world", storedValue)

		// Cleanup
		_, err = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", setting.Key)
		require.NoError(t, err)
	})

	t.Run("successful insertion with empty string", func(t *testing.T) {
		// Arrange
		setting := &domain.Setting[string]{
			Key:       "test_empty_string",
			Value:     "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetString(ctx, setting)

		// Assert
		require.NoError(t, err)

		// Verify the record was inserted with empty value
		var storedValue string
		query := "SELECT value FROM settings.plain WHERE key = $1"
		err = testPool.QueryRow(ctx, query, setting.Key).Scan(&storedValue)
		require.NoError(t, err)
		assert.Equal(t, "", storedValue)

		// Cleanup
		_, err = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", setting.Key)
		require.NoError(t, err)
	})

	t.Run("successful insertion with special characters", func(t *testing.T) {
		// Arrange
		setting := &domain.Setting[string]{
			Key:       "test_special_chars",
			Value:     "special!@#$%^&*()_+{}[]|\\:;\"'<>?,./",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetString(ctx, setting)

		// Assert
		require.NoError(t, err)

		// Verify the record was inserted
		var storedValue string
		query := "SELECT value FROM settings.plain WHERE key = $1"
		err = testPool.QueryRow(ctx, query, setting.Key).Scan(&storedValue)
		require.NoError(t, err)
		assert.Equal(t, "special!@#$%^&*()_+{}[]|\\:;\"'<>?,./", storedValue)

		// Cleanup
		_, err = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", setting.Key)
		require.NoError(t, err)
	})

	t.Run("successful insertion with unicode characters", func(t *testing.T) {
		// Arrange
		setting := &domain.Setting[string]{
			Key:       "test_unicode",
			Value:     "Hello 世界 🌍 Café résumé naïve",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetString(ctx, setting)

		// Assert
		require.NoError(t, err)

		// Verify the record was inserted with unicode preserved
		var storedValue string
		query := "SELECT value FROM settings.plain WHERE key = $1"
		err = testPool.QueryRow(ctx, query, setting.Key).Scan(&storedValue)
		require.NoError(t, err)
		assert.Equal(t, "Hello 世界 🌍 Café résumé naïve", storedValue)

		// Cleanup
		_, err = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", setting.Key)
		require.NoError(t, err)
	})

	t.Run("successful insertion with multiline string", func(t *testing.T) {
		// Arrange
		multilineValue := `Line 1
Line 2
Line 3 with	tabs
And some "quotes"`
		setting := &domain.Setting[string]{
			Key:       "test_multiline",
			Value:     multilineValue,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetString(ctx, setting)

		// Assert
		require.NoError(t, err)

		// Verify the record was inserted
		var storedValue string
		query := "SELECT value FROM settings.plain WHERE key = $1"
		err = testPool.QueryRow(ctx, query, setting.Key).Scan(&storedValue)
		require.NoError(t, err)
		assert.Equal(t, multilineValue, storedValue)

		// Cleanup
		_, err = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", setting.Key)
		require.NoError(t, err)
	})

	t.Run("successful insertion with very long string", func(t *testing.T) {
		// Arrange - Create a long string (1MB)
		longValue := strings.Repeat("A", 1024*1024)
		setting := &domain.Setting[string]{
			Key:       "test_long_string",
			Value:     longValue,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetString(ctx, setting)

		// Assert
		require.NoError(t, err)

		// Verify the record was inserted
		var storedValue string
		query := "SELECT value FROM settings.plain WHERE key = $1"
		err = testPool.QueryRow(ctx, query, setting.Key).Scan(&storedValue)
		require.NoError(t, err)
		assert.Equal(t, len(longValue), len(storedValue))
		assert.Equal(t, longValue, storedValue)

		// Cleanup
		_, err = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", setting.Key)
		require.NoError(t, err)
	})

	t.Run("successful insertion with JSON-like string", func(t *testing.T) {
		// Arrange
		jsonValue := `{"name": "test", "value": 123, "nested": {"key": "value"}}`
		setting := &domain.Setting[string]{
			Key:       "test_json_string",
			Value:     jsonValue,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetString(ctx, setting)

		// Assert
		require.NoError(t, err)

		// Verify the record was inserted
		var storedValue string
		query := "SELECT value FROM settings.plain WHERE key = $1"
		err = testPool.QueryRow(ctx, query, setting.Key).Scan(&storedValue)
		require.NoError(t, err)
		assert.Equal(t, jsonValue, storedValue)

		// Cleanup
		_, err = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", setting.Key)
		require.NoError(t, err)
	})

	t.Run("duplicate key insertion should update existing value", func(t *testing.T) {
		// Arrange
		setting1 := &domain.Setting[string]{
			Key:       "duplicate_string_key_test",
			Value:     "first value",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		setting2 := &domain.Setting[string]{
			Key:       "duplicate_string_key_test", // Same key
			Value:     "second value",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act - Insert first setting
		err := repo.SetString(ctx, setting1)
		require.NoError(t, err)

		// Act - Insert duplicate key (should update)
		err = repo.SetString(ctx, setting2)
		require.NoError(t, err)

		// Assert - Verify the value was updated
		retrievedSetting, err := repo.GetString(ctx, "duplicate_string_key_test")
		require.NoError(t, err)
		assert.Equal(t, "second value", retrievedSetting.Value)

		// Cleanup
		_, err = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", setting1.Key)
		require.NoError(t, err)
	})

	t.Run("setting existing migration values should update them", func(t *testing.T) {
		// Test updating initial migration data
		setting := &domain.Setting[string]{
			Key:       "company_name", // This exists in migration
			Value:     "different company",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act - Should succeed and update the existing value
		err := repo.SetString(ctx, setting)
		require.NoError(t, err)

		// Assert - Verify the value was updated
		retrievedSetting, err := repo.GetString(ctx, "company_name")
		require.NoError(t, err)
		assert.Equal(t, "different company", retrievedSetting.Value)

		// Cleanup - Reset to original migration value
		originalSetting := &domain.Setting[string]{
			Key:   "company_name",
			Value: "leviosa",
		}
		err = repo.SetString(ctx, originalSetting)
		require.NoError(t, err)
	})

	t.Run("nil setting should panic or cause error", func(t *testing.T) {
		// This test depends on your implementation - if you don't handle nil settings,
		// it might panic or cause a database error

		// Act & Assert
		assert.Panics(t, func() {
			_ = repo.SetString(ctx, nil)
		})
	})

	t.Run("context cancellation should return error", func(t *testing.T) {
		// Arrange
		canceledCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		setting := &domain.Setting[string]{
			Key:       "test_cancelled_context_string",
			Value:     "should not be inserted",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetString(canceledCtx, setting)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("context timeout should return error", func(t *testing.T) {
		// Arrange
		timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)
		defer cancel()
		time.Sleep(2 * time.Nanosecond) // Ensure timeout

		setting := &domain.Setting[string]{
			Key:       "test_timeout_context_string",
			Value:     "should not be inserted",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetString(timeoutCtx, setting)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout")
	})

	t.Run("empty key should work if allowed by schema", func(t *testing.T) {
		// This test assumes empty keys are allowed - adjust based on your schema constraints
		setting := &domain.Setting[string]{
			Key:       "",
			Value:     "value with empty key",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetString(ctx, setting)

		// Assert - This might error if your schema doesn't allow empty keys
		if err != nil {
			assert.Contains(t, err.Error(), "key")
		} else {
			// Cleanup if successful
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = ''")
		}
	})

	t.Run("whitespace-only key should work if allowed", func(t *testing.T) {
		setting := &domain.Setting[string]{
			Key:       "   ",
			Value:     "value with whitespace key",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetString(ctx, setting)

		// Assert
		if err != nil {
			// Schema might not allow whitespace-only keys
			assert.Contains(t, err.Error(), "key")
		} else {
			// Cleanup if successful
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = '   '")
		}
	})

	t.Run("database connection closed should return error", func(t *testing.T) {
		// This is harder to test without mocking, but you could create a separate test
		// that uses a closed pool if your test setup allows it

		// For now, we'll skip this test unless you have a way to simulate connection issues
		t.Skip("Database connection failure test requires specific setup")
	})
}

// Helper function to clean up test data
func cleanupTestStringSetting(t *testing.T, ctx context.Context, key string) {
	_, err := testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
	require.NoError(t, err)
}

// Benchmark test for performance measurement
func BenchmarkRepository_SetString(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setting := &domain.Setting[string]{
			Key:       fmt.Sprintf("benchmark_string_key_%d", i),
			Value:     fmt.Sprintf("benchmark_value_%d", i),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.SetString(ctx, setting)
		if err != nil {
			b.Fatalf("SetString failed: %v", err)
		}
	}

	// Cleanup after benchmark
	_, err := testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key LIKE 'benchmark_string_key_%'")
	if err != nil {
		b.Logf("Cleanup failed: %v", err)
	}
}

// Benchmark test for very long strings
func BenchmarkRepository_SetString_LongValues(b *testing.B) {
	ctx := context.Background()
	longValue := strings.Repeat("A", 10*1024) // 10KB string

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setting := &domain.Setting[string]{
			Key:       fmt.Sprintf("benchmark_long_string_%d", i),
			Value:     longValue,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.SetString(ctx, setting)
		if err != nil {
			b.Fatalf("SetString failed: %v", err)
		}
	}

	// Cleanup after benchmark
	_, err := testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key LIKE 'benchmark_long_string_%'")
	if err != nil {
		b.Logf("Cleanup failed: %v", err)
	}
}
