package userRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/core/errs"
)

func (r *Repository) ExistsByAppleID(ctx context.Context, appleID string) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS(
			SELECT 1 FROM %s.users 
			WHERE apple_id_encrypted = $1
		)
	`, r.schema)

	var exists bool
	err := r.pool.QueryRow(ctx, query, appleID).Scan(&exists)
	if err != nil {
		return false, errs.ClassifyPgError("check if user exists by Apple ID", err)
	}

	return exists, nil
}