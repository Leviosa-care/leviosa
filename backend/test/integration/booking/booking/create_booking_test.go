package booking

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBooking(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully create a booking with payment", func(t *testing.T) {
		// Clean state
		clearBookingTables(t, ctx)

		// Create test user (standard role to create booking)
		user, session := authCtx.CreateTestUser(t, ctx, "client@example.com", "John", "Doe", identity.Standard)

		// Create test building
		building := createTestBuilding(t, ctx, "Test Building", "123 Test St")

		// Create test room
		room := createTestRoom(t, ctx, building.ID, "Test Room", "101", 1)

		// Create test partner (practitioner)
		partner, _ := authCtx.CreateTestUser(t, ctx, "partner@example.com", "Jane", "Smith", identity.Partner)

		// Create test allocation (assign partner to room)
		allocation := createTestAllocation(t, ctx, room.ID, partner.ID, time.Now().Add(-24*time.Hour), nil)

		// Create test availability
		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour)
		endTime := startTime.Add(1 * time.Hour)
		availability := createTestAvailability(t, ctx, partner.ID, room.ID, allocation.ID, startTime, endTime, 5000)

		// Prepare request
		requestBody := map[string]interface{}{
			"availability_id": availability.ID.String(),
			"client_id":       user.ID,
			"client_notes":    "Looking forward to this session",
		}
		bodyBytes, err := json.Marshal(requestBody)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", testServerURL+"/bookings", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+session.Token)

		// Execute request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Parse response
		var booking domain.BookingResponse
		err = json.NewDecoder(resp.Body).Decode(&booking)
		require.NoError(t, err)

		// Assert booking details
		assert.NotEqual(t, uuid.Nil, booking.ID)
		assert.Equal(t, availability.ID, booking.AvailabilityID)
		assert.Equal(t, user.ID, booking.ClientID)
		assert.Equal(t, partner.ID, booking.PartnerID)
		assert.Equal(t, room.ID, booking.RoomID)
		assert.Equal(t, 5000, booking.TotalPriceCents)
		assert.Equal(t, "EUR", booking.Currency)
		assert.Equal(t, domain.BookingStatusConfirmed, booking.Status)
		assert.Equal(t, domain.PaymentStatusPending, booking.PaymentStatus)
		assert.NotNil(t, booking.PaymentIntentID)
		assert.NotEmpty(t, *booking.PaymentIntentID)

		// Verify booking exists in database
		savedBooking, err := bookingRepo.GetByID(ctx, booking.ID)
		require.NoError(t, err)
		assert.Equal(t, booking.ID, savedBooking.ID)

		// Verify availability is marked as booked
		savedAvailability, err := availabilityRepo.GetByID(ctx, availability.ID)
		require.NoError(t, err)
		assert.False(t, savedAvailability.IsAvailable)
		assert.True(t, savedAvailability.IsBooked)
	})

	t.Run("should create booking without payment when price is 0", func(t *testing.T) {
		// Clean state
		clearBookingTables(t, ctx)

		// Create test user
		user, session := authCtx.CreateTestUser(t, ctx, "client2@example.com", "Alice", "Johnson", identity.Standard)

		// Create test building, room, partner, allocation
		building := createTestBuilding(t, ctx, "Test Building 2", "456 Test Ave")
		room := createTestRoom(t, ctx, building.ID, "Free Room", "201", 1)
		partner, _ := authCtx.CreateTestUser(t, ctx, "partner2@example.com", "Bob", "Wilson", identity.Partner)
		allocation := createTestAllocation(t, ctx, room.ID, partner.ID, time.Now().Add(-24*time.Hour), nil)

		// Create availability with 0 price
		startTime := time.Now().Add(48 * time.Hour).Truncate(time.Hour)
		endTime := startTime.Add(1 * time.Hour)
		availability := createTestAvailability(t, ctx, partner.ID, room.ID, allocation.ID, startTime, endTime, 0)

		// Prepare request
		requestBody := map[string]interface{}{
			"availability_id": availability.ID.String(),
			"client_id":       user.ID,
		}
		bodyBytes, err := json.Marshal(requestBody)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", testServerURL+"/bookings", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+session.Token)

		// Execute request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Parse response
		var booking domain.BookingResponse
		err = json.NewDecoder(resp.Body).Decode(&booking)
		require.NoError(t, err)

		// Assert no payment intent was created
		assert.Nil(t, booking.PaymentIntentID)
		assert.Equal(t, domain.PaymentStatusNotRequired, booking.PaymentStatus)
	})

	t.Run("should fail when availability is already booked", func(t *testing.T) {
		// Clean state
		clearBookingTables(t, ctx)

		// Create test data
		user1, session1 := authCtx.CreateTestUser(t, ctx, "client3@example.com", "User", "One", identity.Standard)
		user2, session2 := authCtx.CreateTestUser(t, ctx, "client4@example.com", "User", "Two", identity.Standard)

		building := createTestBuilding(t, ctx, "Test Building 3", "789 Test Blvd")
		room := createTestRoom(t, ctx, building.ID, "Shared Room", "301", 1)
		partner, _ := authCtx.CreateTestUser(t, ctx, "partner3@example.com", "Partner", "Three", identity.Partner)
		allocation := createTestAllocation(t, ctx, room.ID, partner.ID, time.Now().Add(-24*time.Hour), nil)

		startTime := time.Now().Add(72 * time.Hour).Truncate(time.Hour)
		endTime := startTime.Add(1 * time.Hour)
		availability := createTestAvailability(t, ctx, partner.ID, room.ID, allocation.ID, startTime, endTime, 3000)

		// First user creates booking successfully
		requestBody1 := map[string]interface{}{
			"availability_id": availability.ID.String(),
			"client_id":       user1.ID,
		}
		bodyBytes1, _ := json.Marshal(requestBody1)

		req1, _ := http.NewRequest("POST", testServerURL+"/bookings", bytes.NewReader(bodyBytes1))
		req1.Header.Set("Content-Type", "application/json")
		req1.Header.Set("Authorization", "Bearer "+session1.Token)

		resp1, err := client.Do(req1)
		require.NoError(t, err)
		defer resp1.Body.Close()
		assert.Equal(t, http.StatusCreated, resp1.StatusCode)

		// Second user tries to book same availability - should fail
		requestBody2 := map[string]interface{}{
			"availability_id": availability.ID.String(),
			"client_id":       user2.ID,
		}
		bodyBytes2, _ := json.Marshal(requestBody2)

		req2, _ := http.NewRequest("POST", testServerURL+"/bookings", bytes.NewReader(bodyBytes2))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Authorization", "Bearer "+session2.Token)

		resp2, err := client.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()

		// Assert conflict or bad request
		assert.Contains(t, []int{http.StatusConflict, http.StatusBadRequest}, resp2.StatusCode)
	})
}

// Helper functions for test data creation

func clearBookingTables(t *testing.T, ctx context.Context) {
	t.Helper()

	// Clear in correct order due to foreign keys
	_, err := testPool.Exec(ctx, "TRUNCATE TABLE bookingschema.bookings CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE TABLE bookingschema.availabilities CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE TABLE bookingschema.room_allocations CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE TABLE bookingschema.rooms CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE TABLE bookingschema.buildings CASCADE")
	require.NoError(t, err)
}

func createTestBuilding(t *testing.T, ctx context.Context, name, address string) *domain.Building {
	t.Helper()

	// Encrypt building data
	nameEnc, nameHash, err := crypto.Encrypt(ctx, []byte(name))
	require.NoError(t, err)

	addressEnc, addressHash, err := crypto.Encrypt(ctx, []byte(address))
	require.NoError(t, err)

	buildingID := uuid.New()
	query := `
		INSERT INTO bookingschema.buildings (
			id, name_encrypted, name_hash, address_encrypted, address_hash,
			is_active, created_at, updated_at,
			dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	now := time.Now()
	_, err = testPool.Exec(ctx, query,
		buildingID,
		nameEnc.Ciphertext,
		nameHash,
		addressEnc.Ciphertext,
		addressHash,
		true,
		now,
		now,
		nameEnc.DEK,
		nameEnc.KeyVersion,
		nameEnc.Metadata,
	)
	require.NoError(t, err)

	return &domain.Building{
		ID:        buildingID,
		Name:      name,
		Address:   address,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func createTestRoom(t *testing.T, ctx context.Context, buildingID uuid.UUID, name, roomNumber string, capacity int) *domain.Room {
	t.Helper()

	// Encrypt room data
	nameEnc, nameHash, err := crypto.Encrypt(ctx, []byte(name))
	require.NoError(t, err)

	roomNumEnc, roomNumHash, err := crypto.Encrypt(ctx, []byte(roomNumber))
	require.NoError(t, err)

	roomID := uuid.New()
	query := `
		INSERT INTO bookingschema.rooms (
			id, building_id, name_encrypted, name_hash,
			room_number_encrypted, room_number_hash, capacity,
			operating_start_time, operating_end_time,
			is_active, created_at, updated_at,
			dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	now := time.Now()
	_, err = testPool.Exec(ctx, query,
		roomID,
		buildingID,
		nameEnc.Ciphertext,
		nameHash,
		roomNumEnc.Ciphertext,
		roomNumHash,
		capacity,
		time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC), // 08:00
		time.Date(0, 1, 1, 20, 0, 0, 0, time.UTC), // 20:00
		true,
		now,
		now,
		nameEnc.DEK,
		nameEnc.KeyVersion,
		nameEnc.Metadata,
	)
	require.NoError(t, err)

	return &domain.Room{
		ID:                 roomID,
		BuildingID:         buildingID,
		Name:               name,
		RoomNumber:         roomNumber,
		Capacity:           capacity,
		OperatingStartTime: time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC),
		OperatingEndTime:   time.Date(0, 1, 1, 20, 0, 0, 0, time.UTC),
		IsActive:           true,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

func createTestAllocation(t *testing.T, ctx context.Context, roomID, partnerID uuid.UUID, startTime time.Time, endTime *time.Time) *domain.RoomAllocation {
	t.Helper()

	allocationID := uuid.New()
	query := `
		INSERT INTO bookingschema.room_allocations (
			id, room_id, partner_id, allocation_type, start_time, end_time,
			is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	now := time.Now()
	allocationType := "shared"
	if endTime != nil {
		allocationType = "dedicated"
	}

	_, err := testPool.Exec(ctx, query,
		allocationID,
		roomID,
		partnerID,
		allocationType,
		startTime,
		endTime,
		true,
		now,
		now,
	)
	require.NoError(t, err)

	return &domain.RoomAllocation{
		ID:             allocationID,
		RoomID:         roomID,
		PartnerID:      partnerID,
		AllocationType: domain.AllocationTypeFromString(allocationType),
		StartTime:      startTime,
		EndTime:        endTime,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func createTestAvailability(t *testing.T, ctx context.Context, partnerID, roomID, allocationID uuid.UUID, startTime, endTime time.Time, priceCents int) *domain.Availability {
	t.Helper()

	availabilityID := uuid.New()
	query := `
		INSERT INTO bookingschema.availabilities (
			id, partner_id, room_id, allocation_id, start_time, end_time,
			price_cents, is_available, is_booked, is_cancelled,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	now := time.Now()
	_, err := testPool.Exec(ctx, query,
		availabilityID,
		partnerID,
		roomID,
		allocationID,
		startTime,
		endTime,
		priceCents,
		true,  // is_available
		false, // is_booked
		false, // is_cancelled
		now,
		now,
	)
	require.NoError(t, err)

	return &domain.Availability{
		ID:           availabilityID,
		PartnerID:    partnerID,
		RoomID:       roomID,
		AllocationID: allocationID,
		StartTime:    startTime,
		EndTime:      endTime,
		PriceCents:   &priceCents,
		IsAvailable:  true,
		IsBooked:     false,
		IsCancelled:  false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
