package webhookHandler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

const (
	// StripeSignatureHeader is the header name for the Stripe webhook signature
	StripeSignatureHeader = "Stripe-Signature"

	// MaxWebhookPayloadSize limits the webhook payload to 64KB to prevent abuse
	MaxWebhookPayloadSize = 65536
)

// mapStripeAccountStatus derives the StripeAccountStatus from Stripe capability flags.
//
// Mapping rules:
//   - charges_enabled && payouts_enabled → active
//   - charges_enabled && !payouts_enabled → restricted
//   - !charges_enabled && !payouts_enabled → disabled
//   - !charges_enabled && payouts_enabled → disabled (unusual, treat as disabled)
func mapStripeAccountStatus(chargesEnabled, payoutsEnabled bool) domain.StripeAccountStatus {
	if chargesEnabled && payoutsEnabled {
		return domain.StripeAccountStatusActive
	}
	if chargesEnabled && !payoutsEnabled {
		return domain.StripeAccountStatusRestricted
	}
	return domain.StripeAccountStatusDisabled
}

// HandleStripeConnectWebhook handles incoming Stripe Connect webhook events
// for account status updates.
//
// This endpoint:
//  1. Reads and validates the raw request body
//  2. Verifies the Stripe webhook signature using the Connect-specific secret
//  3. Processes account.updated events to update partner status
//  4. Returns appropriate HTTP status codes for Stripe retry logic
//
// Security:
//   - No authentication required (webhook is signed by Stripe)
//   - Signature verification prevents spoofed requests
//   - Payload size limited to prevent abuse
//
// Stripe retry behavior:
//   - 2xx: Event processed successfully, no retry
//   - 4xx: Event invalid, no retry (except 408, 429)
//   - 5xx: Temporary failure, Stripe will retry
func (h *handler) HandleStripeConnectWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Read the raw body for signature verification
	body, err := io.ReadAll(io.LimitReader(r.Body, MaxWebhookPayloadSize))
	if err != nil {
		slog.ErrorContext(ctx, "failed to read connect webhook payload",
			"error", err,
		)
		httpx.RespondWithError(w, errors.New("failed to read webhook payload"), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Get the Stripe signature header
	signature := r.Header.Get(StripeSignatureHeader)
	if signature == "" {
		slog.ErrorContext(ctx, "missing Stripe signature header")
		httpx.RespondWithError(w, errors.New("missing Stripe signature header"), http.StatusBadRequest)
		return
	}

	// Verify signature and parse event
	accountID, chargesEnabled, payoutsEnabled, err := h.stripe.VerifyConnectWebhookSignature(body, signature)
	if err != nil {
		slog.ErrorContext(ctx, "connect webhook signature verification failed",
			"error", err,
		)
		httpx.RespondWithError(w, errors.New("webhook signature verification failed"), http.StatusBadRequest)
		return
	}

	// Empty accountID means unknown event type — acknowledge but don't process
	if accountID == "" {
		slog.InfoContext(ctx, "ignoring unknown Stripe Connect webhook event type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Map Stripe capabilities to domain status
	status := mapStripeAccountStatus(chargesEnabled, payoutsEnabled)

	slog.InfoContext(ctx, "processing Stripe Connect account.updated",
		"stripe_account_id", accountID,
		"charges_enabled", chargesEnabled,
		"payouts_enabled", payoutsEnabled,
		"mapped_status", status,
	)

	// Update the partner's stripe status
	_, err = h.svc.UpdateStripeAccountStatus(ctx, accountID, status)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			// No partner found for this Stripe account — might be a test event
			slog.WarnContext(ctx, "no partner found for Stripe Connect account",
				"stripe_account_id", accountID,
			)
			w.WriteHeader(http.StatusOK)
			return
		}
		// Transient error — return 500 so Stripe retries
		slog.ErrorContext(ctx, "failed to update partner stripe status",
			"stripe_account_id", accountID,
			"error", err,
		)
		httpx.RespondWithError(w, errors.New("internal error processing webhook"), http.StatusInternalServerError)
		return
	}

	// Success
	w.WriteHeader(http.StatusOK)
}
