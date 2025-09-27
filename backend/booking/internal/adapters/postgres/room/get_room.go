package roomRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	query := fmt.Sprintf(`
		SELECT
			id, building_id, name_encrypted, description_encrypted,
			room_number_encrypted, capacity, equipment_encrypted,
			hourly_rate_cents, is_active, created_at, updated_at
		FROM %s.rooms
		WHERE id = $1
	`, r.schema)

	room := &domain.Room{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&room.ID,
		&room.BuildingID,
		&room.NameEncrypted,
		&room.DescriptionEncrypted,
		&room.RoomNumberEncrypted,
		&room.Capacity,
		&room.EquipmentEncrypted,
		&room.HourlyRateCents,
		&room.IsActive,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get room by id", err)
	}

	// Decrypt sensitive fields
	if err := r.crypto.DecryptStruct(ctx, room); err != nil {
		return nil, fmt.Errorf("decrypt room data: %w", err)
	}

	return room, nil
}

func (r *Repository) GetByBuildingID(ctx context.Context, buildingID uuid.UUID, activeOnly bool) ([]*domain.Room, error) {
	query := fmt.Sprintf(`
		SELECT
			id, building_id, name_encrypted, description_encrypted,
			room_number_encrypted, capacity, equipment_encrypted,
			hourly_rate_cents, is_active, created_at, updated_at
		FROM %s.rooms
		WHERE building_id = $1
	`, r.schema)

	args := []interface{}{buildingID}
	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY name_encrypted ASC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errs.ClassifyPgError("get rooms by building id", err)
	}
	defer rows.Close()

	var rooms []*domain.Room
	for rows.Next() {
		room := &domain.Room{}
		err := rows.Scan(
			&room.ID,
			&room.BuildingID,
			&room.NameEncrypted,
			&room.DescriptionEncrypted,
			&room.RoomNumberEncrypted,
			&room.Capacity,
			&room.EquipmentEncrypted,
			&room.HourlyRateCents,
			&room.IsActive,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan room row", err)
		}

		// Decrypt sensitive fields
		if err := r.crypto.DecryptStruct(ctx, room); err != nil {
			return nil, fmt.Errorf("decrypt room data: %w", err)
		}

		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate room rows", err)
	}

	return rooms, nil
}