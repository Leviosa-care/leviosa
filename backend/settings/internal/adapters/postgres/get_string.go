package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/Leviosa-care/core/errs"
	"github.com/jackc/pgx/v5"
)

func (r *repository) GetString(ctx context.Context, key string) (*domain.Setting[string], error) {
	var res domain.Setting[string]
	query := `
	SELECT
	id,
	value,
	created_at,
	updated_at
	FROM settings.plain
	WHERE key = $1;`

	err := r.pool.QueryRow(ctx, query, key).Scan(
		&res.ID,
		&res.Value,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NewRepositoryNotFoundErr(err, fmt.Sprintf("string value for key '%s'", key))
		}
		return nil, errs.ClassifyPgError(fmt.Sprintf("get string value for key '%s'", key), err)
	}
	res.Key = key
	return &res, nil
}
