package buildingRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Create(ctx context.Context, building *domain.BuildingEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.buildings (
			id, name_encrypted, address_encrypted, city_encrypted,
			postal_code_encrypted, country_encrypted, description_encrypted,
			phone_encrypted, email_encrypted, dek_encrypted, key_version,
			is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		building.ID,
		building.NameEncrypted,
		building.AddressEncrypted,
		building.CityEncrypted,
		building.PostalCodeEncrypted,
		building.CountryEncrypted,
		building.DescriptionEncrypted,
		building.PhoneEncrypted,
		building.EmailEncrypted,
		building.DEKEncrypted,
		building.KeyVersion,
		building.IsActive,
		building.CreatedAt,
		building.UpdatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError("create building", err)
	}

	return nil
}
