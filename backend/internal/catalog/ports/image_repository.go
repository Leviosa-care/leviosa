package ports

import (
	"context"
	"io"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/google/uuid"
)

type ImageRepository interface {
	// reader
	// GetImageByID retrieves a single image record by its ID.
	GetImageByID(ctx context.Context, imageID uuid.UUID) (*domain.Image, error)
	// GetImagesByParentID retrieves all image records for a parent entity.
	GetImagesByParentID(ctx context.Context, parentID uuid.UUID, parentType domain.ParentType) ([]*domain.Image, error)
	// GetActiveImage retrieves the single active image record for a parent entity.
	GetActiveImage(ctx context.Context, parentID uuid.UUID, parentType domain.ParentType) (*domain.Image, error)
	// GetActiveImage retrieves all active images given a parent type ('category', 'product').
	GetAllActiveImages(ctx context.Context, parentType domain.ParentType) ([]*domain.Image, error)
	// writer
	// CreateImage inserts a new image record into the database.
	CreateImage(ctx context.Context, image *domain.Image) error
	// SetImageActive updates a specific image record to be active, and others inactive.
	SetActiveImage(ctx context.Context, imageID uuid.UUID, parentID uuid.UUID, parentType domain.ParentType) error
	// DeleteImage removes an image record from the database.
	DeleteImage(ctx context.Context, imageID uuid.UUID) error
	// DeleteImagesByParentID removes all images record from the database for a given parent ID and parent type.
	DeleteImagesByParentID(ctx context.Context, parentID uuid.UUID, parentType domain.ParentType) (int64, error)
}

// ImageMediaRepository defines the contract for all file storage operations (e.g., S3).
type ImageMedia interface {
	// reader
	// GetFilePresignedURL generates a signed URL for a file, allowing temporary public access.
	// GetFilePresignedURL(ctx context.Context, key string) (string, error)

	// writer
	UploadFile(ctx context.Context, key string, file io.Reader, size int64, contentType string) (string, error)
	// DeleteFile removes a file from S3 using its key.
	DeleteFile(ctx context.Context, key string) error
}
