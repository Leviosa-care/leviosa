package specializationRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) GetAllSpecializations(ctx context.Context) ([]*domain.SpecializationEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, name_encrypted, display_name_encrypted, description_encrypted,
			is_active, dek_encrypted, key_version, created_at, updated_at
		FROM %s.specializations
		ORDER BY created_at ASC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errs.ClassifyPgError("get all specializations", err)
	}
	defer rows.Close()

	var specializations []*domain.SpecializationEncx
	for rows.Next() {
		specialization := &domain.SpecializationEncx{}
		err := rows.Scan(
			&specialization.ID,
			&specialization.NameEncrypted,
			&specialization.DisplayNameEncrypted,
			&specialization.DescriptionEncrypted,
			&specialization.IsActive,
			&specialization.DEKEncrypted,
			&specialization.KeyVersion,
			&specialization.CreatedAt,
			&specialization.UpdatedAt,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan specialization", err)
		}
		specializations = append(specializations, specialization)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate specializations", err)
	}

	return specializations, nil
}