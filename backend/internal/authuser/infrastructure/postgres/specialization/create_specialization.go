package specializationRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) CreateSpecialization(ctx context.Context, specialization *domain.SpecializationEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.specializations (
			id, name_encrypted, display_name_encrypted, description_encrypted,
			is_active, dek_encrypted, key_version
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		specialization.ID,
		specialization.NameEncrypted,
		specialization.DisplayNameEncrypted,
		specialization.DescriptionEncrypted,
		specialization.IsActive,
		specialization.DEKEncrypted,
		specialization.KeyVersion,
	)
	if err != nil {
		return errs.ClassifyPgError("create specialization", err)
	}

	return nil
}

