package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) RemovePartnerSpecialization(ctx context.Context, partnerID, specializationID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.partner_specializations
		WHERE partner_id = $1 AND specialization_id = $2
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, partnerID, specializationID)
	if err != nil {
		return errs.ClassifyPgError("remove partner specialization", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}