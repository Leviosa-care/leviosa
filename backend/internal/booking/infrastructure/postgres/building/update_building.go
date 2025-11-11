package buildingRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Update(ctx context.Context, buildingEncx *domain.BuildingEncx) error {
	query := fmt.Sprintf(`
		UPDATE %s.buildings SET
			name_encrypted = $2,
			name_hash = $3,
			address_encrypted = $4,
			address_hash = $5,
			city_encrypted = $6,
			city_hash = $7,
			postal_code_encrypted = $8,
			country_encrypted = $9,
			country_hash = $10,
			description_encrypted = $11,
			phone_encrypted = $12,
			email_encrypted = $13,
			dek_encrypted = $14,
			key_version = $15,
			is_active = $16,
			updated_at = $17
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
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
		buildingEncx.DEKEncrypted,
		buildingEncx.KeyVersion,
		buildingEncx.IsActive,
		buildingEncx.UpdatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError("update building", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}
