package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) GetAllPartners(ctx context.Context) ([]*domain.PartnerEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			user_id, bio_encrypted, experience_encrypted, certifications_encrypted,
			category_ids_encrypted, product_ids_encrypted,
			stripe_connected_account_id_encrypted, stripe_account_status, stripe_onboarding_complete,
			dek_encrypted, key_version, created_at, updated_at
		FROM %s.partners
		ORDER BY created_at DESC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errs.ClassifyPgError("get all partners", err)
	}
	defer rows.Close()

	var partners []*domain.PartnerEncx
	for rows.Next() {
		partner := &domain.PartnerEncx{}
		err := rows.Scan(
			&partner.UserID,
			&partner.BioEncrypted,
			&partner.ExperienceEncrypted,
			&partner.CertificationsEncrypted,
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
			return nil, errs.ClassifyPgError("scan partner", err)
		}
		partners = append(partners, partner)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate partners", err)
	}

	if len(partners) == 0 {
		return []*domain.PartnerEncx{}, nil
	}

	return partners, nil
}
