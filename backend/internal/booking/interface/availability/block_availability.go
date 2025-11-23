package availabilityHandler

import (
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) BlockAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract availability ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid availability ID in path"), http.StatusBadRequest)
		return
	}

	availabilityID, err := uuid.Parse(pathParts[1])
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid availability ID format"), http.StatusBadRequest)
		return
	}

	// Call service to block availability
	err = h.svc.BlockAvailability(ctx, availabilityID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "block availability")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
