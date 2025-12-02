package bookingHelpers

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearBookingsTable removes all test data from the bookings table
func ClearBookingsTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE booking.bookings RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

// EnsureBookingForeignKeys ensures all foreign key dependencies exist for a booking
func EnsureBookingForeignKeys(t *testing.T, ctx context.Context, pool *pgxpool.Pool, bookingEncx *domain.BookingEncx) {
	t.Helper()

	// Check if availability already exists
	var availExists bool
	err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM booking.availabilities WHERE id = $1)", bookingEncx.AvailabilityID).Scan(&availExists)
	require.NoError(t, err)

	if availExists {
		// Foreign keys already exist
		return
	}

	// Check if room already exists
	var roomExists bool
	err = pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM booking.rooms WHERE id = $1)", bookingEncx.RoomID).Scan(&roomExists)
	require.NoError(t, err)

	if !roomExists {
		// Create a stub building
		buildingID := uuid.New()
		_, err = pool.Exec(ctx, `
			INSERT INTO booking.buildings (id, name_encrypted, name_hash, address_encrypted, address_hash,
				city_encrypted, city_hash, postal_code_encrypted, country_encrypted, country_hash, dek_encrypted, key_version)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, buildingID, []byte("name"), "name_hash", []byte("address"), "address_hash",
			[]byte("city"), "city_hash", []byte("postal"), []byte("country"), "country_hash",
			[]byte("dek"), 1)
		require.NoError(t, err)

		// Create a stub room
		_, err = pool.Exec(ctx, `
			INSERT INTO booking.rooms (id, building_id, name_encrypted, name_hash, dek_encrypted, key_version)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, bookingEncx.RoomID, buildingID, []byte("room_name"), "room_hash", []byte("dek"), 1)
		require.NoError(t, err)
	}

	// Create a stub availability with future start time
	_, err = pool.Exec(ctx, `
		INSERT INTO booking.availabilities (id, user_id, room_id, start_time, end_time, dek_encrypted, key_version)
		VALUES ($1, $2, $3, NOW() + INTERVAL '1 hour', NOW() + INTERVAL '2 hours', $4, $5)
	`, bookingEncx.AvailabilityID, bookingEncx.PartnerID, bookingEncx.RoomID, []byte("dek"), 1)
	require.NoError(t, err)
}

// InsertBookingEncx inserts a booking into the database directly for testing
func InsertBookingEncx(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	bookingEncx *domain.BookingEncx,
) error {
	t.Helper()

	// Ensure foreign key dependencies exist
	EnsureBookingForeignKeys(t, ctx, pool, bookingEncx)

	query := `
		INSERT INTO booking.bookings (
			id, availability_id, client_id, user_id, room_id,
			product_id_encrypted, slot_start_time_encrypted, slot_end_time_encrypted,
			client_notes_encrypted, partner_notes_encrypted,
			total_price_cents, currency, payment_status, payment_intent_id,
			status, cancelled_at, cancellation_reason_encrypted, completed_at,
			created_at, updated_at,
			dek_encrypted, key_version, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
		)
	`

	_, err := pool.Exec(ctx, query,
		bookingEncx.ID,
		bookingEncx.AvailabilityID,
		bookingEncx.ClientID,
		bookingEncx.PartnerID,
		bookingEncx.RoomID,
		bookingEncx.ProductIDEncrypted,
		bookingEncx.SlotStartTimeEncrypted,
		bookingEncx.SlotEndTimeEncrypted,
		bookingEncx.ClientNotesEncrypted,
		bookingEncx.PartnerNotesEncrypted,
		bookingEncx.TotalPriceCents,
		bookingEncx.Currency,
		bookingEncx.PaymentStatus,
		bookingEncx.PaymentIntentID,
		bookingEncx.Status,
		bookingEncx.CancelledAt,
		bookingEncx.CancellationReasonEncrypted,
		bookingEncx.CompletedAt,
		bookingEncx.CreatedAt,
		bookingEncx.UpdatedAt,
		bookingEncx.DEKEncrypted,
		bookingEncx.KeyVersion,
		bookingEncx.Metadata,
	)

	return err
}

// GetBookingEncxByID retrieves a booking from the database by ID
func GetBookingEncxByID(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	bookingID uuid.UUID,
) (*domain.BookingEncx, error) {
	t.Helper()

	query := `
		SELECT
			id, availability_id, client_id, user_id, room_id,
			product_id_encrypted, slot_start_time_encrypted, slot_end_time_encrypted,
			client_notes_encrypted, partner_notes_encrypted,
			total_price_cents, currency, payment_status, payment_intent_id,
			status, cancelled_at, cancellation_reason_encrypted, completed_at,
			created_at, updated_at,
			dek_encrypted, key_version, metadata
		FROM booking.bookings
		WHERE id = $1
	`

	var bookingEncx domain.BookingEncx
	err := pool.QueryRow(ctx, query, bookingID).Scan(
		&bookingEncx.ID,
		&bookingEncx.AvailabilityID,
		&bookingEncx.ClientID,
		&bookingEncx.PartnerID,
		&bookingEncx.RoomID,
		&bookingEncx.ProductIDEncrypted,
		&bookingEncx.SlotStartTimeEncrypted,
		&bookingEncx.SlotEndTimeEncrypted,
		&bookingEncx.ClientNotesEncrypted,
		&bookingEncx.PartnerNotesEncrypted,
		&bookingEncx.TotalPriceCents,
		&bookingEncx.Currency,
		&bookingEncx.PaymentStatus,
		&bookingEncx.PaymentIntentID,
		&bookingEncx.Status,
		&bookingEncx.CancelledAt,
		&bookingEncx.CancellationReasonEncrypted,
		&bookingEncx.CompletedAt,
		&bookingEncx.CreatedAt,
		&bookingEncx.UpdatedAt,
		&bookingEncx.DEKEncrypted,
		&bookingEncx.KeyVersion,
		&bookingEncx.Metadata,
	)

	return &bookingEncx, err
}

// GetBookingsByAvailabilityID retrieves all bookings for a specific availability
func GetBookingsByAvailabilityID(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	availabilityID uuid.UUID,
) ([]*domain.BookingEncx, error) {
	t.Helper()

	query := `
		SELECT
			id, availability_id, client_id, user_id, room_id,
			product_id_encrypted, slot_start_time_encrypted, slot_end_time_encrypted,
			client_notes_encrypted, partner_notes_encrypted,
			total_price_cents, currency, payment_status, payment_intent_id,
			status, cancelled_at, cancellation_reason_encrypted, completed_at,
			created_at, updated_at,
			dek_encrypted, key_version, metadata
		FROM booking.bookings
		WHERE availability_id = $1
		ORDER BY created_at ASC
	`

	rows, err := pool.Query(ctx, query, availabilityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookingsEncx []*domain.BookingEncx
	for rows.Next() {
		var bookingEncx domain.BookingEncx
		err := rows.Scan(
			&bookingEncx.ID,
			&bookingEncx.AvailabilityID,
			&bookingEncx.ClientID,
			&bookingEncx.PartnerID,
			&bookingEncx.RoomID,
			&bookingEncx.ProductIDEncrypted,
			&bookingEncx.SlotStartTimeEncrypted,
			&bookingEncx.SlotEndTimeEncrypted,
			&bookingEncx.ClientNotesEncrypted,
			&bookingEncx.PartnerNotesEncrypted,
			&bookingEncx.TotalPriceCents,
			&bookingEncx.Currency,
			&bookingEncx.PaymentStatus,
			&bookingEncx.PaymentIntentID,
			&bookingEncx.Status,
			&bookingEncx.CancelledAt,
			&bookingEncx.CancellationReasonEncrypted,
			&bookingEncx.CompletedAt,
			&bookingEncx.CreatedAt,
			&bookingEncx.UpdatedAt,
			&bookingEncx.DEKEncrypted,
			&bookingEncx.KeyVersion,
			&bookingEncx.Metadata,
		)
		if err != nil {
			return nil, err
		}
		bookingsEncx = append(bookingsEncx, &bookingEncx)
	}

	return bookingsEncx, rows.Err()
}

// GetBookingsByClientID retrieves all bookings for a specific client
func GetBookingsByClientID(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	clientID uuid.UUID,
) ([]*domain.BookingEncx, error) {
	t.Helper()

	query := `
		SELECT
			id, availability_id, client_id, user_id, room_id,
			product_id_encrypted, slot_start_time_encrypted, slot_end_time_encrypted,
			client_notes_encrypted, partner_notes_encrypted,
			total_price_cents, currency, payment_status, payment_intent_id,
			status, cancelled_at, cancellation_reason_encrypted, completed_at,
			created_at, updated_at,
			dek_encrypted, key_version, metadata
		FROM booking.bookings
		WHERE client_id = $1
		ORDER BY created_at DESC
	`

	rows, err := pool.Query(ctx, query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookingsEncx []*domain.BookingEncx
	for rows.Next() {
		var bookingEncx domain.BookingEncx
		err := rows.Scan(
			&bookingEncx.ID,
			&bookingEncx.AvailabilityID,
			&bookingEncx.ClientID,
			&bookingEncx.PartnerID,
			&bookingEncx.RoomID,
			&bookingEncx.ProductIDEncrypted,
			&bookingEncx.SlotStartTimeEncrypted,
			&bookingEncx.SlotEndTimeEncrypted,
			&bookingEncx.ClientNotesEncrypted,
			&bookingEncx.PartnerNotesEncrypted,
			&bookingEncx.TotalPriceCents,
			&bookingEncx.Currency,
			&bookingEncx.PaymentStatus,
			&bookingEncx.PaymentIntentID,
			&bookingEncx.Status,
			&bookingEncx.CancelledAt,
			&bookingEncx.CancellationReasonEncrypted,
			&bookingEncx.CompletedAt,
			&bookingEncx.CreatedAt,
			&bookingEncx.UpdatedAt,
			&bookingEncx.DEKEncrypted,
			&bookingEncx.KeyVersion,
			&bookingEncx.Metadata,
		)
		if err != nil {
			return nil, err
		}
		bookingsEncx = append(bookingsEncx, &bookingEncx)
	}

	return bookingsEncx, rows.Err()
}

// GetBookingsByPartnerID retrieves all bookings for a specific partner
func GetBookingsByPartnerID(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	partnerID uuid.UUID,
) ([]*domain.BookingEncx, error) {
	t.Helper()

	query := `
		SELECT
			id, availability_id, client_id, user_id, room_id,
			product_id_encrypted, slot_start_time_encrypted, slot_end_time_encrypted,
			client_notes_encrypted, partner_notes_encrypted,
			total_price_cents, currency, payment_status, payment_intent_id,
			status, cancelled_at, cancellation_reason_encrypted, completed_at,
			created_at, updated_at,
			dek_encrypted, key_version, metadata
		FROM booking.bookings
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := pool.Query(ctx, query, partnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookingsEncx []*domain.BookingEncx
	for rows.Next() {
		var bookingEncx domain.BookingEncx
		err := rows.Scan(
			&bookingEncx.ID,
			&bookingEncx.AvailabilityID,
			&bookingEncx.ClientID,
			&bookingEncx.PartnerID,
			&bookingEncx.RoomID,
			&bookingEncx.ProductIDEncrypted,
			&bookingEncx.SlotStartTimeEncrypted,
			&bookingEncx.SlotEndTimeEncrypted,
			&bookingEncx.ClientNotesEncrypted,
			&bookingEncx.PartnerNotesEncrypted,
			&bookingEncx.TotalPriceCents,
			&bookingEncx.Currency,
			&bookingEncx.PaymentStatus,
			&bookingEncx.PaymentIntentID,
			&bookingEncx.Status,
			&bookingEncx.CancelledAt,
			&bookingEncx.CancellationReasonEncrypted,
			&bookingEncx.CompletedAt,
			&bookingEncx.CreatedAt,
			&bookingEncx.UpdatedAt,
			&bookingEncx.DEKEncrypted,
			&bookingEncx.KeyVersion,
			&bookingEncx.Metadata,
		)
		if err != nil {
			return nil, err
		}
		bookingsEncx = append(bookingsEncx, &bookingEncx)
	}

	return bookingsEncx, rows.Err()
}

// GetBookingByPaymentIntentID retrieves a booking by Stripe payment intent ID
func GetBookingByPaymentIntentID(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	paymentIntentID string,
) (*domain.BookingEncx, error) {
	t.Helper()

	query := `
		SELECT
			id, availability_id, client_id, user_id, room_id,
			product_id_encrypted, slot_start_time_encrypted, slot_end_time_encrypted,
			client_notes_encrypted, partner_notes_encrypted,
			total_price_cents, currency, payment_status, payment_intent_id,
			status, cancelled_at, cancellation_reason_encrypted, completed_at,
			created_at, updated_at,
			dek_encrypted, key_version, metadata
		FROM booking.bookings
		WHERE payment_intent_id = $1
	`

	var bookingEncx domain.BookingEncx
	err := pool.QueryRow(ctx, query, paymentIntentID).Scan(
		&bookingEncx.ID,
		&bookingEncx.AvailabilityID,
		&bookingEncx.ClientID,
		&bookingEncx.PartnerID,
		&bookingEncx.RoomID,
		&bookingEncx.ProductIDEncrypted,
		&bookingEncx.SlotStartTimeEncrypted,
		&bookingEncx.SlotEndTimeEncrypted,
		&bookingEncx.ClientNotesEncrypted,
		&bookingEncx.PartnerNotesEncrypted,
		&bookingEncx.TotalPriceCents,
		&bookingEncx.Currency,
		&bookingEncx.PaymentStatus,
		&bookingEncx.PaymentIntentID,
		&bookingEncx.Status,
		&bookingEncx.CancelledAt,
		&bookingEncx.CancellationReasonEncrypted,
		&bookingEncx.CompletedAt,
		&bookingEncx.CreatedAt,
		&bookingEncx.UpdatedAt,
		&bookingEncx.DEKEncrypted,
		&bookingEncx.KeyVersion,
		&bookingEncx.Metadata,
	)

	return &bookingEncx, err
}

// CountBookings counts all bookings in the database
func CountBookings(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
) (int, error) {
	t.Helper()

	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM booking.bookings").Scan(&count)
	return count, err
}

// CountBookingsByStatus counts bookings with a specific status
func CountBookingsByStatus(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	status domain.BookingStatus,
) (int, error) {
	t.Helper()

	var count int
	err := pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM booking.bookings WHERE status = $1",
		status,
	).Scan(&count)
	return count, err
}

// UpdateBookingStatus updates the status of a booking
func UpdateBookingStatus(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	bookingID uuid.UUID,
	status domain.BookingStatus,
) error {
	t.Helper()

	_, err := pool.Exec(ctx,
		"UPDATE booking.bookings SET status = $1, updated_at = NOW() WHERE id = $2",
		status, bookingID,
	)
	return err
}

// BookingExists checks if a booking exists in the database
func BookingExists(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	bookingID uuid.UUID,
) (bool, error) {
	t.Helper()

	var exists bool
	err := pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM booking.bookings WHERE id = $1)",
		bookingID,
	).Scan(&exists)
	return exists, err
}
