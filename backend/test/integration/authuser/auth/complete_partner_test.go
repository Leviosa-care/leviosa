package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	authEndpoints "github.com/Leviosa-care/leviosa/backend/internal/authuser/interface/auth"
	catalogDomain "github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCompletePartner TEST_PATH=test/integration/authuser/auth/complete_partner_test.go

// clearAllTestData clears all test data from auth and catalog tables
func clearAllTestData(t *testing.T, ctx context.Context, pool *pgxpool.Pool, redisClient *redis.Client) {
	t.Helper()
	td.ClearAuthTestData(t, ctx, pool, redisClient)
	td.ClearCategoriesTable(t, ctx, pool)
	td.ClearProductsTable(t, ctx, pool)
}

func TestCompletePartner(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	validEmail := "partner@example.com"

	// Create valid complete partner request
	validRequest := domain.CompletePartnerRequest{
		// User fields
		Password:  "qPDAR0.4Z8{vpCO]",
		FirstName: "Jane",
		LastName:  "Partner",
		BirthDate: time.Now().AddDate(-30, 0, 0),
		Gender: domain.GenderInput{
			Gender: domain.GenderWoman,
		},
		Telephone:  "0687654321",
		PostalCode: "75001",
		City:       "Paris",
		Address1:   "123 Rue de Rivoli",
		Address2:   "Ap 4B",
		// Partner fields
		Bio:         "Experienced healthcare professional with 10 years of experience",
		Experience:  "10 years in home healthcare services",
		CategoryIDs: []uuid.UUID{}, // Will be populated with valid IDs in tests
		ProductIDs:  []uuid.UUID{}, // Will be populated with valid IDs in tests
	}

	t.Run("should successfully complete partner registration with pending session", func(t *testing.T) {
		// Clean state
		clearAllTestData(t, ctx, testPool, redisClient)

		// Create pending user
		pendingUser := newPendingUser(validEmail)
		pendingUserEncx, err := domain.ProcessUserEncx(ctx, crypto, pendingUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, pendingUserEncx, testPool)
		require.NoError(t, err)

		// Create pending session for this user
		pendingSession, err := td.NewTestSession(t, crypto)
		require.NoError(t, err)
		pendingSession.UserID = pendingUser.ID
		pendingSession.State = session.SessionPending

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		// Populate with valid catalog IDs from services
		validCategoryIDs, validProductIDs := getValidCatalogIDsFromServices(t, ctx)
		request := validRequest
		request.CategoryIDs = validCategoryIDs
		request.ProductIDs = validProductIDs

		// Make HTTP request with session access token
		req := td.NewCompletePartnerRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		message := td.ParseCompletePartnerResponse(t, resp)
		assert.Equal(t, "Partner registration completed successfully. Awaiting admin approval.", message)

		// Verify user was completed in database
		userEncx, err := td.GetUserEnxByID(t, ctx, pendingUser.ID, testPool)
		require.NoError(t, err)
		user, err := domain.DecryptUserEncx(ctx, crypto, userEncx)
		require.NoError(t, err)
		assert.Equal(t, request.FirstName, user.FirstName)
		assert.Equal(t, request.LastName, user.LastName)
		assert.Equal(t, domain.Pending, user.State) // Should still be pending awaiting admin approval
		assert.NotEmpty(t, user.StripeCustomerID)

		// Verify partner was created in database
		partnerEncx, err := td.GetPartnerEncxByUserID(t, ctx, pendingUser.ID, testPool)
		require.NoError(t, err)
		partner, err := domain.DecryptPartnerEncx(ctx, crypto, partnerEncx)
		require.NoError(t, err)
		assert.Equal(t, pendingUser.ID, partner.UserID)
		assert.Equal(t, request.Bio, partner.Bio)
		assert.Equal(t, request.Experience, partner.Experience)
		assert.ElementsMatch(t, request.CategoryIDs, partner.CategoryIDs)
		assert.ElementsMatch(t, request.ProductIDs, partner.ProductIDs)
		// assert.Nil(t, partner.VerifiedByUserID)

		// Verify session was removed after completion
		sessionValueExists, err := redisClient.Exists(ctx, session.SessionKeyPrefix+"*").Result()
		assert.Equal(t, int64(0), sessionValueExists, "All sessions should be removed after completion")
	})

	t.Run("should return 401 when session cookie is missing", func(t *testing.T) {
		// Clean state
		clearAllTestData(t, ctx, testPool, redisClient)

		// Populate with valid catalog IDs
		validCategoryIDs, validProductIDs := getValidCatalogIDsFromServices(t, ctx)
		request := validRequest
		request.CategoryIDs = validCategoryIDs
		request.ProductIDs = validProductIDs

		// Make HTTP request without session token
		req := td.NewCompletePartnerRequest(t, ctx, testServerURL, request, "")
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Contains(t, errorMsg, errs.ErrUnauthorized.Error())
	})

	t.Run("should return 409 when session is already active", func(t *testing.T) {
		// Clean state
		clearAllTestData(t, ctx, testPool, redisClient)

		// Create an active session
		activeSession, err := td.NewTestSession(t, crypto)
		require.NoError(t, err)
		activeSession.State = session.SessionActive

		sessionEncx, err := session.ProcessSessionEncx(ctx, crypto, activeSession)
		require.NoError(t, err)
		td.InsertSessionEncx(t, ctx, redisClient, sessionEncx, time.Hour)

		// Populate with valid catalog IDs
		validCategoryIDs, validProductIDs := getValidCatalogIDsFromServices(t, ctx)
		request := validRequest
		request.CategoryIDs = validCategoryIDs
		request.ProductIDs = validProductIDs

		// Make HTTP request with active session token
		req := td.NewCompletePartnerRequest(t, ctx, testServerURL, request, activeSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("should return 400 for invalid JSON request body", func(t *testing.T) {
		// Clean state
		clearAllTestData(t, ctx, testPool, redisClient)

		// Create a pending session
		pendingSession, err := td.NewTestSession(t, crypto)
		require.NoError(t, err)
		pendingSession.State = session.SessionPending

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		// Make HTTP request with invalid JSON (manually crafted)
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			testServerURL+authEndpoints.CompletePartnerEndpoint,
			nil, // No body - will cause JSON decode error
		)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: pendingSession.AccessToken,
		}
		req.AddCookie(cookie)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Contains(t, errorMsg, "invalid request body")
	})

	t.Run("should return 400 for invalid password", func(t *testing.T) {
		// Clean state
		clearAllTestData(t, ctx, testPool, redisClient)

		// Create a pending session
		pendingSession, err := td.NewTestSession(t, crypto)
		require.NoError(t, err)
		pendingSession.State = session.SessionPending

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		// Populate with valid catalog IDs
		validCategoryIDs, validProductIDs := getValidCatalogIDsFromServices(t, ctx)
		request := validRequest
		request.CategoryIDs = validCategoryIDs
		request.ProductIDs = validProductIDs
		request.Password = "weak"

		// Make HTTP request
		req := td.NewCompletePartnerRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should filter out invalid category IDs and create partner successfully", func(t *testing.T) {
		// Clean state
		clearAllTestData(t, ctx, testPool, redisClient)

		// Create pending user
		pendingUser := newPendingUser(validEmail)
		pendingUserEncx, err := domain.ProcessUserEncx(ctx, crypto, pendingUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, pendingUserEncx, testPool)
		require.NoError(t, err)

		// Create a pending session
		pendingSession, err := td.NewTestSession(t, crypto)
		require.NoError(t, err)
		pendingSession.UserID = pendingUser.ID
		pendingSession.State = session.SessionPending

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		// Use invalid category ID (non-existent in catalog)
		_, _ = getValidCatalogIDsFromServices(t, ctx)
		request := validRequest
		request.CategoryIDs = []uuid.UUID{uuid.New()} // Random UUID not in catalog
		request.ProductIDs = []uuid.UUID{}

		// Make HTTP request
		req := td.NewCompletePartnerRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response - should succeed with invalid IDs filtered out
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		message := td.ParseCompletePartnerResponse(t, resp)
		assert.Equal(t, "Partner registration completed successfully. Awaiting admin approval.", message)

		// Verify partner was created with empty category IDs (invalid ones filtered out)
		partnerEncx, err := td.GetPartnerEncxByUserID(t, ctx, pendingUser.ID, testPool)
		require.NoError(t, err)
		partner, err := domain.DecryptPartnerEncx(ctx, crypto, partnerEncx)
		require.NoError(t, err)
		assert.Equal(t, pendingUser.ID, partner.UserID)
		assert.Empty(t, partner.CategoryIDs, "Invalid category IDs should be filtered out")
		assert.Empty(t, partner.ProductIDs)
	})

	t.Run("should filter out invalid product IDs and create partner with valid categories", func(t *testing.T) {
		// Clean state
		clearAllTestData(t, ctx, testPool, redisClient)

		// Create pending user
		pendingUser := newPendingUser(validEmail)
		pendingUserEncx, err := domain.ProcessUserEncx(ctx, crypto, pendingUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, pendingUserEncx, testPool)
		require.NoError(t, err)

		// Create a pending session
		pendingSession, err := td.NewTestSession(t, crypto)
		require.NoError(t, err)
		pendingSession.UserID = pendingUser.ID
		pendingSession.State = session.SessionPending

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		// Use valid category but invalid product ID
		validCategoryIDs, _ := getValidCatalogIDsFromServices(t, ctx)
		request := validRequest
		request.CategoryIDs = validCategoryIDs
		request.ProductIDs = []uuid.UUID{uuid.New()} // Random UUID not in catalog

		// Make HTTP request
		req := td.NewCompletePartnerRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response - should succeed with invalid product IDs filtered out
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		message := td.ParseCompletePartnerResponse(t, resp)
		assert.Equal(t, "Partner registration completed successfully. Awaiting admin approval.", message)

		// Verify partner was created with valid categories but empty products
		partnerEncx, err := td.GetPartnerEncxByUserID(t, ctx, pendingUser.ID, testPool)
		require.NoError(t, err)
		partner, err := domain.DecryptPartnerEncx(ctx, crypto, partnerEncx)
		require.NoError(t, err)
		assert.Equal(t, pendingUser.ID, partner.UserID)
		assert.ElementsMatch(t, validCategoryIDs, partner.CategoryIDs, "Valid category IDs should be preserved")
		assert.Empty(t, partner.ProductIDs, "Invalid product IDs should be filtered out")
	})

	t.Run("should filter out invalid IDs and keep only valid ones when mixed", func(t *testing.T) {
		// Clean state
		clearAllTestData(t, ctx, testPool, redisClient)

		// Create pending user
		pendingUser := newPendingUser("mixed-ids@example.com")
		pendingUserEncx, err := domain.ProcessUserEncx(ctx, crypto, pendingUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, pendingUserEncx, testPool)
		require.NoError(t, err)

		// Create a pending session
		pendingSession, err := td.NewTestSession(t, crypto)
		require.NoError(t, err)
		pendingSession.UserID = pendingUser.ID
		pendingSession.State = session.SessionPending

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		// Get valid IDs from catalog
		validCategoryIDs, validProductIDs := getValidCatalogIDsFromServices(t, ctx)

		// Mix valid and invalid IDs
		request := validRequest
		request.CategoryIDs = append(validCategoryIDs, uuid.New(), uuid.New()) // 1 valid + 2 invalid
		request.ProductIDs = append(validProductIDs, uuid.New())               // 1 valid + 1 invalid

		// Make HTTP request
		req := td.NewCompletePartnerRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response - should succeed
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		message := td.ParseCompletePartnerResponse(t, resp)
		assert.Equal(t, "Partner registration completed successfully. Awaiting admin approval.", message)

		// Verify partner was created with only valid IDs
		partnerEncx, err := td.GetPartnerEncxByUserID(t, ctx, pendingUser.ID, testPool)
		require.NoError(t, err)
		partner, err := domain.DecryptPartnerEncx(ctx, crypto, partnerEncx)
		require.NoError(t, err)
		assert.Equal(t, pendingUser.ID, partner.UserID)
		assert.ElementsMatch(t, validCategoryIDs, partner.CategoryIDs, "Should only contain valid category IDs")
		assert.ElementsMatch(t, validProductIDs, partner.ProductIDs, "Should only contain valid product IDs")
		assert.Len(t, partner.CategoryIDs, 1, "Should have filtered out 2 invalid category IDs")
		assert.Len(t, partner.ProductIDs, 1, "Should have filtered out 1 invalid product ID")
	})

	t.Run("should return 400 for missing required partner fields", func(t *testing.T) {
		// Clean state
		clearAllTestData(t, ctx, testPool, redisClient)

		// Create a pending session
		pendingSession, err := td.NewTestSession(t, crypto)
		require.NoError(t, err)
		pendingSession.State = session.SessionPending

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		// Create partner request with missing user fields
		request := domain.CompletePartnerRequest{
			Password: validRequest.Password,
			// Missing all other required fields
		}

		// Make HTTP request
		req := td.NewCompletePartnerRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should successfully complete partner registration with minimal partner data", func(t *testing.T) {
		// Clean state
		clearAllTestData(t, ctx, testPool, redisClient)

		// Create pending user
		pendingUser := newPendingUser("minimal-partner@example.com")
		pendingUserEncx, err := domain.ProcessUserEncx(ctx, crypto, pendingUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, pendingUserEncx, testPool)
		require.NoError(t, err)

		// Create pending session
		pendingSession, err := td.NewTestSession(t, crypto)
		require.NoError(t, err)
		pendingSession.UserID = pendingUser.ID
		pendingSession.State = session.SessionPending

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		// Create partner request with minimal partner fields (bio, experience, certifications all optional)
		request := domain.CompletePartnerRequest{
			Password:  "qPDAR0.4Z8{vpCO]",
			FirstName: "Min",
			LastName:  "Partner",
			BirthDate: time.Now().AddDate(-25, 0, 0),
			Gender: domain.GenderInput{
				Gender: domain.GenderPreferNotToSay,
			},
			Telephone:   "0612345678",
			PostalCode:  "75001",
			City:        "Paris",
			Address1:    "123 Rue de Rivoli",
			Bio:         "",            // Optional
			Experience:  "",            // Optional
			CategoryIDs: []uuid.UUID{}, // Empty is valid
			ProductIDs:  []uuid.UUID{}, // Empty is valid
		}

		// Make HTTP request
		req := td.NewCompletePartnerRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify partner was created
		partnerEncx, err := td.GetPartnerEncxByUserID(t, ctx, pendingUser.ID, testPool)
		assert.NoError(t, err)
		partner, err := domain.DecryptPartnerEncx(ctx, crypto, partnerEncx)
		assert.NoError(t, err)
		assert.Equal(t, pendingUser.ID, partner.UserID)
	})
}

// getValidCatalogIDsFromServices retrieves valid category and product IDs by creating test data.
// This ensures tests have real catalog data to validate against.
func getValidCatalogIDsFromServices(t *testing.T, ctx context.Context) (categoryIDs []uuid.UUID, productIDs []uuid.UUID) {
	t.Helper()

	// Create a published category for testing
	testCategory := td.NewValidCategory("Test Partner Category")
	testCategory.Status = catalogDomain.Published // Set to Published for partner validation
	td.InsertCategory(t, ctx, testCategory, testPool)

	// Create a published product for testing
	testProduct := td.NewValidProduct("Test Partner Product", testCategory.ID)
	testProduct.Status = catalogDomain.Published // Set to Published for partner validation
	td.InsertProduct(t, ctx, testPool, testProduct)

	// Return the IDs for testing
	categoryIDs = []uuid.UUID{testCategory.ID}
	productIDs = []uuid.UUID{testProduct.ID}

	return categoryIDs, productIDs
}
