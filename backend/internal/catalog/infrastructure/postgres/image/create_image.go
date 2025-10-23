package imageRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *ImageRepository) CreateImage(ctx context.Context, image *domain.Image) error {
	query := `
		INSERT INTO catalog.images (
			id,
			parent_id,
			parent_type,
			title,
			s3_key,
			size,
			content_type,
			is_active,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pool.Exec(
		ctx,
		query,
		image.ID,
		image.ParentID,
		image.ParentType,
		image.Title,
		image.S3Key,
		image.Size,
		image.ContentType,
		image.IsActive,
		image.CreatedAt,
	)

	if err != nil {
		return errs.ClassifyPgError("create image", err)
	}
	return nil
}
