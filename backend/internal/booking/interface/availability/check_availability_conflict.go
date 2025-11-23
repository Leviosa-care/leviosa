package availabilityHandler

import (
	"net/http"
	"strings"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) CheckAvailabilityConflict(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract partner ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID in path"), http.StatusBadRequest)
		return
	}

	partnerID, err := uuid.Parse(pathParts[1])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
		return
	}

	// Parse query parameters
	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")
	excludeIDStr := r.URL.Query().Get("exclude_id")

	if startTimeStr == "" || endTimeStr == "" {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("start_time and end_time query parameters required"), http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid start_time format"), http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid end_time format"), http.StatusBadRequest)
		return
	}

	var excludeID *uuid.UUID
	if excludeIDStr != "" {
		if parsed, err := uuid.Parse(excludeIDStr); err == nil {
			excludeID = &parsed
		}
	}

	// Call service to check conflict
	hasConflict, err := h.svc.CheckAvailabilityConflict(ctx, partnerID, startTime, endTime, excludeID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "check availability conflict")
		return
	}

	response := struct {
		HasConflict bool `json:"has_conflict"`
	}{
		HasConflict: hasConflict,
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
