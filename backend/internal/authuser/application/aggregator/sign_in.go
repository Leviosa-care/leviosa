package aggregator

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s AuthAggregatorService) SignIn(ctx context.Context, request *domain.SignInRequest) (*domain.CreateSessionResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Get user by email
	user, err := s.user.GetUserByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.NewUnauthorizedErr("invalid credentials")
		}
		return nil, err
	}

	// Check if user is in active state (not pending or unverified)
	if user.State != domain.Active {
		return nil, errs.NewUnauthorizedErr("account not activated")
	}

	// Verify password
	if err := s.user.VerifyUserPassword(ctx, user.ID, request.Password); err != nil {
		if errors.Is(err, errs.ErrInvalidValue) {
			return nil, errs.NewUnauthorizedErr("invalid credentials")
		}
		return nil, err
	}

	// Convert role
	role, ok := identity.ConvertToRole(user.Role)
	if !ok {
		return nil, errs.NewInternalErr(fmt.Errorf("invalid user role: %s", user.Role))
	}

	// Create session
	token, err := s.session.CreateSession(ctx, &domain.CreateSessionRequest{
		UserID: user.ID.String(),
		Role:   role,
	})
	if err != nil {
	}

	if err := s.user.UpdateLastLoginTime(ctx, user.ID); err != nil {
		return nil, err
	}

	return token, nil
}
