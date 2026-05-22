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

// make test-func TEST_NAME=TestGetByUserID TEST_PATH=internal/booking/infrastructure/postgres/availability/get_by_user_id_test.go

func TestGetByUserID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return empty list when partner has no availabilities", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies: building → room → partner (but no availabilities)
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

		// Execute
		filter := ports.AvailabilityFilter{}
		availabilities, err := repo.GetByUserID(ctx, partnerEncx.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Empty(t, availabilities)
	})

	t.Run("should return all availabilities for a partner", func(t *testing.T) {
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

		// Create multiple availabilities for the same partner
		avail1 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		avail2 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		avail3 := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail2, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, avail3, testPool)

		// Execute
		filter := ports.AvailabilityFilter{}
		result, err := repo.GetByUserID(ctx, partnerEncx.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 3)

		// Verify all results belong to the correct partner
		for _, avail := range result {
			assert.Equal(t, partnerEncx.ID, avail.UserID)
		}
	})

	t.Run("should only return availabilities for the specified partner", func(t *testing.T) {
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

		// Execute - get availabilities for partner1
		filter := ports.AvailabilityFilter{}
		result, err := repo.GetByUserID(ctx, partnerEncx1.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify all results belong to partner1
		for _, avail := range result {
			assert.Equal(t, partnerEncx1.ID, avail.UserID)
		}
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
		result, err := repo.GetByUserID(ctx, partnerEncx.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify all results have available status and correct partner
		for _, avail := range result {
			assert.Equal(t, partnerEncx.ID, avail.UserID)
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
		result, err := repo.GetByUserID(ctx, partnerEncx.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify all results belong to the correct partner
		for _, avail := range result {
			assert.Equal(t, partnerEncx.ID, avail.UserID)
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
		recurring.UserID = partnerEncx.ID
		recurring.RoomID = roomEncx.ID

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, nonRecurring1, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, nonRecurring2, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, recurring, testPool)

		// Execute
		filter := ports.AvailabilityFilter{}
		result, err := repo.GetByUserID(ctx, partnerEncx.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 3)

		// Verify all results belong to the correct partner
		recurringCount := 0
		nonRecurringCount := 0
		for _, avail := range result {
			assert.Equal(t, partnerEncx.ID, avail.UserID)
			if avail.IsRecurring {
				recurringCount++
			} else {
				nonRecurringCount++
			}
		}
		assert.Equal(t, 1, recurringCount)
		assert.Equal(t, 2, nonRecurringCount)
	})

	t.Run("should apply limit and offset correctly", func(t *testing.T) {
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

		// Create multiple availabilities for the partner
		availabilities := make([]*domain.AvailabilityEncx, 5)
		for i := 0; i < 5; i++ {
			availabilities[i] = availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
			availabilityHelpers.InsertAvailabilityEncx(t, ctx, availabilities[i], testPool)
		}

		// Execute - get first 2 availabilities
		filter := ports.AvailabilityFilter{
			Limit:  2,
			Offset: 0,
		}
		result, err := repo.GetByUserID(ctx, partnerEncx.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify all results belong to the correct partner
		for _, avail := range result {
			assert.Equal(t, partnerEncx.ID, avail.UserID)
		}
	})

	t.Run("should handle availabilities with null prices", func(t *testing.T) {
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

		// Create availability with null price (free session)
		freeAvail := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		freeAvail.PriceCents = nil

		// Create availability with price
		paidAvail := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)

		availabilityHelpers.InsertAvailabilityEncx(t, ctx, freeAvail, testPool)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, paidAvail, testPool)

		// Execute
		filter := ports.AvailabilityFilter{}
		result, err := repo.GetByUserID(ctx, partnerEncx.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify we have both free and paid availabilities
		freeCount := 0
		paidCount := 0
		for _, avail := range result {
			assert.Equal(t, partnerEncx.ID, avail.UserID)
			if avail.PriceCents == nil {
				freeCount++
			} else {
				paidCount++
			}
		}
		assert.Equal(t, 1, freeCount)
		assert.Equal(t, 1, paidCount)
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
		filter := ports.AvailabilityFilter{}
		result, err := repo.GetByUserID(ctx, partnerEncx.ID, filter)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 1)

		retrieved := result[0]
		assert.Equal(t, original.ID, retrieved.ID)
		assert.Equal(t, original.ServiceTypeEncrypted, retrieved.ServiceTypeEncrypted)
		assert.Equal(t, original.NotesEncrypted, retrieved.NotesEncrypted)
		assert.Equal(t, original.DEKEncrypted, retrieved.DEKEncrypted)
		assert.Equal(t, original.KeyVersion, retrieved.KeyVersion)
	})
}
