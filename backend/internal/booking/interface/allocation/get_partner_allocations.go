package allocationHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) GetPartnerAllocations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get partner allocations request",
		"operation", "get_partner_allocations",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Extract partner ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		logger.WarnContext(ctx, "Handler: Invalid partner ID in path",
			"error", "invalid partner ID in path",
			"path", r.URL.Path,
			"operation", "get_partner_allocations",
			"method", r.Method)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID in path"), http.StatusBadRequest)
		return
	}

	partnerID, err := uuid.Parse(pathParts[1])
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid partner ID format",
			"error", err,
			"raw_partner_id", pathParts[1],
			"operation", "get_partner_allocations",
			"method", r.Method)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
		return
	}

	// Parse query parameters
	activeOnly := r.URL.Query().Get("active_only") != "false" // Default to active only

	// Call service to get partner allocations
	allocations, err := h.svc.GetPartnerAllocations(ctx, partnerID, activeOnly)
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

	// Convert to response DTOs
	var responses []domain.RoomAllocationResponse
	for _, allocation := range allocations {
		responses = append(responses, domain.RoomAllocationResponse{
			ID:             allocation.ID,
			RoomID:         allocation.RoomID,
			UserID:         allocation.UserID,
			AllocationType: allocation.AllocationType,
			StartDate:      allocation.StartDate,
			EndDate:        allocation.EndDate,
			IsActive:       allocation.IsActive,
		})
	}

	logger.InfoContext(ctx, "Handler: Partner allocations retrieved successfully",
		"partner_id", partnerID,
		"active_only", activeOnly,
		"allocation_count", len(responses),
		"operation", "get_partner_allocations")

	httpx.RespondWithJSON(w, responses, http.StatusOK)
}
