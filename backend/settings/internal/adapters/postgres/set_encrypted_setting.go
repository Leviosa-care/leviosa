package postgres

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/Leviosa-care/core/errs"
)

func (r *repository) SetEncryptedSetting(ctx context.Context, setting *domain.SettingEncrypted[string]) error {
	query := `
	INSERT INTO settings.encrypted (key, value_encrypted, dek_encrypted, key_version)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (key) DO UPDATE SET
		value_encrypted = EXCLUDED.value_encrypted,
		dek_encrypted = EXCLUDED.dek_encrypted,
		key_version = EXCLUDED.key_version,
		updated_at = NOW()
	`
	commandTag, err := r.pool.Exec(
		ctx,
		query,
		setting.Key,
		setting.ValueEncrypted,
		setting.DEKEncrypted,
		setting.KeyVersion,
	)
	if err != nil {
		return errs.ClassifyPgError("set encrypted setting", err)
	}

	if commandTag.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(fmt.Errorf("failed to set encrypted setting"), setting.Key)
	}

	return nil
}
