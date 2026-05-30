package partner_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	webhookEndpoints "github.com/Leviosa-care/leviosa/backend/internal/authuser/interface/webhook"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	stripeWebhook "github.com/stripe/stripe-go/v82/webhook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// signStripePayload generates a valid Stripe-Signature header for the given payload and secret.
func signStripePayload(payload []byte, secret string) string {
	t := time.Now()
	signed := stripeWebhook.GenerateTestSignedPayload(&stripeWebhook.UnsignedPayload{
		Payload:   payload,
		Secret:    secret,
		Timestamp: t,
	})
	return signed.Header
}

// buildAccountUpdatedPayload constructs a minimal account.updated webhook event JSON.
func buildAccountUpdatedPayload(accountID string, chargesEnabled, payoutsEnabled bool) []byte {
	payload := fmt.Sprintf(`{
		"id": "evt_test_%s",
		"object": "event",
		"api_version": "%s",
		"type": "account.updated",
		"data": {
			"object": {
				"id": "%s",
				"object": "account",
				"charges_enabled": %t,
				"payouts_enabled": %t
			}
		}
	}`, accountID, stripe.APIVersion, accountID, chargesEnabled, payoutsEnabled)
	return []byte(payload)
}

func TestStripeConnectWebhook(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}
	webhookURL := testServerURL + webhookEndpoints.HandleStripeConnectWebhookEndpoint

	t.Run("should update partner status to active for valid account.updated event", func(t *testing.T) {
		td.ClearPartnersTable(t, ctx, testPool)

		stripeAccountID := "acct_integration_active_" + uuid.New().String()[:8]

		// Seed a partner with pending status and a known Stripe account ID
		testUser := td.NewTestUser(t, "webhook_partner_active@example.com", "Webhook", "Partner")
		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		testPartner := &domain.Partner{
			ID:                       uuid.New(),
			UserID:                   testUser.ID,
			Bio:                      "Integration test partner",
			Experience:               "3 years",
			CategoryIDs:              []uuid.UUID{},
			ProductIDs:               []uuid.UUID{},
			StripeConnectedAccountID: stripeAccountID,
			StripeAccountStatus:      domain.StripeAccountStatusPending,
			StripeOnboardingComplete: false,
		}
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, testPartner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Build and sign the webhook payload
		payload := buildAccountUpdatedPayload(stripeAccountID, true, true)
		signature := signStripePayload(payload, testConnectWebhookSecret)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewReader(payload))
		require.NoError(t, err)
		req.Header.Set("Stripe-Signature", signature)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Assert the partner's stripe_account_status is now active in the database
		updated, err := td.GetPartnerEncxByID(t, ctx, testPartner.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, domain.StripeAccountStatusActive, updated.StripeAccountStatus)
		assert.True(t, updated.StripeOnboardingComplete)
	})

	t.Run("should return 400 for missing Stripe-Signature header", func(t *testing.T) {
		payload := buildAccountUpdatedPayload("acct_missing_sig", true, true)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewReader(payload))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		// No Stripe-Signature header

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid Stripe-Signature header", func(t *testing.T) {
		payload := buildAccountUpdatedPayload("acct_invalid_sig", true, true)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewReader(payload))
		require.NoError(t, err)
		req.Header.Set("Stripe-Signature", "t=12345,v1=invalidsignature")
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 200 for unknown event type without touching the database", func(t *testing.T) {
		td.ClearPartnersTable(t, ctx, testPool)

		payload := []byte(fmt.Sprintf(`{
			"id": "evt_unknown",
			"object": "event",
			"api_version": "%s",
			"type": "some.unknown.event",
			"data": {"object": {}}
		}`, stripe.APIVersion))
		signature := signStripePayload(payload, testConnectWebhookSecret)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewReader(payload))
		require.NoError(t, err)
		req.Header.Set("Stripe-Signature", signature)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		count, err := td.CountPartners(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("should return 200 when account ID does not match any partner", func(t *testing.T) {
		td.ClearPartnersTable(t, ctx, testPool)

		payload := buildAccountUpdatedPayload("acct_no_match_partner", true, true)
		signature := signStripePayload(payload, testConnectWebhookSecret)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewReader(payload))
		require.NoError(t, err)
		req.Header.Set("Stripe-Signature", signature)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// No partner found — webhook acknowledges without error
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
