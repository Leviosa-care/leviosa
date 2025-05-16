package settingsRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/hengadev/leviosa/internal/domain/settings"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (r *repository) GetInt(ctx context.Context, key string) (*settings.Setting[int], error) {
	var res settings.Setting[int]
	var valueStr string
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
		&valueStr,
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
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert retrieved string value into int")
	}
	res.Value = value
	return &res, nil
}
