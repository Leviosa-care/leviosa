package specializationRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/errs"
)

func (r *Repository) GetSpecializationByName(ctx context.Context, name string) (*domain.SpecializationEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, name_encrypted, display_name_encrypted, description_encrypted,
			is_active, dek_encrypted, key_version, created_at, updated_at
		FROM %s.specializations
		WHERE name_encrypted = $1
	`, r.schema)

	specialization := &domain.SpecializationEncx{}
	err := r.pool.QueryRow(ctx, query, name).Scan(
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
		return nil, errs.ClassifyPgError("get specialization by name", err)
	}

	return specialization, nil
}