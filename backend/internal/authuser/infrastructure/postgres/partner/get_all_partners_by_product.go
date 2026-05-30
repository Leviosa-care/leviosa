package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// GetAllPartnersByProduct retrieves all partners that offer a specific product.
// It searches for partners whose product_ids array contains the given productID.
func (r *Repository) GetAllPartnersByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.PartnerEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, user_id, bio, experience, occupation, quote, tags,
			category_ids, product_ids,
			stripe_connected_account_id_encrypted, stripe_account_status, stripe_onboarding_complete,
			dek_encrypted, key_version, created_at, updated_at
		FROM %s.partners
		WHERE $1 = ANY(product_ids)
		ORDER BY created_at DESC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, productID)
	if err != nil {
		return nil, errs.ClassifyPgError("get partners by product", err)
	}
	defer rows.Close()

	var partners []*domain.PartnerEncx
	for rows.Next() {
		partner := &domain.PartnerEncx{}
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
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan partner by product", err)
		}
		partners = append(partners, partner)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate partners by product", err)
	}

	if len(partners) == 0 {
		return []*domain.PartnerEncx{}, nil
	}

	return partners, nil
}
