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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdate TEST_PATH=internal/booking/infrastructure/postgres/availability/update_availability_test.go

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	t.Run("should update availability successfully", func(t *testing.T) {
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

		// Prepare updated availability
		updated := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		updated.ID = original.ID
		updated.PartnerID = original.PartnerID
		updated.RoomID = original.RoomID
		updated.StartTime = original.StartTime.Add(1 * time.Hour)
		updated.EndTime = original.EndTime.Add(1 * time.Hour)
		updated.PriceCents = &[]int{25000}[0] // $250.00
		updated.MaxCapacity = 3
		updated.Status = domain.AvailabilityStatusBooked

		// Execute
		err = repo.Update(ctx, updated)

		// Assert
		require.NoError(t, err)
		retrieved := availabilityHelpers.GetAvailabilityEncxFromDB(t, ctx, original.ID, testPool)
		assert.WithinDuration(t, updated.StartTime, retrieved.StartTime, time.Second)
		assert.WithinDuration(t, updated.EndTime, retrieved.EndTime, time.Second)
		assert.Equal(t, updated.PriceCents, retrieved.PriceCents)
		assert.Equal(t, updated.MaxCapacity, retrieved.MaxCapacity)
		assert.Equal(t, updated.Status, retrieved.Status)
	})

	t.Run("should return not found when updating non-existent availability", func(t *testing.T) {
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

		nonExistent := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)

		// Execute
		err = repo.Update(ctx, nonExistent)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should update recurring availability", func(t *testing.T) {
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
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, original, testPool)

		// Update to non-recurring
		updated := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		updated.ID = original.ID
		updated.PartnerID = original.PartnerID
		updated.RoomID = original.RoomID
		updated.StartTime = original.StartTime
		updated.EndTime = original.EndTime
		updated.IsRecurring = false
		updated.RecurrencePatternEncrypted = nil

		// Execute
		err = repo.Update(ctx, updated)

		// Assert
		require.NoError(t, err)
		retrieved := availabilityHelpers.GetAvailabilityEncxFromDB(t, ctx, original.ID, testPool)
		assert.False(t, retrieved.IsRecurring)
		assert.Nil(t, retrieved.RecurrencePatternEncrypted)
	})

	t.Run("should update encrypted fields", func(t *testing.T) {
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

		// Update with new encrypted content
		updated := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		updated.ID = original.ID
		updated.PartnerID = original.PartnerID
		updated.RoomID = original.RoomID
		updated.ServiceTypeEncrypted = []byte("encrypted_new_service_type")
		updated.NotesEncrypted = []byte("encrypted_new_notes")
		updated.DEKEncrypted = []byte("new_mock_dek_data")
		updated.KeyVersion = 2

		// Execute
		err = repo.Update(ctx, updated)

		// Assert
		require.NoError(t, err)
		retrieved := availabilityHelpers.GetAvailabilityEncxFromDB(t, ctx, original.ID, testPool)
		assert.Equal(t, updated.ServiceTypeEncrypted, retrieved.ServiceTypeEncrypted)
		assert.Equal(t, updated.NotesEncrypted, retrieved.NotesEncrypted)
		assert.Equal(t, updated.DEKEncrypted, retrieved.DEKEncrypted)
		assert.Equal(t, updated.KeyVersion, retrieved.KeyVersion)
	})

	t.Run("should update to different status", func(t *testing.T) {
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
		original.Status = domain.AvailabilityStatusAvailable
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, original, testPool)

		// Update to booked status
		updated := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		updated.ID = original.ID
		updated.PartnerID = original.PartnerID
		updated.RoomID = original.RoomID
		updated.Status = domain.AvailabilityStatusBooked

		// Execute
		err = repo.Update(ctx, updated)

		// Assert
		require.NoError(t, err)
		retrieved := availabilityHelpers.GetAvailabilityEncxFromDB(t, ctx, original.ID, testPool)
		assert.Equal(t, domain.AvailabilityStatusBooked, retrieved.Status)
	})
}
