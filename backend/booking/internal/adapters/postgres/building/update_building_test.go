package buildingRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/require"
)

func TestUpdateBuilding(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully update existing building", func(t *testing.T) {
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
			[]byte("encrypted_original_name"),
			[]byte("encrypted_original_address"),
			[]byte("encrypted_original_city"),
			[]byte("encrypted_original_postal"),
			[]byte("encrypted_original_country"),
			[]byte("encrypted_original_description"),
			[]byte("encrypted_original_phone"),
			[]byte("encrypted_original_email"),
			true,
			now,
			now,
			[]byte("mock_dek_data"),
			1,
			`{"kek_alias":"test","encryption_time":12345}`,
		)
		require.NoError(t, err)

		// Prepare updated building data
		updatedTime := now.Add(time.Hour)
		updatedBuilding := &domain.BuildingEncx{
			ID:                     buildingID,
			NameEncrypted:          []byte("encrypted_updated_name"),
			AddressEncrypted:       []byte("encrypted_updated_address"),
			CityEncrypted:          []byte("encrypted_updated_city"),
			PostalCodeEncrypted:    []byte("encrypted_updated_postal"),
			CountryEncrypted:       []byte("encrypted_updated_country"),
			DescriptionEncrypted:   []byte("encrypted_updated_description"),
			PhoneEncrypted:         []byte("encrypted_updated_phone"),
			EmailEncrypted:         []byte("encrypted_updated_email"),
			IsActive:               false,
			CreatedAt:              now,
			UpdatedAt:              updatedTime,
			DEKEncrypted:           []byte("mock_dek_data"),
			KeyVersion:             1,
			Metadata:               encx.EncryptionMetadata{},
		}

		// Test repository Update method
		err = repo.Update(ctx, updatedBuilding)
		require.NoError(t, err)

		// Verify the building was updated by querying directly
		var nameEncrypted, addressEncrypted, isActive bool
		var updatedAt time.Time
		err = testPool.QueryRow(ctx,
			"SELECT name_encrypted, address_encrypted, is_active, updated_at FROM booking.buildings WHERE id = $1",
			buildingID).Scan(&nameEncrypted, &addressEncrypted, &isActive, &updatedAt)
		require.NoError(t, err)
		require.Equal(t, []byte("encrypted_updated_name"), nameEncrypted)
		require.Equal(t, []byte("encrypted_updated_address"), addressEncrypted)
		require.False(t, isActive)
		require.Equal(t, updatedTime.Unix(), updatedAt.Unix()) // Compare timestamps approximately
	})

	t.Run("should return error when updating non-existent building", func(t *testing.T) {
		// Clean up before test
		clearBuildingsTable(t, ctx)

		// Try to update non-existent building
		nonExistentID := uuid.New()
		updatedBuilding := &domain.BuildingEncx{
			ID:                     nonExistentID,
			NameEncrypted:          []byte("encrypted_updated_name"),
			AddressEncrypted:       []byte("encrypted_updated_address"),
			CityEncrypted:          []byte("encrypted_updated_city"),
			PostalCodeEncrypted:    []byte("encrypted_updated_postal"),
			CountryEncrypted:       []byte("encrypted_updated_country"),
			DescriptionEncrypted:   []byte("encrypted_updated_description"),
			PhoneEncrypted:         []byte("encrypted_updated_phone"),
			EmailEncrypted:         []byte("encrypted_updated_email"),
			IsActive:               true,
			CreatedAt:              time.Now(),
			UpdatedAt:              time.Now(),
			DEKEncrypted:           []byte("mock_dek_data"),
			KeyVersion:             1,
			Metadata:               encx.EncryptionMetadata{},
		}

		err := repo.Update(ctx, updatedBuilding)
		require.Error(t, err)
	})

	t.Run("should handle partial updates correctly", func(t *testing.T) {
		// Clean up before test
		clearBuildingsTable(t, ctx)

		// Insert test building
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
			[]byte("encrypted_original_name"),
			[]byte("encrypted_original_address"),
			[]byte("encrypted_original_city"),
			[]byte("encrypted_original_postal"),
			[]byte("encrypted_original_country"),
			[]byte("encrypted_original_description"),
			[]byte("encrypted_original_phone"),
			[]byte("encrypted_original_email"),
			true,
			now,
			now,
			[]byte("mock_dek_data"),
			1,
			`{"kek_alias":"test","encryption_time":12345}`,
		)
		require.NoError(t, err)

		// Update only some fields
		updatedBuilding := &domain.BuildingEncx{
			ID:                     buildingID,
			NameEncrypted:          []byte("encrypted_partially_updated_name"),
			AddressEncrypted:       []byte("encrypted_original_address"), // unchanged
			CityEncrypted:          []byte("encrypted_original_city"),     // unchanged
			PostalCodeEncrypted:    []byte("encrypted_original_postal"),   // unchanged
			CountryEncrypted:       []byte("encrypted_original_country"),  // unchanged
			DescriptionEncrypted:   []byte("encrypted_updated_description"),
			PhoneEncrypted:         []byte("encrypted_original_phone"),    // unchanged
			EmailEncrypted:         []byte("encrypted_updated_email"),
			IsActive:               true,
			CreatedAt:              now,
			UpdatedAt:              now.Add(time.Hour),
			DEKEncrypted:           []byte("mock_dek_data"),
			KeyVersion:             1,
			Metadata:               encx.EncryptionMetadata{},
		}

		// Test repository Update method
		err = repo.Update(ctx, updatedBuilding)
		require.NoError(t, err)

		// Verify specific fields were updated
		var nameEncrypted, descriptionEncrypted, emailEncrypted, addressEncrypted []byte
		err = testPool.QueryRow(ctx,
			"SELECT name_encrypted, description_encrypted, email_encrypted, address_encrypted FROM booking.buildings WHERE id = $1",
			buildingID).Scan(&nameEncrypted, &descriptionEncrypted, &emailEncrypted, &addressEncrypted)
		require.NoError(t, err)
		require.Equal(t, []byte("encrypted_partially_updated_name"), nameEncrypted)
		require.Equal(t, []byte("encrypted_updated_description"), descriptionEncrypted)
		require.Equal(t, []byte("encrypted_updated_email"), emailEncrypted)
		require.Equal(t, []byte("encrypted_original_address"), addressEncrypted) // unchanged
	})
}