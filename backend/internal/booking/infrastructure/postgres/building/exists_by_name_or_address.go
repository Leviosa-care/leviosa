package buildingRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) ExistsByNameOrAddress(ctx context.Context, nameHash, addressHash string) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS(
			SELECT 1 FROM %s.buildings
			WHERE name_hash = $1 OR address_hash = $2
		)
	`, r.schema)

	var exists bool
	err := r.pool.QueryRow(ctx, query, nameHash, addressHash).Scan(&exists)
	if err != nil {
		return false, errs.ClassifyPgError("check building exists by name or address", err)
	}

	return exists, nil
}
