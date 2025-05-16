package settingsRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/settings"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (r *repository) GetString(ctx context.Context, key string) (*settings.Setting[string], error) {
	var res settings.Setting[string]
	query := `
        SELECT
			id,
            value,
            created_at,
            updated_at
        FROM settings
        WHERE key = ?;`
	err := r.DB.QueryRowContext(ctx, query, key).Scan(
		&res.ID,
		&res.Value,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, rp.NewNotFoundErr(err, fmt.Sprintf("settings value %s", key))
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return nil, rp.NewContextErr(err)
		default:
			return nil, rp.NewDatabaseErr(err)
		}
	}
	return &res, nil
}
