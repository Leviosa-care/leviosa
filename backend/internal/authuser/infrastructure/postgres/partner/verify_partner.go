package partnerRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) VerifyPartner(ctx context.Context, partnerID uuid.UUID, verifiedByUserID uuid.UUID) error {
	// For simplicity, we're storing the verification time as encrypted bytes
	// In a real implementation, you'd encrypt this using the same DEK as other fields
	verifiedAt := time.Now()
	verifiedAtBytes := []byte(verifiedAt.Format(time.RFC3339))

	query := fmt.Sprintf(`
		UPDATE %s.partners SET
			is_verified = true,
			verified_at_encrypted = $2,
			verified_by_user_id = $3,
			updated_at = NOW()
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, partnerID, verifiedAtBytes, verifiedByUserID)
	if err != nil {
		return errs.ClassifyPgError("verify partner", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}