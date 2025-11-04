package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// GetAllPartnersByCategories retrieves all partners that offer services for any of the specified categories.
// It searches for partners whose category_ids array overlaps with the given categoryIDs slice.
// Uses the PostgreSQL && operator to check for array overlap.
func (r *Repository) GetAllPartnersByCategories(ctx context.Context, categoryIDs []uuid.UUID) ([]*domain.PartnerEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			user_id, bio, experience,
			category_ids, product_ids,
			created_at, updated_at
		FROM %s.partners
		WHERE category_ids && $1::uuid[]
		ORDER BY created_at DESC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, categoryIDs)
	if err != nil {
		return nil, errs.ClassifyPgError("get partners by categories", err)
	}
	defer rows.Close()

	var partners []*domain.PartnerEncx
	for rows.Next() {
		partner := &domain.PartnerEncx{}
		err := rows.Scan(
			&partner.UserID,
			&partner.Bio,
			&partner.Experience,
			&partner.CategoryIDs,
			&partner.ProductIDs,
			&partner.CreatedAt,
			&partner.UpdatedAt,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan partner by categories", err)
		}
		partners = append(partners, partner)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate partners by categories", err)
	}

	if len(partners) == 0 {
		return []*domain.PartnerEncx{}, nil
	}

	return partners, nil
}
