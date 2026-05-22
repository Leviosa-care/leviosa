package availabilityRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"
	availabilityHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreate TEST_PATH=internal/booking/infrastructure/postgres/availability/create_availability_test.go

func TestCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("should create availability successfully", func(t *testing.T) {
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

		// Create test encrypted availability with valid foreign keys
		availabilityEncx := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)

		// Execute
		err = repo.Create(ctx, availabilityEncx)

		// Assert
		require.NoError(t, err)
		assert.True(t, availabilityHelpers.AvailabilityExistsInTable(t, ctx, availabilityEncx.ID, testPool))
		assert.Equal(t, 1, availabilityHelpers.CountAvailabilitiesInTable(t, ctx, testPool))
	})

	t.Run("should fail when ID already exists", func(t *testing.T) {
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

		availabilityEncx := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)

		// Insert first time - should succeed
		err = repo.Create(ctx, availabilityEncx)
		require.NoError(t, err)

		// Try to insert same ID again - should fail
		err = repo.Create(ctx, availabilityEncx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unique constraint")
	})

	t.Run("should create recurring availability successfully", func(t *testing.T) {
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

		// Create test recurring encrypted availability
		availabilityEncx := availabilityHelpers.NewTestRecurringAvailabilityEncx(t)
		availabilityEncx.UserID = partnerEncx.ID
		availabilityEncx.RoomID = roomEncx.ID

		// Execute
		err = repo.Create(ctx, availabilityEncx)

		// Assert
		require.NoError(t, err)
		assert.True(t, availabilityHelpers.AvailabilityExistsInTable(t, ctx, availabilityEncx.ID, testPool))
		assert.Equal(t, 1, availabilityHelpers.CountRecurringAvailabilities(t, ctx, testPool))
	})

	t.Run("should create availability with different statuses", func(t *testing.T) {
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

		testCases := []struct {
			name   string
			status domain.AvailabilityStatus
		}{
			{"available", domain.AvailabilityStatusAvailable},
			{"booked", domain.AvailabilityStatusBooked},
			{"cancelled", domain.AvailabilityStatusCancelled},
			{"blocked", domain.AvailabilityStatusBlocked},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				availabilityEncx := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
				availabilityEncx.Status = tc.status

				err := repo.Create(ctx, availabilityEncx)
				require.NoError(t, err)

				// Verify it was created with correct status
				retrieved := availabilityHelpers.GetAvailabilityEncxFromDB(t, ctx, availabilityEncx.ID, testPool)
				assert.Equal(t, tc.status, retrieved.Status)
			})
		}
	})

	t.Run("should create availability with null price", func(t *testing.T) {
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

		availabilityEncx := availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		availabilityEncx.PriceCents = nil // Free session

		// Execute
		err = repo.Create(ctx, availabilityEncx)

		// Assert
		require.NoError(t, err)
		retrieved := availabilityHelpers.GetAvailabilityEncxFromDB(t, ctx, availabilityEncx.ID, testPool)
		assert.Nil(t, retrieved.PriceCents)
	})

	t.Run("should create multiple availabilities successfully", func(t *testing.T) {
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

		// Create multiple availabilities
		availabilities := make([]*domain.AvailabilityEncx, 3)
		for i := 0; i < 3; i++ {
			availabilities[i] = availabilityHelpers.NewTestAvailabilityEncxWithPartnerAndRoom(t, partnerEncx.ID, roomEncx.ID)
		}

		// Execute
		for _, availability := range availabilities {
			err := repo.Create(ctx, availability)
			require.NoError(t, err)
		}

		// Assert
		assert.Equal(t, 3, availabilityHelpers.CountAvailabilitiesInTable(t, ctx, testPool))
	})

	t.Run("should preserve all encrypted fields", func(t *testing.T) {
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

		// Execute
		err = repo.Create(ctx, original)
		require.NoError(t, err)

		// Assert - verify encrypted fields are preserved
		retrieved := availabilityHelpers.GetAvailabilityEncxFromDB(t, ctx, original.ID, testPool)
		assert.Equal(t, original.ServiceTypeEncrypted, retrieved.ServiceTypeEncrypted)
		assert.Equal(t, original.NotesEncrypted, retrieved.NotesEncrypted)
		assert.Equal(t, original.DEKEncrypted, retrieved.DEKEncrypted)
		assert.Equal(t, original.KeyVersion, retrieved.KeyVersion)
	})
}
