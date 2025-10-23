package categoryRepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *CategoryRepository) CategoryExistsByName(ctx context.Context, name string) (bool, error) {
	query := `
	SELECT 1
	FROM catalog.categories
	WHERE name = $1
	LIMIT 1;
	`

	var exists int
	err := r.pool.QueryRow(ctx, query, name).Scan(&exists)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			// Category doesn't exist — not an error
			return false, nil
		}
		// Actual query error — wrap as database error
		return false, errs.ClassifyPgError("check category existence ", err)
	}

	return true, nil
}
