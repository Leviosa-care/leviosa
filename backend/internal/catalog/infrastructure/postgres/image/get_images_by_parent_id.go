package imageRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetImagesByParentID retrieves all image records for a specific parent entity.
// It returns an empty slice if no images are found for the given parent.
func (r *ImageRepository) GetImagesByParentID(ctx context.Context, parentID uuid.UUID, parentType domain.ParentType) ([]*domain.Image, error) {
	query := `
	SELECT 
	id, parent_id, parent_type, title, s3_key, size, content_type, is_active, created_at, updated_at
	FROM catalog.images 
	WHERE parent_id = $1 AND parent_type = $2
	ORDER BY created_at ASC;` // Order by creation time for consistency

	rows, err := r.pool.Query(ctx, query, parentID, parentType)
	if err != nil {
		return nil, errs.ClassifyPgError(fmt.Sprintf("get images by parent ID %s and type %s", parentID, parentType), err)
	}
	defer rows.Close()

	var images []*domain.Image
	for rows.Next() {
		var image domain.Image
		err := rows.Scan(
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
			return nil, errs.NewDBQueryErr(fmt.Errorf("failed to scan image row for parent ID %s: %w", parentID, err))
		}
		images = append(images, &image)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.NewDBQueryErr(fmt.Errorf("error iterating over image rows for parent ID %s: %w", parentID, err))
	}

	if len(images) == 0 {
		return []*domain.Image{}, nil
	}

	return images, nil
}
