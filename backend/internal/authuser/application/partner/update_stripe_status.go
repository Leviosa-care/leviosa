package partner

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// UpdateStripeAccountStatus finds the partner whose StripeConnectedAccountID matches
// the given accountID and updates their stripe_account_status.
// It loads all partners with a Stripe account, decrypts the IDs, matches, then persists.
// Returns the partner ID of the updated partner.
func (s *PartnerService) UpdateStripeAccountStatus(ctx context.Context, stripeAccountID string, status domain.StripeAccountStatus) (uuid.UUID, error) {
	partners, err := s.partnerRepo.GetAllPartnersWithStripeAccount(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("get partners with stripe account: %w", err)
	}

	for _, partnerEncx := range partners {
		partner, err := domain.DecryptPartnerEncx(ctx, s.crypto, partnerEncx)
		if err != nil {
			slog.ErrorContext(ctx, "failed to decrypt partner during stripe status update",
				"partner_id", partnerEncx.ID,
				"error", err,
			)
			continue
		}

		if partner.StripeConnectedAccountID == stripeAccountID {
			if err := s.partnerRepo.UpdatePartnerStripeStatus(ctx, partner.ID, status); err != nil {
				return uuid.Nil, fmt.Errorf("update partner stripe status: %w", err)
			}
			slog.InfoContext(ctx, "updated partner stripe account status",
				"partner_id", partner.ID,
				"stripe_account_id", stripeAccountID,
				"new_status", status,
			)
			return partner.ID, nil
		}
	}

	return uuid.Nil, errs.ErrRepositoryNotFound
}
