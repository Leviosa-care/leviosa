package roomRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.RoomEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, building_id, name_encrypted, name_hash, description_encrypted,
			room_number_encrypted, room_number_hash, capacity, equipment_encrypted,
			operating_start_time, operating_end_time,
			is_active, created_at, updated_at,
			dek_encrypted, key_version, metadata
		FROM %s.rooms
		WHERE id = $1
	`, r.schema)

	roomEncx := &domain.RoomEncx{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&roomEncx.ID,
		&roomEncx.BuildingID,
		&roomEncx.NameEncrypted,
		&roomEncx.NameHash,
		&roomEncx.DescriptionEncrypted,
		&roomEncx.RoomNumberEncrypted,
		&roomEncx.RoomNumberHash,
		&roomEncx.Capacity,
		&roomEncx.EquipmentEncrypted,
		&roomEncx.OperatingStartTime,
		&roomEncx.OperatingEndTime,
		&roomEncx.IsActive,
		&roomEncx.CreatedAt,
		&roomEncx.UpdatedAt,
		&roomEncx.DEKEncrypted,
		&roomEncx.KeyVersion,
		&roomEncx.Metadata,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get room by id", err)
	}

	return roomEncx, nil
}
