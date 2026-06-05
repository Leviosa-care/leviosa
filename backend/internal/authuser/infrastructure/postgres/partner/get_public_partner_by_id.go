package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetPublicPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.PublicPartnerRow, error) {
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
		WHERE p.id = $1
		  AND u.state = 'active'
		  AND p.stripe_account_status != 'disabled'
	`, r.schema, r.schema)

	partner := &domain.PartnerEncx{}
	row := &domain.PublicPartnerRow{PartnerEncx: partner}

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
		&row.FirstNameEncrypted,
		&row.LastNameEncrypted,
		&row.PictureEncrypted,
		&row.UserDEKEncrypted,
		&row.UserKeyVersion,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get public partner by ID", err)
	}

	return row, nil
}
