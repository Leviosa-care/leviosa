package auth

import (
	"net/http"

	"github.com/Leviosa-care/core/contracts/identity"
)

// AuthMiddleware provides session-based authentication and role-based authorization
type AuthMiddleware interface {
	// RequireSession validates session and makes session available in context
	RequireSession(next http.Handler) http.Handler

	// RequireMinimumRole validates session and ensures user has at least the specified role
	RequireMinimumRole(minRole identity.Role) func(http.Handler) http.Handler

	// RequireAnyRole validates session and ensures user has one of the specified roles
	RequireAnyRole(roles ...identity.Role) func(http.Handler) http.Handler

	// RequireAdmin validates session and ensures user has admin role
	RequireAdmin(next http.Handler) http.Handler
}

