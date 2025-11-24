package availability_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	catalogDomain "github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAvailability_DurationValidation(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should succeed with valid duration for single product in shared room", func(t *testing.T) {
		// Clean state
		clearTestData(t, ctx)

		// Create partner user
		partner, session := authCtx.CreateTestUser(t, ctx, "partner1@example.com", "Partner", "One", identity.Partner)

		// Create building and room
		building := createTestBuilding(t, ctx, "Building A", "123 Main St")
		room := createTestRoom(t, ctx, building.ID, "Shared Room 1", "101", 2)

		// Create SHARED allocation
		allocation := createTestSharedAllocation(t, ctx, room.ID, partner.ID)

		// Create product: 30 min service + 10 min buffer = 40 min total
		product := createTestProduct(t, ctx, "Swedish Massage", 30, 10)

		// Create 80-minute availability (2 × 40min sessions) - VALID
		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Minute)
		endTime := startTime.Add(80 * time.Minute)

		requestBody := map[string]interface{}{
			"room_id":    room.ID.String(),
			"start_time": startTime.Format(time.RFC3339),
			"end_time":   endTime.Format(time.RFC3339),
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", testServerURL+"/availabilities", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+session.Token)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should succeed - 80min = 2 sessions of 40min
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("should fail with invalid duration for single product in shared room", func(t *testing.T) {
		// Clean state
		clearTestData(t, ctx)

		// Create partner user
		partner, session := authCtx.CreateTestUser(t, ctx, "partner2@example.com", "Partner", "Two", identity.Partner)

		// Create building and room
		building := createTestBuilding(t, ctx, "Building B", "456 Main St")
		room := createTestRoom(t, ctx, building.ID, "Shared Room 2", "102", 2)

		// Create SHARED allocation
		allocation := createTestSharedAllocation(t, ctx, room.ID, partner.ID)

		// Create product: 30 min service + 10 min buffer = 40 min total
		product := createTestProduct(t, ctx, "Swedish Massage", 30, 10)

		// Create 65-minute availability - INVALID (not a multiple of 40)
		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Minute)
		endTime := startTime.Add(65 * time.Minute)

		requestBody := map[string]interface{}{
			"room_id":    room.ID.String(),
			"start_time": startTime.Format(time.RFC3339),
			"end_time":   endTime.Format(time.RFC3339),
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", testServerURL+"/availabilities", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+session.Token)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should fail with validation error
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Parse error response
		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)

		// Verify structured error response
		assert.Equal(t, "Availability duration does not align with product offerings", errorResponse["error"])
		assert.Equal(t, float64(65), errorResponse["requested_duration"])
		assert.NotEmpty(t, errorResponse["suggested_durations"])

		// Verify suggestions include valid durations
		suggestions := errorResponse["suggested_durations"].([]interface{})
		assert.NotEmpty(t, suggestions)
	})

	t.Run("should succeed with valid duration for multiple products in shared room", func(t *testing.T) {
		// Clean state
		clearTestData(t, ctx)

		// Create partner user
		partner, session := authCtx.CreateTestUser(t, ctx, "partner3@example.com", "Partner", "Three", identity.Partner)

		// Create building and room
		building := createTestBuilding(t, ctx, "Building C", "789 Main St")
		room := createTestRoom(t, ctx, building.ID, "Shared Room 3", "103", 2)

		// Create SHARED allocation
		allocation := createTestSharedAllocation(t, ctx, room.ID, partner.ID)

		// Create two products:
		// Product 1: 30 min + 10 min = 40 min total
		// Product 2: 45 min + 15 min = 60 min total
		product1 := createTestProduct(t, ctx, "Swedish Massage", 30, 10)
		product2 := createTestProduct(t, ctx, "Deep Tissue Massage", 45, 15)

		// Create 120-minute availability - VALID (3×40min OR 2×60min)
		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Minute)
		endTime := startTime.Add(120 * time.Minute)

		requestBody := map[string]interface{}{
			"room_id":    room.ID.String(),
			"start_time": startTime.Format(time.RFC3339),
			"end_time":   endTime.Format(time.RFC3339),
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", testServerURL+"/availabilities", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+session.Token)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should succeed - 120min works for both products
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("should succeed with any duration for dedicated room", func(t *testing.T) {
		// Clean state
		clearTestData(t, ctx)

		// Create partner user
		partner, session := authCtx.CreateTestUser(t, ctx, "partner4@example.com", "Partner", "Four", identity.Partner)

		// Create building and room
		building := createTestBuilding(t, ctx, "Building D", "321 Main St")
		room := createTestRoom(t, ctx, building.ID, "Dedicated Room", "201", 1)

		// Create DEDICATED allocation (time-bounded)
		startDate := time.Now().Add(-24 * time.Hour)
		endDate := time.Now().Add(7 * 24 * time.Hour)
		allocation := createTestDedicatedAllocation(t, ctx, room.ID, partner.ID, &startDate, &endDate)

		// Create product
		product := createTestProduct(t, ctx, "Private Session", 30, 10)

		// Create 73-minute availability - odd duration, but should work for dedicated
		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Minute)
		endTime := startTime.Add(73 * time.Minute)

		requestBody := map[string]interface{}{
			"room_id":    room.ID.String(),
			"start_time": startTime.Format(time.RFC3339),
			"end_time":   endTime.Format(time.RFC3339),
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", testServerURL+"/availabilities", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+session.Token)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should succeed - dedicated rooms bypass validation
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("should handle product with buffer time correctly", func(t *testing.T) {
		// Clean state
		clearTestData(t, ctx)

		// Create partner user
		partner, session := authCtx.CreateTestUser(t, ctx, "partner5@example.com", "Partner", "Five", identity.Partner)

		// Create building and room
		building := createTestBuilding(t, ctx, "Building E", "654 Main St")
		room := createTestRoom(t, ctx, building.ID, "Shared Room 4", "104", 2)

		// Create SHARED allocation
		allocation := createTestSharedAllocation(t, ctx, room.ID, partner.ID)

		// Create product: 50 min service + 20 min buffer = 70 min total
		product := createTestProduct(t, ctx, "Extended Massage", 50, 20)

		// Create 140-minute availability (2 × 70min) - VALID
		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Minute)
		endTime := startTime.Add(140 * time.Minute)

		requestBody := map[string]interface{}{
			"room_id":    room.ID.String(),
			"start_time": startTime.Format(time.RFC3339),
			"end_time":   endTime.Format(time.RFC3339),
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", testServerURL+"/availabilities", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+session.Token)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should succeed - buffer time is included in calculation
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})
}

// Helper functions

func clearTestData(t *testing.T, ctx context.Context) {
	t.Helper()

	// Clear in correct order due to foreign keys
	_, err := testPool.Exec(ctx, "TRUNCATE TABLE bookingschema.availabilities CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE TABLE bookingschema.room_allocations CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE TABLE bookingschema.rooms CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE TABLE bookingschema.buildings CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE TABLE catalogschema.products CASCADE")
	require.NoError(t, err)
}

func createTestSharedAllocation(t *testing.T, ctx context.Context, roomID, partnerID uuid.UUID) *domain.RoomAllocation {
	t.Helper()

	allocationID := uuid.New()
	query := `
		INSERT INTO bookingschema.room_allocations (
			id, room_id, partner_id, allocation_type, start_time, end_time,
			is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, NULL, NULL, $5, $6, $7)
	`

	now := time.Now()
	_, err := testPool.Exec(ctx, query,
		allocationID,
		roomID,
		partnerID,
		"shared",
		true,
		now,
		now,
	)
	require.NoError(t, err)

	return &domain.RoomAllocation{
		ID:             allocationID,
		RoomID:         roomID,
		PartnerID:      partnerID,
		AllocationType: domain.AllocationTypeShared,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func createTestDedicatedAllocation(t *testing.T, ctx context.Context, roomID, partnerID uuid.UUID, startDate, endDate *time.Time) *domain.RoomAllocation {
	t.Helper()

	allocationID := uuid.New()
	query := `
		INSERT INTO bookingschema.room_allocations (
			id, room_id, partner_id, allocation_type, start_time, end_time,
			is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	now := time.Now()
	_, err := testPool.Exec(ctx, query,
		allocationID,
		roomID,
		partnerID,
		"dedicated",
		startDate,
		endDate,
		true,
		now,
		now,
	)
	require.NoError(t, err)

	return &domain.RoomAllocation{
		ID:             allocationID,
		RoomID:         roomID,
		PartnerID:      partnerID,
		AllocationType: domain.AllocationTypeDedicated,
		StartDate:      startDate,
		EndDate:        endDate,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func createTestProduct(t *testing.T, ctx context.Context, name string, durationMinutes, bufferMinutes int) *catalogDomain.Product {
	t.Helper()

	productID := uuid.New()
	categoryID := uuid.New() // Simplified - in real tests you'd create a category

	query := `
		INSERT INTO catalogschema.products (
			id, name, description, category_id, duration, buffer_time,
			published_status, availability, cancellation_hours,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	now := time.Now()
	_, err := testPool.Exec(ctx, query,
		productID,
		name,
		"Test product description",
		categoryID,
		durationMinutes,
		bufferMinutes,
		"published", // published status
		"available", // availability
		24,          // cancellation hours
		now,
		now,
	)
	require.NoError(t, err)

	return &catalogDomain.Product{
		ID:                productID,
		Name:              name,
		Duration:          durationMinutes,
		BufferTime:        bufferMinutes,
		Status:            catalogDomain.PublishedStatusPublished,
		Availability:      catalogDomain.AvailabilityTypeAvailable,
		CancellationHours: 24,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}
