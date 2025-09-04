package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (s *UserService) CompleteUser(ctx context.Context, userID uuid.UUID, request *domain.CompleteUserRequest) error {
	// Validate the request
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(fmt.Sprintf("user completion validation failed: %s", err.Error()))
	}

	// Get the existing pending user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "user")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		default:
			return errs.NewInternalErr(fmt.Errorf("failed to get user: %w", err))
		}
	}

	// Verify user is in unverified state
	if user.State != domain.Unverified {
		return errs.NewConflictErr(fmt.Errorf("user is not in unverified state: %s", user.State))
	}

	// Create Stripe customer for the user
	stripeCustomer, err := s.stripe.CreateCustomer(ctx, userID, user.Email, request.FirstName, request.LastName)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			return errs.NewInvalidValueErr(fmt.Sprintf("stripe customer creation failed: %s", err.Error()))
		case errors.Is(err, errs.ErrPermissionDenied):
			return errs.NewPermissionErr(fmt.Sprintf("stripe customer creation failed: %s", err.Error()))
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrResourceExhausted):
			return errs.NewExternalServiceErr(err, "stripe service unavailable")
		default:
			return errs.NewInternalErr(fmt.Errorf("failed to create stripe customer: %w", err))
		}
	}

	// Update user with new information
	user.Password = request.Password
	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.BirthDate = request.BirthDate
	user.Telephone = request.Telephone
	user.PostalCode = request.PostalCode
	user.City = request.City
	user.Address1 = request.Address1
	user.Address2 = request.Address2
	user.StripeCustomerID = stripeCustomer.ID
	user.State = domain.Pending

	// Handle gender (convert from GenderInput to string)
	if request.Gender.Gender == domain.GenderCustom {
		user.Gender = request.Gender.CustomGender
	} else {
		user.Gender = string(request.Gender.Gender)
	}

	// Process encryption on the updated user
	if err := s.crypto.ProcessStruct(ctx, user); err != nil {
		return errs.NewNotEncryptedErr("user", err)
	}

	// Update user in repository
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "user")
		case errors.Is(err, errs.ErrConflict):
			return errs.NewConflictErr(fmt.Errorf("%w: user already completed", err))
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		default:
			return errs.NewInternalErr(fmt.Errorf("failed to update user: %w", err))
		}
	}

	return nil
}
