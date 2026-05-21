package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	tsetup "github.com/Leviosa-care/leviosa/backend/test/helpers/booking"
	tavail "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tbooking "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/booking"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"
	tcatalog "github.com/Leviosa-care/leviosa/backend/test/helpers/catalog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAdminDashboardStats(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("empty platform returns all zeros", func(t *testing.T) {
		tbooking.ClearBookingsTable(t, ctx, testPool)
		tavail.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tcatalog.ClearProductsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminToken := tu.SetupAdminUser(t, ctx, authCtx)

		req, err := http.NewRequest("GET", testServerURL+"/admin/dashboard/stats", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: adminToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var stats domain.DashboardStats
		err = json.NewDecoder(resp.Body).Decode(&stats)
		require.NoError(t, err)

		assert.Equal(t, 0, stats.BookingsThisWeek)
		assert.Equal(t, 0, stats.RevenueThisWeek)
		assert.Equal(t, 0, stats.UpcomingBookingsCount)
		assert.Equal(t, 0, stats.PendingBookingsCount)
		assert.Equal(t, 0, stats.ActiveProductsCount)
	})

	t.Run("ISO week boundary excludes last Sunday bookings", func(t *testing.T) {
		tbooking.ClearBookingsTable(t, ctx, testPool)
		tavail.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tcatalog.ClearProductsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminToken := tu.SetupAdminUser(t, ctx, authCtx)

		now := time.Now()
		startOfISOWeek := getStartOfISOWeek(now)

		lastSunday := startOfISOWeek.Add(-1 * time.Second)
		thisMonday := startOfISOWeek

		bookingEncx := tbooking.NewTestBookingEncx(t)
		bookingEncx.CreatedAt = lastSunday
		bookingEncx.TotalPriceCents = 5000
		bookingEncx.PaymentStatus = domain.PaymentStatusPaid
		tbooking.InsertBookingEncx(t, ctx, testPool, bookingEncx)

		thisWeekBooking := tbooking.NewTestBookingEncx(t)
		thisWeekBooking.CreatedAt = thisMonday
		thisWeekBooking.TotalPriceCents = 3000
		thisWeekBooking.PaymentStatus = domain.PaymentStatusPaid
		tbooking.InsertBookingEncx(t, ctx, testPool, thisWeekBooking)

		req, err := http.NewRequest("GET", testServerURL+"/admin/dashboard/stats", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: adminToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var stats domain.DashboardStats
		err = json.NewDecoder(resp.Body).Decode(&stats)
		require.NoError(t, err)

		assert.Equal(t, 1, stats.BookingsThisWeek, "should count only bookings from this ISO week")
		assert.Equal(t, 3000, stats.RevenueThisWeek, "should sum revenue only from this ISO week")
	})

	t.Run("mix of booking statuses counts correctly", func(t *testing.T) {
		tbooking.ClearBookingsTable(t, ctx, testPool)
		tavail.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tcatalog.ClearProductsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminToken := tu.SetupAdminUser(t, ctx, authCtx)

		now := time.Now()
		startOfISOWeek := getStartOfISOWeek(now)

		upcomingBooking := tbooking.NewTestBookingEncx(t)
		upcomingBooking.CreatedAt = startOfISOWeek.Add(2 * time.Hour)
		upcomingBooking.SlotStartTimeEncrypted = []byte("future_slot")
		upcomingBooking.Status = domain.BookingStatusConfirmed
		tbooking.InsertBookingEncx(t, ctx, testPool, upcomingBooking)

		pastConfirmedBooking := tbooking.NewTestBookingEncx(t)
		pastConfirmedBooking.CreatedAt = startOfISOWeek.Add(1 * time.Hour)
		pastConfirmedBooking.SlotStartTimeEncrypted = []byte("past_slot")
		pastConfirmedBooking.Status = domain.BookingStatusConfirmed
		tbooking.InsertBookingEncx(t, ctx, testPool, pastConfirmedBooking)

		cancelledBooking := tbooking.NewTestBookingEncx(t)
		cancelledBooking.CreatedAt = startOfISOWeek.Add(3 * time.Hour)
		cancelledBooking.Status = domain.BookingStatusCancelled
		tbooking.InsertBookingEncx(t, ctx, testPool, cancelledBooking)

		pendingBooking := tbooking.NewTestBookingEncx(t)
		pendingBooking.CreatedAt = startOfISOWeek.Add(4 * time.Hour)
		pendingBooking.PaymentStatus = domain.PaymentStatusPending
		tbooking.InsertBookingEncx(t, ctx, testPool, pendingBooking)

		req, err := http.NewRequest("GET", testServerURL+"/admin/dashboard/stats", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: adminToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var stats domain.DashboardStats
		err = json.NewDecoder(resp.Body).Decode(&stats)
		require.NoError(t, err)

		assert.Equal(t, 4, stats.BookingsThisWeek)
		assert.Equal(t, 1, stats.UpcomingBookingsCount, "should count only confirmed future bookings")
		assert.Equal(t, 1, stats.PendingBookingsCount, "should count pending payment bookings")
	})

	t.Run("non-admin returns 403", func(t *testing.T) {
		tbooking.ClearBookingsTable(t, ctx, testPool)
		tavail.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

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

		partnerToken, _ := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@test.com", room.ID, testPool, authCtx.Redis, crypto)

		req, err := http.NewRequest("GET", testServerURL+"/admin/dashboard/stats", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("unauthenticated returns 401", func(t *testing.T) {
		req, err := http.NewRequest("GET", testServerURL+"/admin/dashboard/stats", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func getStartOfISOWeek(t time.Time) time.Time {
	weekday := t.Weekday()
	var daysSinceMonday int
	if weekday == time.Sunday {
		daysSinceMonday = 6
	} else {
		daysSinceMonday = int(weekday - time.Monday)
	}
	return time.Date(t.Year(), t.Month(), t.Day()-daysSinceMonday, 0, 0, 0, 0, t.Location())
}
