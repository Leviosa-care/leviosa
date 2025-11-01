package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"
)

func (r *repository) SetInt(ctx context.Context, setting *domain.Setting[int]) error {
	query := `
	INSERT INTO settings.plain (key, value)
	VALUES ($1, $2)
	ON CONFLICT (key) DO UPDATE SET
		value = EXCLUDED.value,
		updated_at = NOW()
	`
	valueStr := fmt.Sprintf("%d", setting.Value)
	commandTag, err := r.pool.Exec(
		ctx,
		query,
		setting.Key,
		valueStr,
	)
	if err != nil {
		return errs.ClassifyPgError(fmt.Sprintf("set int value for key '%s'", setting.Key), err)
	}

	if commandTag.RowsAffected() == 0 {
		return errs.NewNotCreatedErr(errors.New("no rows affected by insertion statement"), fmt.Sprintf("int setting for key '%s'", setting.Key))
	}
	return nil
}
