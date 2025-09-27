package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (r *Repository) DeletePartner(ctx context.Context, partnerID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.partners
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, partnerID)
	if err != nil {
		return errs.ClassifyPgError("delete partner", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}