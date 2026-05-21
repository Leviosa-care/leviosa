package bookingHandler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

const (
	// StripeSignatureHeader is the header name for the Stripe webhook signature
	StripeSignatureHeader = "Stripe-Signature"

	// MaxWebhookPayloadSize limits the webhook payload to 64KB to prevent abuse
	MaxWebhookPayloadSize = 65536
)

// HandleStripeWebhook handles incoming Stripe webhook events for payment status updates.
//
// This endpoint:
// 1. Reads and validates the raw request body
// 2. Verifies the Stripe webhook signature
// 3. Processes the event through the booking service
// 4. Returns appropriate HTTP status codes for Stripe retry logic
//
// Security:
// - No authentication required (webhook is signed by Stripe)
// - Signature verification prevents spoofed requests
// - Payload size limited to prevent abuse
//
// Stripe retry behavior:
// - 2xx: Event processed successfully, no retry
// - 4xx: Event invalid, no retry (except 408, 429)
// - 5xx: Temporary failure, Stripe will retry
func (h *handler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Read the raw body for signature verification
	// Limit payload size to prevent abuse
	body, err := io.ReadAll(io.LimitReader(r.Body, MaxWebhookPayloadSize))
	if err != nil {
		slog.ErrorContext(ctx, "failed to read webhook payload",
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
	event, err := h.paymentService.VerifyWebhookSignature(body, signature)
	if err != nil {
		// Signature verification failed - could be spoofed request or configuration issue
		slog.ErrorContext(ctx, "webhook signature verification failed",
			"error", err,
		)
		httpx.RespondWithError(w, errors.New("webhook signature verification failed"), http.StatusBadRequest)
		return
	}

	// Process the webhook event
	if err := h.svc.HandlePaymentWebhook(ctx, event); err != nil {
		// Determine appropriate response based on error type
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			// Booking not found - don't retry, it won't be found later either
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrInvalidValue):
			// Invalid data in webhook - don't retry
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		default:
			// Transient error - return 500 so Stripe retries
			httpx.RespondWithError(w, errors.New("internal error processing webhook"), http.StatusInternalServerError)
		}
		return
	}

	// Success - acknowledge receipt
	w.WriteHeader(http.StatusOK)
}
