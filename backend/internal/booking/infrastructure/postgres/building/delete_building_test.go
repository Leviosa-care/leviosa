package buildingRepository_test

import (
	"context"
	"testing"
	"time"

	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestDeleteBuilding TEST_PATH=internal/booking/infrastructure/postgres/building/delete_building_test.go

func TestDeleteBuilding(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully soft delete existing building", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert test building directly into database
		buildingEncx := tb.NewTestBuildingEncx(t)
		buildingEncx.IsActive = true
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Verify building is initially active
		var isActive bool
		err = testPool.QueryRow(ctx,
			"SELECT is_active FROM booking.buildings WHERE id = $1",
			buildingEncx.ID).Scan(&isActive)
		require.NoError(t, err)
		require.True(t, isActive)

		// Test repository Delete method (soft delete)
		err = repo.Delete(ctx, buildingEncx.ID)
		require.NoError(t, err)

		// Verify building is now inactive (soft deleted)
		err = testPool.QueryRow(ctx,
			"SELECT is_active FROM booking.buildings WHERE id = $1",
			buildingEncx.ID).Scan(&isActive)
		require.NoError(t, err)
		require.False(t, isActive)

		// Verify building still exists in database (soft delete, not hard delete)
		var count int
		err = testPool.QueryRow(ctx,
			"SELECT COUNT(*) FROM booking.buildings WHERE id = $1",
			buildingEncx.ID).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("should return error when deleting non-existent building", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Try to delete non-existent building
		nonExistentID := uuid.New()
		err := repo.Delete(ctx, nonExistentID)
		require.Error(t, err)
	})

	t.Run("should handle multiple soft deletes correctly", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert multiple test buildings
		now := time.Now()
		buildingIDs := []uuid.UUID{
			uuid.New(),
			uuid.New(),
			uuid.New(),
		}

		for i, buildingID := range buildingIDs {
			_, err := testPool.Exec(ctx, `
				INSERT INTO booking.buildings (
					id, name_encrypted, address_encrypted, city_encrypted,
					postal_code_encrypted, country_encrypted, description_encrypted,
					phone_encrypted, email_encrypted, is_active, created_at, updated_at,
					dek_encrypted, key_version, metadata
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
			`,
				buildingID,
				[]byte("encrypted_building_"+string(rune(i))),
				[]byte("encrypted_test_address"),
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
		}

		// Soft delete first two buildings
		err := repo.Delete(ctx, buildingIDs[0])
		require.NoError(t, err)

		err = repo.Delete(ctx, buildingIDs[1])
		require.NoError(t, err)

		// Verify building states
		var isActive1, isActive2, isActive3 bool
		err = testPool.QueryRow(ctx,
			"SELECT is_active FROM booking.buildings WHERE id = $1",
			buildingIDs[0]).Scan(&isActive1)
		require.NoError(t, err)
		require.False(t, isActive1)

		err = testPool.QueryRow(ctx,
			"SELECT is_active FROM booking.buildings WHERE id = $1",
			buildingIDs[1]).Scan(&isActive2)
		require.NoError(t, err)
		require.False(t, isActive2)

		err = testPool.QueryRow(ctx,
			"SELECT is_active FROM booking.buildings WHERE id = $1",
			buildingIDs[2]).Scan(&isActive3)
		require.NoError(t, err)
		require.True(t, isActive3) // Still active

		// Delete the third building
		err = repo.Delete(ctx, buildingIDs[2])
		require.NoError(t, err)

		err = testPool.QueryRow(ctx,
			"SELECT is_active FROM booking.buildings WHERE id = $1",
			buildingIDs[2]).Scan(&isActive3)
		require.NoError(t, err)
		require.False(t, isActive3) // Now inactive
	})

	t.Run("should update timestamp on soft delete", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert test building
		buildingID := uuid.New()
		originalTime := time.Now().Add(-1 * time.Hour) // 1 hour ago

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
			originalTime,
			originalTime,
			[]byte("mock_dek_data"),
			1,
			`{"kek_alias":"test","encryption_time":12345}`,
		)
		require.NoError(t, err)

		// Verify original timestamp
		var originalUpdatedAt time.Time
		err = testPool.QueryRow(ctx,
			"SELECT updated_at FROM booking.buildings WHERE id = $1",
			buildingID).Scan(&originalUpdatedAt)
		require.NoError(t, err)
		require.Equal(t, originalTime.Unix(), originalUpdatedAt.Unix())

		// Wait a moment to ensure timestamp difference
		time.Sleep(10 * time.Millisecond)

		// Soft delete the building
		err = repo.Delete(ctx, buildingID)
		require.NoError(t, err)

		// Verify timestamp was updated
		var newUpdatedAt time.Time
		err = testPool.QueryRow(ctx,
			"SELECT updated_at FROM booking.buildings WHERE id = $1",
			buildingID).Scan(&newUpdatedAt)
		require.NoError(t, err)
		require.True(t, newUpdatedAt.After(originalUpdatedAt), "Updated timestamp should be after original")
	})
}
