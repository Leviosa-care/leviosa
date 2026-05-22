package image

import (
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

const (
	ProductPrefix  = "public/assets/products"
	CategoryPrefix = "public/assets/categories"
	MaxImageSize   = 5 * 1024 * 1024 // 5 MB
)

func CreateParentPrefix(parentID uuid.UUID, parentType domain.ParentType) (string, error) {
	var prefix string
	switch parentType {
	case domain.CategoryType:
		prefix = CategoryPrefix
	case domain.ProductType:
		prefix = ProductPrefix
	}
	return fmt.Sprintf("%s/%s", prefix, parentID), nil
}

func CreateParentImagePrefix(parentID, imageID uuid.UUID, parentType domain.ParentType, contentType string) (string, error) {
	ext := getFileExtensionFromContentType(contentType) // Helper function needed
	if ext == "" {
		return "", errs.NewInvalidValueErr("unsupported file extension from content type")
	}
	bucketPrefix, err := CreateParentPrefix(parentID, parentType)
	if err != nil {
		// TODO: do something in that error case
	}
	return fmt.Sprintf("%s/%s%s", bucketPrefix, imageID, ext), nil
}

func isValidImageContentType(contentType string) bool {
	return contentType == "image/jpeg" || contentType == "image/png" || contentType == "image/webp"
}

func getFileExtensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}
