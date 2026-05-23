package bookingHandler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) GetAdminBookings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter, err := parseAdminBookingsFilter(r)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	result, err := h.svc.GetAdminBookings(ctx, filter)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	httpx.RespondWithJSON(w, result, http.StatusOK)
}

func parseAdminBookingsFilter(r *http.Request) (ports.AdminBookingsFilter, error) {
	filter := ports.AdminBookingsFilter{
		Page:  1,
		Limit: 20,
	}

	q := r.URL.Query()

	if statusStr := q.Get("status"); statusStr != "" {
		status := domain.BookingStatus(statusStr)
		switch status {
		case domain.BookingStatusConfirmed, domain.BookingStatusCompleted,
			domain.BookingStatusCancelled, domain.BookingStatusNoShow:
			filter.Status = &status
		default:
			return filter, errs.NewInvalidValueErr("invalid status filter")
		}
	}

	if partnerIDStr := q.Get("partner_id"); partnerIDStr != "" {
		partnerID, err := uuid.Parse(partnerIDStr)
		if err != nil {
			return filter, errs.NewInvalidValueErr("invalid partner_id")
		}
		filter.PartnerID = &partnerID
	}

	if fromStr := q.Get("from"); fromStr != "" {
		from, err := time.Parse(time.DateOnly, fromStr)
		if err != nil {
			return filter, errs.NewInvalidValueErr("invalid from date, expected YYYY-MM-DD")
		}
		filter.From = &from
	}

	if toStr := q.Get("to"); toStr != "" {
		to, err := time.Parse(time.DateOnly, toStr)
		if err != nil {
			return filter, errs.NewInvalidValueErr("invalid to date, expected YYYY-MM-DD")
		}
		// End of day
		endOfDay := to.Add(24*time.Hour - time.Nanosecond)
		filter.To = &endOfDay
	}

	if pageStr := q.Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page >= 1 {
			filter.Page = page
		}
	}

	if limitStr := q.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit >= 1 && limit <= 100 {
			filter.Limit = limit
		}
	}

	return filter, nil
}
