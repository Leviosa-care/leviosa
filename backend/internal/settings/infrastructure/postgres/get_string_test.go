package postgres_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-unit-postgres TEST=TestGetString

func TestGetString(t *testing.T) {
	ctx := context.Background()

	t.Run("successful retrieval of string setting", func(t *testing.T) {
		// Arrange
		key := "test_string_key"
		expectedValue := "hello world"

		// Insert test data
		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, expectedValue, now, now)
		require.NoError(t, err)

		// Clean up after test
		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetString(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, key, result.Key)
		assert.Equal(t, expectedValue, result.Value)
		assert.NotZero(t, result.ID)
		assert.False(t, result.CreatedAt.IsZero())
		assert.False(t, result.UpdatedAt.IsZero())
	})

	t.Run("successful retrieval of empty string", func(t *testing.T) {
		// Arrange
		key := "test_empty_string"
		expectedValue := ""

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, expectedValue, now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetString(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedValue, result.Value)
	})

	t.Run("successful retrieval of string with whitespace", func(t *testing.T) {
		// Arrange
		key := "test_whitespace_string"
		expectedValue := "  hello world  "

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, expectedValue, now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetString(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedValue, result.Value)
	})

	t.Run("successful retrieval of string with special characters", func(t *testing.T) {
		// Arrange
		key := "test_special_chars"
		expectedValue := "Hello! @#$%^&*()_+-={}[]|\\:;\"'<>?,./"

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, expectedValue, now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetString(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedValue, result.Value)
	})

	t.Run("successful retrieval of unicode string", func(t *testing.T) {
		// Arrange
		key := "test_unicode_string"
		expectedValue := "Hello 世界 🌍 café naïve résumé"

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, expectedValue, now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetString(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedValue, result.Value)
	})

	t.Run("successful retrieval of multiline string", func(t *testing.T) {
		// Arrange
		key := "test_multiline_string"
		expectedValue := "Line 1\nLine 2\nLine 3\n"

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, expectedValue, now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetString(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedValue, result.Value)
		assert.Contains(t, result.Value, "\n")
	})

	t.Run("successful retrieval of JSON string", func(t *testing.T) {
		// Arrange
		key := "test_json_string"
		expectedValue := `{"name": "John Doe", "age": 30, "active": true}`

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, expectedValue, now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetString(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedValue, result.Value)
	})

	t.Run("successful retrieval of numeric string", func(t *testing.T) {
		// Arrange
		key := "test_numeric_string"
		expectedValue := "12345"

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, expectedValue, now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetString(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedValue, result.Value)
		// Verify it's still a string, not converted to int
		assert.IsType(t, "", result.Value)
	})

	t.Run("successful retrieval of very long string", func(t *testing.T) {
		// Arrange
		key := "test_long_string"
		expectedValue := strings.Repeat("A", 1000) // 1000 character string

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, expectedValue, now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetString(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedValue, result.Value)
		assert.Len(t, result.Value, 1000)
	})

	t.Run("key not found returns repository not found error", func(t *testing.T) {
		// Arrange
		nonExistentKey := "non_existent_string_key_12345"

		// Act
		result, err := repo.GetString(ctx, nonExistentKey)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)

		// Check if it's the expected error type (adjust based on your errs package)
		assert.Contains(t, err.Error(), "string value for key 'non_existent_string_key_12345'")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Arrange
		key := "test_context_cancel_string"

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, "test value", now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Create cancelled context
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel()

		// Act
		result, err := repo.GetString(cancelledCtx, key)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("successful retrieval with tab characters", func(t *testing.T) {
		// Arrange
		key := "test_tab_string"
		expectedValue := "Column1\tColumn2\tColumn3"

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, expectedValue, now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetString(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedValue, result.Value)
		assert.Contains(t, result.Value, "\t")
	})

	t.Run("successful retrieval of single character", func(t *testing.T) {
		// Arrange
		key := "test_single_char"
		expectedValue := "X"

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, expectedValue, now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetString(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedValue, result.Value)
		assert.Len(t, result.Value, 1)
	})

	t.Run("successful retrieval of string with SQL injection attempt", func(t *testing.T) {
		// Arrange
		key := "test_sql_injection"
		expectedValue := "'; DROP TABLE settings.plain; --"

		insertQuery := `
			INSERT INTO settings.plain (key, value, created_at, updated_at) 
			VALUES ($1, $2, $3, $4)`

		now := time.Now()
		_, err := testPool.Exec(ctx, insertQuery, key, expectedValue, now, now)
		require.NoError(t, err)

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
		}()

		// Act
		result, err := repo.GetString(ctx, key)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedValue, result.Value)
		// Verify the table still exists by doing another query
		_, err = testPool.Exec(ctx, "SELECT COUNT(*) FROM settings.plain")
		assert.NoError(t, err, "Table should still exist - SQL injection should be prevented")
	})
}

// Benchmark test to check performance with different string sizes
func BenchmarkRepository_GetString(b *testing.B) {
	ctx := context.Background()

	testCases := []struct {
		name  string
		value string
	}{
		{"short", "short"},
		{"medium", strings.Repeat("medium", 20)},
		{"long", strings.Repeat("long", 250)},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			key := "benchmark_string_key_" + tc.name

			// Setup test data
			insertQuery := `
				INSERT INTO settings.plain (key, value, created_at, updated_at) 
				VALUES ($1, $2, $3, $4)`

			now := time.Now()
			_, err := testPool.Exec(ctx, insertQuery, key, tc.value, now, now)
			require.NoError(b, err)

			defer func() {
				_, _ = testPool.Exec(ctx, "DELETE FROM settings.plain WHERE key = $1", key)
			}()

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := repo.GetString(ctx, key)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
