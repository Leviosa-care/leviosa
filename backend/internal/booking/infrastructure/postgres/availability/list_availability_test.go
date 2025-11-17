package availabilityRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"
	availabilityHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestList TEST_PATH=internal/booking/infrastructure/postgres/availability/list_availability_test.go

func TestList(t *testing.T) {
	ctx := context.Background()

	t.Run("should return empty list when no availabilities exist", func(t *testing.T) {
		// Setup
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)

		// Execute
		filter := ports.AvailabilityFilter{}
		availabilities, err := repo.List(ctx, filter)

		// Assert
		require.NoError(t, err)
		assert.Empty(t, availabilities)
	})

	t.Run("should return all availabilities when no filter is applied", func(t *testing.T) {
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

		// Insert multiple availabilities
		availabilities := make([]*domain.AvailabilityEncx, 3)
		for i := 0; i < 3; i++ {
			availabilities[i] = availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
			availabilityHelpers.InsertAvailabilityEncx(t, ctx, availabilities[i], testPool)
		}

		// Execute
		filter := ports.AvailabilityFilter{}
		result, err := repo.List(ctx, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 3)
	})

	t.Run("should filter by partner ID", func(t *testing.T) {
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
		require.NoError(t, err)
		partnerEncx1 := th.NewTestPartnerEncxWithUserID(t, userEncx1.ID)
		err = th.InsertPartnerEncx(t, ctx, partnerEncx1, testPool)
		require.NoError(t, err)

		// Create partner 2
		userID2 := th.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx2 := th.NewTestPartnerEncxWithUserID(t, userID2)
		err = th.InsertPartnerEncx(t, ctx, partnerEncx2, testPool)
		require.NoError(t, err)

		// Create availabilities for different partners
		avail1 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx1.ID, roomEncx.ID)
		avail2 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx1.ID, roomEncx.ID)
		avail3 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx2.ID, roomEncx.ID)

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail2, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail3, testPool)

		// Execute - filter by partner1
		filter := ports.AvailabilityFilter{UserID: &partnerEncx1.ID}
		result, err := repo.List(ctx, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify all results belong to partner1
		for _, avail := range result {
			assert.Equal(t, partnerEncx1.ID, avail.UserID)
		}
	})

	t.Run("should filter by room ID", func(t *testing.T) {
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

		// Create room 2
		roomEncx2 := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx2)

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

		// Execute - filter by room1
		filter := ports.AvailabilityFilter{RoomID: &roomEncx1.ID}
		result, err := repo.List(ctx, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify all results belong to room1
		for _, avail := range result {
			assert.Equal(t, roomEncx1.ID, avail.RoomID)
		}
	})

	t.Run("should filter by status", func(t *testing.T) {
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
		filter := ports.AvailabilityFilter{Status: []domain.AvailabilityStatus{domain.AvailabilityStatusAvailable}}
		result, err := repo.List(ctx, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify all results have available status
		for _, avail := range result {
			assert.Equal(t, domain.AvailabilityStatusAvailable, avail.Status)
		}
	})

	t.Run("should apply limit and offset", func(t *testing.T) {
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

		// Create multiple availabilities
		for i := 0; i < 5; i++ {
			avail := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
			availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail, testPool)
		}

		// Execute - limit 2, offset 1
		filter := ports.AvailabilityFilter{
			Limit:  2,
			Offset: 1,
		}
		result, err := repo.List(ctx, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)
	})
}
