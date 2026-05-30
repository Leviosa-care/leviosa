package booking

import (
	"context"
	"testing"
	"time"

	bookingService "github.com/Leviosa-care/leviosa/backend/internal/booking/application/booking"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"

	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	tsetup "github.com/Leviosa-care/leviosa/backend/test/helpers/booking"
	tavail "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tbooking "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/booking"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"
	tcatalog "github.com/Leviosa-care/leviosa/backend/test/helpers/catalog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
)

func TestReminderScheduler_OnlyConfirmedUnremindedBookings(t *testing.T) {
	ctx := context.Background()

	// Clean state
	tbooking.ClearBookingsTable(t, ctx, testPool)
	tavail.ClearAvailabilityTable(t, ctx, testPool)
	tr.ClearRoomsTable(t, ctx, testPool)
	tb.ClearBuildingsTable(t, ctx, testPool)
	tcatalog.ClearProductsTable(t, ctx, testPool)
	defer tu.ClearAuthData(t, ctx, authCtx)

	// Setup test category and products
	categoryID := tcatalog.CreateTestCategory(t, ctx, testPool)
	products := tcatalog.CreateDefaultTestProducts(t, ctx, testPool, categoryID)

	// Setup test building and room
	building := tb.NewTestBuilding(t)
	buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
	require.NoError(t, err)
	err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
	require.NoError(t, err)

	room := tr.NewTestRoomWithBuilding(t, building.ID)
	roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
	require.NoError(t, err)
	err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
	require.NoError(t, err)

	// Setup users
	_, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner-reminder@example.com", room.ID, testPool, authCtx.Redis, crypto)

	// Create three bookings:
	// 1. Due within the window, reminded_at IS NULL → should be reminded
	// 2. Outside the window → should NOT be reminded
	// 3. Already reminded → should NOT be reminded again

	now := time.Now()

	// Booking 1: due in 12 hours (within 24h window)
	startTime1 := now.Add(12 * time.Hour).Truncate(time.Hour)
	endTime1 := startTime1.Add(2 * time.Hour)
	availability1 := tavail.NewTestAvailability(t, partnerID, room.ID, startTime1, endTime1, 0)
	availabilityEncx1, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability1)
	require.NoError(t, err)
	tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx1, testPool)

	_, clientID1 := tsetup.SetupStandardUser(t, ctx, "client-reminder1@example.com", room.ID, testPool, authCtx.Redis, crypto)

	booking1 := createBookingViaService(t, ctx, availability1.ID, clientID1, partnerID, room.ID, products[0].ID, startTime1.Add(30*time.Minute))

	// Booking 2: due in 48 hours (outside 24h window)
	startTime2 := now.Add(48 * time.Hour).Truncate(time.Hour)
	endTime2 := startTime2.Add(2 * time.Hour)
	availability2 := tavail.NewTestAvailability(t, partnerID, room.ID, startTime2, endTime2, 0)
	availabilityEncx2, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability2)
	require.NoError(t, err)
	tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx2, testPool)

	_, clientID2 := tsetup.SetupStandardUser(t, ctx, "client-reminder2@example.com", room.ID, testPool, authCtx.Redis, crypto)

	booking2 := createBookingViaService(t, ctx, availability2.ID, clientID2, partnerID, room.ID, products[0].ID, startTime2.Add(30*time.Minute))

	// Booking 3: due in 6 hours but already reminded
	startTime3 := now.Add(6 * time.Hour).Truncate(time.Hour)
	endTime3 := startTime3.Add(2 * time.Hour)
	availability3 := tavail.NewTestAvailability(t, partnerID, room.ID, startTime3, endTime3, 0)
	availabilityEncx3, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability3)
	require.NoError(t, err)
	tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx3, testPool)

	_, clientID3 := tsetup.SetupStandardUser(t, ctx, "client-reminder3@example.com", room.ID, testPool, authCtx.Redis, crypto)

	booking3 := createBookingViaService(t, ctx, availability3.ID, clientID3, partnerID, room.ID, products[0].ID, startTime3.Add(30*time.Minute))

	// Mark booking 3 as already reminded
	err = bookingRepo.MarkReminderSent(ctx, booking3.ID)
	require.NoError(t, err)

	// Reset spy
	notificationSpy.BookingReminders = nil

	// Create scheduler with 24h window and run one tick
	scheduler := bookingService.NewReminderScheduler(
		bookingRepo,
		notificationSpy,
		crypto,
		bookingService.WithReminderWindow(24*time.Hour),
	)
	scheduler.TickOnce(ctx)

	// Assert: exactly one reminder sent (for booking 1)
	assert.Len(t, notificationSpy.BookingReminders, 1, "expected exactly one reminder notification")
	if len(notificationSpy.BookingReminders) == 1 {
		assert.Equal(t, booking1.ID, notificationSpy.BookingReminders[0].BookingID)
	}

	// Assert: booking 1 has reminded_at set
	encx1, err := bookingRepo.GetByID(ctx, booking1.ID)
	require.NoError(t, err)
	require.NotNil(t, encx1.RemindedAt, "booking 1 should have reminded_at set")

	// Assert: booking 2 does NOT have reminded_at set
	encx2, err := bookingRepo.GetByID(ctx, booking2.ID)
	require.NoError(t, err)
	assert.Nil(t, encx2.RemindedAt, "booking 2 should NOT have reminded_at set (outside window)")

	// Assert: booking 3 still has its original reminded_at (not updated again)
	encx3, err := bookingRepo.GetByID(ctx, booking3.ID)
	require.NoError(t, err)
	assert.NotNil(t, encx3.RemindedAt, "booking 3 should still have reminded_at set")
}

func TestReminderScheduler_MarkSentOnNotificationError(t *testing.T) {
	ctx := context.Background()

	// Clean state
	tbooking.ClearBookingsTable(t, ctx, testPool)
	tavail.ClearAvailabilityTable(t, ctx, testPool)
	tr.ClearRoomsTable(t, ctx, testPool)
	tb.ClearBuildingsTable(t, ctx, testPool)
	tcatalog.ClearProductsTable(t, ctx, testPool)
	defer tu.ClearAuthData(t, ctx, authCtx)

	// Setup test category and products
	categoryID := tcatalog.CreateTestCategory(t, ctx, testPool)
	products := tcatalog.CreateDefaultTestProducts(t, ctx, testPool, categoryID)

	// Setup test building and room
	building := tb.NewTestBuilding(t)
	buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
	require.NoError(t, err)
	err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
	require.NoError(t, err)

	room := tr.NewTestRoomWithBuilding(t, building.ID)
	roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
	require.NoError(t, err)
	err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
	require.NoError(t, err)

	_, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner-reminder-err@example.com", room.ID, testPool, authCtx.Redis, crypto)

	now := time.Now()
	startTime := now.Add(12 * time.Hour).Truncate(time.Hour)
	endTime := startTime.Add(2 * time.Hour)
	availability := tavail.NewTestAvailability(t, partnerID, room.ID, startTime, endTime, 0)
	availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
	require.NoError(t, err)
	tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

	_, clientID := tsetup.SetupStandardUser(t, ctx, "client-reminder-err@example.com", room.ID, testPool, authCtx.Redis, crypto)

	booking := createBookingViaService(t, ctx, availability.ID, clientID, partnerID, room.ID, products[0].ID, startTime.Add(30*time.Minute))

	// Create a notification service that always errors
	errorNotif := &errorNotificationService{}
	scheduler := bookingService.NewReminderScheduler(
		bookingRepo,
		errorNotif,
		crypto,
		bookingService.WithReminderWindow(24*time.Hour),
	)
	scheduler.TickOnce(ctx)

	// Assert: reminded_at is set even though notification failed
	encx, err := bookingRepo.GetByID(ctx, booking.ID)
	require.NoError(t, err)
	assert.NotNil(t, encx.RemindedAt, "reminded_at should be set even when notification fails")
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// createBookingViaService creates a booking through the service layer and returns
// the decrypted booking domain object.
func createBookingViaService(t *testing.T, ctx context.Context, availabilityID, clientID, partnerID, roomID, productID uuid.UUID, slotStart time.Time) *domain.Booking {
	t.Helper()
	b, err := service.CreateBooking(ctx, availabilityID, &clientID, productID, slotStart, "", "", "", "", "")
	require.NoError(t, err)
	return b
}

// errorNotificationService always returns an error for every method.
type errorNotificationService struct{}

func (e *errorNotificationService) SendBookingConfirmation(_ context.Context, _ ports.BookingNotificationData) error {
	return assert.AnError
}
func (e *errorNotificationService) SendBookingCancellation(_ context.Context, _ ports.BookingNotificationData) error {
	return assert.AnError
}
func (e *errorNotificationService) SendBookingReminder(_ context.Context, _ ports.BookingNotificationData) error {
	return assert.AnError
}
func (e *errorNotificationService) SendPaymentConfirmation(_ context.Context, _ ports.BookingNotificationData) error {
	return assert.AnError
}
func (e *errorNotificationService) SendPaymentFailed(_ context.Context, _ ports.BookingNotificationData) error {
	return assert.AnError
}

// compile-time check
var _ ports.BookingNotificationService = (*errorNotificationService)(nil)
