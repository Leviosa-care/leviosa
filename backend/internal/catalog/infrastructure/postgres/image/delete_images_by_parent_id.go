package imageRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// DeleteImagesByParentID deletes all image records associated with a specific parent.
// It returns the number of rows affected.
func (r *ImageRepository) DeleteImagesByParentID(ctx context.Context, parentID uuid.UUID, parentType domain.ParentType) (int64, error) {
	query := `DELETE FROM catalog.images WHERE parent_id = $1 AND parent_type = $2;`

	commandTag, err := r.pool.Exec(ctx, query, parentID, parentType)
	if err != nil {
		return 0, errs.ClassifyPgError(fmt.Sprintf("delete images by parent ID %s and type %s", parentID, parentType), err)
	}

	return commandTag.RowsAffected(), nil
}
