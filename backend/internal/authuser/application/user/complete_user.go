package user

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *UserService) CompleteUser(ctx context.Context, userID uuid.UUID, request *domain.CompleteUserRequest) error {
	// Validate the request
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(fmt.Sprintf("user completion validation failed: %s", err.Error()))
	}

	// Get the existing pending user
	userEncx, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Decrypt the user data using the new generated function
	user, err := domain.DecryptUserEncx(ctx, s.crypto, userEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("user for completion", err)
	}

	// Verify user is in unverified state
	if user.State != domain.Unverified {
		return errs.NewConflictErr(fmt.Errorf("user is not in unverified state: %s", user.State))
	}

	// Create Stripe customer for the user
	stripeCustomer, err := s.stripe.CreateCustomer(ctx, userID, user.Email, request.FirstName, request.LastName)
	if err != nil {
		return fmt.Errorf("failed to create stripe customer: %w", err)
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

	// Encrypt the user data using the new generated function
	updatedUserEncx, err := domain.ProcessUserEncx(ctx, s.crypto, user)
	if err != nil {
		return errs.NewNotEncryptedErr("user", err)
	}

	// Update user in repository
	if err := s.repo.UpdateUser(ctx, updatedUserEncx); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
