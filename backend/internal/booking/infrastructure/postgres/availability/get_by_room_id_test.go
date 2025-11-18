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

// make test-func TEST_NAME=TestGetByRoomID TEST_PATH=internal/booking/infrastructure/postgres/availability/get_by_room_id_test.go

func TestGetByRoomID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return empty list when room has no availabilities", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create dependencies: building → room (but no availabilities)
		buildingEncx := tb.NewTestBuildingEncx(t)
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		// Execute
		filter := ports.AvailabilityFilter{}
		availabilities, err := repo.GetByRoomID(ctx, roomEncx.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Empty(t, availabilities)
	})

	t.Run("should return all availabilities for a room", func(t *testing.T) {
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

		// Create multiple availabilities for the same room
		avail1 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		avail2 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		avail3 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail2, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail3, testPool)

		// Execute
		filter := ports.AvailabilityFilter{}
		result, err := repo.GetByRoomID(ctx, roomEncx.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 3)

		// Verify all results belong to the correct room
		for _, avail := range result {
			assert.Equal(t, roomEncx.ID, avail.RoomID)
		}
	})

	t.Run("should only return availabilities for the specified room", func(t *testing.T) {
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

		// Create availabilities for different rooms
		avail1 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx1.ID)
		avail2 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx1.ID)
		avail3 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx2.ID)

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail2, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail3, testPool)

		// Execute - get availabilities for room1
		filter := ports.AvailabilityFilter{}
		result, err := repo.GetByRoomID(ctx, roomEncx1.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify all results belong to room1
		for _, avail := range result {
			assert.Equal(t, roomEncx1.ID, avail.RoomID)
		}
	})

	t.Run("should return availabilities from different partners for the same room", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building → room
		buildingEncx := tb.NewTestBuildingEncx(t)
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		// Create partner 1
		userEncx1 := th.NewTestUserEncx(t)
		userEncx1.EmailHash = "email_hash_1"
		err = th.InsertUserEncx(t, ctx, userEncx1, testPool)
		partnerEncx1 := th.NewTestPartnerEncxWithUserID(t, userEncx1.ID)
		err = th.InsertPartnerEncx(t, ctx, partnerEncx1, testPool)
		require.NoError(t, err)

		// Create partner 2
		userID2 := th.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx2 := th.NewTestPartnerEncxWithUserID(t, userID2)
		err = th.InsertPartnerEncx(t, ctx, partnerEncx2, testPool)
		require.NoError(t, err)

		// Create availabilities for different partners in the same room
		avail1 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx1.ID, roomEncx.ID)
		avail2 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx2.ID, roomEncx.ID)
		avail3 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx1.ID, roomEncx.ID)

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail2, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail3, testPool)

		// Execute
		filter := ports.AvailabilityFilter{}
		result, err := repo.GetByRoomID(ctx, roomEncx.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 3)

		// Verify all results belong to the correct room but may have different partners
		partner1Count := 0
		partner2Count := 0
		for _, avail := range result {
			assert.Equal(t, roomEncx.ID, avail.RoomID)
			if avail.UserID == partnerEncx1.ID {
				partner1Count++
			} else if avail.UserID == partnerEncx2.ID {
				partner2Count++
			}
		}
		assert.Equal(t, 2, partner1Count)
		assert.Equal(t, 1, partner2Count)
	})

	t.Run("should apply additional filters correctly", func(t *testing.T) {
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
		avail1 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		avail2 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		avail3 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)

		avail1.Status = domain.AvailabilityStatusAvailable
		avail2.Status = domain.AvailabilityStatusBooked
		avail3.Status = domain.AvailabilityStatusAvailable

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail2, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail3, testPool)

		// Execute - filter by available status
		filter := ports.AvailabilityFilter{
			Status: []domain.AvailabilityStatus{domain.AvailabilityStatusAvailable},
		}
		result, err := repo.GetByRoomID(ctx, roomEncx.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify all results have available status and correct room
		for _, avail := range result {
			assert.Equal(t, roomEncx.ID, avail.RoomID)
			assert.Equal(t, domain.AvailabilityStatusAvailable, avail.Status)
		}
	})

	t.Run("should handle time range filtering", func(t *testing.T) {
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
		avail2 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		avail3 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)

		avail1.StartTime = now.Add(1 * time.Hour)
		avail1.EndTime = now.Add(2 * time.Hour)

		avail2.StartTime = now.Add(3 * time.Hour)
		avail2.EndTime = now.Add(4 * time.Hour)

		avail3.StartTime = now.Add(10 * time.Hour)
		avail3.EndTime = now.Add(11 * time.Hour)

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail2, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail3, testPool)

		// Execute - filter by time range that includes avail1 and avail2
		filter := ports.AvailabilityFilter{
			StartTime: &[]time.Time{now}[0],
			EndTime:   &[]time.Time{now.Add(5 * time.Hour)}[0],
		}
		result, err := repo.GetByRoomID(ctx, roomEncx.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify all results belong to the correct room and time range
		for _, avail := range result {
			assert.Equal(t, roomEncx.ID, avail.RoomID)
			assert.True(t, avail.StartTime.Before(now.Add(5*time.Hour)))
		}
	})

	t.Run("should handle recurring and non-recurring availabilities", func(t *testing.T) {
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

		// Create both recurring and non-recurring availabilities
		nonRecurring1 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		nonRecurring2 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		recurring := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		recurring.RoomID = roomEncx.ID
		recurring.UserID = partnerEncx.ID

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, nonRecurring1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, nonRecurring2, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, recurring, testPool)

		// Execute
		filter := ports.AvailabilityFilter{}
		result, err := repo.GetByRoomID(ctx, roomEncx.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 3)

		// Verify all results belong to the correct room
		recurringCount := 0
		nonRecurringCount := 0
		for _, avail := range result {
			assert.Equal(t, roomEncx.ID, avail.RoomID)
			if avail.IsRecurring {
				recurringCount++
			} else {
				nonRecurringCount++
			}
		}
		assert.Equal(t, 1, recurringCount)
		assert.Equal(t, 2, nonRecurringCount)
	})

	t.Run("should handle different capacities in the same room", func(t *testing.T) {
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

		// Create availabilities with different capacities
		avail1 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		avail2 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		avail3 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)

		avail1.MaxCapacity = 1
		avail2.MaxCapacity = 3
		avail3.MaxCapacity = 5

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail2, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail3, testPool)

		// Execute
		filter := ports.AvailabilityFilter{}
		result, err := repo.GetByRoomID(ctx, roomEncx.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 3)

		// Verify all results belong to the correct room and have expected capacities
		capacities := make(map[int]int)
		for _, avail := range result {
			assert.Equal(t, roomEncx.ID, avail.RoomID)
			capacities[avail.MaxCapacity]++
		}
		assert.Equal(t, 1, capacities[1])
		assert.Equal(t, 1, capacities[3])
		assert.Equal(t, 1, capacities[5])
	})

	t.Run("should apply ordering correctly", func(t *testing.T) {
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
		avail2 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		avail3 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)

		avail1.StartTime = now.Add(3 * time.Hour) // Latest
		avail2.StartTime = now.Add(1 * time.Hour) // Earliest
		avail3.StartTime = now.Add(2 * time.Hour) // Middle

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail2, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail3, testPool)

		// Execute - order by start_time descending
		filter := ports.AvailabilityFilter{
			OrderBy:        "start_time",
			OrderDirection: "desc",
		}
		result, err := repo.GetByRoomID(ctx, roomEncx.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 3)

		// Verify results are ordered by start_time descending
		assert.True(t, result[0].StartTime.After(result[1].StartTime))
		assert.True(t, result[1].StartTime.After(result[2].StartTime))
	})
}
