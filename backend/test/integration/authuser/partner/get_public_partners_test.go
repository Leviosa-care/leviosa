package partner_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME='^TestGetPublicPartners$$' TEST_PATH=test/integration/authuser/partner/get_public_partners_test.go

func TestGetPublicPartners(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return only active non-disabled partners unauthenticated", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// 1) Active partner with stripe_account_status = active → included
		user1 := td.NewTestUser(t, "active@example.com", "Alice", "Martin")
		user1.State = domain.Active
		user1.Role = identity.PartnerStr
		user1.Picture = "https://example.com/alice.jpg"
		user1Encx, err := domain.ProcessUserEncx(ctx, crypto, user1)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, user1Encx, testPool)
		require.NoError(t, err)

		partner1 := td.NewTestPartner(t, user1.ID)
		partner1.Occupation = "Kinésithérapeute du sport"
		partner1.Quote = "Le mouvement est la vie"
		partner1.Tags = []string{"sport", "rééducation", "blessures"}
		partner1.StripeAccountStatus = domain.StripeAccountStatusActive
		partner1.CategoryIDs = []uuid.UUID{uuid.New()}
		partner1Encx, err := domain.ProcessPartnerEncx(ctx, crypto, partner1)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partner1Encx, testPool)
		require.NoError(t, err)

		// 2) Partner with stripe_account_status = disabled → excluded
		user2 := td.NewTestUser(t, "disabled@example.com", "Bob", "Durand")
		user2.State = domain.Active
		user2.Role = identity.PartnerStr
		user2Encx, err := domain.ProcessUserEncx(ctx, crypto, user2)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, user2Encx, testPool)
		require.NoError(t, err)

		partner2 := td.NewTestPartner(t, user2.ID)
		partner2.StripeAccountStatus = domain.StripeAccountStatusDisabled
		partner2Encx, err := domain.ProcessPartnerEncx(ctx, crypto, partner2)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partner2Encx, testPool)
		require.NoError(t, err)

		// 3) Partner with user state = pending → excluded
		user3 := td.NewTestUser(t, "pending@example.com", "Charlie", "Lambert")
		user3.State = domain.Pending
		user3.Role = identity.PartnerStr
		user3Encx, err := domain.ProcessUserEncx(ctx, crypto, user3)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, user3Encx, testPool)
		require.NoError(t, err)

		partner3 := td.NewTestPartner(t, user3.ID)
		partner3.StripeAccountStatus = domain.StripeAccountStatusActive
		partner3Encx, err := domain.ProcessPartnerEncx(ctx, crypto, partner3)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partner3Encx, testPool)
		require.NoError(t, err)

		// Act: call GET /partners unauthenticated
		url := testServerURL + "/partners"
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var partners []domain.PublicPartnerResponse
		err = json.NewDecoder(resp.Body).Decode(&partners)
		require.NoError(t, err)

		// Only the active partner should be returned
		require.Len(t, partners, 1, "Expected exactly 1 public partner (active, non-disabled)")

		p := partners[0]
		assert.Equal(t, partner1.ID, p.ID)
		assert.Equal(t, "Alice", p.FirstName)
		assert.Equal(t, "Martin", p.LastName)
		assert.Equal(t, "Kinésithérapeute du sport", p.Occupation)
		assert.Equal(t, "Le mouvement est la vie", p.Quote)
		assert.Equal(t, []string{"sport", "rééducation", "blessures"}, p.Tags)
		assert.Equal(t, partner1.CategoryIDs, p.CategoryIDs)
	})

	t.Run("should return empty array when no public partners exist", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		url := testServerURL + "/partners"
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var partners []domain.PublicPartnerResponse
		err = json.NewDecoder(resp.Body).Decode(&partners)
		require.NoError(t, err)
		assert.Empty(t, partners)
		assert.NotNil(t, partners)
	})
}
