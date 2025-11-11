package buildingRepository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdateBuilding TEST_PATH=internal/booking/infrastructure/postgres/building/update_building_test.go

func TestUpdateBuilding(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully update existing building", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert test building directly into database
		buildingEncx := tb.NewTestBuildingEncx(t)
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Get original building to compare UpdatedAt later
		originalBuilding, err := tb.GetBuildingEncxByID(t, ctx, testPool, buildingEncx.ID)
		require.NoError(t, err)

		// Small delay to ensure UpdatedAt will be different
		time.Sleep(10 * time.Millisecond)

		updatedNameEncrypted := []byte("encrypted_updated_name")
		updatedAddressEncrypted := []byte("encrypted_updated_address")
		updatedCityEncrypted := []byte("encrypted_updated_city")
		updatedPostalCodeEncrypted := []byte("encrypted_updated_postal")
		updatedCountryEncrypted := []byte("encrypted_updated_country")
		updatedDescriptionEncrypted := []byte("encrypted_updated_description")
		updatedPhoneEncrypted := []byte("encrypted_updated_phone")
		updatedEmailEncrypted := []byte("encrypted_updated_email")
		updatedIsActive := false

		// Prepare updated building data
		// Note: UpdatedAt will be set by database trigger, not by this struct
		updatedBuilding := &domain.BuildingEncx{
			ID:                   buildingEncx.ID,
			NameEncrypted:        updatedNameEncrypted,
			AddressEncrypted:     updatedAddressEncrypted,
			CityEncrypted:        updatedCityEncrypted,
			PostalCodeEncrypted:  updatedPostalCodeEncrypted,
			CountryEncrypted:     updatedCountryEncrypted,
			DescriptionEncrypted: updatedDescriptionEncrypted,
			PhoneEncrypted:       updatedPhoneEncrypted,
			EmailEncrypted:       updatedEmailEncrypted,
			IsActive:             updatedIsActive,
			CreatedAt:            buildingEncx.CreatedAt,
			UpdatedAt:            time.Now(), // This will be overridden by DB trigger
			DEKEncrypted:         []byte("mock_dek_data"),
			KeyVersion:           1,
			Metadata:             encx.EncryptionMetadata{},
		}

		// Test repository Update method
		err = repo.Update(ctx, updatedBuilding)
		require.NoError(t, err)

		// Verify the building was updated by querying directly
		retrievedBuildingEncx, err := tb.GetBuildingEncxByID(t, ctx, testPool, buildingEncx.ID)
		require.NoError(t, err)
		assert.Equal(t, updatedNameEncrypted, retrievedBuildingEncx.NameEncrypted)
		assert.Equal(t, updatedAddressEncrypted, retrievedBuildingEncx.AddressEncrypted)
		assert.Equal(t, updatedCityEncrypted, retrievedBuildingEncx.CityEncrypted)
		assert.Equal(t, updatedPostalCodeEncrypted, retrievedBuildingEncx.PostalCodeEncrypted)
		assert.Equal(t, updatedCountryEncrypted, retrievedBuildingEncx.CountryEncrypted)
		assert.Equal(t, updatedDescriptionEncrypted, retrievedBuildingEncx.DescriptionEncrypted)
		assert.Equal(t, updatedPhoneEncrypted, retrievedBuildingEncx.PhoneEncrypted)
		assert.Equal(t, updatedEmailEncrypted, retrievedBuildingEncx.EmailEncrypted)
		assert.Equal(t, updatedIsActive, retrievedBuildingEncx.IsActive)

		// Verify UpdatedAt was changed by the database trigger
		assert.True(t, retrievedBuildingEncx.UpdatedAt.After(originalBuilding.UpdatedAt),
			fmt.Sprintf("UpdatedAt should have changed: original=%v, retrieved=%v",
				originalBuilding.UpdatedAt, retrievedBuildingEncx.UpdatedAt))
	})

	t.Run("should return error when updating non-existent building", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Try to update non-existent building
		nonExistentID := uuid.New()
		updatedBuilding := &domain.BuildingEncx{
			ID:                   nonExistentID,
			NameEncrypted:        []byte("encrypted_updated_name"),
			AddressEncrypted:     []byte("encrypted_updated_address"),
			CityEncrypted:        []byte("encrypted_updated_city"),
			PostalCodeEncrypted:  []byte("encrypted_updated_postal"),
			CountryEncrypted:     []byte("encrypted_updated_country"),
			DescriptionEncrypted: []byte("encrypted_updated_description"),
			PhoneEncrypted:       []byte("encrypted_updated_phone"),
			EmailEncrypted:       []byte("encrypted_updated_email"),
			IsActive:             true,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			DEKEncrypted:         []byte("mock_dek_data"),
			KeyVersion:           1,
			Metadata:             encx.EncryptionMetadata{},
		}

		err := repo.Update(ctx, updatedBuilding)
		require.Error(t, err)
	})

	t.Run("should handle partial updates correctly", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert test building
		buildingEncx := tb.NewTestBuildingEncx(t)
		buildingEncx.IsActive = true
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		originalAddress := []byte("encrypted_original_address")

		updatedName := []byte("encrypted_partially_updated_name")
		updatedDescription := []byte("encrypted_updated_description")
		updatedEmail := []byte("encrypted_updated_email")

		// Update only some fields
		now := time.Now()
		updatedBuilding := &domain.BuildingEncx{
			ID:                   buildingEncx.ID,
			NameEncrypted:        updatedName,
			AddressEncrypted:     originalAddress,                      // unchanged
			CityEncrypted:        []byte("encrypted_original_city"),    // unchanged
			PostalCodeEncrypted:  []byte("encrypted_original_postal"),  // unchanged
			CountryEncrypted:     []byte("encrypted_original_country"), // unchanged
			DescriptionEncrypted: updatedDescription,
			PhoneEncrypted:       []byte("encrypted_original_phone"), // unchanged
			EmailEncrypted:       updatedEmail,
			IsActive:             true,
			CreatedAt:            now,
			UpdatedAt:            now.Add(time.Hour),
			DEKEncrypted:         []byte("mock_dek_data"),
			KeyVersion:           1,
			Metadata:             encx.EncryptionMetadata{},
		}

		// Test repository Update method
		err = repo.Update(ctx, updatedBuilding)
		require.NoError(t, err)

		// Verify specific fields were updated
		retrievedBuildingEncx, err := tb.GetBuildingEncxByID(t, ctx, testPool, buildingEncx.ID)
		require.NoError(t, err)

		assert.Equal(t, updatedName, retrievedBuildingEncx.NameEncrypted)
		assert.Equal(t, updatedDescription, retrievedBuildingEncx.DescriptionEncrypted)
		assert.Equal(t, updatedEmail, retrievedBuildingEncx.EmailEncrypted)
		assert.Equal(t, originalAddress, retrievedBuildingEncx.AddressEncrypted) // unchanged
	})
}
