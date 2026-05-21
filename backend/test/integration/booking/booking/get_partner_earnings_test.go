package booking

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPartnerEarnings(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("partner with no bookings should return zeroed summary and empty transactions", func(t *testing.T) {
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

		partnerToken, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@test.com", room.ID, testPool, authCtx.Redis, crypto)

		req, err := http.NewRequest("GET", testServerURL+"/partners/"+partnerID.String()+"/earnings", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var summary domain.EarningsSummary
		err = json.NewDecoder(resp.Body).Decode(&summary)
		require.NoError(t, err)

		assert.Equal(t, 0, summary.CurrentMonthCents)
		assert.Equal(t, 0, summary.LastMonthCents)
		assert.Equal(t, 0, summary.PendingCents)
		assert.Equal(t, 0, summary.NextPayoutCents)
		assert.NotEmpty(t, summary.NextPayoutDate)
		assert.Empty(t, summary.Transactions)
	})

	t.Run("partner with only pending bookings should populate pendingCents", func(t *testing.T) {
		tbooking.ClearBookingsTable(t, ctx, testPool)
		tavail.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tcatalog.ClearProductsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		categoryID := tcatalog.CreateTestCategory(t, ctx, testPool)
		products := tcatalog.CreateDefaultTestProducts(t, ctx, testPool, categoryID)

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

		partnerToken, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@test.com", room.ID, testPool, authCtx.Redis, crypto)
		clientToken, clientID := tsetup.SetupStandardUser(t, ctx, "client@test.com", room.ID, testPool, authCtx.Redis, crypto)

		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour)
		endTime := startTime.Add(2 * time.Hour)
		availability := tavail.NewTestAvailability(t, partnerID, room.ID, startTime, endTime, 5000)
		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)
		require.NoError(t, err)

		slotStartTime := startTime.Add(30 * time.Minute).Truncate(10 * time.Minute)
		requestBody := map[string]interface{}{
			"availability_id": availability.ID.String(),
			"client_id":       clientID,
			"product_id":      products[0].ID.String(),
			"slot_start_time": slotStartTime.Format(time.RFC3339),
		}
		bodyBytes, err := json.Marshal(requestBody)
		require.NoError(t, err)

		reqCreate, err := http.NewRequest("POST", testServerURL+"/bookings", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		reqCreate.Header.Set("Content-Type", "application/json")
		reqCreate.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

		respCreate, err := client.Do(reqCreate)
		require.NoError(t, err)
		respCreate.Body.Close()
		require.Equal(t, http.StatusCreated, respCreate.StatusCode)

		req, err := http.NewRequest("GET", testServerURL+"/partners/"+partnerID.String()+"/earnings", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var summary domain.EarningsSummary
		err = json.NewDecoder(resp.Body).Decode(&summary)
		require.NoError(t, err)

		assert.Equal(t, 0, summary.CurrentMonthCents, "currentMonthCents should be 0 for pending bookings")
		assert.Equal(t, 5000, summary.PendingCents, "pendingCents should reflect pending booking amount")
		assert.Equal(t, 0, summary.NextPayoutCents, "nextPayoutCents should be 0 for pending bookings")
		assert.Len(t, summary.Transactions, 1, "should have 1 transaction")
		assert.Equal(t, domain.PaymentStatusPending, summary.Transactions[0].PaymentStatus)
	})

	t.Run("partner with mix of paid/pending/refunded bookings should have correct buckets", func(t *testing.T) {
		tbooking.ClearBookingsTable(t, ctx, testPool)
		tavail.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tcatalog.ClearProductsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		categoryID := tcatalog.CreateTestCategory(t, ctx, testPool)
		products := tcatalog.CreateDefaultTestProducts(t, ctx, testPool, categoryID)

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

		partnerToken, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@test.com", room.ID, testPool, authCtx.Redis, crypto)
		clientToken, clientID := tsetup.SetupStandardUser(t, ctx, "client@test.com", room.ID, testPool, authCtx.Redis, crypto)

		createBooking := func(baseTime time.Time, price int) uuid.UUID {
			startTime := baseTime.Truncate(time.Hour)
			endTime := startTime.Add(2 * time.Hour)
			availability := tavail.NewTestAvailability(t, partnerID, room.ID, startTime, endTime, price)
			availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
			require.NoError(t, err)
			tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)
			require.NoError(t, err)

			slotStartTime := startTime.Add(30 * time.Minute).Truncate(10 * time.Minute)
			requestBody := map[string]interface{}{
				"availability_id": availability.ID.String(),
				"client_id":       clientID,
				"product_id":      products[0].ID.String(),
				"slot_start_time": slotStartTime.Format(time.RFC3339),
			}
			bodyBytes, err := json.Marshal(requestBody)
			require.NoError(t, err)

			reqCreate, err := http.NewRequest("POST", testServerURL+"/bookings", bytes.NewReader(bodyBytes))
			require.NoError(t, err)
			reqCreate.Header.Set("Content-Type", "application/json")
			reqCreate.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

			respCreate, err := client.Do(reqCreate)
			require.NoError(t, err)
			defer respCreate.Body.Close()
			require.Equal(t, http.StatusCreated, respCreate.StatusCode)

			var createdBooking domain.BookingResponse
			err = json.NewDecoder(respCreate.Body).Decode(&createdBooking)
			require.NoError(t, err)
			return createdBooking.ID
		}

		now := time.Now()

		paidBooking1 := createBooking(now.Add(24*time.Hour), 10000)
		paidBooking2 := createBooking(now.Add(48*time.Hour), 15000)
		_ = createBooking(now.Add(72*time.Hour), 8000)

		updatePaymentStatus := func(bookingID uuid.UUID, paymentStatus domain.PaymentStatus) {
			req, err := http.NewRequest("POST", testServerURL+"/bookings/"+bookingID.String()+"/payment", nil)
			require.NoError(t, err)
			req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

			paymentReq := map[string]string{"payment_intent_id": "pi_test_" + bookingID.String()}
			bodyBytes, _ := json.Marshal(paymentReq)
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			require.NoError(t, err)
			resp.Body.Close()
			require.Equal(t, http.StatusOK, resp.StatusCode)

			if paymentStatus == domain.PaymentStatusRefunded {
				reqRefund, err := http.NewRequest("POST", testServerURL+"/bookings/"+bookingID.String()+"/refund", nil)
				require.NoError(t, err)
				reqRefund.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

				respRefund, err := client.Do(reqRefund)
				require.NoError(t, err)
				respRefund.Body.Close()
				require.Equal(t, http.StatusOK, respRefund.StatusCode)
			}
		}

		updatePaymentStatus(paidBooking1, domain.PaymentStatusPaid)
		updatePaymentStatus(paidBooking2, domain.PaymentStatusPaid)

		req, err := http.NewRequest("GET", testServerURL+"/partners/"+partnerID.String()+"/earnings", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var summary domain.EarningsSummary
		err = json.NewDecoder(resp.Body).Decode(&summary)
		require.NoError(t, err)

		assert.Equal(t, 25000, summary.CurrentMonthCents, "should sum paid bookings in current month")
		assert.Equal(t, 8000, summary.PendingCents, "should sum pending bookings")
		assert.Len(t, summary.Transactions, 3, "should have 3 transactions")
	})

	t.Run("cross-month boundary bookings should count correctly", func(t *testing.T) {
		tbooking.ClearBookingsTable(t, ctx, testPool)
		tavail.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tcatalog.ClearProductsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		categoryID := tcatalog.CreateTestCategory(t, ctx, testPool)
		products := tcatalog.CreateDefaultTestProducts(t, ctx, testPool, categoryID)

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

		partnerToken, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@test.com", room.ID, testPool, authCtx.Redis, crypto)
		clientToken, clientID := tsetup.SetupStandardUser(t, ctx, "client@test.com", room.ID, testPool, authCtx.Redis, crypto)

		now := time.Now()
		lastMonth := now.AddDate(0, -1, 0)

		createPaidBooking := func(bookingTime time.Time, price int) {
			startTime := bookingTime.Truncate(time.Hour)
			endTime := startTime.Add(2 * time.Hour)
			availability := tavail.NewTestAvailability(t, partnerID, room.ID, startTime, endTime, price)
			availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
			require.NoError(t, err)
			tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)
			require.NoError(t, err)

			slotStartTime := startTime.Add(30 * time.Minute).Truncate(10 * time.Minute)
			requestBody := map[string]interface{}{
				"availability_id": availability.ID.String(),
				"client_id":       clientID,
				"product_id":      products[0].ID.String(),
				"slot_start_time": slotStartTime.Format(time.RFC3339),
			}
			bodyBytes, err := json.Marshal(requestBody)
			require.NoError(t, err)

			reqCreate, err := http.NewRequest("POST", testServerURL+"/bookings", bytes.NewReader(bodyBytes))
			require.NoError(t, err)
			reqCreate.Header.Set("Content-Type", "application/json")
			reqCreate.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

			respCreate, err := client.Do(reqCreate)
			require.NoError(t, err)
			defer respCreate.Body.Close()
			require.Equal(t, http.StatusCreated, respCreate.StatusCode)

			var createdBooking domain.BookingResponse
			err = json.NewDecoder(respCreate.Body).Decode(&createdBooking)
			require.NoError(t, err)

			reqPayment, err := http.NewRequest("POST", testServerURL+"/bookings/"+createdBooking.ID.String()+"/payment", nil)
			require.NoError(t, err)
			reqPayment.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

			paymentReq := map[string]string{"payment_intent_id": "pi_test_" + createdBooking.ID.String()}
			bodyBytes, _ = json.Marshal(paymentReq)
			reqPayment.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			reqPayment.Header.Set("Content-Type", "application/json")

			respPayment, err := client.Do(reqPayment)
			require.NoError(t, err)
			respPayment.Body.Close()
			require.Equal(t, http.StatusOK, respPayment.StatusCode)
		}

		createPaidBooking(lastMonth.Add(24*time.Hour), 5000)
		createPaidBooking(now.Add(24*time.Hour), 7000)

		req, err := http.NewRequest("GET", testServerURL+"/partners/"+partnerID.String()+"/earnings", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var summary domain.EarningsSummary
		err = json.NewDecoder(resp.Body).Decode(&summary)
		require.NoError(t, err)

		assert.Equal(t, 7000, summary.CurrentMonthCents, "should count current month paid booking")
		assert.Equal(t, 5000, summary.LastMonthCents, "should count last month paid booking")
	})

	t.Run("partner A cannot access partner B earnings - should return 403", func(t *testing.T) {
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

		partnerAToken, _ := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partnerA@test.com", room.ID, testPool, authCtx.Redis, crypto)
		_, partnerBID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partnerB@test.com", room.ID, testPool, authCtx.Redis, crypto)

		req, err := http.NewRequest("GET", testServerURL+"/partners/"+partnerBID.String()+"/earnings", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerAToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode, "partner A should not access partner B's earnings")
	})

	t.Run("admin can access any partner earnings", func(t *testing.T) {
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

		adminToken, _ := tsetup.SetupAdminWithAllocation(t, ctx, room.ID, testPool, authCtx.Redis, crypto)
		_, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@test.com", room.ID, testPool, authCtx.Redis, crypto)

		req, err := http.NewRequest("GET", testServerURL+"/partners/"+partnerID.String()+"/earnings", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: adminToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "admin should be able to access partner earnings")

		var summary domain.EarningsSummary
		err = json.NewDecoder(resp.Body).Decode(&summary)
		require.NoError(t, err)
		assert.NotNil(t, summary)
	})

	t.Run("unauthenticated request should return 401", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		partnerID := uuid.New()

		req, err := http.NewRequest("GET", testServerURL+"/partners/"+partnerID.String()+"/earnings", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "unauthenticated request should return 401")
	})

	t.Run("next payout date should always be next Monday from current date", func(t *testing.T) {
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

		partnerToken, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@test.com", room.ID, testPool, authCtx.Redis, crypto)

		req, err := http.NewRequest("GET", testServerURL+"/partners/"+partnerID.String()+"/earnings", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var summary domain.EarningsSummary
		err = json.NewDecoder(resp.Body).Decode(&summary)
		require.NoError(t, err)

		parsedDate, err := time.Parse(time.RFC3339, summary.NextPayoutDate)
		require.NoError(t, err)

		now := time.Now()
		expectedNextMonday := nextMonday(now)

		assert.Equal(t, expectedNextMonday.Year(), parsedDate.Year())
		assert.Equal(t, expectedNextMonday.Month(), parsedDate.Month())
		assert.Equal(t, expectedNextMonday.Day(), parsedDate.Day())
		assert.Equal(t, time.Monday, parsedDate.Weekday(), "next payout date should be a Monday")
	})
}

func nextMonday(t time.Time) time.Time {
	daysUntilMonday := (7 - int(t.Weekday())) % 7
	if daysUntilMonday == 0 {
		daysUntilMonday = 7
	}
	return t.AddDate(0, 0, daysUntilMonday)
}
