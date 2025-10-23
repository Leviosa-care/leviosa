package imageRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *ImageRepository) SetActiveImage(ctx context.Context, imageID uuid.UUID, parentID uuid.UUID, parentType domain.ParentType) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return errs.NewDBQueryErr(fmt.Errorf("failed to begin transaction: %w", err))
	}
	defer tx.Rollback(ctx)

	// Step 1: Deactivate the currently active image for this parent (if any).
	deactivateQuery := `
		UPDATE catalog.images
		SET is_active = FALSE
		WHERE parent_id = $1 AND parent_type = $2 AND is_active = TRUE;`

	_, err = tx.Exec(ctx, deactivateQuery, parentID, parentType)
	if err != nil {
		return errs.ClassifyPgError(fmt.Sprintf("deactivate old active image for %s %s", parentType, parentID), err)
	}

	// Step 2: Activate the new image.
	activateQuery := `
		UPDATE catalog.images
		SET is_active = TRUE
		WHERE id = $1 AND parent_id = $2 AND parent_type = $3;`

	res, err := tx.Exec(ctx, activateQuery, imageID, parentID, parentType)
	if err != nil {
		return errs.ClassifyPgError(fmt.Sprintf("activate new image %s for %s %s", imageID, parentType, parentID), err)
	}

	rowsAffected := res.RowsAffected()

	if rowsAffected == 0 {
		// This means the imageID to be activated was not found or didn't match the parent.
		// This is a "not found" scenario at the repository level.
		return errs.NewRepositoryNotFoundErr(nil, fmt.Sprintf("image with ID %s for %s %s", imageID, parentType, parentID))
	}

	// Commit the transaction if all operations were successful.
	if err := tx.Commit(ctx); err != nil {
		return errs.NewDBQueryErr(fmt.Errorf("failed to commit transaction: %w", err))
	}

	return nil
}
