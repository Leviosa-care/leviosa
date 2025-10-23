package imageRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// GetAllActiveImages retrieves all active images for a given parent type.
// It returns an empty slice if no active images are found, and handles database errors.
func (r *ImageRepository) GetAllActiveImages(ctx context.Context, parentType domain.ParentType) ([]*domain.Image, error) {
	query := `
		SELECT 
			id, parent_id, parent_type, title, s3_key, size, content_type, is_active, created_at, updated_at
		FROM 
			catalog.images
		WHERE
			parent_type = $1 AND is_active = TRUE;
	`

	rows, err := r.pool.Query(ctx, query, parentType)
	if err != nil {
		return nil, errs.ClassifyPgError(fmt.Sprintf("get active images for parent type '%s'", parentType), err)
	}
	defer rows.Close()

	images := []*domain.Image{}
	for rows.Next() {
		var img domain.Image
		err := rows.Scan(
			&img.ID,
			&img.ParentID,
			&img.ParentType,
			&img.Title,
			&img.S3Key,
			&img.Size,
			&img.ContentType,
			&img.IsActive,
			&img.CreatedAt,
			&img.UpdatedAt,
		)
		if err != nil {

			return nil, errs.NewDBQueryErr(fmt.Errorf("failed to scan image row for parent type '%s': %w", parentType, err))
		}
		images = append(images, &img)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.NewDBQueryErr(fmt.Errorf("error iterating over image rows for parent type '%s': %w", parentType, err))
	}

	if len(images) == 0 {
		return []*domain.Image{}, nil
	}

	// If no rows are found, return an empty slice, not a "not found" error.
	return images, nil
}
