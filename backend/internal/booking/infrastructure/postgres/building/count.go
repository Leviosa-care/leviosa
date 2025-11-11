package buildingRepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Count(ctx context.Context, filter ports.BuildingFilter) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s.buildings", r.schema)

	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	// Apply filters
	if filter.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}

	if filter.CityHash != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("city_hash = $%d", argIndex))
		args = append(args, *filter.CityHash)
		argIndex++
	}

	if filter.CountryHash != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("country_hash = $%d", argIndex))
		args = append(args, *filter.CountryHash)
		argIndex++
	}

	// Add WHERE clause if we have conditions
	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	var count int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, errs.ClassifyPgError("count buildings", err)
	}

	return count, nil
}
