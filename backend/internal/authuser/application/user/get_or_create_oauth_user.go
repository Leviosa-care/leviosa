package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetOrCreateOAuthUser retrieves an existing OAuth user or creates a new one
func (s *UserService) GetOrCreateOAuthUser(ctx context.Context, provider, oauthUserID, email, firstName, lastName string) (*domain.UserResponse, bool, error) {
	// Validate input parameters
	if provider == "" || oauthUserID == "" || email == "" {
		return nil, false, errs.NewInvalidValueErr("provider, OAuth user ID, and email are required")
	}

	// First, check if a user with this OAuth ID already exists
	existingOAuthUser, err := s.findUserByOAuthID(ctx, provider, oauthUserID)
	if err != nil {
		// Check if it's a "not found" error by examining if the error contains the expected types
		isNotFound := false
		// Try to determine if this is a "not found" scenario by checking error chain
		var currentErr error = err
		for currentErr != nil {
			if errors.Is(currentErr, errs.ErrRepositoryNotFound) {
				isNotFound = true
				break
			}
			currentErr = errors.Unwrap(currentErr)
		}
		
		if !isNotFound {
			return nil, false, fmt.Errorf("failed to check for existing OAuth user: %w", err)
		}
		// If it's a not found error, continue with the logic
	}

	if existingOAuthUser != nil {
		// OAuth user already exists, return it
		return existingOAuthUser, false, nil
	}

	// OAuth user doesn't exist, check if a user with this email exists
	existingEmailUser, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		// Check if it's a "not found" error
		isNotFound := false
		var currentErr error = err
		for currentErr != nil {
			if errors.Is(currentErr, errs.ErrRepositoryNotFound) {
				isNotFound = true
				break
			}
			currentErr = errors.Unwrap(currentErr)
		}
		
		if !isNotFound {
			return nil, false, fmt.Errorf("failed to check for existing email user: %w", err)
		}
		// If it's a not found error, continue with the logic
	}

	if existingEmailUser != nil {
		// User exists with this email but not linked to this OAuth provider
		// Link the OAuth account to the existing user
		return s.linkOAuthAccount(ctx, existingEmailUser, provider, oauthUserID)
	}

	// Neither OAuth user nor email user exists, create new OAuth user
	return s.createNewOAuthUser(ctx, provider, oauthUserID, email, firstName, lastName)
}

// findUserByOAuthID finds a user by their OAuth provider ID
func (s *UserService) findUserByOAuthID(ctx context.Context, provider, oauthUserID string) (*domain.UserResponse, error) {
	switch provider {
	case "google":
		return s.GetUserByGoogleID(ctx, oauthUserID)
	case "apple":
		return s.GetUserByAppleID(ctx, oauthUserID)
	default:
		return nil, errs.NewInvalidValueErr("unsupported provider: " + provider)
	}
}

// linkOAuthAccount links an OAuth provider to an existing user account
func (s *UserService) linkOAuthAccount(ctx context.Context, user *domain.UserResponse, provider, oauthUserID string) (*domain.UserResponse, bool, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(user.ID.String())
	if err != nil {
		return nil, false, errs.NewInternalErr(fmt.Errorf("invalid user ID format: %w", err))
	}

	// Create update request to link OAuth account
	updateRequest := &domain.UpdateUserRequest{}
	switch provider {
	case "google":
		updateRequest.GoogleID = &oauthUserID
	case "apple":
		updateRequest.AppleID = &oauthUserID
	default:
		return nil, false, errs.NewInvalidValueErr("unsupported provider: " + provider)
	}

	// Update the user
	updatedUser, err := s.UpdateUser(ctx, userUUID, updateRequest)
	if err != nil {
		return nil, false, fmt.Errorf("failed to link OAuth account: %w", err)
	}

	return updatedUser, false, nil
}

// createNewOAuthUser creates a new user account from OAuth data
func (s *UserService) createNewOAuthUser(ctx context.Context, provider, oauthUserID, email, firstName, lastName string) (*domain.UserResponse, bool, error) {
	// Generate new user ID
	userID := uuid.New()

	// Create new user
	user := &domain.User{
		ID:        userID,
		State:     domain.Active, // OAuth users are automatically active
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      identity.Standard.String(), // Default role for new OAuth users
		CreatedAt: time.Now(),
	}

	// Set OAuth provider ID
	switch provider {
	case "google":
		user.GoogleID = oauthUserID
	case "apple":
		user.AppleID = oauthUserID
	default:
		return nil, false, errs.NewInvalidValueErr("unsupported provider: " + provider)
	}

	// Encrypt the user data using the new generated function
	userEncx, err := domain.ProcessUserEncx(ctx, s.crypto, user)
	if err != nil {
		return nil, false, fmt.Errorf("failed to encrypt user data: %w", err)
	}

	// Create user in repository
	err = s.repo.CreateUser(ctx, userEncx)
	if err != nil {
		if errors.Is(err, errs.ErrUniqueViolation) {
			return nil, false, errs.NewConflictErr(fmt.Errorf("user with this email already exists"))
		}
		return nil, false, fmt.Errorf("failed to create OAuth user: %w", err)
	}

	// Create Stripe customer for new user (use the plain user object for this)
	stripeCustomer, err := s.stripe.CreateCustomer(ctx, userID, user.Email, user.FirstName, user.LastName)
	if err != nil {
		// Log error but don't fail user creation
		// The Stripe customer can be created later if needed
		// TODO: Consider implementing a retry mechanism or async processing
	} else {
		// Update user with Stripe customer ID
		user.StripeCustomerID = stripeCustomer.ID
		// Encrypt the updated user data
		updatedUserEncx, err := domain.ProcessUserEncx(ctx, s.crypto, user)
		if err == nil {
			s.repo.UpdateUser(ctx, updatedUserEncx) // Best effort update
		}
	}

	// Convert to response format (use the plain user object)
	return user.ToResponse(), true, nil
}