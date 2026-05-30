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
			occupation = $4,
			quote = $5,
			tags = $6,
			category_ids = $7,
			product_ids = $8,
			stripe_connected_account_id_encrypted = $9,
			stripe_account_status = $10,
			stripe_onboarding_complete = $11,
			dek_encrypted = $12,
			key_version = $13,
			updated_at = NOW()
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		partner.ID,
		partner.Bio,
		partner.Experience,
		partner.Occupation,
		partner.Quote,
		partner.Tags,
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
