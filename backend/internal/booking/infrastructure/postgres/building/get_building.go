package buildingRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.BuildingEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, name_encrypted, name_hash, address_encrypted, address_hash,
			city_encrypted, city_hash, postal_code_encrypted,
			country_encrypted, country_hash, description_encrypted,
			phone_encrypted, email_encrypted, is_active, created_at, updated_at,
			dek_encrypted, key_version, metadata
		FROM %s.buildings
		WHERE id = $1
	`, r.schema)

	buildingEncx := &domain.BuildingEncx{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
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
		&buildingEncx.IsActive,
		&buildingEncx.CreatedAt,
		&buildingEncx.UpdatedAt,
		&buildingEncx.DEKEncrypted,
		&buildingEncx.KeyVersion,
		&buildingEncx.Metadata,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get building by id", err)
	}

	return buildingEncx, nil
}
