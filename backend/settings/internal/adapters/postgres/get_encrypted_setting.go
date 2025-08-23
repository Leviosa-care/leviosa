package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/Leviosa-care/core/errs"
)

func (r *repository) GetEncryptedSetting(ctx context.Context, key string) (*domain.SettingEncrypted[string], error) {
	var res domain.SettingEncrypted[string]
	query := `
	SELECT
	id,
	value_encrypted,
	created_at,
	updated_at,
	dek_encrypted
	FROM settings.encrypted
	WHERE key = $1;`

	err := r.pool.QueryRow(ctx, query, key).Scan(
		&res.ID,
		&res.ValueEncrypted,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.DEKEncrypted,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundErr(err, "encrypted setting")
		}
		return nil, errs.ClassifyPgError("get encrypted setting", err)
	}

	res.Key = key

	return &res, nil
}
