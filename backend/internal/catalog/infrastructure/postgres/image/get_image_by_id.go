package imageRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GetImageByID retrieves a single image record by its ID.
// It returns errs.NewRepositoryNotFoundErr if the image does not exist.
func (r *ImageRepository) GetImageByID(ctx context.Context, imageID uuid.UUID) (*domain.Image, error) {
	query := `
	SELECT 
	id, parent_id, parent_type, title, s3_key, size, content_type, is_active, created_at, updated_at
	FROM catalog.images 
	WHERE id = $1;`

	var image domain.Image
	err := r.pool.QueryRow(ctx, query, imageID).Scan(
		&image.ID,
		&image.ParentID,
		&image.ParentType,
		&image.Title,
		&image.S3Key,
		&image.Size,
		&image.ContentType,
		&image.IsActive,
		&image.CreatedAt,
		&image.UpdatedAt,
	)

	if err != nil {
		// if err == sql.ErrNoRows {
		if err == pgx.ErrNoRows {
			return nil, errs.NewRepositoryNotFoundErr(err, fmt.Sprintf("image with ID %s", imageID))
		}
		return nil, errs.ClassifyPgError(fmt.Sprintf("get image by ID %s", imageID), err)
	}

	return &image, nil
}
