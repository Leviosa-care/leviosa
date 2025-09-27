package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (r *Repository) AddPartnerSpecialization(ctx context.Context, partnerID, specializationID uuid.UUID) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.partner_specializations (id, partner_id, specialization_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (partner_id, specialization_id) DO NOTHING
	`, r.schema)

	_, err := r.pool.Exec(ctx, query, uuid.New(), partnerID, specializationID)
	if err != nil {
		return errs.ClassifyPgError("add partner specialization", err)
	}

	return nil
}