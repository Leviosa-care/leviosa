package user

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *UserService) VerifyUserPassword(ctx context.Context, userID uuid.UUID, password string) error {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get logger: %w", err)
	}

	logger.DebugContext(ctx, "Service: VerifyUserPassword called",
		"user_id", userID.String(),
		"password_length", len(password))

	// Get the user from repository to access the stored password hash
	userEncx, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		logger.DebugContext(ctx, "Service: Error getting user by ID during password verification",
			"user_id", userID.String(),
			"error", err)
		return fmt.Errorf("failed to get user by ID: %w", err)
	}

	logger.DebugContext(ctx, "Service: User retrieved for password verification",
		"user_id", userID.String(),
		"has_stored_hash", userEncx.PasswordHashSecure != "",
		"stored_hash_length", len(userEncx.PasswordHashSecure))

	ok, err := s.crypto.CompareSecureHashAndValue(ctx, password, userEncx.PasswordHashSecure)
	if err != nil {
		logger.DebugContext(ctx, "Service: Error comparing secure hash",
			"user_id", userID.String(),
			"error", err)
		return errs.NewUnexpectedError(err)
	}
	if !ok {
		logger.DebugContext(ctx, "Service: Password comparison failed - passwords do not match",
			"user_id", userID.String(),
			"provided_password_length", len(password),
			"stored_hash_length", len(userEncx.PasswordHashSecure))
		return errs.NewInvalidValueErr("password verification failed: provided password does not match stored hash")
	}

	logger.DebugContext(ctx, "Service: Password verification successful",
		"user_id", userID.String())

	return nil
}
