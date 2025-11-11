package buildingRepository_test

import (
	"context"
	"testing"

	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateBuilding TEST_PATH=internal/booking/infrastructure/postgres/building/create_building_test.go

func TestCreateBuilding(t *testing.T) {
	ctx := context.Background()

	// Create test building data using encrypted struct
	buildingEncx := tb.NewTestBuildingEncx(t)

	t.Run("should successfully create a building", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Test repository Create method
		err := repo.Create(ctx, buildingEncx)
		require.NoError(t, err)

		// Verify the building was inserted by querying directly
		var count int
		err = testPool.QueryRow(ctx,
			"SELECT COUNT(*) FROM booking.buildings WHERE id = $1",
			buildingEncx.ID).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count, "Building should be inserted in database")
	})

	t.Run("should handle duplicate building ID", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert first building
		err := repo.Create(ctx, buildingEncx)
		require.NoError(t, err)

		// Try to insert building with same ID
		err = repo.Create(ctx, buildingEncx)
		assert.Error(t, err, "Should fail on duplicate ID")
	})
}
