package aggregator

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s AuthAggregatorService) SignIn(ctx context.Context, request *domain.SignInRequest) (*domain.CreateSessionResponse, error) {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		return nil, errs.NewInternalErr(fmt.Errorf("failed to get logger: %w", err))
	}

	logger.DebugContext(ctx, "Service: SignIn called",
		"email", request.Email,
		"request_id", ctx.Value("request_id"))

	if err := request.Valid(ctx); err != nil {
		logger.DebugContext(ctx, "Service: SignIn request validation failed",
			"error", err,
			"email", request.Email)
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Get user by email
	user, err := s.user.GetUserByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			logger.DebugContext(ctx, "Service: User not found in database",
				"email", request.Email,
				"error", err)
			return nil, errs.NewUnauthorizedErr("invalid credentials")
		}
		logger.DebugContext(ctx, "Service: Error getting user by email",
			"email", request.Email,
			"error", err)
		return nil, err
	}

	logger.DebugContext(ctx, "Service: User found",
		"user_id", user.ID.String(),
		"email", user.Email,
		"state", user.State,
		"role", user.Role,
		"is_active", user.State == domain.Active)

	// Check if user is in active state (not pending or unverified)
	if user.State != domain.Active {
		logger.DebugContext(ctx, "Service: User account not activated",
			"user_id", user.ID.String(),
			"current_state", user.State,
			"expected_state", domain.Active)
		return nil, errs.NewUnauthorizedErr("account not activated")
	}

	// Verify password
	logger.DebugContext(ctx, "Service: Verifying password",
		"user_id", user.ID.String())

	if err := s.user.VerifyUserPassword(ctx, user.ID, request.Password); err != nil {
		if errors.Is(err, errs.ErrInvalidValue) {
			logger.DebugContext(ctx, "Service: Password verification failed",
				"user_id", user.ID.String(),
				"error", err)
			return nil, errs.NewUnauthorizedErr("invalid credentials")
		}
		logger.DebugContext(ctx, "Service: Unexpected error during password verification",
			"user_id", user.ID.String(),
			"error", err)
		return nil, err
	}

	logger.DebugContext(ctx, "Service: Password verified successfully",
		"user_id", user.ID.String())

	// Convert role
	role, ok := identity.ConvertToRole(user.Role)
	if !ok {
		logger.DebugContext(ctx, "Service: Invalid user role",
			"user_id", user.ID.String(),
			"user_role", user.Role)
		return nil, errs.NewInternalErr(fmt.Errorf("invalid user role: %s", user.Role))
	}

	// Create session
	logger.DebugContext(ctx, "Service: Creating session",
		"user_id", user.ID.String(),
		"role", role)

	token, err := s.session.CreateSession(ctx, &domain.CreateSessionRequest{
		UserID: user.ID.String(),
		Role:   role,
		State:  session.SessionActive,
	})
	if err != nil {
		logger.DebugContext(ctx, "Service: Error creating session",
			"user_id", user.ID.String(),
			"error", err)
		return nil, err
	}

	if err := s.user.UpdateLastLoginTime(ctx, user.ID); err != nil {
		logger.DebugContext(ctx, "Service: Error updating last login time",
			"user_id", user.ID.String(),
			"error", err)
		return nil, err
	}

	logger.DebugContext(ctx, "Service: SignIn completed successfully",
		"user_id", user.ID.String(),
		"role", role)

	return token, nil
}
