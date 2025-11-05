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
			bio = $2,
			experience = $3,
			category_ids = $4,
			product_ids = $5,
			stripe_connected_account_id_encrypted = $6,
			stripe_account_status = $7,
			stripe_onboarding_complete = $8,
			dek_encrypted = $9,
			key_version = $10,
			updated_at = NOW()
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		partner.ID,
		partner.Bio,
		partner.Experience,
		partner.CategoryIDs,
		partner.ProductIDs,
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
