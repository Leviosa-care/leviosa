package auth

import (
	"net/http"

	"github.com/Leviosa-care/core/contracts/identity"
)

// AuthMiddleware provides session-based authentication and role-based authorization
type AuthMiddleware interface {
	// Dual-token authentication
	RequireAccessToken(next http.Handler) http.Handler
	RequireRefreshToken(next http.Handler) http.Handler

	// Role-based authorization (works with access token authentication)
	RequireMinimumRole(minRole identity.Role) func(http.Handler) http.Handler
	RequireAnyRole(roles ...identity.Role) func(http.Handler) http.Handler
	RequireAdmin(next http.Handler) http.Handler
}

