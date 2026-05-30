package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// GetAllPartnersWithStripeAccount returns all partners that have a Stripe connected account
// (i.e. stripe_connected_account_id_encrypted IS NOT NULL).
// Used by the Connect webhook handler to match an incoming account.updated event to a partner.
func (r *Repository) GetAllPartnersWithStripeAccount(ctx context.Context) ([]*domain.PartnerEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, user_id, bio, experience, occupation, quote, tags,
			category_ids, product_ids,
			stripe_connected_account_id_encrypted, stripe_account_status, stripe_onboarding_complete,
			dek_encrypted, key_version, created_at, updated_at
		FROM %s.partners
		WHERE stripe_connected_account_id_encrypted IS NOT NULL
	`, r.schema)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errs.ClassifyPgError("get all partners with stripe account", err)
	}
	defer rows.Close()

	var partners []*domain.PartnerEncx
	for rows.Next() {
		partner := &domain.PartnerEncx{}
		if err := rows.Scan(
			&partner.ID,
			&partner.UserID,
			&partner.Bio,
			&partner.Experience,
			&partner.Occupation,
			&partner.Quote,
			&partner.Tags,
			&partner.CategoryIDs,
			&partner.ProductIDs,
			&partner.StripeConnectedAccountIDEncrypted,
			&partner.StripeAccountStatus,
			&partner.StripeOnboardingComplete,
			&partner.DEKEncrypted,
			&partner.KeyVersion,
			&partner.CreatedAt,
			&partner.UpdatedAt,
		); err != nil {
			return nil, errs.ClassifyPgError("scan partner with stripe account", err)
		}
		partners = append(partners, partner)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate partners with stripe account", err)
	}

	return partners, nil
}
