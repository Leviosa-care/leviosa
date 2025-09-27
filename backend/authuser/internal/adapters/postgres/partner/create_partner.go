package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/errs"
)

func (r *Repository) CreatePartner(ctx context.Context, partner *domain.Partner) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.partners (
			id, user_id, bio_encrypted, experience_encrypted, certifications_encrypted,
			is_verified, dek_encrypted, key_version
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		partner.ID,
		partner.UserID,
		partner.BioEncrypted,
		partner.ExperienceEncrypted,
		partner.CertificationsEncrypted,
		partner.IsVerified,
		partner.DEKEncrypted,
		partner.KeyVersion,
	)
	if err != nil {
		return errs.ClassifyPgError("create partner", err)
	}

	return nil
}