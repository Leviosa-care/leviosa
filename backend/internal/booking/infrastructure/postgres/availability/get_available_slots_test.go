package availabilityRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"
	availabilityHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAvailableSlots TEST_PATH=internal/booking/infrastructure/postgres/availability/get_available_slots_test.go

func TestGetAvailableSlots(t *testing.T) {
	ctx := context.Background()

	t.Run("should return empty list when no available slots exist", func(t *testing.T) {
		// Setup
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)

		// Execute
		filter := ports.AvailabilityFilter{}
		availableSlots, err := repo.GetAvailableSlots(ctx, filter)

		// Assert
		require.NoError(t, err)
		assert.Empty(t, availableSlots)
	})

	t.Run("should return only available status availabilities", func(t *testing.T) {
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

		// Create availabilities with different statuses
		available1 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		available1.Status = domain.AvailabilityStatusAvailable

		booked := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		booked.Status = domain.AvailabilityStatusBooked

		cancelled := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		cancelled.Status = domain.AvailabilityStatusCancelled

		blocked := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		blocked.Status = domain.AvailabilityStatusBlocked

		available2 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		available2.Status = domain.AvailabilityStatusAvailable

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, available1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, booked, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, cancelled, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, blocked, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, available2, testPool)

		// Execute
		filter := ports.AvailabilityFilter{}
		result, err := repo.GetAvailableSlots(ctx, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify all results have available status
		for _, avail := range result {
			assert.Equal(t, domain.AvailabilityStatusAvailable, avail.Status)
		}
	})

	// Note: "should filter out past availabilities" test removed because the database
	// CHECK constraint (start_time >= NOW()) prevents creating past availabilities,
	// making this test scenario impossible and redundant.

	t.Run("should include availabilities that start now", func(t *testing.T) {
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

		now := time.Now()

		// Create availability that starts now
		currentAvail := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		currentAvail.Status = domain.AvailabilityStatusAvailable
		currentAvail.StartTime = now.Add(1 * time.Minute) // Slightly in the future to account for test timing
		currentAvail.EndTime = now.Add(2 * time.Hour)

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, currentAvail, testPool)

		// Execute
		filter := ports.AvailabilityFilter{}
		result, err := repo.GetAvailableSlots(ctx, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 1)

		// Verify the result includes the current availability
		retrieved := result[0]
		assert.Equal(t, currentAvail.ID, retrieved.ID)
		assert.True(t, retrieved.StartTime.After(now.Add(-1*time.Minute))) // Allow for some test timing variance
	})

	t.Run("should apply additional filters correctly", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building
		buildingEncx := tb.NewTestBuildingEncx(t)
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Create room 1
		roomEncx1 := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx1)
		require.NoError(t, err)

		// Create room 2
		roomEncx2 := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx2)
		require.NoError(t, err)

		// Create partner
		userID := th.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := th.NewTestPartnerEncxWithUserID(t, userID)
		err = th.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Create available availabilities for different rooms
		avail1 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx1.ID)
		avail1.Status = domain.AvailabilityStatusAvailable

		avail2 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx2.ID)
		avail2.Status = domain.AvailabilityStatusAvailable

		// Create booked availability for the same room (should be filtered out)
		booked := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx1.ID)
		booked.Status = domain.AvailabilityStatusBooked

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail2, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, booked, testPool)

		// Execute - filter by room1
		filter := ports.AvailabilityFilter{
			RoomID: &roomEncx1.ID,
		}
		result, err := repo.GetAvailableSlots(ctx, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 1)

		// Verify result has available status and correct room
		retrieved := result[0]
		assert.Equal(t, domain.AvailabilityStatusAvailable, retrieved.Status)
		assert.Equal(t, roomEncx1.ID, retrieved.RoomID)
	})

	t.Run("should handle time range filtering with available filter", func(t *testing.T) {
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

		now := time.Now()

		// Create availabilities at different times
		avail1 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		avail1.Status = domain.AvailabilityStatusAvailable
		avail1.StartTime = now.Add(1 * time.Hour)
		avail1.EndTime = now.Add(2 * time.Hour)

		avail2 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		avail2.Status = domain.AvailabilityStatusAvailable
		avail2.StartTime = now.Add(3 * time.Hour)
		avail2.EndTime = now.Add(4 * time.Hour)

		// Create availability outside time range but available (should be filtered out by time)
		avail3 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		avail3.Status = domain.AvailabilityStatusAvailable
		avail3.StartTime = now.Add(10 * time.Hour)
		avail3.EndTime = now.Add(11 * time.Hour)

		// Create availability in time range but booked (should be filtered out by status)
		booked := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		booked.Status = domain.AvailabilityStatusBooked
		booked.StartTime = now.Add(2 * time.Hour)
		booked.EndTime = now.Add(3 * time.Hour)

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail2, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail3, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, booked, testPool)

		// Execute - filter by time range that includes avail1, avail2, and booked
		filter := ports.AvailabilityFilter{
			StartTime: &[]time.Time{now}[0],
			EndTime:   &[]time.Time{now.Add(5 * time.Hour)}[0],
		}
		result, err := repo.GetAvailableSlots(ctx, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify all results are available and in time range
		for _, avail := range result {
			assert.Equal(t, domain.AvailabilityStatusAvailable, avail.Status)
			assert.True(t, avail.StartTime.Before(now.Add(5*time.Hour)))
		}
	})

	t.Run("should handle recurring available slots", func(t *testing.T) {
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

		// Create recurring availability that is available
		recurring := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		recurring.PartnerID = partnerEncx.ID
		recurring.RoomID = roomEncx.ID
		recurring.Status = domain.AvailabilityStatusAvailable
		// Ensure it starts in the future
		recurring.StartTime = time.Now().Add(24 * time.Hour)
		recurring.EndTime = recurring.StartTime.Add(1 * time.Hour)

		// Create non-recurring availability that is available
		nonRecurring := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		nonRecurring.Status = domain.AvailabilityStatusAvailable

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, recurring, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, nonRecurring, testPool)

		// Execute
		filter := ports.AvailabilityFilter{}
		result, err := repo.GetAvailableSlots(ctx, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify both recurring and non-recurring are included
		recurringCount := 0
		nonRecurringCount := 0
		for _, avail := range result {
			assert.Equal(t, domain.AvailabilityStatusAvailable, avail.Status)
			if avail.IsRecurring {
				recurringCount++
			} else {
				nonRecurringCount++
			}
		}
		assert.Equal(t, 1, recurringCount)
		assert.Equal(t, 1, nonRecurringCount)
	})

	t.Run("should handle limit and offset correctly", func(t *testing.T) {
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

		// Create multiple available availabilities
		for i := 0; i < 5; i++ {
			avail := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
			avail.Status = domain.AvailabilityStatusAvailable
			availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail, testPool)
		}

		// Execute - get first 2 available slots
		filter := ports.AvailabilityFilter{
			Limit: 2,
		}
		result, err := repo.GetAvailableSlots(ctx, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify all results are available
		for _, avail := range result {
			assert.Equal(t, domain.AvailabilityStatusAvailable, avail.Status)
		}
	})

	t.Run("should preserve encrypted fields in available slots", func(t *testing.T) {
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
		original.PartnerID = partnerEncx.ID
		original.RoomID = roomEncx.ID
		original.Status = domain.AvailabilityStatusAvailable
		// Ensure it starts in the future
		original.StartTime = time.Now().Add(24 * time.Hour)
		original.EndTime = original.StartTime.Add(1 * time.Hour)

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, original, testPool)

		// Execute
		filter := ports.AvailabilityFilter{}
		result, err := repo.GetAvailableSlots(ctx, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 1)

		retrieved := result[0]
		assert.Equal(t, original.ID, retrieved.ID)
		assert.Equal(t, domain.AvailabilityStatusAvailable, retrieved.Status)
		assert.Equal(t, original.ServiceTypeEncrypted, retrieved.ServiceTypeEncrypted)
		assert.Equal(t, original.NotesEncrypted, retrieved.NotesEncrypted)
		assert.Equal(t, original.RecurrencePatternEncrypted, retrieved.RecurrencePatternEncrypted)
	})
}
