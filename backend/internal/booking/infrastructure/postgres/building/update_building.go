package buildingRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Update(ctx context.Context, building *domain.BuildingEncx) error {
	query := fmt.Sprintf(`
		UPDATE %s.buildings SET
			name_encrypted = $2,
			address_encrypted = $3,
			city_encrypted = $4,
			postal_code_encrypted = $5,
			country_encrypted = $6,
			description_encrypted = $7,
			phone_encrypted = $8,
			email_encrypted = $9,
			is_active = $10,
			updated_at = $11
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		building.ID,
		building.NameEncrypted,
		building.AddressEncrypted,
		building.CityEncrypted,
		building.PostalCodeEncrypted,
		building.CountryEncrypted,
		building.DescriptionEncrypted,
		building.PhoneEncrypted,
		building.EmailEncrypted,
		building.IsActive,
		building.UpdatedAt,
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
