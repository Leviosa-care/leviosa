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
			address_encrypted = $3,
			city_encrypted = $4,
			city_hash = $5,
			postal_code_encrypted = $6,
			country_encrypted = $7,
			country_hash = $8,
			description_encrypted = $9,
			phone_encrypted = $10,
			email_encrypted = $11,
			dek_encrypted = $12,
			key_version = $13,
			is_active = $14,
			updated_at = $15
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		buildingEncx.ID,
		buildingEncx.NameEncrypted,
		buildingEncx.AddressEncrypted,
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
