package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// DeletePartner deletes a partner by user ID (since partners are identified by user_id in the new domain)
func (r *Repository) DeletePartner(ctx context.Context, userID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.partners
		WHERE user_id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return errs.ClassifyPgError("delete partner", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}
