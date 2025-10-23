package userRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.users 
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return errs.ClassifyPgError("delete user", err)
	}

	// Check if any row was actually deleted
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.NewRepositoryNotFoundErr(err, "user for deletion")
	}

	return nil
}

