package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// GetAllPartnersByCategory retrieves all partners that offer services for a specific category.
// It searches for partners whose category_ids array contains the given categoryID.
func (r *Repository) GetAllPartnersByCategory(ctx context.Context, categoryID uuid.UUID) ([]*domain.PartnerEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, user_id, bio, experience,
			category_ids, product_ids,
			stripe_connected_account_id_encrypted, stripe_account_status, stripe_onboarding_complete,
			dek_encrypted, key_version, created_at, updated_at
		FROM %s.partners
		WHERE $1 = ANY(category_ids)
		ORDER BY created_at DESC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, categoryID)
	if err != nil {
		return nil, errs.ClassifyPgError("get partners by category", err)
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
			return nil, errs.ClassifyPgError("scan partner by category", err)
		}
		partners = append(partners, partner)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate partners by category", err)
	}

	if len(partners) == 0 {
		return []*domain.PartnerEncx{}, nil
	}

	return partners, nil
}
