package user

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/validation"

	"github.com/hengadev/encx"
)

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.UserResponse, error) {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get logger: %w", err)
	}

	logger.DebugContext(ctx, "Service: GetUserByEmail called",
		"email", email)

	if err := validation.ValidateEmail(email); err != nil {
		logger.DebugContext(ctx, "Service: Email validation failed",
			"email", email,
			"error", err)
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	emailBytes, err := encx.SerializeValue(email)
	if err != nil {
		logger.DebugContext(ctx, "Service: Failed to serialize email",
			"email", email,
			"error", err)
		return nil, errs.NewInvalidValueErr(err.Error())
	}
	emailHash := s.crypto.HashBasic(ctx, emailBytes)

	logger.DebugContext(ctx, "Service: Generated email hash",
		"email", email,
		"email_hash", fmt.Sprintf("%x", emailHash))

	userEncx, err := s.repo.GetUserByEmailHash(ctx, emailHash)
	if err != nil {
		logger.DebugContext(ctx, "Service: User not found in database",
			"email", email,
			"email_hash", fmt.Sprintf("%x", emailHash),
			"error", err)
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	logger.DebugContext(ctx, "Service: User found, decrypting data",
		"user_id", userEncx.ID.String(),
		"email", email,
		"key_version", userEncx.KeyVersion,
		"has_dek", len(userEncx.DEKEncrypted) > 0)

	// Decrypt user data
	user, err := domain.DecryptUserEncx(ctx, s.crypto, userEncx)
	if err != nil {
		logger.DebugContext(ctx, "Service: Failed to decrypt user data",
			"user_id", userEncx.ID.String(),
			"email", email,
			"error", err)
		return nil, errs.NewNotDecryptedErr("user retrieved by email", err)
	}

	logger.DebugContext(ctx, "Service: User decrypted successfully",
		"user_id", user.ID.String(),
		"email", user.Email,
		"state", user.State,
		"role", user.Role)

	return user.ToResponse(), nil
}
