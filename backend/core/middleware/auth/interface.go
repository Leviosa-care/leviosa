package auth

import (
	"github.com/Leviosa-care/core/contracts/identity"
	mw "github.com/Leviosa-care/core/middleware"
)

// // AuthMiddleware provides session-based authentication and role-based authorization
// type AuthMiddleware interface {
// 	// Dual-token authentication
// 	RequireAccessToken(next http.Handler) http.Handler
// 	RequireRefreshToken(next http.Handler) http.Handler
//
// 	// Role-based authorization (works with access token authentication)
// 	RequireMinimumRole(minRole identity.Role) func(http.Handler) http.Handler
// 	RequireAnyRole(roles ...identity.Role) func(http.Handler) http.Handler
// 	RequireAdmin(next http.Handler) http.Handler
// }

// AuthMiddleware provides session-based authentication and role-based authorization
type AuthMiddleware interface {
	// Dual-token authentication
	RequireAccessToken(next mw.Handler) mw.Handler
	RequireRefreshToken(next mw.Handler) mw.Handler

	// Role-based authorization (works with access token authentication)
	RequireMinimumRole(minRole identity.Role) func(mw.Handler) mw.Handler
	RequireAnyRole(roles ...identity.Role) func(mw.Handler) mw.Handler
	RequireAdmin(next mw.Handler) mw.Handler
}
