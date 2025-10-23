package auth

import (
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

// RequireAdmin validates session and ensures user has admin role
func (m *SessionAuthMiddleware) RequireAdmin(next mw.Handler) mw.Handler {
	return m.RequireMinimumRole(identity.Administrator)(next)
}
