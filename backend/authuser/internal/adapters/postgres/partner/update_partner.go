package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/errs"
)

func (r *Repository) UpdatePartner(ctx context.Context, partner *domain.Partner) error {
	query := fmt.Sprintf(`
		UPDATE %s.partners SET
			bio_encrypted = $2,
			experience_encrypted = $3,
			certifications_encrypted = $4,
			dek_encrypted = $5,
			key_version = $6,
			updated_at = NOW()
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		partner.ID,
		partner.BioEncrypted,
		partner.ExperienceEncrypted,
		partner.CertificationsEncrypted,
		partner.DEKEncrypted,
		partner.KeyVersion,
	)
	if err != nil {
		return errs.ClassifyPgError("update partner", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}