package roomHelpers

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearRoomsTable removes all test data from the rooms table
func ClearRoomsTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE booking.rooms RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

// GetRoomEncxByID retrieves a room from the database by ID
func GetRoomEncxByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, roomID interface{}) (*domain.RoomEncx, error) {
	t.Helper()

	query := `
		SELECT id, building_id, name_encrypted, name_hash, description_encrypted,
		       room_number_encrypted, room_number_hash, capacity, equipment_encrypted,
		       is_active, created_at, updated_at,
		       dek_encrypted, key_version, metadata
		FROM booking.rooms
		WHERE id = $1
	`

	var roomEncx domain.RoomEncx
	err := pool.QueryRow(ctx, query, roomID).Scan(
		&roomEncx.ID,
		&roomEncx.BuildingID,
		&roomEncx.NameEncrypted,
		&roomEncx.NameHash,
		&roomEncx.DescriptionEncrypted,
		&roomEncx.RoomNumberEncrypted,
		&roomEncx.RoomNumberHash,
		&roomEncx.Capacity,
		&roomEncx.EquipmentEncrypted,
		&roomEncx.IsActive,
		&roomEncx.CreatedAt,
		&roomEncx.UpdatedAt,
		&roomEncx.DEKEncrypted,
		&roomEncx.KeyVersion,
		&roomEncx.Metadata,
	)

	return &roomEncx, err
}

// InsertRoomEncx inserts a room into the database directly for testing
func InsertRoomEncx(t *testing.T, ctx context.Context, pool *pgxpool.Pool, roomEncx *domain.RoomEncx) error {
	t.Helper()
	_, err := pool.Exec(ctx, `
			INSERT INTO booking.rooms (
				id, building_id, name_encrypted, name_hash, description_encrypted,
				room_number_encrypted, room_number_hash, capacity, equipment_encrypted,
				is_active, created_at, updated_at,
				dek_encrypted, key_version, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		`,
		roomEncx.ID,
		roomEncx.BuildingID,
		roomEncx.NameEncrypted,
		roomEncx.NameHash,
		roomEncx.DescriptionEncrypted,
		roomEncx.RoomNumberEncrypted,
		roomEncx.RoomNumberHash,
		roomEncx.Capacity,
		roomEncx.EquipmentEncrypted,
		roomEncx.IsActive,
		roomEncx.CreatedAt,
		roomEncx.UpdatedAt,
		roomEncx.DEKEncrypted,
		roomEncx.KeyVersion,
		roomEncx.Metadata,
	)
	return err
}

// GetRoomsByBuildingID retrieves all rooms for a specific building
func GetRoomsByBuildingID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, buildingID interface{}, activeOnly bool) ([]*domain.RoomEncx, error) {
	t.Helper()

	query := `
		SELECT id, building_id, name_encrypted, name_hash, description_encrypted,
		       room_number_encrypted, room_number_hash, capacity, equipment_encrypted,
		       is_active, created_at, updated_at,
		       dek_encrypted, key_version, metadata
		FROM booking.rooms
		WHERE building_id = $1
	`

	args := []interface{}{buildingID}

	if activeOnly {
		query += " AND is_active = $2"
		args = append(args, true)
	}

	query += " ORDER BY created_at"

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*domain.RoomEncx
	for rows.Next() {
		var roomEncx domain.RoomEncx
		err := rows.Scan(
			&roomEncx.ID,
			&roomEncx.BuildingID,
			&roomEncx.NameEncrypted,
			&roomEncx.NameHash,
			&roomEncx.DescriptionEncrypted,
			&roomEncx.RoomNumberEncrypted,
			&roomEncx.RoomNumberHash,
			&roomEncx.Capacity,
			&roomEncx.EquipmentEncrypted,
			&roomEncx.IsActive,
			&roomEncx.CreatedAt,
			&roomEncx.UpdatedAt,
			&roomEncx.DEKEncrypted,
			&roomEncx.KeyVersion,
			&roomEncx.Metadata,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, &roomEncx)
	}

	return rooms, rows.Err()
}

// CountRoomsByFilter counts rooms matching the given filter criteria
func CountRoomsByFilter(t *testing.T, ctx context.Context, pool *pgxpool.Pool, filter struct {
	BuildingID     interface{}
	IsActive       *bool
	MinCapacity    *int
	MaxCapacity    *int
	MinHourlyRate  *int
	MaxHourlyRate  *int
	NameHash       *string
	RoomNumberHash *string
}) (int, error) {
	t.Helper()

	query := `SELECT COUNT(*) FROM booking.rooms WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filter.BuildingID != nil {
		query += ` AND building_id = $` + string(rune('0'+argIndex))
		args = append(args, filter.BuildingID)
		argIndex++
	}

	if filter.IsActive != nil {
		query += ` AND is_active = $` + string(rune('0'+argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}

	if filter.MinCapacity != nil {
		query += ` AND capacity >= $` + string(rune('0'+argIndex))
		args = append(args, *filter.MinCapacity)
		argIndex++
	}

	if filter.MaxCapacity != nil {
		query += ` AND capacity <= $` + string(rune('0'+argIndex))
		args = append(args, *filter.MaxCapacity)
		argIndex++
	}

	if filter.MinHourlyRate != nil {
		query += ` AND hourly_rate_cents >= $` + string(rune('0'+argIndex))
		args = append(args, *filter.MinHourlyRate)
		argIndex++
	}

	if filter.MaxHourlyRate != nil {
		query += ` AND hourly_rate_cents <= $` + string(rune('0'+argIndex))
		args = append(args, *filter.MaxHourlyRate)
		argIndex++
	}

	if filter.NameHash != nil {
		query += ` AND name_hash = $` + string(rune('0'+argIndex))
		args = append(args, *filter.NameHash)
		argIndex++
	}

	if filter.RoomNumberHash != nil {
		query += ` AND room_number_hash = $` + string(rune('0'+argIndex))
		args = append(args, *filter.RoomNumberHash)
		argIndex++
	}

	var count int
	err := pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// UpdateRoomIsActive updates the active status of a room
func UpdateRoomIsActive(t *testing.T, ctx context.Context, pool *pgxpool.Pool, roomID interface{}, isActive bool) error {
	t.Helper()
	_, err := pool.Exec(ctx, "UPDATE booking.rooms SET is_active = $1, updated_at = NOW() WHERE id = $2", isActive, roomID)
	return err
}

// RoomExistsByBuildingAndNumber checks if a room exists for a building with a specific room number hash
func RoomExistsByBuildingAndNumber(t *testing.T, ctx context.Context, pool *pgxpool.Pool, buildingID interface{}, roomNumberHash string) (bool, error) {
	t.Helper()
	var exists bool
	err := pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM booking.rooms WHERE building_id = $1 AND room_number_hash = $2)",
		buildingID, roomNumberHash).Scan(&exists)
	return exists, err
}
