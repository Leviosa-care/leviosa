package allocationRepository

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/require"
)

// NewTestSharedAllocationEncx creates a pre-encrypted shared allocation for repository tests
func NewTestSharedAllocationEncx(t *testing.T, crypto encx.CryptoService, roomID, userID uuid.UUID) *domain.RoomAllocationEncx {
	t.Helper()

	ctx := context.Background()

	// Create domain entity
	allocation, err := domain.NewSharedAllocation(roomID, userID)
	require.NoError(t, err, "failed to create shared allocation")

	// Encrypt
	allocationEncx, err := domain.ProcessRoomAllocationEncx(ctx, crypto, allocation)
	require.NoError(t, err, "failed to encrypt allocation")

	return allocationEncx
}

// NewTestDedicatedAllocationEncx creates a pre-encrypted dedicated allocation for repository tests
func NewTestDedicatedAllocationEncx(t *testing.T, crypto encx.CryptoService, roomID, userID uuid.UUID, startDate, endDate time.Time) *domain.RoomAllocationEncx {
	t.Helper()

	ctx := context.Background()

	// Create domain entity
	allocation, err := domain.NewDedicatedAllocation(roomID, userID, &startDate, &endDate)
	require.NoError(t, err, "failed to create dedicated allocation")

	// Encrypt
	allocationEncx, err := domain.ProcessRoomAllocationEncx(ctx, crypto, allocation)
	require.NoError(t, err, "failed to encrypt allocation")

	return allocationEncx
}

// NewTestDedicatedAllocationEncxWithNilEndDate creates a pre-encrypted dedicated allocation with nil end date
func NewTestDedicatedAllocationEncxWithNilEndDate(t *testing.T, crypto encx.CryptoService, roomID, userID uuid.UUID, startDate time.Time) *domain.RoomAllocationEncx {
	t.Helper()

	ctx := context.Background()

	// Create domain entity
	allocation, err := domain.NewDedicatedAllocation(roomID, userID, &startDate, nil)
	require.NoError(t, err, "failed to create dedicated allocation with nil end date")

	// Encrypt
	allocationEncx, err := domain.ProcessRoomAllocationEncx(ctx, crypto, allocation)
	require.NoError(t, err, "failed to encrypt allocation")

	return allocationEncx
}

// NewTestInactiveSharedAllocationEncx creates a pre-encrypted inactive shared allocation
func NewTestInactiveSharedAllocationEncx(t *testing.T, crypto encx.CryptoService, roomID, userID uuid.UUID) *domain.RoomAllocationEncx {
	t.Helper()

	ctx := context.Background()

	// Create domain entity
	allocation, err := domain.NewSharedAllocation(roomID, userID)
	require.NoError(t, err, "failed to create shared allocation")

	// Deactivate
	allocation.Deactivate()

	// Encrypt
	allocationEncx, err := domain.ProcessRoomAllocationEncx(ctx, crypto, allocation)
	require.NoError(t, err, "failed to encrypt allocation")

	return allocationEncx
}
