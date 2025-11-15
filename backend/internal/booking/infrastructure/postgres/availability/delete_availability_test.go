package availabilityRepository_test

import (
	"context"
	"testing"

	th "github.com/Leviosa-care/leviosa/backend/test/helpers"
	availabilityHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestDelete TEST_PATH=internal/booking/infrastructure/postgres/availability/delete_availability_test.go

func TestDelete(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete availability successfully", func(t *testing.T) {
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

		availability := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, availability, testPool)

		// Verify it exists
		require.True(t, availabilityHelpers.AvailabilityExistsInTable(t, ctx, availability.ID, testPool))
		require.Equal(t, 1, availabilityHelpers.CountAvailabilitiesInTable(t, ctx, testPool))

		// Execute
		err = repo.Delete(ctx, availability.ID)

		// Assert
		require.NoError(t, err)
		assert.False(t, availabilityHelpers.AvailabilityExistsInTable(t, ctx, availability.ID, testPool))
		assert.Equal(t, 0, availabilityHelpers.CountAvailabilitiesInTable(t, ctx, testPool))
	})

	t.Run("should return not found when deleting non-existent availability", func(t *testing.T) {
		// Setup - clear all tables
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)

		nonExistentID := uuid.New()

		// Execute
		err := repo.Delete(ctx, nonExistentID)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should delete recurring availability", func(t *testing.T) {
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

		availability := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		availability.PartnerID = partnerEncx.ID
		availability.RoomID = roomEncx.ID
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, availability, testPool)

		// Verify it exists and is recurring
		require.True(t, availabilityHelpers.AvailabilityExistsInTable(t, ctx, availability.ID, testPool))
		require.Equal(t, 1, availabilityHelpers.CountRecurringAvailabilities(t, ctx, testPool))

		// Execute
		err = repo.Delete(ctx, availability.ID)

		// Assert
		require.NoError(t, err)
		assert.False(t, availabilityHelpers.AvailabilityExistsInTable(t, ctx, availability.ID, testPool))
		assert.Equal(t, 0, availabilityHelpers.CountRecurringAvailabilities(t, ctx, testPool))
	})
}
