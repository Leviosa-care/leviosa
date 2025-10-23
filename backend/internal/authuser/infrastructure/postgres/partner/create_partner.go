package partnerRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) CreatePartner(ctx context.Context, partner *domain.PartnerEncx) error {
	// Marshal category and product IDs to JSONB
	categoryIDsJSON, err := json.Marshal(partner.CategoryIDs)
	if err != nil {
		return fmt.Errorf("marshal category IDs: %w", err)
	}

	productIDsJSON, err := json.Marshal(partner.ProductIDs)
	if err != nil {
		return fmt.Errorf("marshal product IDs: %w", err)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s.partners (
			id, user_id, bio_encrypted, experience_encrypted, certifications_encrypted,
			category_ids, product_ids, is_verified, dek_encrypted, key_version
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
	`, r.schema)

	if _, err := r.pool.Exec(ctx, query,
		partner.ID,
		partner.UserID,
		partner.BioEncrypted,
		partner.ExperienceEncrypted,
		partner.CertificationsEncrypted,
		categoryIDsJSON,
		productIDsJSON,
		partner.IsVerified,
		partner.DEKEncrypted,
		partner.KeyVersion,
	); err != nil {
		return errs.ClassifyPgError("create partner", err)
	}

	return nil
}
