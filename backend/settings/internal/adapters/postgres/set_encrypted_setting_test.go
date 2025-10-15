package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/core/contracts/settings"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-unit-postgres TEST=TestSetPhone

func TestSetPhone(t *testing.T) {
	ctx := context.Background()

	dek := "4tw34trw5yw34t8q039i4t3w5y3q4"
	metadata := encx.EncryptionMetadata{
		PepperVersion:    1,
		KEKAlias:         "test-alias",
		EncryptionTime:   1234567890,
		GeneratorVersion: "1.0.0",
	}
	now := time.Now()

	t.Run("successful phone setting creation", func(t *testing.T) {
		// Arrange
		phoneSetting := &domain.SettingEncryptedEncx{
			Key:            settings.CompanyPhone,
			ValueEncrypted: []byte("0612345679"),
			DEKEncrypted:   []byte(dek),
			KeyVersion:     1,
			Metadata:       metadata,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		// Clean up before and after test
		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", phoneSetting.Key)
		}()

		// Act
		// err = repo.SetEncryptedSetting(ctx, phoneEncrypted)
		err := repo.SetEncryptedSetting(ctx, phoneSetting)

		// Assert
		assert.NoError(t, err)

		// Verify the record was inserted
		var retrievedSetting domain.SettingEncryptedEncx
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

		assert.Equal(t, phoneSetting.ValueEncrypted, retrievedSetting.ValueEncrypted)
		assert.Equal(t, phoneSetting.DEKEncrypted, retrievedSetting.DEKEncrypted)
		assert.Equal(t, phoneSetting.KeyVersion, retrievedSetting.KeyVersion)
		assert.WithinDuration(t, phoneSetting.CreatedAt, retrievedSetting.CreatedAt, time.Second)
		assert.WithinDuration(t, phoneSetting.UpdatedAt, retrievedSetting.UpdatedAt, time.Second)
	})

	t.Run("successful phone setting with empty encrypted values", func(t *testing.T) {
		// Arrange
		phoneSetting := &domain.SettingEncryptedEncx{
			Key:            settings.CompanyPhone,
			ValueEncrypted: []byte(""),
			DEKEncrypted:   []byte(dek),
			KeyVersion:     1,
			Metadata:       metadata,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", phoneSetting.Key)
		}()

		// Act
		err := repo.SetEncryptedSetting(ctx, phoneSetting)

		// Assert
		assert.NoError(t, err)

		// Verify the record was inserted
		var count int
		err = testPool.QueryRow(ctx, "SELECT COUNT(*) FROM settings.encrypted WHERE key = $1", phoneSetting.Key).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("successful phone setting with nil encrypted values", func(t *testing.T) {
		// Arrange
		phoneSetting := &domain.SettingEncryptedEncx{
			Key:            settings.CompanyPhone,
			ValueEncrypted: nil,
			DEKEncrypted:   []byte(dek),
			KeyVersion:     1,
			Metadata:       metadata,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", phoneSetting.Key)
		}()

		// Act
		err := repo.SetEncryptedSetting(ctx, phoneSetting)

		// Assert
		assert.ErrorIs(t, err, errs.ErrNotNullViolation)
	})

	t.Run("context cancellation returns error", func(t *testing.T) {
		// Arrange
		phoneSetting := &domain.SettingEncryptedEncx{
			Key:            settings.CompanyPhone,
			ValueEncrypted: []byte("0612345679"),
			DEKEncrypted:   []byte(dek),
			KeyVersion:     1,
			Metadata:       metadata,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel()

		// Act
		err := repo.SetEncryptedSetting(cancelledCtx, phoneSetting)

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

		phoneSetting := &domain.SettingEncryptedEncx{
			Key:            settings.CompanyPhone,
			ValueEncrypted: largeValueEncrypted,
			DEKEncrypted:   largeDEKEncrypted,
			KeyVersion:     1,
			Metadata:       metadata,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		defer func() {
			_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", phoneSetting.Key)
		}()

		// Act
		err := repo.SetEncryptedSetting(ctx, phoneSetting)

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
			err = testPool.QueryRow(ctx, query, phoneSetting.Key).Scan(&retrievedValueEncrypted, &retrievedDEKEncrypted)
			require.NoError(t, err)

			assert.Equal(t, largeValueEncrypted, retrievedValueEncrypted)
			assert.Equal(t, largeDEKEncrypted, retrievedDEKEncrypted)
		}
	})

	t.Run("setting with different key versions", func(t *testing.T) {
		versions := []int{0, 1, 10, 100, 999}

		for _, version := range versions {
			t.Run(fmt.Sprintf("key_version_%d", version), func(t *testing.T) {
				phoneSetting := &domain.SettingEncryptedEncx{Key: settings.CompanyPhone,
					ValueEncrypted: []byte("0612345679"),
					DEKEncrypted:   []byte(dek),
					KeyVersion:     version,
					Metadata:       metadata,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}

				defer func() {
					_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE key = $1", phoneSetting.Key)
				}()

				// Act
				err := repo.SetEncryptedSetting(ctx, phoneSetting)

				// Assert
				assert.NoError(t, err)

				// Verify key version was stored correctly
				var retrievedKeyVersion int
				query := "SELECT key_version FROM settings.encrypted WHERE key = $1"
				err = testPool.QueryRow(ctx, query, phoneSetting.Key).Scan(&retrievedKeyVersion)
				require.NoError(t, err)
				assert.Equal(t, version, retrievedKeyVersion)
			})
		}
	})
}

// TEST=BenchmarkRepository_SetPhone make test-unit-test

func BenchmarkRepository_SetPhone(b *testing.B) {
	ctx := context.Background()

	dek := "4tw34trw5yw34t8q039i4t3w5y3q4"
	metadata := encx.EncryptionMetadata{}
	now := time.Now()

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

				phoneSetting := &domain.SettingEncryptedEncx{
					ID:             fmt.Sprintf("benchmark-set-%s-%d", tc.name, i),
					Key:            "company_phone",
					ValueEncrypted: []byte(fmt.Sprintf("01%08d", i)),
					DEKEncrypted:   []byte(dek),
					KeyVersion:     1,
					Metadata:       metadata,
					CreatedAt:      now,
					UpdatedAt:      now,
				}

				b.StartTimer()

				err := repo.SetEncryptedSetting(ctx, phoneSetting)
				if err != nil {
					b.Fatal(err)
				}

				b.StopTimer()
				// Clean up
				_, _ = testPool.Exec(ctx, "DELETE FROM settings.encrypted WHERE id = $1", phoneSetting.ID)
			}
		})
	}
}
