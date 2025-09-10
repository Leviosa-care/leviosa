package sessionRepository

import (
	"context"

	"github.com/Leviosa-care/core/errs"
)

func (r *SessionRepository) ValidateResetSession(ctx context.Context, tokenHash string) (string, error) {
	key := FormatResetSessionKey(tokenHash)

	// Get the user email hash associated with this reset token
	userEmailHash, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", errs.ClassifyRedisError("validate reset session", err)
	}

	// Delete the token immediately after successful validation (single-use)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		// Log the deletion error but don't fail the validation
		// The token was valid and we got the email hash
		return userEmailHash, errs.ClassifyRedisError("delete consumed reset session", err)
	}

	return userEmailHash, nil
}

