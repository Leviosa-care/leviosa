package allocationHandler

import (
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
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing check partner room access request",
		"operation", "check_partner_room_access",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Extract partner ID and room ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 5 {
		logger.WarnContext(ctx, "Handler: Invalid URL path",
			"error", "invalid URL path",
			"path", r.URL.Path,
			"operation", "check_partner_room_access",
			"method", r.Method)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid URL path"), http.StatusBadRequest)
		return
	}

	partnerID, err := uuid.Parse(pathParts[1])
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid partner ID format",
			"error", err,
			"raw_partner_id", pathParts[1],
			"operation", "check_partner_room_access",
			"method", r.Method)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
		return
	}

	roomID, err := uuid.Parse(pathParts[3])
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid room ID format",
			"error", err,
			"raw_room_id", pathParts[3],
			"operation", "check_partner_room_access",
			"method", r.Method)
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
		httpx.RespondWithServiceError(w, logger, ctx, err, "check partner room access")
		return
	}

	response := struct {
		HasAccess bool      `json:"has_access"`
		CheckedAt time.Time `json:"checked_at"`
	}{
		HasAccess: hasAccess,
		CheckedAt: checkTime,
	}

	logger.InfoContext(ctx, "Handler: Partner room access check completed",
		"partner_id", partnerID,
		"room_id", roomID,
		"has_access", hasAccess,
		"check_time", checkTime,
		"operation", "check_partner_room_access")

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
