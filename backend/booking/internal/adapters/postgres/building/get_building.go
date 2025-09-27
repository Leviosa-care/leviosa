package buildingRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Building, error) {
	query := fmt.Sprintf(`
		SELECT
			id, name_encrypted, address_encrypted, city_encrypted,
			postal_code_encrypted, country_encrypted, description_encrypted,
			phone_encrypted, email_encrypted, is_active, created_at, updated_at
		FROM %s.buildings
		WHERE id = $1
	`, r.schema)

	building := &domain.Building{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&building.ID,
		&building.NameEncrypted,
		&building.AddressEncrypted,
		&building.CityEncrypted,
		&building.PostalCodeEncrypted,
		&building.CountryEncrypted,
		&building.DescriptionEncrypted,
		&building.PhoneEncrypted,
		&building.EmailEncrypted,
		&building.IsActive,
		&building.CreatedAt,
		&building.UpdatedAt,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get building by id", err)
	}

	// Decrypt sensitive fields
	if err := r.crypto.DecryptStruct(ctx, building); err != nil {
		return nil, fmt.Errorf("decrypt building data: %w", err)
	}

	return building, nil
}