package ctxutil

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
)

func ValidateRoleInContext(ctx context.Context, expectedRole identity.Role) error {
	role, ok := ctx.Value(RoleKey).(identity.Role)
	if !ok {
		return fmt.Errorf("role not found in context")
	}
	if role != expectedRole {
		return fmt.Errorf("expected role %q, got %q", expectedRole, role)
	}
	return nil
}
