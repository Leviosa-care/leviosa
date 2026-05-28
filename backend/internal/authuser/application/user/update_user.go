package user

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *UserService) UpdateUser(ctx context.Context, userID uuid.UUID, request *domain.UpdateUserRequest) (*domain.UserResponse, error) {
	// Get existing user from repository
	userEncx, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user for update: %w", err)
	}

	// Decrypt user data to allow field updates using the new generated function
	user, err := domain.DecryptUserEncx(ctx, s.crypto, userEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("user for update", err)
	}

	// Update only non-nil fields from request
	if request.Picture != nil {
		user.Picture = *request.Picture
	}
	if request.FirstName != nil {
		user.FirstName = *request.FirstName
	}
	if request.LastName != nil {
		user.LastName = *request.LastName
	}
	if request.BirthDate != nil {
		user.BirthDate = *request.BirthDate
	}
	if request.Gender != nil {
		user.Gender = *request.Gender
	}
	if request.Email != nil {
		user.Email = *request.Email
	}
	if request.Telephone != nil {
		user.Telephone = *request.Telephone
	}
	if request.PostalCode != nil {
		user.PostalCode = *request.PostalCode
	}
	if request.City != nil {
		user.City = *request.City
	}
	if request.Address1 != nil {
		user.Address1 = *request.Address1
	}
	if request.Address2 != nil {
		user.Address2 = *request.Address2
	}
	if request.GoogleID != nil {
		user.GoogleID = *request.GoogleID
	}
	if request.AppleID != nil {
		user.AppleID = *request.AppleID
	}

	// Auto-clear profile_incomplete once all required fields are filled.
	if user.ProfileIncomplete {
		if user.Gender != "" && !user.BirthDate.IsZero() && user.Address1 != "" {
			user.ProfileIncomplete = false
		}
	}

	// Encrypt the user data using the new generated function
	updatedUserEncx, err := domain.ProcessUserEncx(ctx, s.crypto, user)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("user for update", err)
	}

	// Save updated user to repository
	if err := s.repo.UpdateUser(ctx, updatedUserEncx); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	// Convert to response format (use the plain user object)
	response := user.ToResponse()
	return response, nil
}
