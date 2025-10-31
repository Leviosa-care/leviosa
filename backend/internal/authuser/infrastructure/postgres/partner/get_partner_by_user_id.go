package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*domain.PartnerEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			user_id, bio_encrypted, experience_encrypted, 
			category_ids_encrypted, product_ids_encrypted,
			stripe_connected_account_id_encrypted, stripe_account_status, stripe_onboarding_complete,
			dek_encrypted, key_version, created_at, updated_at
		FROM %s.partners
		WHERE user_id = $1
	`, r.schema)

	partner := &domain.PartnerEncx{}
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&partner.UserID,
		&partner.BioEncrypted,
		&partner.ExperienceEncrypted,
		&partner.CategoryIDsEncrypted,
		&partner.ProductIDsEncrypted,
		&partner.StripeConnectedAccountIDEncrypted,
		&partner.StripeAccountStatus,
		&partner.StripeOnboardingComplete,
		&partner.DEKEncrypted,
		&partner.KeyVersion,
		&partner.CreatedAt,
		&partner.UpdatedAt,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get partner by user ID", err)
	}

	return partner, nil
}
