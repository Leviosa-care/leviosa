package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) UpdatePartner(ctx context.Context, partner *domain.PartnerEncx) error {
	query := fmt.Sprintf(`
		UPDATE %s.partners SET
			bio_encrypted = $2,
			experience_encrypted = $3,
			certifications_encrypted = $4,
			category_ids_encrypted = $5,
			product_ids_encrypted = $6,
			stripe_connected_account_id_encrypted = $7,
			stripe_account_status = $8,
			stripe_onboarding_complete = $9,
			dek_encrypted = $10,
			key_version = $11,
			updated_at = NOW()
		WHERE user_id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		partner.UserID,
		partner.BioEncrypted,
		partner.ExperienceEncrypted,
		partner.CertificationsEncrypted,
		partner.CategoryIDsEncrypted,
		partner.ProductIDsEncrypted,
		partner.StripeConnectedAccountIDEncrypted,
		partner.StripeAccountStatus,
		partner.StripeOnboardingComplete,
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
