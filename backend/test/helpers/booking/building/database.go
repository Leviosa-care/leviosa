package buildingHelpers

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearBuildingsTable removes all test data from the buildings table
func ClearBuildingsTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE booking.buildings RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

// GetBuildingEncxByID retrieves a building from the database by ID
func GetBuildingEncxByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, buildingID interface{}) (*domain.BuildingEncx, error) {
	t.Helper()

	query := `
		SELECT id, name_encrypted, address_encrypted, city_encrypted, postal_code_encrypted, country_encrypted, description_encrypted, phone_encrypted, email_encrypted, dek_encrypted, key_version, is_active, created_at, updated_at
		FROM booking.buildings
		WHERE id = $1
	`

	var buildingEncx domain.BuildingEncx
	err := pool.QueryRow(ctx, query, buildingID).Scan(
		&buildingEncx.ID,
		&buildingEncx.NameEncrypted,
		&buildingEncx.AddressEncrypted,
		&buildingEncx.CityEncrypted,
		&buildingEncx.PostalCodeEncrypted,
		&buildingEncx.CountryEncrypted,
		&buildingEncx.DescriptionEncrypted,
		&buildingEncx.PhoneEncrypted,
		&buildingEncx.EmailEncrypted,
		&buildingEncx.DEKEncrypted,
		&buildingEncx.KeyVersion,
		&buildingEncx.IsActive,
		&buildingEncx.CreatedAt,
		&buildingEncx.UpdatedAt,
	)

	return &buildingEncx, err
}
