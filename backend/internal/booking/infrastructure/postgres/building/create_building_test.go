package buildingRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/require"
)

func TestCreateBuilding(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully create a building", func(t *testing.T) {
		// Clean up before test
		clearBuildingsTable(t, ctx)

		// Create test building data using encrypted struct
		buildingID := uuid.New()
		now := time.Now()

		buildingEncx := &domain.BuildingEncx{
			ID:                   buildingID,
			NameEncrypted:        []byte("encrypted_test_building_name"),
			AddressEncrypted:     []byte("encrypted_123_test_street"),
			CityEncrypted:        []byte("encrypted_test_city"),
			PostalCodeEncrypted:  []byte("encrypted_12345"),
			CountryEncrypted:     []byte("encrypted_test_country"),
			DescriptionEncrypted: []byte("encrypted_test_description"),
			PhoneEncrypted:       []byte("encrypted_+1234567890"),
			EmailEncrypted:       []byte("encrypted_test@example.com"),
			IsActive:             true,
			CreatedAt:            now,
			UpdatedAt:            now,
			DEKEncrypted:         []byte("mock_dek_data"),
			KeyVersion:           1,
			Metadata:             encx.EncryptionMetadata{},
		}

		// Test repository Create method
		err := repo.Create(ctx, buildingEncx)
		require.NoError(t, err)

		// Verify the building was inserted by querying directly
		var count int
		err = testPool.QueryRow(ctx,
			"SELECT COUNT(*) FROM booking.buildings WHERE id = $1",
			buildingID).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count, "Building should be inserted in database")
	})

	t.Run("should handle duplicate building ID", func(t *testing.T) {
		// Clean up before test
		clearBuildingsTable(t, ctx)

		// Create test building data
		buildingID := uuid.New()
		now := time.Now()

		buildingEncx := &domain.BuildingEncx{
			ID:                   buildingID,
			NameEncrypted:        []byte("encrypted_test_building_name"),
			AddressEncrypted:     []byte("encrypted_123_test_street"),
			CityEncrypted:        []byte("encrypted_test_city"),
			PostalCodeEncrypted:  []byte("encrypted_12345"),
			CountryEncrypted:     []byte("encrypted_test_country"),
			DescriptionEncrypted: []byte("encrypted_test_description"),
			PhoneEncrypted:       []byte("encrypted_+1234567890"),
			EmailEncrypted:       []byte("encrypted_test@example.com"),
			IsActive:             true,
			CreatedAt:            now,
			UpdatedAt:            now,
			DEKEncrypted:         []byte("mock_dek_data"),
			KeyVersion:           1,
			Metadata:             encx.EncryptionMetadata{},
		}

		// Insert first building
		err := repo.Create(ctx, buildingEncx)
		require.NoError(t, err)

		// Try to insert building with same ID
		err = repo.Create(ctx, buildingEncx)
		require.Error(t, err, "Should fail on duplicate ID")
	})
}

// clearBuildingsTable removes all test data from the buildings table
func clearBuildingsTable(t *testing.T, ctx context.Context) {
	t.Helper()
	_, err := testPool.Exec(ctx, "DELETE FROM booking.buildings")
	require.NoError(t, err)
}
