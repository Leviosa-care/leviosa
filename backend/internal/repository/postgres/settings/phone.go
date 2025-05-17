package settingsRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/settings"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (r *repository) GetPhone(ctx context.Context) (*settings.SettingEncrypted[string], error) {
	var res settings.SettingEncrypted[string]
	query := `
        SELECT
			id,
            value_encrypted,
            created_at,
            updated_at
            dek_encrypted
        FROM settings_encrypted
        WHERE key = $1;`
	err := r.DB.QueryRowContext(ctx, query, settings.CompanyPhoneKey).Scan(
		&res.ID,
		&res.ValueEncrypted,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.DEKEncrypted,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, rp.NewNotFoundErr(err, "res")
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return nil, rp.NewContextErr(err)
		default:
			return nil, rp.NewDatabaseErr(err)
		}
	}
	return &res, nil
}

func (r *repository) SetPhone(ctx context.Context, setting *settings.SettingEncrypted[string]) error {
	query := `
		INSERT INTO settings (id, key, value_encrypted, created_at, updated_at, dek_encrypted)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.DB.ExecContext(
		ctx,
		query,
		setting.ID,
		setting.Key,
		setting.ValueEncrypted,
		setting.CreatedAt,
		setting.UpdatedAt,
		setting.DEKEncrypted,
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
		return rp.NewNotCreatedErr(errors.New("no rows affected by insertion statement"), fmt.Sprintf("setting for key '%s'", setting.Key))
	}
	return nil
}
