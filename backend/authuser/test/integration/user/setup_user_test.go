package user_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	td "github.com/Leviosa-care/authuser/test/helpers"
	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func setupUser(t *testing.T, ctx context.Context, role identity.Role) string {
	// create admin user
	now := time.Now()
	user := &domain.User{
		ID:         uuid.New(),
		Role:       role.String(),
		State:      domain.Active,
		Email:      fmt.Sprintf("%s@leviosa.care", role.String()),
		FirstName:  role.String(),
		LastName:   role.String(),
		Password:   "bMPSrxQK#?rPO.[<",
		Telephone:  "0612345678",
		CreatedAt:  now,
		LoggedInAt: now,
	}
	// encrypt admin user
	err := crypto.ProcessStruct(ctx, user)
	require.NoError(t, err)
	// insert admin user
	td.InsertUser(t, ctx, user, testPool)

	// create an active session for that user
	sessionID := uuid.New()

	// Generate valid base64 tokens for testing
	accessToken, err := session.GenerateToken()
	require.NoError(t, err)

	refreshToken, err := session.GenerateToken()
	require.NoError(t, err)

	activeSession := &session.Session{
		ID:           sessionID,
		UserID:       user.ID,
		Role:         role,
		State:        session.SessionPending,
		CreatedAt:    now,
		ExpiresAt:    now.Add(24 * time.Hour),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// encrypt admin session
	err = crypto.ProcessStruct(context.Background(), activeSession)
	require.NoError(t, err, "Failed to re-process session struct for state update")

	// insert admin session in database
	td.InsertSessionDirectly(t, ctx, testClient, activeSession, time.Hour)
	// return admin access token hash to attach it to the request

	return activeSession.AccessToken
}

func setupStandardUser(t *testing.T, ctx context.Context) string {
	return setupUser(t, ctx, identity.Standard)
}

func setupAdminUser(t *testing.T, ctx context.Context) string {
	return setupUser(t, ctx, identity.Administrator)
}
