package postgres

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *repository) SetEncryptedSetting(ctx context.Context, settingEncx *domain.SettingEncryptedEncx) error {
	query := `
	INSERT INTO settings.encrypted (key, value_encrypted, dek_encrypted, key_version, metadata)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (key) DO UPDATE SET
		value_encrypted = EXCLUDED.value_encrypted,
		dek_encrypted = EXCLUDED.dek_encrypted,
		key_version = EXCLUDED.key_version,
		metadata = EXCLUDED.metadata,
		updated_at = NOW()
	`
	commandTag, err := r.pool.Exec(
		ctx,
		query,
		settingEncx.Key,
		settingEncx.ValueEncrypted,
		settingEncx.DEKEncrypted,
		settingEncx.KeyVersion,
		settingEncx.Metadata,
	)
	if err != nil {
		return errs.ClassifyPgError("set encrypted setting", err)
	}

	if commandTag.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(fmt.Errorf("failed to set encrypted setting"), settingEncx.Key)
	}

	return nil
}
