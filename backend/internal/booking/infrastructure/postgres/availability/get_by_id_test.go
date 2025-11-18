package availabilityRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"
	availabilityHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetByID TEST_PATH=internal/booking/infrastructure/postgres/availability/get_by_id_test.go

func TestGetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should get availability by ID successfully", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building → room → partner
		buildingEncx := tb.NewTestBuildingEncx(t)
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		userID := th.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := th.NewTestPartnerEncxWithUserID(t, userID)
		err = th.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		original := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, original, testPool)

		// Execute
		retrieved, err := repo.GetByID(ctx, original.ID)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, original.ID, retrieved.ID)
		assert.Equal(t, original.UserID, retrieved.UserID)
		assert.Equal(t, original.RoomID, retrieved.RoomID)
		assert.WithinDuration(t, original.StartTime, retrieved.StartTime, time.Second)
		assert.WithinDuration(t, original.EndTime, retrieved.EndTime, time.Second)
		assert.Equal(t, original.PriceCents, retrieved.PriceCents)
		assert.Equal(t, original.MaxCapacity, retrieved.MaxCapacity)
		assert.Equal(t, original.IsRecurring, retrieved.IsRecurring)
		assert.Equal(t, original.Status, retrieved.Status)
		assert.Equal(t, original.ServiceTypeEncrypted, retrieved.ServiceTypeEncrypted)
		assert.Equal(t, original.NotesEncrypted, retrieved.NotesEncrypted)
		assert.Equal(t, original.RecurrencePatternEncrypted, retrieved.RecurrencePatternEncrypted)
	})

	t.Run("should return not found when availability does not exist", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)

		nonExistentID := uuid.New()

		// Execute
		retrieved, err := repo.GetByID(ctx, nonExistentID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, retrieved)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should get recurring availability by ID", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building → room → partner
		buildingEncx := tb.NewTestBuildingEncx(t)
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		userID := th.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := th.NewTestPartnerEncxWithUserID(t, userID)
		err = th.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		original := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		original.UserID = partnerEncx.ID
		original.RoomID = roomEncx.ID
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, original, testPool)

		// Execute
		retrieved, err := repo.GetByID(ctx, original.ID)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.True(t, retrieved.IsRecurring)
		assert.NotNil(t, retrieved.RecurrencePatternEncrypted)
		assert.Equal(t, original.RecurrencePatternEncrypted, retrieved.RecurrencePatternEncrypted)
	})

	t.Run("should get availability with different statuses", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building → room → partner
		buildingEncx := tb.NewTestBuildingEncx(t)
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		userID := th.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := th.NewTestPartnerEncxWithUserID(t, userID)
		err = th.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		testCases := []domain.AvailabilityStatus{
			domain.AvailabilityStatusAvailable,
			domain.AvailabilityStatusBooked,
			domain.AvailabilityStatusCancelled,
			domain.AvailabilityStatusBlocked,
		}

		for _, status := range testCases {
			t.Run(string(status), func(t *testing.T) {
				original := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
				original.ID = uuid.New() // Ensure unique ID for each test
				original.Status = status
				availabilityHelpers.InsertAvailabilityEncx(t, ctx, original, testPool)

				// Execute
				retrieved, err := repo.GetByID(ctx, original.ID)

				// Assert
				require.NoError(t, err)
				assert.NotNil(t, retrieved)
				assert.Equal(t, status, retrieved.Status)

				// Clean up this specific test data
				availabilityHelpers.DeleteAvailabilityEncx(t, ctx, original.ID, testPool)
			})
		}
	})

	t.Run("should get availability with null price", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building → room → partner
		buildingEncx := tb.NewTestBuildingEncx(t)
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)

		userID := th.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := th.NewTestPartnerEncxWithUserID(t, userID)
		err = th.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		original := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		original.PriceCents = nil // Free session
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, original, testPool)

		// Execute
		retrieved, err := repo.GetByID(ctx, original.ID)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Nil(t, retrieved.PriceCents)
	})

	t.Run("should preserve encrypted fields", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building → room → partner
		buildingEncx := tb.NewTestBuildingEncx(t)
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		userID := th.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := th.NewTestPartnerEncxWithUserID(t, userID)
		err = th.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		original := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		original.UserID = partnerEncx.ID
		original.RoomID = roomEncx.ID
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, original, testPool)

		// Execute
		retrieved, err := repo.GetByID(ctx, original.ID)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, original.ServiceTypeEncrypted, retrieved.ServiceTypeEncrypted)
		assert.Equal(t, original.NotesEncrypted, retrieved.NotesEncrypted)
		assert.Equal(t, original.RecurrencePatternEncrypted, retrieved.RecurrencePatternEncrypted)
		assert.Equal(t, original.DEKEncrypted, retrieved.DEKEncrypted)
		assert.Equal(t, original.KeyVersion, retrieved.KeyVersion)
		assert.Equal(t, original.Metadata, retrieved.Metadata)
	})

	t.Run("should get availability from multiple records", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building → room → partner
		buildingEncx := tb.NewTestBuildingEncx(t)
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)

		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		userID := th.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := th.NewTestPartnerEncxWithUserID(t, userID)
		err = th.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Insert multiple availabilities
		availabilities := make([]*domain.AvailabilityEncx, 5)
		for i := 0; i < 5; i++ {
			availabilities[i] = availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
			availabilityHelpers.InsertAvailabilityEncx(t, ctx, availabilities[i], testPool)
		}

		// Execute - get the third availability
		targetAvailability := availabilities[2]
		retrieved, err := repo.GetByID(ctx, targetAvailability.ID)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, targetAvailability.ID, retrieved.ID)
		assert.Equal(t, targetAvailability.UserID, retrieved.UserID)
		assert.Equal(t, targetAvailability.RoomID, retrieved.RoomID)
	})

	t.Run("should handle empty encrypted fields", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building → room → partner
		buildingEncx := tb.NewTestBuildingEncx(t)
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		userID := th.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := th.NewTestPartnerEncxWithUserID(t, userID)
		err = th.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		original := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		original.ServiceTypeEncrypted = []byte{}
		original.NotesEncrypted = []byte{}
		original.RecurrencePatternEncrypted = nil // Non-recurring
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, original, testPool)

		// Execute
		retrieved, err := repo.GetByID(ctx, original.ID)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Empty(t, retrieved.ServiceTypeEncrypted)
		assert.Empty(t, retrieved.NotesEncrypted)
		assert.Nil(t, retrieved.RecurrencePatternEncrypted)
	})
}
