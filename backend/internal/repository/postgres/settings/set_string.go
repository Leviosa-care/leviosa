package settingsRepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/settings"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

func (r *repository) SetString(ctx context.Context, setting *settings.Setting[string]) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, key, value, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)
	`, pg.QualifiedTable(r.schema, "settings"))
	result, err := r.DB.ExecContext(
		ctx,
		query,
		setting.ID,
		setting.Key,
		setting.Value,
		setting.CreatedAt,
		setting.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return rp.NewContextErr(err)
		default:
			return rp.NewDatabaseErr(err)
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return rp.NewDatabaseErr(err)
	}
	if rowsAffected == 0 {
		return rp.NewNotCreatedErr(errors.New("no rows affected by insertion statement"), fmt.Sprintf("string setting for key '%s'", setting.Key))
	}
	return nil
}
