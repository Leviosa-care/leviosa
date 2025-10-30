package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// VerifyPartner updates partner's Stripe account status to reflect verification completion
// In the new Partner domain, verification is handled through Stripe Connect account status
// rather than separate verification fields
func (r *Repository) VerifyPartner(ctx context.Context, userID uuid.UUID, verifiedByUserID uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.partners SET
			stripe_account_status = '%s',
			stripe_onboarding_complete = true,
			updated_at = NOW()
		WHERE user_id = $1
	`, r.schema, domain.StripeAccountStatusActive)

	result, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return errs.ClassifyPgError("verify partner", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}
