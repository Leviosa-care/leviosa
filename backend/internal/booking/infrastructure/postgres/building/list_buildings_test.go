package buildingRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestListBuildings(t *testing.T) {
	ctx := context.Background()

	t.Run("should return empty list when no buildings exist", func(t *testing.T) {
		// Clean up before test
		clearBuildingsTable(t, ctx)

		// Test repository List method with empty filter
		filter := ports.BuildingFilter{}
		buildings, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Empty(t, buildings)
	})

	t.Run("should list all buildings without filters", func(t *testing.T) {
		// Clean up before test
		clearBuildingsTable(t, ctx)

		// Insert test buildings
		now := time.Now()
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
			_, err := testPool.Exec(ctx, `
				INSERT INTO booking.buildings (
					id, name_encrypted, address_encrypted, city_encrypted,
					postal_code_encrypted, country_encrypted, description_encrypted,
					phone_encrypted, email_encrypted, is_active, created_at, updated_at,
					dek_encrypted, key_version, metadata
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
			`,
				b.id,
				b.name,
				[]byte("encrypted_test_address"),
				[]byte("encrypted_test_city"),
				[]byte("encrypted_12345"),
				[]byte("encrypted_test_country"),
				[]byte("encrypted_test_description"),
				[]byte("encrypted_+1234567890"),
				[]byte("encrypted_test@example.com"),
				b.isActive,
				now,
				now,
				[]byte("mock_dek_data"),
				1,
				`{"kek_alias":"test","encryption_time":12345}`,
			)
			require.NoError(t, err)
		}

		// Test repository List method
		filter := ports.BuildingFilter{}
		result, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Len(t, result, 3)

		// Verify the results are BuildingEncx structs
		for _, building := range result {
			require.NotNil(t, building.ID)
			require.NotNil(t, building.NameEncrypted)
		}
	})

	t.Run("should filter buildings by active status", func(t *testing.T) {
		// Clean up before test
		clearBuildingsTable(t, ctx)

		// Insert test buildings
		now := time.Now()
		active := uuid.New()
		inactive := uuid.New()

		// Insert active building
		_, err := testPool.Exec(ctx, `
			INSERT INTO booking.buildings (
				id, name_encrypted, address_encrypted, city_encrypted,
				postal_code_encrypted, country_encrypted, description_encrypted,
				phone_encrypted, email_encrypted, is_active, created_at, updated_at,
				dek_encrypted, key_version, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		`,
			active,
			[]byte("encrypted_active_building"),
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

		// Insert inactive building
		_, err = testPool.Exec(ctx, `
			INSERT INTO booking.buildings (
				id, name_encrypted, address_encrypted, city_encrypted,
				postal_code_encrypted, country_encrypted, description_encrypted,
				phone_encrypted, email_encrypted, is_active, created_at, updated_at,
				dek_encrypted, key_version, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		`,
			inactive,
			[]byte("encrypted_inactive_building"),
			[]byte("encrypted_test_address"),
			[]byte("encrypted_test_city"),
			[]byte("encrypted_12345"),
			[]byte("encrypted_test_country"),
			[]byte("encrypted_test_description"),
			[]byte("encrypted_+1234567890"),
			[]byte("encrypted_test@example.com"),
			false,
			now,
			now,
			[]byte("mock_dek_data"),
			1,
			`{"kek_alias":"test","encryption_time":12345}`,
		)
		require.NoError(t, err)

		// Test filtering by active status
		activeFilter := ports.BuildingFilter{
			IsActive: &[]bool{true}[0], // Active buildings only
		}
		activeBuildings, err := repo.List(ctx, activeFilter)
		require.NoError(t, err)
		require.Len(t, activeBuildings, 1)

		inactiveFilter := ports.BuildingFilter{
			IsActive: &[]bool{false}[0], // Inactive buildings only
		}
		inactiveBuildings, err := repo.List(ctx, inactiveFilter)
		require.NoError(t, err)
		require.Len(t, inactiveBuildings, 1)
	})

	t.Run("should apply ordering correctly", func(t *testing.T) {
		// Clean up before test
		clearBuildingsTable(t, ctx)

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
			_, err := testPool.Exec(ctx, `
				INSERT INTO booking.buildings (
					id, name_encrypted, address_encrypted, city_encrypted,
					postal_code_encrypted, country_encrypted, description_encrypted,
					phone_encrypted, email_encrypted, is_active, created_at, updated_at,
					dek_encrypted, key_version, metadata
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
			`,
				b.id,
				b.name,
				[]byte("encrypted_test_address"),
				[]byte("encrypted_test_city"),
				[]byte("encrypted_12345"),
				[]byte("encrypted_test_country"),
				[]byte("encrypted_test_description"),
				[]byte("encrypted_+1234567890"),
				[]byte("encrypted_test@example.com"),
				true,
				b.time,
				b.time,
				[]byte("mock_dek_data"),
				1,
				`{"kek_alias":"test","encryption_time":12345}`,
			)
			require.NoError(t, err)
		}

		// Test ordering by created_at DESC (default)
		filter := ports.BuildingFilter{}
		result, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Len(t, result, 3)
		// Should be ordered by created_at DESC (newest first)
		require.Equal(t, "encrypted_building_c", string(result[0].NameEncrypted))
		require.Equal(t, "encrypted_building_b", string(result[1].NameEncrypted))
		require.Equal(t, "encrypted_building_a", string(result[2].NameEncrypted))
	})

	t.Run("should apply pagination correctly", func(t *testing.T) {
		// Clean up before test
		clearBuildingsTable(t, ctx)

		// Insert 5 test buildings
		now := time.Now()
		for i := 1; i <= 5; i++ {
			buildingID := uuid.New()
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
				now.Add(time.Duration(i)*time.Hour),
				now.Add(time.Duration(i)*time.Hour),
				[]byte("mock_dek_data"),
				1,
				`{"kek_alias":"test","encryption_time":12345}`,
			)
			require.NoError(t, err)
		}

		// Test pagination
		filter := ports.BuildingFilter{
			Limit:  2,
			Offset: 1,
		}
		result, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Len(t, result, 2)
	})
}

func TestCountBuildings(t *testing.T) {
	ctx := context.Background()

	t.Run("should count zero when no buildings exist", func(t *testing.T) {
		// Clean up before test
		clearBuildingsTable(t, ctx)

		// Test repository Count method
		filter := ports.BuildingFilter{}
		count, err := repo.Count(ctx, filter)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("should count all buildings without filters", func(t *testing.T) {
		// Clean up before test
		clearBuildingsTable(t, ctx)

		// Insert test buildings
		now := time.Now()
		for i := 1; i <= 3; i++ {
			buildingID := uuid.New()
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

		// Test repository Count method
		filter := ports.BuildingFilter{}
		count, err := repo.Count(ctx, filter)
		require.NoError(t, err)
		require.Equal(t, 3, count)
	})

	t.Run("should count buildings by active status", func(t *testing.T) {
		// Clean up before test
		clearBuildingsTable(t, ctx)

		// Insert test buildings
		now := time.Now()
		for i := 0; i < 5; i++ {
			buildingID := uuid.New()
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
				i < 3, // First 3 are active, last 2 are inactive
				now,
				now,
				[]byte("mock_dek_data"),
				1,
				`{"kek_alias":"test","encryption_time":12345}`,
			)
			require.NoError(t, err)
		}

		// Count active buildings
		activeFilter := ports.BuildingFilter{
			IsActive: &[]bool{true}[0],
		}
		activeCount, err := repo.Count(ctx, activeFilter)
		require.NoError(t, err)
		require.Equal(t, 3, activeCount)

		// Count inactive buildings
		inactiveFilter := ports.BuildingFilter{
			IsActive: &[]bool{false}[0],
		}
		inactiveCount, err := repo.Count(ctx, inactiveFilter)
		require.NoError(t, err)
		require.Equal(t, 2, inactiveCount)
	})
}
