package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.PartnerEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, user_id, bio_encrypted, experience_encrypted, certifications_encrypted,
			is_verified, verified_at_encrypted, verified_by_user_id,
			dek_encrypted, key_version, created_at, updated_at
		FROM %s.partners
		WHERE id = $1
	`, r.schema)

	partner := &domain.PartnerEncx{}
	err := r.pool.QueryRow(ctx, query, partnerID).Scan(
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
		return nil, errs.ClassifyPgError("get partner by ID", err)
	}

	return partner, nil
}