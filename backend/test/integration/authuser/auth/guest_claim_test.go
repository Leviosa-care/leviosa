package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Integration tests for the guest claim flow:
//   - POST /auth/guest-claim        → validates input, sends OTP, returns 202
//   - POST /auth/guest-claim/verify → validates OTP, creates account, returns session

func TestGuestClaim(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Setup RabbitMQ for OTP message verification
	ch, err := testMQConn.Channel()
	require.NoError(t, err)
	defer ch.Close()

	td.SetupOTPQueue(t, ch)
	msgs := td.ConsumeOTPMessages(t, ch)

	t.Run("should send OTP for valid guest claim and return 202", func(t *testing.T) {
		td.ClearAuthTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		req := td.NewGuestClaimRequest(t, ctx, testServerURL,
			"guest@example.com", "0612345678", td.GenerateStrongPassword(t),
			"Jean", "Dupont",
		)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusAccepted, resp.StatusCode)
		message, status := td.ParseGuestClaimResponse(t, resp)
		assert.Equal(t, "OTP sent to email for guest claim verification", message)
		assert.Equal(t, "sent", status)

		// Verify OTP was created in Redis
		emailBytes, err := encx.SerializeValue("guest@example.com")
		require.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)
		exists := td.CheckOTPExists(t, ctx, emailHash, redisClient)
		assert.True(t, exists, "OTP should exist in Redis")

		// Verify RabbitMQ message was sent
		delivery := td.WaitForOTPMessage(t, msgs, 3*time.Second)
		assert.NotNil(t, delivery)
	})

	t.Run("should return 409 for already registered email", func(t *testing.T) {
		td.ClearAuthTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		// Insert existing user
		existingUser := td.NewTestUser(t, "taken@example.com", "Alice", "Martin")
		existingUserEncx, err := domain.ProcessUserEncx(ctx, crypto, existingUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, existingUserEncx, testPool)
		require.NoError(t, err)

		req := td.NewGuestClaimRequest(t, ctx, testServerURL,
			"taken@example.com", "0612345678", td.GenerateStrongPassword(t),
			"Bob", "Durand",
		)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, errorMsg, "already registered")
		assert.Equal(t, http.StatusConflict, statusCode)
	})

	t.Run("should return 400 for invalid input", func(t *testing.T) {
		td.ClearAuthTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		// Empty email, short password
		req := td.NewGuestClaimRequest(t, ctx, testServerURL,
			"", "0612345678", "short",
			"J", "D",
		)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestGuestClaimVerify(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Setup RabbitMQ for OTP
	ch, err := testMQConn.Channel()
	require.NoError(t, err)
	defer ch.Close()

	td.SetupOTPQueue(t, ch)
	msgs := td.ConsumeOTPMessages(t, ch)

	t.Run("happy path: verify OTP and create guest account", func(t *testing.T) {
		td.ClearAuthTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		email := "newguest@example.com"
		password := td.GenerateStrongPassword(t)

		// Step 1: guest-claim sends OTP
		claimReq := td.NewGuestClaimRequest(t, ctx, testServerURL,
			email, "0612345678", password, "Marie", "Curie",
		)
		resp, err := client.Do(claimReq)
		require.NoError(t, err)
		resp.Body.Close()
		require.Equal(t, http.StatusAccepted, resp.StatusCode)

		// Consume OTP from RabbitMQ to get the code
		delivery := td.WaitForOTPMessage(t, msgs, 3*time.Second)
		require.NotNil(t, delivery)

		// Extract OTP code from Redis
		emailBytes, err := encx.SerializeValue(email)
		require.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)
		otpEncx, err := td.GetOTPEncxByEmailHash(t, ctx, emailHash, redisClient)
		require.NoError(t, err)
		otp, err := domain.DecryptOTPEncx(ctx, crypto, otpEncx)
		require.NoError(t, err)

		// Step 2: guest-claim/verify
		verifyReq := td.NewGuestClaimVerifyRequest(t, ctx, testServerURL,
			email, otp.Code, password, "Marie", "Curie", "0612345678",
		)
		resp, err = client.Do(verifyReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		message, status := td.ParseGuestClaimVerifyResponse(t, resp)
		assert.Equal(t, "Guest account created successfully", message)
		assert.Equal(t, "created", status)

		// Verify session cookies are set
		var hasAccessToken, hasRefreshToken bool
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "leviosa_access_token" {
				hasAccessToken = true
			}
			if cookie.Name == "leviosa_refresh_token" {
				hasRefreshToken = true
			}
		}
		assert.True(t, hasAccessToken, "access token cookie should be set")
		assert.True(t, hasRefreshToken, "refresh token cookie should be set")

		// Verify user was created with correct state
		userEncx, err := td.GetUserEncxByEmailHash(t, ctx, emailHash, testPool, crypto)
		require.NoError(t, err)
		user, err := domain.DecryptUserEncx(ctx, crypto, userEncx)
		require.NoError(t, err)

		assert.Equal(t, domain.Active, user.State, "user should be active")
		assert.Equal(t, "standard", user.Role, "user should have standard role")
		assert.True(t, user.ProfileIncomplete, "profile should be incomplete")
		assert.Equal(t, "Marie", user.FirstName)
		assert.Equal(t, "Curie", user.LastName)
		assert.Equal(t, email, user.Email)
	})

	t.Run("wrong OTP returns 400", func(t *testing.T) {
		td.ClearAuthTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		email := "wrongotp@example.com"
		password := td.GenerateStrongPassword(t)

		// Step 1: send OTP
		claimReq := td.NewGuestClaimRequest(t, ctx, testServerURL,
			email, "0612345678", password, "Test", "User",
		)
		resp, err := client.Do(claimReq)
		require.NoError(t, err)
		resp.Body.Close()
		require.Equal(t, http.StatusAccepted, resp.StatusCode)

		// Consume OTP message
		delivery := td.WaitForOTPMessage(t, msgs, 3*time.Second)
		require.NotNil(t, delivery)

		// Step 2: verify with wrong OTP
		verifyReq := td.NewGuestClaimVerifyRequest(t, ctx, testServerURL,
			email, "000000", password, "Test", "User", "0612345678",
		)
		resp, err = client.Do(verifyReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("expired OTP returns 401", func(t *testing.T) {
		td.ClearAuthTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		email := "expiredotp@example.com"
		password := td.GenerateStrongPassword(t)

		// Insert an already-expired OTP
		expiredOTP := td.NewExpiredOTP(email)
		expiredOTPEncx, err := domain.ProcessOTPEncx(ctx, crypto, expiredOTP)
		require.NoError(t, err)
		err = td.InsertOTPEncx(t, ctx, expiredOTPEncx, redisClient, 10*time.Minute)
		require.NoError(t, err)

		// Step 2: verify with the expired OTP's code
		verifyReq := td.NewGuestClaimVerifyRequest(t, ctx, testServerURL,
			email, expiredOTP.Code, password, "Test", "User", "0612345678",
		)
		resp, err := client.Do(verifyReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("duplicate email on verify returns 409", func(t *testing.T) {
		td.ClearAuthTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		email := "duplicate@example.com"
		password := td.GenerateStrongPassword(t)

		// Insert existing user with same email
		existingUser := td.NewTestUser(t, email, "First", "User")
		existingUserEncx, err := domain.ProcessUserEncx(ctx, crypto, existingUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, existingUserEncx, testPool)
		require.NoError(t, err)

		// Insert a valid OTP
		validOTP := td.NewTestOTP(email)
		validOTPEncx, err := domain.ProcessOTPEncx(ctx, crypto, validOTP)
		require.NoError(t, err)
		err = td.InsertOTPEncx(t, ctx, validOTPEncx, redisClient, 10*time.Minute)
		require.NoError(t, err)

		// Verify with valid OTP but email already taken
		verifyReq := td.NewGuestClaimVerifyRequest(t, ctx, testServerURL,
			email, validOTP.Code, password, "Second", "User", "0612345678",
		)
		resp, err := client.Do(verifyReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})
}
