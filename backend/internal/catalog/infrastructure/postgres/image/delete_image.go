package imageRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// DeleteImage deletes an image record by its ID.
// It returns errs.ErrRepositoryNotFound if the image does not exist.
func (r *ImageRepository) DeleteImage(ctx context.Context, imageID uuid.UUID) error {
	query := "DELETE FROM catalog.images WHERE id = $1;"

	commandTag, err := r.pool.Exec(ctx, query, imageID)
	if err != nil {
		return errs.ClassifyPgError(fmt.Sprintf("delete image with ID %s", imageID), err)
	}

	if commandTag.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(nil, fmt.Sprintf("image with ID %s", imageID))
	}
	return nil
}
