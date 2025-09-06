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
	router.HandleFunc("GET /admin/auth/users/pending", RequireAdmin(mw.EnableCORS(h.GetPendingUsers)))

	// Retrieves all registered users (admin only).
	router.HandleFunc("GET /admin/users", RequireAdmin(mw.EnableCORS(h.GetAllUsers)))

	// Approves a pending user by setting their role and status to active (admin only).
	router.HandleFunc("PATCH /admin/users/approve", RequireAdmin(mw.EnableCORS(h.ApproveUser)))

	// Retrieves the profile of the currently authenticated user.
	router.HandleFunc("GET /users/me", RequireStandard(mw.EnableCORS(h.GetUser)))

	// TODO: ==============================
	// Suggested additional handlers:
	// ==============================

	// Updates the profile of the currently authenticated user.
	// router.HandleFunc("PATCH /users/me", mw.EnableCORS(h.UpdateUser))

	// Deletes the profile of the currently authenticated user (account removal).
	// router.HandleFunc("DELETE /users/me", mw.EnableCORS(h.DeleteUser))

	// Permanently deletes or suspends a user profile by ID (admin only).
	// router.HandleFunc("DELETE /admin/users/{id}", mw.EnableCORS(h.BanUser))

	// Retrieves details of a specific user by ID (admin only).
	// router.HandleFunc("GET /admin/users/{id}", mw.EnableCORS(h.GetUserByID))

	// Changes the password of the authenticated user (requires old password).
	// router.HandleFunc("PATCH /users/me/password", mw.EnableCORS(h.ChangePassword))

	// Endpoint for admins to udpate roles of existing users.
	// router.HandleFunc("PATCH /admin/users/{id}/role", mw.EnableCORS(h.UpdateUserRole))
}
