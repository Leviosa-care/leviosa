package userRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/core/errs"
)

func (r *Repository) ExistsByGoogleID(ctx context.Context, googleID string) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS(
			SELECT 1 FROM %s.users 
			WHERE google_id_encrypted = $1
		)
	`, r.schema)

	var exists bool
	err := r.pool.QueryRow(ctx, query, googleID).Scan(&exists)
	if err != nil {
		return false, errs.ClassifyPgError("check if user exists by Google ID", err)
	}

	return exists, nil
}