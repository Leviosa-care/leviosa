package availabilityRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	helpers "github.com/Leviosa-care/leviosa/backend/test/helpers"
	availabilityHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	buildingHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	roomHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetRecurringAvailabilities TEST_PATH=internal/booking/infrastructure/postgres/availability/get_recurring_availabilities_test.go

func TestGetRecurringAvailabilities(t *testing.T) {
	ctx := context.Background()

	t.Run("should return empty list when no recurring availabilities exist", func(t *testing.T) {
		// Setup
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)

		until := time.Now().Add(30 * 24 * time.Hour) // 30 days from now

		// Execute
		recurringAvailabilities, err := repo.GetRecurringAvailabilities(ctx, until)

		// Assert
		require.NoError(t, err)
		assert.Empty(t, recurringAvailabilities)
	})

	t.Run("should return only recurring availabilities", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		buildingHelpers.ClearBuildingsTable(t, ctx, testPool)
		roomHelpers.ClearRoomsTable(t, ctx, testPool)
		helpers.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building → room → partner
		buildingEncx := buildingHelpers.NewTestBuildingEncx(t)
		err := buildingHelpers.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		roomEncx := roomHelpers.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = roomHelpers.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		userID := helpers.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := helpers.NewTestPartnerEncxWithUserID(t, userID)
		err = helpers.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		now := time.Now()
		until := now.Add(30 * 24 * time.Hour)

		// Create recurring availability
		recurring := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		recurring.UserID = partnerEncx.ID
		recurring.RoomID = roomEncx.ID
		recurring.StartTime = now.Add(24 * time.Hour) // Start tomorrow
		recurring.Status = domain.AvailabilityStatusAvailable

		// Create non-recurring availabilities
		nonRecurring1 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		nonRecurring1.Status = domain.AvailabilityStatusAvailable

		nonRecurring2 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		nonRecurring2.Status = domain.AvailabilityStatusBooked

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, recurring, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, nonRecurring1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, nonRecurring2, testPool)

		// Execute
		result, err := repo.GetRecurringAvailabilities(ctx, until)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 1)

		// Verify the result is the recurring availability
		retrieved := result[0]
		assert.Equal(t, recurring.ID, retrieved.ID)
		assert.True(t, retrieved.IsRecurring)
	})

	t.Run("should return only available recurring availabilities", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		buildingHelpers.ClearBuildingsTable(t, ctx, testPool)
		roomHelpers.ClearRoomsTable(t, ctx, testPool)
		helpers.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building → room → partner
		buildingEncx := buildingHelpers.NewTestBuildingEncx(t)
		err := buildingHelpers.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		roomEncx := roomHelpers.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = roomHelpers.InsertRoomEncx(t, ctx, testPool, roomEncx)

		userID := helpers.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := helpers.NewTestPartnerEncxWithUserID(t, userID)
		err = helpers.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		now := time.Now()
		until := now.Add(30 * 24 * time.Hour)

		// Create available recurring availability
		availableRecurring := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		availableRecurring.UserID = partnerEncx.ID
		availableRecurring.RoomID = roomEncx.ID
		availableRecurring.StartTime = now.Add(24 * time.Hour)
		availableRecurring.EndTime = now.Add(25 * time.Hour)
		availableRecurring.Status = domain.AvailabilityStatusAvailable

		// Create booked recurring availability
		bookedRecurring := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		bookedRecurring.UserID = partnerEncx.ID
		bookedRecurring.RoomID = roomEncx.ID
		bookedRecurring.StartTime = now.Add(48 * time.Hour)
		bookedRecurring.EndTime = now.Add(49 * time.Hour)
		bookedRecurring.Status = domain.AvailabilityStatusBooked

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, availableRecurring, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, bookedRecurring, testPool)

		// Execute
		result, err := repo.GetRecurringAvailabilities(ctx, until)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 1)

		// Verify the result has available status
		retrieved := result[0]
		assert.Equal(t, availableRecurring.ID, retrieved.ID)
		assert.Equal(t, domain.AvailabilityStatusAvailable, retrieved.Status)
		assert.True(t, retrieved.IsRecurring)
	})

	t.Run("should filter by start time using until parameter", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		buildingHelpers.ClearBuildingsTable(t, ctx, testPool)
		roomHelpers.ClearRoomsTable(t, ctx, testPool)
		helpers.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building → room → partner
		buildingEncx := buildingHelpers.NewTestBuildingEncx(t)
		err := buildingHelpers.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		roomEncx := roomHelpers.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = roomHelpers.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		userID := helpers.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := helpers.NewTestPartnerEncxWithUserID(t, userID)
		err = helpers.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		now := time.Now()

		// Create recurring availability that starts before the until date
		earlyRecurring := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		earlyRecurring.UserID = partnerEncx.ID
		earlyRecurring.RoomID = roomEncx.ID
		earlyRecurring.StartTime = now.Add(1 * time.Hour) // Start in 1 hour
		earlyRecurring.EndTime = now.Add(2 * time.Hour)   // End in 1 hour
		earlyRecurring.Status = domain.AvailabilityStatusAvailable

		// Create recurring availability that starts after the until date
		lateRecurring := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		lateRecurring.UserID = partnerEncx.ID
		lateRecurring.RoomID = roomEncx.ID
		lateRecurring.StartTime = now.Add(48 * time.Hour) // Start in 2 days
		lateRecurring.EndTime = now.Add(49 * time.Hour)   // End in 2 days
		lateRecurring.Status = domain.AvailabilityStatusAvailable

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, earlyRecurring, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, lateRecurring, testPool)

		// Execute - set until to 24 hours from now (should only include earlyRecurring)
		until := now.Add(24 * time.Hour)
		result, err := repo.GetRecurringAvailabilities(ctx, until)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 1)

		// Verify the result is the early recurring availability
		retrieved := result[0]
		assert.Equal(t, earlyRecurring.ID, retrieved.ID)
		assert.True(t, retrieved.StartTime.Before(until) || retrieved.StartTime.Equal(until))
	})

	t.Run("should include recurring availabilities that start exactly at until time", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		buildingHelpers.ClearBuildingsTable(t, ctx, testPool)
		roomHelpers.ClearRoomsTable(t, ctx, testPool)
		helpers.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building → room → partner
		buildingEncx := buildingHelpers.NewTestBuildingEncx(t)
		err := buildingHelpers.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		roomEncx := roomHelpers.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = roomHelpers.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		userID := helpers.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := helpers.NewTestPartnerEncxWithUserID(t, userID)
		err = helpers.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		now := time.Now()
		until := now.Add(24 * time.Hour)

		// Create recurring availability that starts exactly at until time
		exactRecurring := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		exactRecurring.UserID = partnerEncx.ID
		exactRecurring.RoomID = roomEncx.ID
		exactRecurring.StartTime = until
		exactRecurring.Status = domain.AvailabilityStatusAvailable

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, exactRecurring, testPool)

		// Execute
		result, err := repo.GetRecurringAvailabilities(ctx, until)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 1)

		// Verify the result includes the exact match
		retrieved := result[0]
		assert.Equal(t, exactRecurring.ID, retrieved.ID)
		assert.WithinDuration(t, until, retrieved.StartTime, time.Second)
	})

	t.Run("should return multiple recurring availabilities ordered by start time", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		buildingHelpers.ClearBuildingsTable(t, ctx, testPool)
		roomHelpers.ClearRoomsTable(t, ctx, testPool)
		helpers.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building → room → partner
		buildingEncx := buildingHelpers.NewTestBuildingEncx(t)
		err := buildingHelpers.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		roomEncx := roomHelpers.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = roomHelpers.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		userID := helpers.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := helpers.NewTestPartnerEncxWithUserID(t, userID)
		err = helpers.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		now := time.Now()
		until := now.Add(30 * 24 * time.Hour)

		// Create multiple recurring availabilities at different times
		recurring1 := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		recurring1.UserID = partnerEncx.ID
		recurring1.RoomID = roomEncx.ID
		recurring1.StartTime = now.Add(3 * time.Hour)
		recurring1.Status = domain.AvailabilityStatusAvailable

		recurring2 := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		recurring2.UserID = partnerEncx.ID
		recurring2.RoomID = roomEncx.ID
		recurring2.StartTime = now.Add(1 * time.Hour) // Earliest
		recurring2.Status = domain.AvailabilityStatusAvailable

		recurring3 := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		recurring3.UserID = partnerEncx.ID
		recurring3.RoomID = roomEncx.ID
		recurring3.StartTime = now.Add(2 * time.Hour)
		recurring3.Status = domain.AvailabilityStatusAvailable

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, recurring1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, recurring2, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, recurring3, testPool)

		// Execute
		result, err := repo.GetRecurringAvailabilities(ctx, until)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 3)

		// Verify results are ordered by start_time ascending
		assert.WithinDuration(t, recurring2.StartTime, result[0].StartTime, time.Second)
		assert.WithinDuration(t, recurring3.StartTime, result[1].StartTime, time.Second)
		assert.WithinDuration(t, recurring1.StartTime, result[2].StartTime, time.Second)

		// Verify all are recurring and available
		for _, avail := range result {
			assert.Equal(t, domain.AvailabilityStatusAvailable, avail.Status)
			assert.True(t, avail.IsRecurring)
		}
	})

	t.Run("should preserve encrypted fields in recurring availabilities", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		buildingHelpers.ClearBuildingsTable(t, ctx, testPool)
		roomHelpers.ClearRoomsTable(t, ctx, testPool)
		helpers.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building → room → partner
		buildingEncx := buildingHelpers.NewTestBuildingEncx(t)
		err := buildingHelpers.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		roomEncx := roomHelpers.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = roomHelpers.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		userID := helpers.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := helpers.NewTestPartnerEncxWithUserID(t, userID)
		err = helpers.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		now := time.Now()
		until := now.Add(30 * 24 * time.Hour)

		original := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		original.UserID = partnerEncx.ID
		original.RoomID = roomEncx.ID
		original.StartTime = now.Add(24 * time.Hour)
		original.Status = domain.AvailabilityStatusAvailable

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, original, testPool)

		// Execute
		result, err := repo.GetRecurringAvailabilities(ctx, until)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 1)

		retrieved := result[0]
		assert.Equal(t, original.ID, retrieved.ID)
		assert.Equal(t, original.ServiceTypeEncrypted, retrieved.ServiceTypeEncrypted)
		assert.Equal(t, original.NotesEncrypted, retrieved.NotesEncrypted)
		assert.Equal(t, original.DEKEncrypted, retrieved.DEKEncrypted)
		assert.Equal(t, original.KeyVersion, retrieved.KeyVersion)
		assert.Equal(t, original.Metadata, retrieved.Metadata)
	})

	t.Run("should handle recurring availabilities from different partners and rooms", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		buildingHelpers.ClearBuildingsTable(t, ctx, testPool)
		roomHelpers.ClearRoomsTable(t, ctx, testPool)
		helpers.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building
		buildingEncx := buildingHelpers.NewTestBuildingEncx(t)
		err := buildingHelpers.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Create room 1
		roomEncx1 := roomHelpers.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = roomHelpers.InsertRoomEncx(t, ctx, testPool, roomEncx1)
		require.NoError(t, err)

		// Create room 2
		roomEncx2 := roomHelpers.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = roomHelpers.InsertRoomEncx(t, ctx, testPool, roomEncx2)
		require.NoError(t, err)

		// Create partner 1
		userEncx1 := helpers.NewTestUserEncx(t)
		userEncx1.EmailHash = "email_hash_1"
		helpers.InsertUserEncx(t, ctx, userEncx1, testPool)
		partnerEncx1 := helpers.NewTestPartnerEncxWithUserID(t, userEncx1.ID)
		err = helpers.InsertPartnerEncx(t, ctx, partnerEncx1, testPool)
		require.NoError(t, err)

		// Create partner 2
		userID2 := helpers.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx2 := helpers.NewTestPartnerEncxWithUserID(t, userID2)
		err = helpers.InsertPartnerEncx(t, ctx, partnerEncx2, testPool)
		require.NoError(t, err)

		now := time.Now()
		until := now.Add(30 * 24 * time.Hour)

		// Create recurring availabilities for different partners and rooms
		recurring1 := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		recurring1.UserID = partnerEncx1.ID
		recurring1.RoomID = roomEncx1.ID
		recurring1.StartTime = now.Add(1 * time.Hour)
		recurring1.Status = domain.AvailabilityStatusAvailable

		recurring2 := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		recurring2.UserID = partnerEncx2.ID
		recurring2.RoomID = roomEncx2.ID
		recurring2.StartTime = now.Add(2 * time.Hour)
		recurring2.Status = domain.AvailabilityStatusAvailable

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, recurring1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, recurring2, testPool)

		// Execute
		result, err := repo.GetRecurringAvailabilities(ctx, until)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify both are returned
		partnerIDs := make(map[uuid.UUID]bool)
		roomIDs := make(map[uuid.UUID]bool)
		for _, avail := range result {
			assert.Equal(t, domain.AvailabilityStatusAvailable, avail.Status)
			assert.True(t, avail.IsRecurring)
			partnerIDs[avail.UserID] = true
			roomIDs[avail.RoomID] = true
		}
		assert.Len(t, partnerIDs, 2)
		assert.Len(t, roomIDs, 2)
	})

	t.Run("should handle cancelled recurring availabilities correctly", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		buildingHelpers.ClearBuildingsTable(t, ctx, testPool)
		roomHelpers.ClearRoomsTable(t, ctx, testPool)
		helpers.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building → room → partner
		buildingEncx := buildingHelpers.NewTestBuildingEncx(t)
		err := buildingHelpers.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		roomEncx := roomHelpers.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = roomHelpers.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		userID := helpers.CreateTestUserForPartner(t, ctx, testPool)
		partnerEncx := helpers.NewTestPartnerEncxWithUserID(t, userID)
		err = helpers.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		now := time.Now()
		until := now.Add(30 * 24 * time.Hour)

		// Create available recurring availability
		availableRecurring := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		availableRecurring.UserID = partnerEncx.ID
		availableRecurring.RoomID = roomEncx.ID
		availableRecurring.StartTime = now.Add(24 * time.Hour)
		availableRecurring.Status = domain.AvailabilityStatusAvailable

		// Create cancelled recurring availability (should not be returned)
		cancelledRecurring := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		cancelledRecurring.UserID = partnerEncx.ID
		cancelledRecurring.RoomID = roomEncx.ID
		cancelledRecurring.StartTime = now.Add(25 * time.Hour)
		cancelledRecurring.Status = domain.AvailabilityStatusCancelled

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, availableRecurring, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, cancelledRecurring, testPool)

		// Execute
		result, err := repo.GetRecurringAvailabilities(ctx, until)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 1)

		// Verify only the available recurring availability is returned
		retrieved := result[0]
		assert.Equal(t, availableRecurring.ID, retrieved.ID)
		assert.Equal(t, domain.AvailabilityStatusAvailable, retrieved.Status)
		assert.True(t, retrieved.IsRecurring)
	})
}
