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
		SELECT id, name_encrypted, name_hash, address_encrypted, address_hash,
		       city_encrypted, city_hash, postal_code_encrypted,
		       country_encrypted, country_hash, description_encrypted,
		       phone_encrypted, email_encrypted,
		       dek_encrypted, key_version, is_active, created_at, updated_at
		FROM booking.buildings
		WHERE id = $1
	`

	var buildingEncx domain.BuildingEncx
	err := pool.QueryRow(ctx, query, buildingID).Scan(
		&buildingEncx.ID,
		&buildingEncx.NameEncrypted,
		&buildingEncx.NameHash,
		&buildingEncx.AddressEncrypted,
		&buildingEncx.AddressHash,
		&buildingEncx.CityEncrypted,
		&buildingEncx.CityHash,
		&buildingEncx.PostalCodeEncrypted,
		&buildingEncx.CountryEncrypted,
		&buildingEncx.CountryHash,
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

func InsertBuildingEncx(t *testing.T, ctx context.Context, pool *pgxpool.Pool, buildingEncx *domain.BuildingEncx) error {
	t.Helper()
	_, err := pool.Exec(ctx, `
			INSERT INTO booking.buildings (
				id, name_encrypted, name_hash, address_encrypted, address_hash,
				city_encrypted, city_hash, postal_code_encrypted,
				country_encrypted, country_hash, description_encrypted,
				phone_encrypted, email_encrypted, is_active, created_at, updated_at,
				dek_encrypted, key_version, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		`,
		buildingEncx.ID,
		buildingEncx.NameEncrypted,
		buildingEncx.NameHash,
		buildingEncx.AddressEncrypted,
		buildingEncx.AddressHash,
		buildingEncx.CityEncrypted,
		buildingEncx.CityHash,
		buildingEncx.PostalCodeEncrypted,
		buildingEncx.CountryEncrypted,
		buildingEncx.CountryHash,
		buildingEncx.DescriptionEncrypted,
		buildingEncx.PhoneEncrypted,
		buildingEncx.EmailEncrypted,
		buildingEncx.IsActive,
		buildingEncx.CreatedAt,
		buildingEncx.UpdatedAt,
		buildingEncx.DEKEncrypted,
		buildingEncx.KeyVersion,
		buildingEncx.Metadata,
	)
	return err
}
