package buildingRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Create(ctx context.Context, buildingEncx *domain.BuildingEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.buildings (
			id, name_encrypted, address_encrypted, city_encrypted, city_hash,
			postal_code_encrypted, country_encrypted, country_hash,
			description_encrypted, phone_encrypted, email_encrypted,
			dek_encrypted, key_version, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
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
		buildingEncx.CreatedAt,
		buildingEncx.UpdatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError("create building", err)
	}

	return nil
}
