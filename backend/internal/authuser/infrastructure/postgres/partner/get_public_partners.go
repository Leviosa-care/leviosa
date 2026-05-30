package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) GetPublicPartners(ctx context.Context) ([]*domain.PublicPartnerRow, error) {
	query := fmt.Sprintf(`
		SELECT
			p.id, p.user_id, p.bio, p.experience, p.occupation, p.quote, p.tags,
			p.category_ids, p.product_ids,
			p.stripe_connected_account_id_encrypted, p.stripe_account_status, p.stripe_onboarding_complete,
			p.dek_encrypted, p.key_version, p.created_at, p.updated_at,
			u.first_name_encrypted, u.last_name_encrypted, u.picture_encrypted,
			u.dek_encrypted, u.key_version
		FROM %s.partners p
		JOIN %s.users u ON p.user_id = u.id
		WHERE u.state = 'active'
		  AND p.stripe_account_status != 'disabled'
		ORDER BY p.created_at DESC
	`, r.schema, r.schema)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errs.ClassifyPgError("get public partners", err)
	}
	defer rows.Close()

	var result []*domain.PublicPartnerRow
	for rows.Next() {
		partner := &domain.PartnerEncx{}
		row := &domain.PublicPartnerRow{PartnerEncx: partner}

		err := rows.Scan(
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
			&row.FirstNameEncrypted,
			&row.LastNameEncrypted,
			&row.PictureEncrypted,
			&row.UserDEKEncrypted,
			&row.UserKeyVersion,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan public partner", err)
		}
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate public partners", err)
	}

	if len(result) == 0 {
		return []*domain.PublicPartnerRow{}, nil
	}

	return result, nil
}
