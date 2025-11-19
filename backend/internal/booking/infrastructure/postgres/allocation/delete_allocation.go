package allocationRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	// Soft delete by marking as inactive
	query := fmt.Sprintf(`
		UPDATE %s.room_allocations
		SET is_active = false, updated_at = NOW()
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return errs.ClassifyPgError("delete room allocation", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}
