package partnerRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) UpdatePartner(ctx context.Context, partner *domain.PartnerEncx) error {
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
		UPDATE %s.partners SET
			bio_encrypted = $2,
			experience_encrypted = $3,
			certifications_encrypted = $4,
			category_ids = $5,
			product_ids = $6,
			dek_encrypted = $7,
			key_version = $8,
			updated_at = NOW()
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		partner.ID,
		partner.BioEncrypted,
		partner.ExperienceEncrypted,
		partner.CertificationsEncrypted,
		categoryIDsJSON,
		productIDsJSON,
		partner.DEKEncrypted,
		partner.KeyVersion,
	)
	if err != nil {
		return errs.ClassifyPgError("update partner", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}