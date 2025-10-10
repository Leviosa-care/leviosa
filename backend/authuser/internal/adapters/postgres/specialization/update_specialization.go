package specializationRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/errs"
)

func (r *Repository) UpdateSpecialization(ctx context.Context, specialization *domain.SpecializationEncx) error {
	query := fmt.Sprintf(`
		UPDATE %s.specializations SET
			display_name_encrypted = $2,
			description_encrypted = $3,
			is_active = $4,
			dek_encrypted = $5,
			key_version = $6,
			updated_at = NOW()
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		specialization.ID,
		specialization.DisplayNameEncrypted,
		specialization.DescriptionEncrypted,
		specialization.IsActive,
		specialization.DEKEncrypted,
		specialization.KeyVersion,
	)
	if err != nil {
		return errs.ClassifyPgError("update specialization", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}