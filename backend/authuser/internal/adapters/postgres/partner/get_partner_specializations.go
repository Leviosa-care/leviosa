package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetPartnerSpecializations(ctx context.Context, partnerID uuid.UUID) ([]*domain.Specialization, error) {
	query := fmt.Sprintf(`
		SELECT
			s.id, s.name_encrypted, s.display_name_encrypted, s.description_encrypted,
			s.is_active, s.dek_encrypted, s.key_version, s.created_at, s.updated_at
		FROM %s.specializations s
		INNER JOIN %s.partner_specializations ps ON s.id = ps.specialization_id
		WHERE ps.partner_id = $1 AND s.is_active = true
		ORDER BY s.display_name_encrypted ASC
	`, r.schema, r.schema)

	rows, err := r.pool.Query(ctx, query, partnerID)
	if err != nil {
		return nil, errs.ClassifyPgError("get partner specializations", err)
	}
	defer rows.Close()

	var specializations []*domain.Specialization
	for rows.Next() {
		specialization := &domain.Specialization{}
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
			return nil, errs.ClassifyPgError("scan partner specialization", err)
		}
		specializations = append(specializations, specialization)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate partner specializations", err)
	}

	return specializations, nil
}