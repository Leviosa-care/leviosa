package authPayment

import (
	"encoding/json"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

// VerifyConnectWebhookSignature verifies the Stripe webhook signature using the
// Connect-specific webhook secret and parses the account.updated event.
// Returns the Stripe account ID and the charges/payouts enabled flags.
func (s *service) VerifyConnectWebhookSignature(payload []byte, signature string) (accountID string, chargesEnabled bool, payoutsEnabled bool, err error) {
	event, err := webhook.ConstructEvent(payload, signature, s.connectWebhookSecret)
	if err != nil {
		return "", false, false, errs.NewInvalidValueErr(fmt.Sprintf("connect webhook signature verification failed: %s", err.Error()))
	}

	// We only process account.updated events
	switch event.Type {
	case "account.updated":
		var account stripe.Account
		if err := json.Unmarshal(event.Data.Raw, &account); err != nil {
			return "", false, false, errs.NewInvalidValueErr(fmt.Sprintf("failed to parse account from webhook: %s", err.Error()))
		}

		return account.ID, account.ChargesEnabled, account.PayoutsEnabled, nil

	default:
		// Unknown event type — caller should return 200 without processing
		return "", false, false, nil
	}
}
