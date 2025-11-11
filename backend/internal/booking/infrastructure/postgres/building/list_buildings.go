package buildingRepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) List(ctx context.Context, filter ports.BuildingFilter) ([]*domain.BuildingEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, name_encrypted, address_encrypted, city_encrypted, city_hash,
			postal_code_encrypted, country_encrypted, country_hash,
			description_encrypted, phone_encrypted, email_encrypted,
			is_active, created_at, updated_at,
			dek_encrypted, key_version, metadata
		FROM %s.buildings
	`, r.schema)

	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	// Apply filters
	if filter.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}

	// Apply city filter using hash
	if filter.CityHash != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("city_hash = $%d", argIndex))
		args = append(args, *filter.CityHash)
		argIndex++
	}

	// Apply country filter using hash
	if filter.CountryHash != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("country_hash = $%d", argIndex))
		args = append(args, *filter.CountryHash)
		argIndex++
	}

	// Add WHERE clause if we have conditions
	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Add ordering
	orderBy := "created_at"
	if filter.OrderBy != "" {
		switch filter.OrderBy {
		case "name", "created_at", "city":
			orderBy = filter.OrderBy
		}
	}

	orderDirection := "DESC"
	if filter.OrderDirection == "asc" {
		orderDirection = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", orderBy, orderDirection)

	// Add pagination
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errs.ClassifyPgError("list buildings", err)
	}
	defer rows.Close()

	var buildings []*domain.BuildingEncx
	for rows.Next() {
		buildingEncx := &domain.BuildingEncx{}
		err := rows.Scan(
			&buildingEncx.ID,
			&buildingEncx.NameEncrypted,
			&buildingEncx.AddressEncrypted,
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
			return nil, errs.ClassifyPgError("scan building row", err)
		}

		buildings = append(buildings, buildingEncx)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate building rows", err)
	}

	return buildings, nil
}
