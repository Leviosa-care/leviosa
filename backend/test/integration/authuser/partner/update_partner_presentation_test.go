package partner_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdatePartnerPresentationFields TEST_PATH=test/integration/authuser/partner/update_partner_presentation_test.go

func TestUpdatePartnerPresentationFields(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	setupPartnerSession := func(t *testing.T) (string, uuid.UUID) {
		t.Helper()

		user := td.NewTestUser(t, "pres-partner@example.com", "Pres", "Partner")
		user.State = domain.Active
		user.Role = identity.PartnerStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		partner := td.NewTestPartner(t, user.ID)
		partner.Occupation = ""
		partner.Quote = ""
		partner.Tags = nil
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		now := time.Now()
		sessionID := uuid.New()
		accessToken, err := session.GenerateToken()
		require.NoError(t, err)
		refreshToken, err := session.GenerateToken()
		require.NoError(t, err)

		sess := &session.Session{
			ID:           sessionID,
			UserID:       user.ID,
			Role:         identity.Partner,
			State:        session.SessionActive,
			CreatedAt:    now,
			ExpiresAt:    now.Add(24 * time.Hour),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		sessEncx, err := session.ProcessSessionEncx(ctx, crypto, sess)
		require.NoError(t, err)
		td.InsertSessionEncx(t, ctx, redisClient, sessEncx, time.Hour)

		return accessToken, user.ID
	}

	t.Run("should update and return occupation, quote, and tags via PUT /partners/me", func(t *testing.T) {
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken, _ := setupPartnerSession(t)

		// Update all presentation fields
		occupation := "Kinésithérapeute du sport"
		quote := "Le mouvement est la vie"
		tags := []string{"sport", "rééducation", "blessures"}
		updateRequest := domain.UpdatePartnerRequest{
			Occupation: &occupation,
			Quote:      &quote,
			Tags:       &tags,
		}

		req := td.NewUpdatePartnerMeRequest(t, ctx, testServerURL, updateRequest, accessToken)
		resp, err := client.Do(req)

		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		partnerResp := td.ParsePartnerResponse(t, resp)
		assert.Equal(t, occupation, partnerResp.Occupation)
		assert.Equal(t, quote, partnerResp.Quote)
		assert.Equal(t, tags, partnerResp.Tags)
	})

	t.Run("should round-trip presentation fields through GET /partners/me", func(t *testing.T) {
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken, _ := setupPartnerSession(t)

		// Set all fields via PUT
		occupation := "Ostéopathe"
		quote := "Toucher le corps, écouter la vie"
		tags := []string{"ostéopathie", "douleurs chroniques", "cranien"}
		updateRequest := domain.UpdatePartnerRequest{
			Occupation: &occupation,
			Quote:      &quote,
			Tags:       &tags,
		}

		req := td.NewUpdatePartnerMeRequest(t, ctx, testServerURL, updateRequest, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// Fetch via GET /partners/me
		getReq := td.NewGetPartnerMeRequest(t, ctx, testServerURL, accessToken)
		getResp, err := client.Do(getReq)
		require.NoError(t, err)
		defer getResp.Body.Close()
		require.Equal(t, http.StatusOK, getResp.StatusCode)

		fetched := td.ParsePartnerResponse(t, getResp)
		assert.Equal(t, occupation, fetched.Occupation)
		assert.Equal(t, quote, fetched.Quote)
		assert.Equal(t, tags, fetched.Tags)
	})

	t.Run("should not clear existing values when fields are omitted", func(t *testing.T) {
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken, _ := setupPartnerSession(t)

		// Set all presentation fields
		occupation := "Infirmier"
		quote := "Prendre soin de chacun"
		tags := []string{"soins", "domicile"}
		setReq := domain.UpdatePartnerRequest{
			Occupation: &occupation,
			Quote:      &quote,
			Tags:       &tags,
		}
		req := td.NewUpdatePartnerMeRequest(t, ctx, testServerURL, setReq, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// Update only bio (omit presentation fields)
		newBio := "Nouvelle bio"
		bioReq := domain.UpdatePartnerRequest{
			Bio: &newBio,
		}
		req = td.NewUpdatePartnerMeRequest(t, ctx, testServerURL, bioReq, accessToken)
		resp, err = client.Do(req)
		require.NoError(t, err)
		resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify presentation fields are unchanged
		getReq := td.NewGetPartnerMeRequest(t, ctx, testServerURL, accessToken)
		getResp, err := client.Do(getReq)
		require.NoError(t, err)
		defer getResp.Body.Close()
		require.Equal(t, http.StatusOK, getResp.StatusCode)

		fetched := td.ParsePartnerResponse(t, getResp)
		assert.Equal(t, newBio, fetched.Bio, "Bio should be updated")
		assert.Equal(t, occupation, fetched.Occupation, "Occupation should remain unchanged")
		assert.Equal(t, quote, fetched.Quote, "Quote should remain unchanged")
		assert.Equal(t, tags, fetched.Tags, "Tags should remain unchanged")
	})

	t.Run("should return 400 for occupation exceeding max length", func(t *testing.T) {
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken, _ := setupPartnerSession(t)

		longOccupation := ""
		for i := 0; i < 201; i++ {
			longOccupation += "a"
		}
		req := td.NewUpdatePartnerMeRequest(t, ctx, testServerURL, domain.UpdatePartnerRequest{
			Occupation: &longOccupation,
		}, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for quote exceeding max length", func(t *testing.T) {
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken, _ := setupPartnerSession(t)

		longQuote := ""
		for i := 0; i < 301; i++ {
			longQuote += "a"
		}
		req := td.NewUpdatePartnerMeRequest(t, ctx, testServerURL, domain.UpdatePartnerRequest{
			Quote: &longQuote,
		}, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
