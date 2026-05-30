package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (r *Repository) GetPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.PartnerEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, user_id, bio, experience, occupation, quote, tags,
			category_ids, product_ids,
			stripe_connected_account_id_encrypted, stripe_account_status, stripe_onboarding_complete,
			dek_encrypted, key_version, created_at, updated_at
		FROM %s.partners
		WHERE id = $1
	`, r.schema)

	partner := &domain.PartnerEncx{}
	err := r.pool.QueryRow(ctx, query, partnerID).Scan(
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
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get partner by ID", err)
	}

	return partner, nil
}
