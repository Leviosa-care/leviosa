package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-unit-postgres TEST=TestSetInt

func TestSetInt(t *testing.T) {
	ctx := context.Background()

	t.Run("successful insertion", func(t *testing.T) {
		// Arrange
		setting := &domain.Setting[int]{
			Key:       "test_int_setting",
			Value:     42,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetInt(ctx, setting)

		// Assert
		require.NoError(t, err)

		// Verify the record was actually inserted
		var storedValue string
		var storedID int
		query := "SELECT id, value FROM settings.plain WHERE key = $1"
		err = testPool.QueryRow(ctx, query, setting.Key).Scan(&storedID, &storedValue)
		require.NoError(t, err)
		assert.Greater(t, storedID, 0) // ID should be auto-generated and positive
		assert.Equal(t, "42", storedValue)

		// Cleanup
		_, err = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", setting.Key)
		require.NoError(t, err)
	})

	t.Run("successful insertion with negative value", func(t *testing.T) {
		// Arrange
		setting := &domain.Setting[int]{
			Key:       "test_negative_int",
			Value:     -123,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetInt(ctx, setting)

		// Assert
		require.NoError(t, err)

		// Verify the record was inserted with correct negative value
		var storedValue string
		query := "SELECT value FROM settings.plain WHERE key = $1"
		err = testPool.QueryRow(ctx, query, setting.Key).Scan(&storedValue)
		require.NoError(t, err)
		assert.Equal(t, "-123", storedValue)

		// Cleanup
		_, err = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", setting.Key)
		require.NoError(t, err)
	})

	t.Run("successful insertion with zero value", func(t *testing.T) {
		// Arrange
		setting := &domain.Setting[int]{
			Key:       "test_zero_int",
			Value:     0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetInt(ctx, setting)

		// Assert
		require.NoError(t, err)

		// Verify the record was inserted
		var storedValue string
		query := "SELECT value FROM settings.plain WHERE key = $1"
		err = testPool.QueryRow(ctx, query, setting.Key).Scan(&storedValue)
		require.NoError(t, err)
		assert.Equal(t, "0", storedValue)

		// Cleanup
		_, err = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", setting.Key)
		require.NoError(t, err)
	})

	t.Run("duplicate key insertion should update existing value", func(t *testing.T) {
		// Arrange
		setting1 := &domain.Setting[int]{
			Key:       "duplicate_key_test",
			Value:     100,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		setting2 := &domain.Setting[int]{
			Key:       "duplicate_key_test", // Same key
			Value:     200,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act - Insert first setting
		err := repo.SetInt(ctx, setting1)
		require.NoError(t, err)

		// Act - Insert duplicate key (should update)
		err = repo.SetInt(ctx, setting2)
		require.NoError(t, err)

		// Assert - Verify the value was updated
		retrievedSetting, err := repo.GetInt(ctx, "duplicate_key_test")
		require.NoError(t, err)
		assert.Equal(t, 200, retrievedSetting.Value)

		// Cleanup
		_, err = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", setting1.Key)
		require.NoError(t, err)
	})

	t.Run("nil setting should panic or cause error", func(t *testing.T) {
		// This test depends on your implementation - if you don't handle nil settings,
		// it might panic or cause a database error

		// Act & Assert
		assert.Panics(t, func() {
			_ = repo.SetInt(ctx, nil)
		})
	})

	t.Run("context cancellation should return error", func(t *testing.T) {
		// Arrange
		canceledCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		setting := &domain.Setting[int]{
			Key:       "test_cancelled_context",
			Value:     999,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetInt(canceledCtx, setting)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("context timeout should return error", func(t *testing.T) {
		// Arrange
		timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)
		defer cancel()
		time.Sleep(2 * time.Nanosecond) // Ensure timeout

		setting := &domain.Setting[int]{
			Key:       "test_timeout_context",
			Value:     999,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetInt(timeoutCtx, setting)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout")
	})

	t.Run("empty key should work if allowed by schema", func(t *testing.T) {
		// This test assumes empty keys are allowed - adjust based on your schema constraints
		setting := &domain.Setting[int]{
			Key:       "",
			Value:     42,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetInt(ctx, setting)

		// Assert - This might error if your schema doesn't allow empty keys
		if err != nil {
			assert.Contains(t, err.Error(), "key")
		} else {
			// Cleanup if successful
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = ''")
		}
	})

	t.Run("very large integer values", func(t *testing.T) {
		// Test with maximum int value
		setting := &domain.Setting[int]{
			Key:       "test_max_int",
			Value:     2147483647, // Max int32 or use math.MaxInt for int
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		err := repo.SetInt(ctx, setting)

		// Assert
		require.NoError(t, err)

		// Verify
		var storedValue string
		query := "SELECT value FROM settings.plain WHERE key = $1"
		err = testPool.QueryRow(ctx, query, setting.Key).Scan(&storedValue)
		require.NoError(t, err)
		assert.Equal(t, "2147483647", storedValue)

		// Cleanup
		_, err = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", setting.Key)
		require.NoError(t, err)
	})

	t.Run("database connection closed should return error", func(t *testing.T) {
		// This is harder to test without mocking, but you could create a separate test
		// that uses a closed pool if your test setup allows it

		// For now, we'll skip this test unless you have a way to simulate connection issues
		t.Skip("Database connection failure test requires specific setup")
	})
}

// Helper function to clean up test data
func cleanupTestSetting(t *testing.T, ctx context.Context, key string) {
	_, err := testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
	require.NoError(t, err)
}

// Benchmark test for performance measurement
func BenchmarkRepository_SetInt(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setting := &domain.Setting[int]{
			Key:       fmt.Sprintf("benchmark_key_%d", i),
			Value:     i,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.SetInt(ctx, setting)
		if err != nil {
			b.Fatalf("SetInt failed: %v", err)
		}
	}

	// Cleanup after benchmark
	_, err := testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key LIKE 'benchmark_key_%'")
	if err != nil {
		b.Logf("Cleanup failed: %v", err)
	}
}
