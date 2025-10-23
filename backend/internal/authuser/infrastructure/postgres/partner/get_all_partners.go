package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) GetAllPartners(ctx context.Context) ([]*domain.PartnerEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, user_id, bio_encrypted, experience_encrypted, certifications_encrypted,
			is_verified, verified_at_encrypted, verified_by_user_id,
			dek_encrypted, key_version, created_at, updated_at
		FROM %s.partners
		ORDER BY created_at DESC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errs.ClassifyPgError("get all partners", err)
	}
	defer rows.Close()

	var partners []*domain.PartnerEncx
	for rows.Next() {
		partner := &domain.PartnerEncx{}
		err := rows.Scan(
			&partner.ID,
			&partner.UserID,
			&partner.BioEncrypted,
			&partner.ExperienceEncrypted,
			&partner.CertificationsEncrypted,
			&partner.IsVerified,
			&partner.VerifiedAtEncrypted,
			&partner.VerifiedByUserID,
			&partner.DEKEncrypted,
			&partner.KeyVersion,
			&partner.CreatedAt,
			&partner.UpdatedAt,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan partner", err)
		}
		partners = append(partners, partner)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate partners", err)
	}

	return partners, nil
}