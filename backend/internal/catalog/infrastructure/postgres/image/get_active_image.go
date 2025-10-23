package imageRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GetActiveImage retrieves the single active image record for a parent entity.
// It returns errs.NewRepositoryNotFoundErr if no active image is found for the given parent.
func (r *ImageRepository) GetActiveImage(ctx context.Context, parentID uuid.UUID, parentType domain.ParentType) (*domain.Image, error) {
	query := `
	SELECT 
	id, parent_id, parent_type, title, s3_key, size, content_type, is_active, created_at, updated_at
	FROM catalog.images
	WHERE parent_id = $1 AND parent_type = $2 AND is_active = TRUE;`

	var image domain.Image
	err := r.pool.QueryRow(ctx, query, parentID, parentType).Scan(
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
		if err == pgx.ErrNoRows { // Correctly check for pgx.ErrNoRows
			return nil, errs.NewRepositoryNotFoundErr(err, fmt.Sprintf("active image for %s with ID %s", parentType, parentID))
		}
		return nil, errs.ClassifyPgError(fmt.Sprintf("get active image for %s with ID %s", parentType, parentID), err)
	}

	return &image, nil
}
