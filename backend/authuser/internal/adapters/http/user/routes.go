package userHandler

import (
	"net/http"

	"github.com/Leviosa-care/core/contracts/identity"
	mw "github.com/Leviosa-care/core/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin
	RequireStandard := h.authmw.RequireMinimumRole(identity.Standard)

	// Retrieves all users currently in pending state (admin only).
	router.HandleFunc("GET "+GetPendingUsersEndpoint, RequireAdmin(mw.EnableCORS(h.GetPendingUsers)))

	// Retrieves all registered users (admin only).
	router.HandleFunc("GET "+GetAllUsersEndpoint, RequireAdmin(mw.EnableCORS(h.GetAllUsers)))

	// Approves a pending user by setting their role and status to active (admin only).
	router.HandleFunc("PATCH "+ApproveUserEndpoint, RequireAdmin(mw.EnableCORS(h.ApproveUser)))

	// Retrieves the profile of the currently authenticated user.
	router.HandleFunc("GET "+GetUserEndpoint, RequireStandard(mw.EnableCORS(h.GetUser)))

	// Retrieves details of a specific user by ID (admin only).
	router.HandleFunc("GET "+GetUserByIDEndpoint, RequireAdmin(mw.EnableCORS(h.GetUserByID)))

	// Updates the role of a specific user by ID (admin only).
	router.HandleFunc("PATCH "+UpdateUserRoleEndpoint, RequireAdmin(mw.EnableCORS(h.UpdateUserRole)))

	// Updates the profile of the currently authenticated user.
	router.HandleFunc("PATCH "+UpdateUserEndpoint, RequireStandard(mw.EnableCORS(h.UpdateUser)))

	// Changes the password of the authenticated user (requires old password).
	router.HandleFunc("PATCH "+ChangePasswordEndpoint, RequireStandard(mw.EnableCORS(h.ChangePassword)))
}
