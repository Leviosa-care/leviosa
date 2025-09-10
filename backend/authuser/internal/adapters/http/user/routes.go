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

	// Retrieves details of a specific user by ID (admin only).
	router.HandleFunc("GET /admin/users/{id}", RequireAdmin(mw.EnableCORS(h.GetUserByID)))

	// Updates the role of a specific user by ID (admin only).
	router.HandleFunc("PATCH /admin/users/{id}/role", RequireAdmin(mw.EnableCORS(h.UpdateUserRole)))

	// Updates the profile of the currently authenticated user.
	router.HandleFunc("PATCH /users/me", RequireStandard(mw.EnableCORS(h.UpdateUser)))

	// TODO: ==============================
	// Suggested additional handlers:
	// ==============================

	// Changes the password of the authenticated user (requires old password).
	// router.HandleFunc("PATCH /users/me/password", mw.EnableCORS(h.ChangePassword))
}
