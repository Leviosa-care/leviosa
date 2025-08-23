package settings

import (
	"fmt"

	"github.com/Leviosa-care/core/errs"
)

const (
	SettingsPrefix = "public/assets/settings"
)

func CreateLogoPrefix(filename string, contentType string) (string, error) {
	ext := getFileExtensionFromContentType(contentType) // Helper function needed
	if ext == "" {
		return "", errs.NewInvalidValueErr("unsupported file extension from content type")
	}
	return fmt.Sprintf("%s/%s%s", SettingsPrefix, filename, ext), nil
}

func isValidImageContentType(contentType string) bool {
	// IsValidImageContentType checks if the provided content type is a valid image type.
	// Basic check. Extend as needed (e.g., image/png, image/gif, image/webp)
	return contentType == "image/jpeg" || contentType == "image/png" || contentType == "image/gif"
}

// GetFileExtensionFromContentType determines a common file extension from a content type.
func getFileExtensionFromContentType(contentType string) string {
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
