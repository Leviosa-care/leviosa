package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) CreatePartner(ctx context.Context, partner *domain.PartnerEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.partners (
			user_id, bio_encrypted, experience_encrypted,
			category_ids_encrypted, product_ids_encrypted,
			stripe_connected_account_id_encrypted, stripe_account_status, stripe_onboarding_complete,
			dek_encrypted, key_version, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
	`, r.schema)

	if _, err := r.pool.Exec(ctx, query,
		partner.UserID,
		partner.BioEncrypted,
		partner.ExperienceEncrypted,
		partner.CategoryIDsEncrypted,
		partner.ProductIDsEncrypted,
		partner.StripeConnectedAccountIDEncrypted,
		partner.StripeAccountStatus,
		partner.StripeOnboardingComplete,
		partner.DEKEncrypted,
		partner.KeyVersion,
		partner.CreatedAt,
		partner.UpdatedAt,
	); err != nil {
		return errs.ClassifyPgError("create partner", err)
	}

	return nil
}
