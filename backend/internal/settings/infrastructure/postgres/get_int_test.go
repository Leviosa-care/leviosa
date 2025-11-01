package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-unit-postgres TEST=TestGetInt

func TestGetInt(t *testing.T) {
	ctx := context.Background()

	t.Run("successful retrieval of int setting", func(t *testing.T) {
		// Arrange
		key := "test_int_key"
		expectedValue := 42

		// Insert test data
		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, "42", now, now)
		require.NoError(t, err)

		// Clean up after test
		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetInt(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, key, result.Key)
		assert.Equal(t, expectedValue, result.Value)
		assert.NotZero(t, result.ID)
		assert.False(t, result.CreatedAt.IsZero())
		assert.False(t, result.UpdatedAt.IsZero())
	})

	t.Run("successful retrieval of negative int", func(t *testing.T) {
		// Arrange
		key := "test_negative_int"
		expectedValue := -123

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, "-123", now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetInt(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedValue, result.Value)
	})

	t.Run("successful retrieval of zero value", func(t *testing.T) {
		// Arrange
		key := "test_zero_int"
		expectedValue := 0

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, "0", now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetInt(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedValue, result.Value)
	})

	t.Run("key not found returns repository not found error", func(t *testing.T) {
		// Arrange
		nonExistentKey := "non_existent_key_12345"

		// Act
		result, err := repo.GetInt(ctx, nonExistentKey)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)

		// Check if it's the expected error type (adjust based on your errs package)
		// You might want to check the error message or type here
		assert.Contains(t, err.Error(), "int value for key 'non_existent_key_12345'")
	})

	t.Run("invalid int value returns conversion error", func(t *testing.T) {
		// Arrange
		key := "test_invalid_int"

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, "not_an_integer", now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetInt(ctx, key)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to convert retrieved string value into int")
	})

	t.Run("float value stored as string returns conversion error", func(t *testing.T) {
		// Arrange
		key := "test_float_as_string"

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, "42.5", now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetInt(ctx, key)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to convert retrieved string value into int")
	})

	t.Run("empty string value returns conversion error", func(t *testing.T) {
		// Arrange
		key := "test_empty_string"

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, "", now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetInt(ctx, key)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to convert retrieved string value into int")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Arrange
		key := "test_context_cancel"

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, "42", now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Create cancelled context
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel()

		// Act
		result, err := repo.GetInt(cancelledCtx, key)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("very large int value", func(t *testing.T) {
		// Arrange
		key := "test_large_int"
		expectedValue := 2147483647 // Max int32

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, "2147483647", now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetInt(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedValue, result.Value)
	})

	t.Run("whitespace in value returns conversion error", func(t *testing.T) {
		// Arrange
		key := "test_whitespace_int"

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, " 42 ", now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetInt(ctx, key)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to convert retrieved string value into int")
	})
}

// Benchmark test to check performance
func BenchmarkRepository_GetInt(b *testing.B) {
	ctx := context.Background()
	key := "benchmark_int_key"

	// Setup test data
	insertQuery := `
		INSERT INTO settings.plain (key, value, created_at, updated_at) 
		VALUES ($1, $2, $3, $4)`

	now := time.Now()
	_, err := testPool.Exec(ctx, insertQuery, key, "42", now, now)
	require.NoError(b, err)

	defer func() {
		_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
	}()

	b.ResetTimer()

	for range b.N {
		_, err := repo.GetInt(ctx, key)
		if err != nil {
			b.Fatal(err)
		}
	}
}
