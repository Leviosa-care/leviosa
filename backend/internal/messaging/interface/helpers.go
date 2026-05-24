package messagingHandler

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/google/uuid"
)

func getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	// First try the session info (set by auth middleware)
	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if ok && sessionInfo.UserID != uuid.Nil {
		return sessionInfo.UserID, nil
	}

	// Fallback to context value
	raw := ctx.Value(ctxutil.UserIDKey)
	id, ok := raw.(uuid.UUID)
	if !ok {
		str, ok := raw.(string)
		if !ok {
			return uuid.Nil, fmt.Errorf("user ID not found in context")
		}
		return uuid.Parse(str)
	}
	return id, nil
}

func getRoleFromContext(ctx context.Context) (identity.Role, error) {
	// First try the session info
	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if ok {
		return sessionInfo.Role, nil
	}

	// Fallback to context value
	role, ok := ctx.Value(ctxutil.RoleKey).(identity.Role)
	if !ok {
		return identity.Visitor, fmt.Errorf("role not found in context")
	}
	return role, nil
}
