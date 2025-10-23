package partnerRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.PartnerEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, user_id, bio_encrypted, experience_encrypted, certifications_encrypted,
			category_ids, product_ids, is_verified, verified_at_encrypted, verified_by_user_id,
			dek_encrypted, key_version, created_at, updated_at
		FROM %s.partners
		WHERE id = $1
	`, r.schema)

	partner := &domain.PartnerEncx{}
	var categoryIDsJSON, productIDsJSON []byte

	err := r.pool.QueryRow(ctx, query, partnerID).Scan(
		&partner.ID,
		&partner.UserID,
		&partner.BioEncrypted,
		&partner.ExperienceEncrypted,
		&partner.CertificationsEncrypted,
		&categoryIDsJSON,
		&productIDsJSON,
		&partner.IsVerified,
		&partner.VerifiedAtEncrypted,
		&partner.VerifiedByUserID,
		&partner.DEKEncrypted,
		&partner.KeyVersion,
		&partner.CreatedAt,
		&partner.UpdatedAt,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get partner by ID", err)
	}

	// Unmarshal JSONB arrays
	if len(categoryIDsJSON) > 0 {
		if err := json.Unmarshal(categoryIDsJSON, &partner.CategoryIDs); err != nil {
			return nil, fmt.Errorf("unmarshal category IDs: %w", err)
		}
	}

	if len(productIDsJSON) > 0 {
		if err := json.Unmarshal(productIDsJSON, &partner.ProductIDs); err != nil {
			return nil, fmt.Errorf("unmarshal product IDs: %w", err)
		}
	}

	return partner, nil
}