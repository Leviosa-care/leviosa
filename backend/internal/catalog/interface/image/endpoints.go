package imageHandler

const (
	// Base paths
	AdminImagesBasePath = "/admin/images"

	// === Image Resource Paths ===
	SetActivePath = "/set-active"

	// === Admin-Only Endpoints ===
	// (All image endpoints are admin-only)

	// Upload image (admin only)
	UploadImageEndpoint = AdminImagesBasePath

	// Remove image (admin only)
	RemoveImageEndpoint = AdminImagesBasePath

	// Set active image (admin only)
	SetActiveImageEndpoint = AdminImagesBasePath + SetActivePath
)
