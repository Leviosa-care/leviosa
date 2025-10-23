package specializationRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetSpecializationByID(ctx context.Context, specializationID uuid.UUID) (*domain.SpecializationEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, name_encrypted, display_name_encrypted, description_encrypted,
			is_active, dek_encrypted, key_version, created_at, updated_at
		FROM %s.specializations
		WHERE id = $1
	`, r.schema)

	specialization := &domain.SpecializationEncx{}
	err := r.pool.QueryRow(ctx, query, specializationID).Scan(
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
		return nil, errs.ClassifyPgError("get specialization by ID", err)
	}

	return specialization, nil
}