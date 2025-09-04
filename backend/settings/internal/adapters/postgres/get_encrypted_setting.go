package postgres

import (
	"context"
	"fmt"

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
	dek_encrypted,
	key_version
	FROM settings.encrypted
	WHERE key = $1;`

	err := r.pool.QueryRow(ctx, query, key).Scan(
		&res.ID,
		&res.ValueEncrypted,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.DEKEncrypted,
		&res.KeyVersion,
	)
	if err != nil {
		return nil, errs.ClassifyPgError(fmt.Sprintf("get encrypted setting for key '%s'", key), err)
	}

	res.Key = key

	return &res, nil
}
