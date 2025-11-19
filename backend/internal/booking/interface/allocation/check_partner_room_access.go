package allocationHandler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) CheckPartnerRoomAccess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	_ = logger
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract partner ID and room ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 5 {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid URL path"), http.StatusBadRequest)
		return
	}

	partnerID, err := uuid.Parse(pathParts[1])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
		return
	}

	roomID, err := uuid.Parse(pathParts[3])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid room ID format"), http.StatusBadRequest)
		return
	}

	// Parse query parameter for time check (default to now)
	checkTime := time.Now()
	if timeStr := r.URL.Query().Get("at"); timeStr != "" {
		if parsedTime, err := time.Parse(time.RFC3339, timeStr); err == nil {
			checkTime = parsedTime
		}
	}

	// Call service to check access
	hasAccess, err := h.svc.CheckPartnerRoomAccess(ctx, partnerID, roomID, checkTime)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
		default:
			statusCode = http.StatusInternalServerError
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	response := struct {
		HasAccess bool      `json:"has_access"`
		CheckedAt time.Time `json:"checked_at"`
	}{
		HasAccess: hasAccess,
		CheckedAt: checkTime,
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
