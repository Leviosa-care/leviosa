package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/settings/internal/domain"
)

func (r *repository) SetString(ctx context.Context, setting *domain.Setting[string]) error {
	query := `
	INSERT INTO settings.plain (key, value)
	VALUES ($1, $2)
	ON CONFLICT (key) DO UPDATE SET
		value = EXCLUDED.value,
		updated_at = NOW()
	`
	commandTag, err := r.pool.Exec(
		ctx,
		query,
		setting.Key,
		setting.Value,
	)
	if err != nil {
		return errs.ClassifyPgError(fmt.Sprintf("set string value for key '%s'", setting.Key), err)
	}

	if commandTag.RowsAffected() == 0 {
		return errs.NewNotCreatedErr(errors.New("no rows affected by insertion statement"), fmt.Sprintf("string setting for key '%s'", setting.Key))
	}
	return nil
}
