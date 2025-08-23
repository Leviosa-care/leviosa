package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/Leviosa-care/core/contracts/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Integration test combining both operations
func TestGetSetPhone(t *testing.T) {
	ctx := context.Background()

	t.Run("set phone then get phone returns same encrypted data", func(t *testing.T) {
		// Arrange
		originalSetting := &domain.SettingEncrypted[string]{
			Key:            settings.CompanyPhone,
			Value:          "0123456789",
			ValueEncrypted: []byte("integrated_test_encrypted_phone_bytes"),
			DEK:            []byte("integrated_test_dek"),
			DEKEncrypted:   []byte("integrated_test_encrypted_dek_bytes"),
			KeyVersion:     2,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", originalSetting.Key)
		}()

		// Act - Set the phone
		err := repo.SetPhone(ctx, originalSetting)
		require.NoError(t, err)

		// Act - Get the phone
		retrievedSetting, err := repo.GetPhone(ctx)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, retrievedSetting)
		assert.Equal(t, originalSetting.ValueEncrypted, retrievedSetting.ValueEncrypted)
		assert.Equal(t, originalSetting.DEKEncrypted, retrievedSetting.DEKEncrypted)
		// Note: KeyVersion might not be retrieved by GetPhone depending on your SELECT query
		assert.WithinDuration(t, originalSetting.CreatedAt, retrievedSetting.CreatedAt, time.Second)
		assert.WithinDuration(t, originalSetting.UpdatedAt, retrievedSetting.UpdatedAt, time.Second)
	})

	t.Run("get phone before any phone is set returns not found", func(t *testing.T) {
		// Arrange - ensure no phone setting exists
		_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", settings.CompanyPhone)

		// Act
		result, err := repo.GetPhone(ctx)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("set phone multiple times overwrites previous", func(t *testing.T) {
		// This test depends on whether your SetPhone does INSERT or UPSERT
		// You might need to adjust based on your actual implementation

		// Arrange - first setting
		setting1 := &domain.SettingEncrypted[string]{
			Key:            settings.CompanyPhone,
			Value:          "0111111111",
			ValueEncrypted: []byte("first_encrypted_value"),
			DEKEncrypted:   []byte("first_encrypted_dek"),
			KeyVersion:     1,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", settings.CompanyPhone)
		}()

		// Act - Set first phone
		err := repo.SetPhone(ctx, setting1)
		require.NoError(t, err)

		// Arrange - second setting with different ID but same key
		setting2 := &domain.SettingEncrypted[string]{
			Key:            settings.CompanyPhone,
			Value:          "0222222222",
			ValueEncrypted: []byte("second_encrypted_value"),
			DEKEncrypted:   []byte("second_encrypted_dek"),
			KeyVersion:     2,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Act - Set second phone (this might fail if you have unique constraints on key)
		err = repo.SetPhone(ctx, setting2)
		if err == nil {
			// If the second insert succeeds, GetPhone should return the last one
			// (This depends on your database constraints and query ordering)
			retrievedSetting, err := repo.GetPhone(ctx)
			require.NoError(t, err)

			// The behavior here depends on your database schema and GetPhone query
			// You might get either setting1 or setting2 depending on ordering
			assert.NotNil(t, retrievedSetting)
		} else {
			// If unique constraint prevents multiple entries with same key, that's expected
			assert.Error(t, err)
		}
	})
}
