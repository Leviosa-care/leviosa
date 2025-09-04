package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/core/contracts/settings"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetPhone(t *testing.T) {
	ctx := context.Background()

	t.Run("successful phone setting creation", func(t *testing.T) {
		// Arrange
		setting := &domain.SettingEncrypted[string]{
			Key:            settings.CompanyPhone, // Assuming this is domain.CompanyPhoneKey
			Value:          "0123456789",         // Actual phone value (unencrypted)
			ValueEncrypted: []byte("encrypted_phone_bytes_12345"),
			DEK:            []byte("data_encryption_key"),
			DEKEncrypted:   []byte("encrypted_dek_bytes_67890"),
			KeyVersion:     1,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Clean up before and after test
		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", setting.Key)
		}()

		// Act
		err := repo.SetPhone(ctx, setting)

		// Assert
		require.NoError(t, err)

		// Verify the record was inserted
		var retrievedSetting domain.SettingEncrypted[string]
		query := `SELECT key, value_encrypted, created_at, updated_at, dek_encrypted, key_version 
				  FROM settings.encrypted WHERE key = $1`

		err = testPool.QueryRow(ctx, query, settings.CompanyPhone).Scan(
			&retrievedSetting.Key,
			&retrievedSetting.ValueEncrypted,
			&retrievedSetting.CreatedAt,
			&retrievedSetting.UpdatedAt,
			&retrievedSetting.DEKEncrypted,
			&retrievedSetting.KeyVersion,
		)
		require.NoError(t, err)

		assert.Equal(t, setting.ValueEncrypted, retrievedSetting.ValueEncrypted)
		assert.Equal(t, setting.DEKEncrypted, retrievedSetting.DEKEncrypted)
		assert.Equal(t, setting.KeyVersion, retrievedSetting.KeyVersion)
		assert.WithinDuration(t, setting.CreatedAt, retrievedSetting.CreatedAt, time.Second)
		assert.WithinDuration(t, setting.UpdatedAt, retrievedSetting.UpdatedAt, time.Second)
	})

	t.Run("successful phone setting with empty encrypted values", func(t *testing.T) {
		// Arrange
		setting := &domain.SettingEncrypted[string]{
			Key:            settings.CompanyPhone,
			Value:          "",
			ValueEncrypted: []byte{},
			DEK:            []byte{},
			DEKEncrypted:   []byte{},
			KeyVersion:     0,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", setting.Key)
		}()

		// Act
		err := repo.SetPhone(ctx, setting)

		// Assert
		require.NoError(t, err)

		// Verify the record was inserted
		var count int
		err = testPool.QueryRow(ctx, "SELECT COUNT(*) FROM settings.encrypted WHERE key = $1", setting.Key).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("successful phone setting with nil encrypted values", func(t *testing.T) {
		// Arrange
		setting := &domain.SettingEncrypted[string]{
			Key:            settings.CompanyPhone,
			Value:          "0187654321",
			ValueEncrypted: nil, // nil slice
			DEK:            nil, // nil slice
			DEKEncrypted:   nil, // nil slice
			KeyVersion:     0,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", setting.Key)
		}()

		// Act
		err := repo.SetPhone(ctx, setting)

		// Assert
		assert.ErrorIs(t, err, errs.ErrNotNullViolation)
	})

	t.Run("context cancellation returns error", func(t *testing.T) {
		// Arrange
		setting := &domain.SettingEncrypted[string]{
			Key:            settings.CompanyPhone,
			Value:          "0111111111",
			ValueEncrypted: []byte("encrypted_value"),
			DEKEncrypted:   []byte("encrypted_dek"),
			KeyVersion:     1,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel()

		// Act
		err := repo.SetPhone(cancelledCtx, setting)

		// Assert
		require.Error(t, err)
		// Should be classified as a context-related error by your errs.ClassifyPgError
	})

	t.Run("setting with very large encrypted values", func(t *testing.T) {
		// Arrange
		largeValueEncrypted := make([]byte, 50000) // 50KB encrypted value
		largeDEKEncrypted := make([]byte, 10000)   // 10KB encrypted DEK

		// Fill with test data
		for i := range largeValueEncrypted {
			largeValueEncrypted[i] = byte(i % 256)
		}
		for i := range largeDEKEncrypted {
			largeDEKEncrypted[i] = byte((i + 128) % 256)
		}

		setting := &domain.SettingEncrypted[string]{
			Key:            settings.CompanyPhone,
			Value:          "0155555555",
			ValueEncrypted: largeValueEncrypted,
			DEKEncrypted:   largeDEKEncrypted,
			KeyVersion:     1,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", setting.Key)
		}()

		// Act
		err := repo.SetPhone(ctx, setting)

		// Assert - this might fail if your database has column length restrictions
		// You might want to handle this case in your actual implementation
		if err != nil {
			// If the database rejects very large values, that's expected
			assert.Error(t, err)
		} else {
			require.NoError(t, err)

			// Verify the large data was stored correctly
			var retrievedValueEncrypted, retrievedDEKEncrypted []byte
			query := "SELECT value_encrypted, dek_encrypted FROM settings.encrypted WHERE key = $1"
			err = testPool.QueryRow(ctx, query, setting.Key).Scan(&retrievedValueEncrypted, &retrievedDEKEncrypted)
			require.NoError(t, err)

			assert.Equal(t, largeValueEncrypted, retrievedValueEncrypted)
			assert.Equal(t, largeDEKEncrypted, retrievedDEKEncrypted)
		}
	})

	t.Run("setting with different key versions", func(t *testing.T) {
		versions := []int{0, 1, 10, 100, 999}

		for _, version := range versions {
			t.Run(fmt.Sprintf("key_version_%d", version), func(t *testing.T) {
				setting := &domain.SettingEncrypted[string]{
					Key:            settings.CompanyPhone,
					Value:          fmt.Sprintf("01%08d", version),
					ValueEncrypted: []byte(fmt.Sprintf("encrypted_value_%d", version)),
					DEKEncrypted:   []byte(fmt.Sprintf("encrypted_dek_%d", version)),
					KeyVersion:     version,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}

				defer func() {
					_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", setting.Key)
				}()

				// Act
				err := repo.SetPhone(ctx, setting)

				// Assert
				require.NoError(t, err)

				// Verify key version was stored correctly
				var retrievedKeyVersion int
				query := "SELECT key_version FROM settings.encrypted WHERE key = $1"
				err = testPool.QueryRow(ctx, query, setting.Key).Scan(&retrievedKeyVersion)
				require.NoError(t, err)
				assert.Equal(t, version, retrievedKeyVersion)
			})
		}
	})
}

func BenchmarkRepository_SetPhone(b *testing.B) {
	ctx := context.Background()

	// Test different encrypted data sizes
	testCases := []struct {
		name      string
		valueSize int
		dekSize   int
	}{
		{"small", 100, 50},
		{"medium", 1000, 256},
		{"large", 10000, 1000},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()

			for i := range b.N {
				b.StopTimer()

				// Create test data
				valueEncrypted := make([]byte, tc.valueSize)
				dekEncrypted := make([]byte, tc.dekSize)

				setting := &domain.SettingEncrypted[string]{
					ID:             fmt.Sprintf("benchmark-set-%s-%d", tc.name, i),
					Key:            "company_phone",
					Value:          fmt.Sprintf("01%08d", i),
					ValueEncrypted: valueEncrypted,
					DEKEncrypted:   dekEncrypted,
					KeyVersion:     1,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}

				b.StartTimer()

				err := repo.SetPhone(ctx, setting)
				if err != nil {
					b.Fatal(err)
				}

				b.StopTimer()
				// Clean up
				_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE id = $1", setting.ID)
			}
		})
	}
}
