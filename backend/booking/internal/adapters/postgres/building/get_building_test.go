package buildingRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetBuildingByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve existing building", func(t *testing.T) {
		// Clean up before test
		clearBuildingsTable(t, ctx)

		// Insert test building directly into database
		buildingID := uuid.New()
		now := time.Now()

		_, err := testPool.Exec(ctx, `
			INSERT INTO booking.buildings (
				id, name_encrypted, address_encrypted, city_encrypted,
				postal_code_encrypted, country_encrypted, description_encrypted,
				phone_encrypted, email_encrypted, is_active, created_at, updated_at,
				dek_encrypted, key_version, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		`,
			buildingID,
			[]byte("encrypted_test_building_name"),
			[]byte("encrypted_123_test_street"),
			[]byte("encrypted_test_city"),
			[]byte("encrypted_12345"),
			[]byte("encrypted_test_country"),
			[]byte("encrypted_test_description"),
			[]byte("encrypted_+1234567890"),
			[]byte("encrypted_test@example.com"),
			true,
			now,
			now,
			[]byte("mock_dek_data"),
			1,
			`{"kek_alias":"test","encryption_time":12345}`,
		)
		require.NoError(t, err)

		// Test repository GetByID method
		buildingEncx, err := repo.GetByID(ctx, buildingID)
		require.NoError(t, err)
		require.NotNil(t, buildingEncx)

		// Verify retrieved data
		require.Equal(t, buildingID, buildingEncx.ID)
		require.Equal(t, []byte("encrypted_test_building_name"), buildingEncx.NameEncrypted)
		require.Equal(t, []byte("encrypted_123_test_street"), buildingEncx.AddressEncrypted)
		require.Equal(t, []byte("encrypted_test_city"), buildingEncx.CityEncrypted)
		require.Equal(t, []byte("encrypted_12345"), buildingEncx.PostalCodeEncrypted)
		require.Equal(t, []byte("encrypted_test_country"), buildingEncx.CountryEncrypted)
		require.Equal(t, []byte("encrypted_test_description"), buildingEncx.DescriptionEncrypted)
		require.Equal(t, []byte("encrypted_+1234567890"), buildingEncx.PhoneEncrypted)
		require.Equal(t, []byte("encrypted_test@example.com"), buildingEncx.EmailEncrypted)
		require.True(t, buildingEncx.IsActive)
		require.Equal(t, 1, buildingEncx.KeyVersion)
		require.Equal(t, []byte("mock_dek_data"), buildingEncx.DEKEncrypted)
	})

	t.Run("should return error for non-existent building ID", func(t *testing.T) {
		// Clean up before test
		clearBuildingsTable(t, ctx)

		// Try to get non-existent building
		nonExistentID := uuid.New()
		buildingEncx, err := repo.GetByID(ctx, nonExistentID)
		require.Error(t, err)
		require.Nil(t, buildingEncx)
	})

	t.Run("should handle database connection errors gracefully", func(t *testing.T) {
		// Test with invalid UUID format
		invalidID := uuid.UUID{}
		buildingEncx, err := repo.GetByID(ctx, invalidID)
		require.Error(t, err)
		require.Nil(t, buildingEncx)
	})
}