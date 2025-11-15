package availabilityRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"
	ta "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCheckConflict TEST_PATH=internal/booking/infrastructure/postgres/availability/check_conflict_test.go

func TestCheckConflict(t *testing.T) {
	ctx := context.Background()

	t.Run("should return no conflict when no overlapping availabilities exist", func(t *testing.T) {
		// Setup
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies
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

		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(1 * time.Hour)

		// Execute
		hasConflict, err := repo.CheckConflict(ctx, partnerEncx.ID, startTime, endTime, nil)

		// Assert
		assert.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("should detect conflict with existing availability", func(t *testing.T) {
		// Setup
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies
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

		// Create existing availability from 10:00 to 12:00
		existingStartTime := time.Now().Add(24 * time.Hour)
		existingEndTime := existingStartTime.Add(2 * time.Hour)

		existing := ta.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		existing.StartTime = existingStartTime
		existing.EndTime = existingEndTime
		existing.Status = domain.AvailabilityStatusAvailable
		ta.InsertAvailabilityEncx(t, ctx, existing, testPool)

		// Try to create overlapping availability from 11:00 to 13:00
		newStartTime := existingStartTime.Add(1 * time.Hour)
		newEndTime := existingEndTime.Add(1 * time.Hour)

		// Execute
		hasConflict, err := repo.CheckConflict(ctx, partnerEncx.ID, newStartTime, newEndTime, nil)

		// Assert
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("should not detect conflict when times do not overlap", func(t *testing.T) {
		// Setup
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies
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

		// Create existing availability from 10:00 to 12:00
		existingStartTime := time.Now().Add(24 * time.Hour)
		existingEndTime := existingStartTime.Add(2 * time.Hour)

		existing := ta.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		existing.StartTime = existingStartTime
		existing.EndTime = existingEndTime
		existing.Status = domain.AvailabilityStatusAvailable
		ta.InsertAvailabilityEncx(t, ctx, existing, testPool)

		// Try to create non-overlapping availability from 14:00 to 16:00
		newStartTime := existingStartTime.Add(4 * time.Hour)
		newEndTime := existingStartTime.Add(6 * time.Hour)

		// Execute
		hasConflict, err := repo.CheckConflict(ctx, partnerEncx.ID, newStartTime, newEndTime, nil)

		// Assert
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("should ignore excluded availability ID", func(t *testing.T) {
		// Setup
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies
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

		// Create existing availability
		existingStartTime := time.Now().Add(24 * time.Hour)
		existingEndTime := existingStartTime.Add(2 * time.Hour)

		existing := ta.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		existing.StartTime = existingStartTime
		existing.EndTime = existingEndTime
		existing.Status = domain.AvailabilityStatusAvailable
		ta.InsertAvailabilityEncx(t, ctx, existing, testPool)

		// Check conflict with the same time range but exclude the existing ID
		// Execute
		hasConflict, err := repo.CheckConflict(ctx, partnerEncx.ID, existingStartTime, existingEndTime, &existing.ID)

		// Assert
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("should not detect conflict with cancelled availability", func(t *testing.T) {
		// Setup
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies
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

		// Create cancelled availability
		cancelledStartTime := time.Now().Add(24 * time.Hour)
		cancelledEndTime := cancelledStartTime.Add(2 * time.Hour)

		cancelled := ta.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		cancelled.StartTime = cancelledStartTime
		cancelled.EndTime = cancelledEndTime
		cancelled.Status = domain.AvailabilityStatusCancelled
		ta.InsertAvailabilityEncx(t, ctx, cancelled, testPool)

		// Try to create availability in the same time
		// Execute
		hasConflict, err := repo.CheckConflict(ctx, partnerEncx.ID, cancelledStartTime, cancelledEndTime, nil)

		// Assert
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("should detect conflict with blocked availability", func(t *testing.T) {
		// Setup
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies
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

		// Create blocked availability
		blockedStartTime := time.Now().Add(24 * time.Hour)
		blockedEndTime := blockedStartTime.Add(2 * time.Hour)

		blocked := ta.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		blocked.StartTime = blockedStartTime
		blocked.EndTime = blockedEndTime
		blocked.Status = domain.AvailabilityStatusBlocked
		ta.InsertAvailabilityEncx(t, ctx, blocked, testPool)

		// Try to create availability in the same time
		// Execute
		hasConflict, err := repo.CheckConflict(ctx, partnerEncx.ID, blockedStartTime, blockedEndTime, nil)

		// Assert
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("should not detect conflict with different partner", func(t *testing.T) {
		// Setup - clear all tables
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies
		buildingEncx := tb.NewTestBuildingEncx(t)
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		// Create partner 1
		userEncx1 := th.NewTestUserEncx(t)
		userEncx1.EmailHash = "email_hash1"
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

		// Create availability for partner 1
		existingStartTime := time.Now().Add(24 * time.Hour)
		existingEndTime := existingStartTime.Add(2 * time.Hour)

		existing := ta.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx1.ID, roomEncx.ID)
		existing.StartTime = existingStartTime
		existing.EndTime = existingEndTime
		existing.Status = domain.AvailabilityStatusAvailable
		ta.InsertAvailabilityEncx(t, ctx, existing, testPool)

		// Check conflict for partner 2 (different partner)
		// Execute
		hasConflict, err := repo.CheckConflict(ctx, partnerEncx2.ID, existingStartTime, existingEndTime, nil)

		// Assert
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("should detect conflict with adjacent times (end time equals start time)", func(t *testing.T) {
		// Setup
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		th.ClearPartnersTable(t, ctx, testPool)

		// Create dependencies
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

		// Create availability from 10:00 to 12:00
		existingStartTime := time.Now().Add(24 * time.Hour)
		existingEndTime := existingStartTime.Add(2 * time.Hour)

		existing := ta.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		existing.StartTime = existingStartTime
		existing.EndTime = existingEndTime
		existing.Status = domain.AvailabilityStatusAvailable
		ta.InsertAvailabilityEncx(t, ctx, existing, testPool)

		// Try to create availability from 12:00 to 14:00 (touches but doesn't overlap)
		newStartTime := existingEndTime
		newEndTime := existingEndTime.Add(2 * time.Hour)

		// Execute
		hasConflict, err := repo.CheckConflict(ctx, partnerEncx.ID, newStartTime, newEndTime, nil)

		// Assert
		require.NoError(t, err)
		assert.False(t, hasConflict) // Touching edges should not be considered conflicting
	})
}
