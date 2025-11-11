package buildingRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestListBuildings TEST_PATH=internal/booking/infrastructure/postgres/building/list_buildings_test.go

func TestListBuildings(t *testing.T) {
	ctx := context.Background()

	t.Run("should return empty list when no buildings exist", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Test repository List method with empty filter
		filter := ports.BuildingFilter{}
		buildings, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Empty(t, buildings)
	})

	t.Run("should list all buildings without filters", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert test buildings
		buildings := []struct {
			id       uuid.UUID
			name     []byte
			isActive bool
		}{
			{uuid.New(), []byte("encrypted_building_1"), true},
			{uuid.New(), []byte("encrypted_building_2"), false},
			{uuid.New(), []byte("encrypted_building_3"), true},
		}

		for _, b := range buildings {
			building := tb.NewTestBuildingEncx(t)
			building.ID = b.id
			building.NameEncrypted = b.name
			err := tb.InsertBuildingEncx(t, ctx, testPool, building)
			require.NoError(t, err)
		}

		// Test repository List method
		filter := ports.BuildingFilter{}
		result, err := repo.List(ctx, filter)
		assert.NoError(t, err)
		assert.Len(t, result, 3)

		// Verify the results are BuildingEncx structs
		for _, building := range result {
			assert.NotNil(t, building.ID)
			assert.NotNil(t, building.NameEncrypted)
		}
	})

	t.Run("should filter buildings by active status", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert test buildings
		active := tb.NewTestBuildingEncx(t)
		active.IsActive = true
		inactive := tb.NewTestBuildingEncx(t)
		inactive.IsActive = false

		err := tb.InsertBuildingEncx(t, ctx, testPool, active)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, inactive)
		require.NoError(t, err)

		// Test filtering by active status
		activeFilter := ports.BuildingFilter{
			IsActive: &[]bool{true}[0], // Active buildings only
		}
		activeBuildings, err := repo.List(ctx, activeFilter)
		assert.NoError(t, err)
		assert.Len(t, activeBuildings, 1)

		inactiveFilter := ports.BuildingFilter{
			IsActive: &[]bool{false}[0], // Inactive buildings only
		}
		inactiveBuildings, err := repo.List(ctx, inactiveFilter)
		assert.NoError(t, err)
		assert.Len(t, inactiveBuildings, 1)
	})

	t.Run("should apply ordering correctly", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert test buildings with different creation times
		baseTime := time.Now()
		buildings := []struct {
			id   uuid.UUID
			name []byte
			time time.Time
		}{
			{uuid.New(), []byte("encrypted_building_c"), baseTime.Add(2 * time.Hour)},
			{uuid.New(), []byte("encrypted_building_a"), baseTime},
			{uuid.New(), []byte("encrypted_building_b"), baseTime.Add(time.Hour)},
		}

		for _, b := range buildings {
			building := tb.NewTestBuildingEncx(t)
			building.ID = b.id
			building.NameEncrypted = b.name
			err := tb.InsertBuildingEncx(t, ctx, testPool, building)
			require.NoError(t, err)
		}

		// Test ordering by created_at DESC (default)
		filter := ports.BuildingFilter{}
		result, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Len(t, result, 3)
		// Should be ordered by created_at DESC (newest first)
		require.Equal(t, "encrypted_building_b", string(result[0].NameEncrypted))
		require.Equal(t, "encrypted_building_a", string(result[1].NameEncrypted))
		require.Equal(t, "encrypted_building_c", string(result[2].NameEncrypted))
	})

	t.Run("should apply pagination correctly", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert 5 test buildings
		for i := 1; i <= 5; i++ {
			buildingEncx := tb.NewTestBuildingEncx(t)
			buildingEncx.ID = uuid.New()
			buildingEncx.IsActive = true
			err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
			require.NoError(t, err)
		}

		// Test pagination
		filter := ports.BuildingFilter{
			Limit:  2,
			Offset: 1,
		}
		result, err := repo.List(ctx, filter)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
	})
}
