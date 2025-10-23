package ports

import (
	"context"
	"io"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

// TODO:
// - make only one interface and use aggregator instead of using dependencies in the handler

type ImageParentService interface {
	// GetActiveImage gets the single active image for a parent entity.
	GetActiveImage(ctx context.Context, parentIDStr string, parentTypeStr string) (*domain.Image, error)
	// GetAllImages gets all images for a parent entity, sorted by creation date.
	GetAllImages(ctx context.Context, parentIDStr string, parentTypeStr string) ([]*domain.Image, error)
	// DeleteImages deletes all image records for a given parent and removes their files from S3.
	DeleteImages(ctx context.Context, parentID string, parentType string) error
	// GetAllImages gets all active images for a parent type.
	GetAllActiveImages(ctx context.Context, parentType string) ([]*domain.Image, error)
}

type ImageCommandService interface {
	// AddImage uploads an image file and creates a new image record in the database.
	// It returns the ID of the newly created image.
	AddImage(ctx context.Context, req *domain.CreateImageRequest, file io.Reader, fileSize int64, contentType string) (string, error)
	// DeleteImage deletes an image record from the database and removes the file from S3.
	DeleteImage(ctx context.Context, request *domain.ImageModifierRequest) error
	// SetImageAsActive sets a specific image as the active image for a parent entity,
	// and deactivates any previous active image.
	SetActiveImage(ctx context.Context, request *domain.ImageModifierRequest) error
}

type ImageService interface {
	ImageParentService
	ImageCommandService
}
