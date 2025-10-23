package specializationRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) DeleteSpecialization(ctx context.Context, specializationID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.specializations
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, specializationID)
	if err != nil {
		return errs.ClassifyPgError("delete specialization", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}