package buildingRepository_test

import (
	"context"
	"testing"

	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetBuildingByID TEST_PATH=internal/booking/infrastructure/postgres/building/get_building_test.go

func TestGetBuildingByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve existing building", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert test building directly into database
		buildingEncx := tb.NewTestBuildingEncx(t)
		buildingEncx.IsActive = true
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Test repository GetByID method
		retrievedBuildingEncx, err := repo.GetByID(ctx, buildingEncx.ID)
		assert.NoError(t, err)
		assert.NotNil(t, buildingEncx)

		// Verify retrieved data
		assert.Equal(t, buildingEncx.ID, retrievedBuildingEncx.ID)
		assert.Equal(t, buildingEncx.NameEncrypted, retrievedBuildingEncx.NameEncrypted)
		assert.Equal(t, buildingEncx.AddressEncrypted, retrievedBuildingEncx.AddressEncrypted)
		assert.Equal(t, buildingEncx.CityEncrypted, retrievedBuildingEncx.CityEncrypted)
		assert.Equal(t, buildingEncx.PostalCodeEncrypted, retrievedBuildingEncx.PostalCodeEncrypted)
		assert.Equal(t, buildingEncx.CountryEncrypted, retrievedBuildingEncx.CountryEncrypted)
		assert.Equal(t, buildingEncx.DescriptionEncrypted, retrievedBuildingEncx.DescriptionEncrypted)
		assert.Equal(t, buildingEncx.PhoneEncrypted, retrievedBuildingEncx.PhoneEncrypted)
		assert.Equal(t, buildingEncx.EmailEncrypted, retrievedBuildingEncx.EmailEncrypted)
		assert.True(t, retrievedBuildingEncx.IsActive)
		assert.Equal(t, 1, retrievedBuildingEncx.KeyVersion)
		assert.Equal(t, buildingEncx.DEKEncrypted, retrievedBuildingEncx.DEKEncrypted)
	})

	t.Run("should return error for non-existent building ID", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Try to get non-existent building
		nonExistentID := uuid.New()
		buildingEncx, err := repo.GetByID(ctx, nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, buildingEncx)
	})

	t.Run("should handle database connection errors gracefully", func(t *testing.T) {
		// Test with invalid UUID format
		invalidID := uuid.UUID{}
		buildingEncx, err := repo.GetByID(ctx, invalidID)
		assert.Error(t, err)
		assert.Nil(t, buildingEncx)
	})
}
