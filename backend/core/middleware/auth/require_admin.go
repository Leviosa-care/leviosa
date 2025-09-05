package auth

import (
	"github.com/Leviosa-care/core/contracts/identity"
	mw "github.com/Leviosa-care/core/middleware"
)

// RequireAdmin validates session and ensures user has admin role
func (m *SessionAuthMiddleware) RequireAdmin(next mw.Handler) mw.Handler {
	return m.RequireMinimumRole(identity.Administrator)(next)
}