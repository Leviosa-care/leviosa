package domain

// IsValidImageContentType checks if the provided content type is a valid image type.
func IsValidImageContentType(contentType string) bool {
	// Basic check. Extend as needed (e.g., image/png, image/gif, image/webp)
	return contentType == "image/jpeg" || contentType == "image/png" || contentType == "image/gif"
}

// GetFileExtensionFromContentType determines a common file extension from a content type.
func GetFileExtensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	// Add more as needed
	default:
		return "" // Or ".bin" for unknown types, or return an error earlier
	}
}
