package userHandler

const (
	// Base paths
	UsersBasePath      = "/users"
	AdminUsersBasePath = "/admin/users"
	AdminAuthBasePath  = "/admin/auth"

	// === User Resource Paths ===

	// User profile
	MePath = "/me"

	// User management
	PendingPath  = "/pending"
	ApprovePath  = "/approve"
	RolePath     = "/role"
	PasswordPath = "/password"

	// === Public User Endpoints ===
	// (None - all user endpoints require authentication)

	// === Authenticated User Endpoints (Standard role) ===

	// Get current user profile
	GetUserEndpoint = UsersBasePath + MePath

	// Update current user
	UpdateUserEndpoint = UsersBasePath + MePath

	// Change password
	ChangePasswordEndpoint = UsersBasePath + MePath + PasswordPath

	// === Admin-Only Endpoints ===

	// Get pending users (users awaiting approval)
	GetPendingUsersEndpoint = AdminAuthBasePath + AdminUsersBasePath + PendingPath

	// Get all users
	GetAllUsersEndpoint = AdminUsersBasePath

	// Approve a pending user
	ApproveUserEndpoint = AdminUsersBasePath + ApprovePath

	// Get user by ID
	GetUserByIDEndpoint = AdminUsersBasePath + "/{id}"

	// Update user role
	UpdateUserRoleEndpoint = AdminUsersBasePath + "/{id}" + RolePath
)
