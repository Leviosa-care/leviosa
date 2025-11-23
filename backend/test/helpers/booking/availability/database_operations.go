package availabilityHelpers

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ClearAvailabilityTable removes all records from the availabilities table
func ClearAvailabilityTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	query := `DELETE FROM booking.availabilities`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		t.Fatalf("Failed to clear availabilities table: %v", err)
	}
}

// InsertAvailabilityEncx directly inserts an encrypted availability into the database
func InsertAvailabilityEncx(t *testing.T, ctx context.Context, availability *domain.AvailabilityEncx, pool *pgxpool.Pool) {
	query := `
		INSERT INTO booking.availabilities (
			id, user_id, room_id, start_time, end_time,
			service_type_encrypted, price_cents, max_capacity,
			notes_encrypted, is_recurring, recurrence_pattern,
			status, created_at, updated_at,
			dek_encrypted, key_version, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
	`

	_, err := pool.Exec(ctx, query,
		availability.ID,
		availability.UserID,
		availability.RoomID,
		availability.StartTime,
		availability.EndTime,
		availability.ServiceTypeEncrypted,
		availability.PriceCents,
		availability.MaxCapacity,
		availability.NotesEncrypted,
		availability.IsRecurring,
		availability.RecurrencePattern,
		availability.Status,
		availability.CreatedAt,
		availability.UpdatedAt,
		availability.DEKEncrypted,
		availability.KeyVersion,
		availability.Metadata,
	)

	if err != nil {
		t.Fatalf("Failed to insert availability: %v", err)
	}
}

// DeleteAvailabilityEncx removes an availability from the database by ID
func DeleteAvailabilityEncx(t *testing.T, ctx context.Context, availabilityID uuid.UUID, pool *pgxpool.Pool) {
	query := `DELETE FROM booking.availabilities WHERE id = $1`
	_, err := pool.Exec(ctx, query, availabilityID)
	if err != nil {
		t.Fatalf("Failed to delete availability: %v", err)
	}
}

// GetAvailabilityEncxFromDB retrieves an encrypted availability directly from the database
func GetAvailabilityEncxFromDB(t *testing.T, ctx context.Context, availabilityID uuid.UUID, pool *pgxpool.Pool) *domain.AvailabilityEncx {
	query := `
		SELECT
			id, user_id, room_id, start_time, end_time,
			service_type_encrypted, price_cents, max_capacity,
			notes_encrypted, is_recurring, recurrence_pattern,
			status, created_at, updated_at,
			dek_encrypted, key_version, metadata
		FROM booking.availabilities
		WHERE id = $1
	`

	availabilityEncx := &domain.AvailabilityEncx{}
	err := pool.QueryRow(ctx, query, availabilityID).Scan(
		&availabilityEncx.ID,
		&availabilityEncx.UserID,
		&availabilityEncx.RoomID,
		&availabilityEncx.StartTime,
		&availabilityEncx.EndTime,
		&availabilityEncx.ServiceTypeEncrypted,
		&availabilityEncx.PriceCents,
		&availabilityEncx.MaxCapacity,
		&availabilityEncx.NotesEncrypted,
		&availabilityEncx.IsRecurring,
		&availabilityEncx.RecurrencePattern,
		&availabilityEncx.Status,
		&availabilityEncx.CreatedAt,
		&availabilityEncx.UpdatedAt,
		&availabilityEncx.DEKEncrypted,
		&availabilityEncx.KeyVersion,
		&availabilityEncx.Metadata,
	)

	if err != nil {
		t.Fatalf("Failed to get availability from database: %v", err)
	}

	return availabilityEncx
}

// CountAvailabilitiesInTable counts the number of records in the availabilities table
func CountAvailabilitiesInTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) int {
	query := `SELECT COUNT(*) FROM booking.availabilities`
	var count int
	err := pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count availabilities: %v", err)
	}
	return count
}

// CountAvailabilitiesByPartnerID counts availabilities for a specific partner
func CountAvailabilitiesByPartnerID(t *testing.T, ctx context.Context, partnerID uuid.UUID, pool *pgxpool.Pool) int {
	query := `SELECT COUNT(*) FROM booking.availabilities WHERE user_id = $1`
	var count int
	err := pool.QueryRow(ctx, query, partnerID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count availabilities by partner ID: %v", err)
	}
	return count
}

// CountAvailabilitiesByRoomID counts availabilities for a specific room
func CountAvailabilitiesByRoomID(t *testing.T, ctx context.Context, roomID uuid.UUID, pool *pgxpool.Pool) int {
	query := `SELECT COUNT(*) FROM booking.availabilities WHERE room_id = $1`
	var count int
	err := pool.QueryRow(ctx, query, roomID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count availabilities by room ID: %v", err)
	}
	return count
}

// CountAvailabilitiesByStatus counts availabilities by status
func CountAvailabilitiesByStatus(t *testing.T, ctx context.Context, status domain.AvailabilityStatus, pool *pgxpool.Pool) int {
	query := `SELECT COUNT(*) FROM booking.availabilities WHERE status = $1`
	var count int
	err := pool.QueryRow(ctx, query, status).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count availabilities by status: %v", err)
	}
	return count
}

// CountRecurringAvailabilities counts recurring availabilities
func CountRecurringAvailabilities(t *testing.T, ctx context.Context, pool *pgxpool.Pool) int {
	query := `SELECT COUNT(*) FROM booking.availabilities WHERE is_recurring = true`
	var count int
	err := pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count recurring availabilities: %v", err)
	}
	return count
}

// AvailabilityExistsInTable checks if an availability exists in the database
func AvailabilityExistsInTable(t *testing.T, ctx context.Context, availabilityID uuid.UUID, pool *pgxpool.Pool) bool {
	query := `SELECT EXISTS(SELECT 1 FROM booking.availabilities WHERE id = $1)`
	var exists bool
	err := pool.QueryRow(ctx, query, availabilityID).Scan(&exists)
	if err != nil {
		t.Fatalf("Failed to check if availability exists: %v", err)
	}
	return exists
}
