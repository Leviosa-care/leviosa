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

// make test-func TEST_NAME=TestCountBuildings TEST_PATH=internal/booking/infrastructure/postgres/building/count_test.go

func TestCountBuildings(t *testing.T) {
	ctx := context.Background()

	t.Run("should count zero when no buildings exist", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Test repository Count method
		filter := ports.BuildingFilter{}
		count, err := repo.Count(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("should count all buildings without filters", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

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
		assert.NoError(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("should count buildings by active status", func(t *testing.T) {
		// Clean up before test
		tb.ClearBuildingsTable(t, ctx, testPool)

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
		assert.NoError(t, err)
		assert.Equal(t, 3, activeCount)

		// Count inactive buildings
		inactiveFilter := ports.BuildingFilter{
			IsActive: &[]bool{false}[0],
		}
		inactiveCount, err := repo.Count(ctx, inactiveFilter)
		assert.NoError(t, err)
		assert.Equal(t, 2, inactiveCount)
	})
}
