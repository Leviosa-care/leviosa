package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/jackc/pgx/v5"
)

func (r *repository) GetInt(ctx context.Context, key string) (*domain.Setting[int], error) {
	var res domain.Setting[int]
	var valueStr string

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
		&valueStr,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NewRepositoryNotFoundErr(err, fmt.Sprintf("int value for key '%s'", key))
		}
		return nil, errs.ClassifyPgError(fmt.Sprintf("get int value for key '%s'", key), err)
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert retrieved string value into int")
	}
	res.Key = key
	res.Value = value

	return &res, nil
}
