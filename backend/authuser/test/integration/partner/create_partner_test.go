package partner_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/authuser/test/helpers"
	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const httpTimeout = 10 * time.Second

func TestCreatePartner(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: httpTimeout}

	t.Run("should successfully create partner with specializations", func(t *testing.T) {
		// Clean state
		helpers.ClearPartnerTestData(t, ctx, testPool)

		// Create test specialization first
		spec := helpers.NewTestSpecialization("physiotherapist", "Physiotherapist", "Physical rehabilitation specialist")
		err := crypto.Encrypt(ctx, spec)
		require.NoError(t, err)
		helpers.InsertSpecialization(t, ctx, spec, testPool)

		// Create admin session for authorization
		adminSession := helpers.CreateTestSession(t, ctx, testClient, identity.Admin)

		// Test HTTP request
		request := domain.CreatePartnerRequest{
			Email:             "partner@test.com",
			FirstName:         "John",
			LastName:          "Doe",
			PhoneNumber:       "+1234567890",
			Bio:               "Experienced physiotherapist with 10 years of practice",
			Experience:        "10 years in musculoskeletal rehabilitation",
			Certifications:    []string{"DPT", "OCS", "FAAOMPT"},
			SpecializationIDs: []uuid.UUID{spec.ID},
		}

		req := helpers.NewCreatePartnerRequest(t, ctx, testServerURL, request)
		req.Header.Set("Cookie", helpers.ToCookieString(adminSession))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.CompletePartnerResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response structure
		assert.NotEqual(t, uuid.Nil, response.ID)
		assert.NotEqual(t, uuid.Nil, response.User.ID)
		assert.Equal(t, request.Email, response.User.Email)
		assert.Equal(t, request.FirstName, response.User.FirstName)
		assert.Equal(t, request.LastName, response.User.LastName)
		assert.Equal(t, request.PhoneNumber, response.User.PhoneNumber)
		assert.Equal(t, identity.Partner, response.User.Role)
		assert.Equal(t, request.Bio, response.Bio)
		assert.Equal(t, request.Experience, response.Experience)
		assert.Equal(t, request.Certifications, response.Certifications)
		assert.False(t, response.IsVerified) // New partners are not verified
		assert.Len(t, response.Specializations, 1)
		assert.Equal(t, spec.ID, response.Specializations[0].ID)

		// Verify database persistence
		dbPartner, err := helpers.GetPartnerFromDB(t, ctx, testPool, response.ID)
		require.NoError(t, err)
		assert.Equal(t, response.ID, dbPartner.ID)
		assert.Equal(t, response.User.ID, dbPartner.UserID)

		// Verify user persistence
		dbUser, err := helpers.GetUserFromDB(t, ctx, testPool, response.User.ID)
		require.NoError(t, err)
		assert.Equal(t, response.User.ID, dbUser.ID)
		assert.Equal(t, identity.Partner, dbUser.Role)

		// Verify specialization association
		exists := helpers.CheckPartnerSpecializationExists(t, ctx, testPool, response.ID, spec.ID)
		assert.True(t, exists)
	})

	t.Run("should fail with invalid email format", func(t *testing.T) {
		// Clean state
		helpers.ClearPartnerTestData(t, ctx, testPool)

		// Create admin session
		adminSession := helpers.CreateTestSession(t, ctx, testClient, identity.Admin)

		request := domain.CreatePartnerRequest{
			Email:       "invalid-email",
			FirstName:   "John",
			LastName:    "Doe",
			PhoneNumber: "+1234567890",
			Bio:         "Test bio",
			Experience:  "Test experience",
		}

		req := helpers.NewCreatePartnerRequest(t, ctx, testServerURL, request)
		req.Header.Set("Cookie", helpers.ToCookieString(adminSession))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail with duplicate email", func(t *testing.T) {
		// Clean state
		helpers.ClearPartnerTestData(t, ctx, testPool)

		// Create existing user with same email
		existingUser := helpers.NewTestUser("partner@test.com", "Existing", "User")
		existingUser.Role = identity.Standard.String()
		err := crypto.Encrypt(ctx, existingUser)
		require.NoError(t, err)
		helpers.InsertUser(t, ctx, existingUser, testPool)

		// Create admin session
		adminSession := helpers.CreateTestSession(t, ctx, testClient, identity.Admin)

		request := domain.CreatePartnerRequest{
			Email:       "partner@test.com", // Same email as existing user
			FirstName:   "John",
			LastName:    "Doe",
			PhoneNumber: "+1234567890",
			Bio:         "Test bio",
			Experience:  "Test experience",
		}

		req := helpers.NewCreatePartnerRequest(t, ctx, testServerURL, request)
		req.Header.Set("Cookie", helpers.ToCookieString(adminSession))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("should fail with non-existent specialization", func(t *testing.T) {
		// Clean state
		helpers.ClearPartnerTestData(t, ctx, testPool)

		// Create admin session
		adminSession := helpers.CreateTestSession(t, ctx, testClient, identity.Admin)

		request := domain.CreatePartnerRequest{
			Email:             "partner@test.com",
			FirstName:         "John",
			LastName:          "Doe",
			PhoneNumber:       "+1234567890",
			Bio:               "Test bio",
			Experience:        "Test experience",
			SpecializationIDs: []uuid.UUID{uuid.New()}, // Non-existent specialization
		}

		req := helpers.NewCreatePartnerRequest(t, ctx, testServerURL, request)
		req.Header.Set("Cookie", helpers.ToCookieString(adminSession))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should require admin role", func(t *testing.T) {
		// Clean state
		helpers.ClearPartnerTestData(t, ctx, testPool)

		// Create standard user session (not admin)
		standardSession := helpers.CreateTestSession(t, ctx, testClient, identity.Standard)

		request := domain.CreatePartnerRequest{
			Email:       "partner@test.com",
			FirstName:   "John",
			LastName:    "Doe",
			PhoneNumber: "+1234567890",
			Bio:         "Test bio",
			Experience:  "Test experience",
		}

		req := helpers.NewCreatePartnerRequest(t, ctx, testServerURL, request)
		req.Header.Set("Cookie", helpers.ToCookieString(standardSession))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should require authentication", func(t *testing.T) {
		// Clean state
		helpers.ClearPartnerTestData(t, ctx, testPool)

		request := domain.CreatePartnerRequest{
			Email:       "partner@test.com",
			FirstName:   "John",
			LastName:    "Doe",
			PhoneNumber: "+1234567890",
			Bio:         "Test bio",
			Experience:  "Test experience",
		}

		req := helpers.NewCreatePartnerRequest(t, ctx, testServerURL, request)
		// No session cookie

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}