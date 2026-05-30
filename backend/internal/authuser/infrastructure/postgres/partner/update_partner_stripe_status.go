package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// UpdatePartnerStripeStatus updates only the stripe_account_status and stripe_onboarding_complete
// fields for a partner. This is used by the Connect webhook handler to reflect Stripe's
// account capability changes without touching encrypted fields.
func (r *Repository) UpdatePartnerStripeStatus(ctx context.Context, partnerID uuid.UUID, status domain.StripeAccountStatus) error {
	query := fmt.Sprintf(`
		UPDATE %s.partners
		SET stripe_account_status = $2,
		    stripe_onboarding_complete = TRUE,
		    updated_at = NOW()
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, partnerID, status)
	if err != nil {
		return errs.ClassifyPgError("update partner stripe status", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}
