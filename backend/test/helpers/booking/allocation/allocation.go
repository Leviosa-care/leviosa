package allocationHelpers

import (
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// NewTestSharedAllocation creates a shared room allocation for testing
func NewTestSharedAllocation(t *testing.T, roomID, userID uuid.UUID) *domain.RoomAllocation {
	t.Helper()

	allocation, err := domain.NewSharedAllocation(roomID, userID)
	require.NoError(t, err)

	return allocation
}

// NewTestDedicatedAllocation creates a dedicated room allocation with specific dates
func NewTestDedicatedAllocation(t *testing.T, roomID, userID uuid.UUID, startDate, endDate time.Time) *domain.RoomAllocation {
	t.Helper()

	allocation, err := domain.NewDedicatedAllocation(roomID, userID, &startDate, &endDate)
	require.NoError(t, err)

	return allocation
}

// NewTestDedicatedAllocationIndefinite creates a dedicated allocation without an end date
func NewTestDedicatedAllocationIndefinite(t *testing.T, roomID, userID uuid.UUID, startDate time.Time) *domain.RoomAllocation {
	t.Helper()

	// Pass nil for endDate to create NULL end_date in database
	allocation, err := domain.NewDedicatedAllocation(roomID, userID, &startDate, nil)
	require.NoError(t, err)

	return allocation
}

// NewTestPastDedicatedAllocation creates a dedicated allocation that has already ended
func NewTestPastDedicatedAllocation(t *testing.T, roomID, userID uuid.UUID) *domain.RoomAllocation {
	t.Helper()

	// Create allocation from 30 days ago to 15 days ago
	startDate := time.Now().AddDate(0, 0, -30).Truncate(24 * time.Hour)
	endDate := time.Now().AddDate(0, 0, -15).Truncate(24 * time.Hour)

	allocation, err := domain.NewDedicatedAllocation(roomID, userID, &startDate, &endDate)
	require.NoError(t, err)

	return allocation
}

// NewTestFutureDedicatedAllocation creates a dedicated allocation that starts in the future
func NewTestFutureDedicatedAllocation(t *testing.T, roomID, userID uuid.UUID) *domain.RoomAllocation {
	t.Helper()

	// Create allocation starting 15 days from now, ending 45 days from now
	startDate := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
	endDate := time.Now().AddDate(0, 0, 45).Truncate(24 * time.Hour)

	allocation, err := domain.NewDedicatedAllocation(roomID, userID, &startDate, &endDate)
	require.NoError(t, err)

	return allocation
}

// NewTestActiveDedicatedAllocation creates a dedicated allocation that is currently active
func NewTestActiveDedicatedAllocation(t *testing.T, roomID, userID uuid.UUID) *domain.RoomAllocation {
	t.Helper()

	// Create allocation starting 15 days ago, ending 15 days from now
	startDate := time.Now().AddDate(0, 0, -15).Truncate(24 * time.Hour)
	endDate := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)

	allocation, err := domain.NewDedicatedAllocation(roomID, userID, &startDate, &endDate)
	require.NoError(t, err)

	return allocation
}

// NewTestInactiveAllocation creates an inactive allocation (for soft delete tests)
func NewTestInactiveAllocation(t *testing.T, roomID, userID uuid.UUID) *domain.RoomAllocation {
	t.Helper()

	allocation, err := domain.NewSharedAllocation(roomID, userID)
	require.NoError(t, err)

	allocation.Deactivate()

	return allocation
}

// NewTestDedicatedAllocationWithID creates a dedicated allocation with a specific ID (for update tests)
func NewTestDedicatedAllocationWithID(t *testing.T, id, roomID, userID uuid.UUID, startDate, endDate time.Time) *domain.RoomAllocation {
	t.Helper()

	allocation, err := domain.NewDedicatedAllocation(roomID, userID, &startDate, &endDate)
	require.NoError(t, err)

	allocation.ID = id

	return allocation
}

// NewTestSharedAllocationWithID creates a shared allocation with a specific ID (for update tests)
func NewTestSharedAllocationWithID(t *testing.T, id, roomID, userID uuid.UUID) *domain.RoomAllocation {
	t.Helper()

	allocation, err := domain.NewSharedAllocation(roomID, userID)
	require.NoError(t, err)

	allocation.ID = id

	return allocation
}
