package aggregatorHandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
	"github.com/google/uuid"
)

func (h *handler) DeleteUserByAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract user ID from URL path
	userIDStr := r.PathValue("id")
	if userIDStr == "" {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("user ID is required"), http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid user ID format: %v", err)), http.StatusBadRequest)
		return
	}

	// Call the aggregator service to delete the user
	err = h.svc.DeleteUserByAdmin(ctx, userID)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrInvalidValue):
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrResourceExhausted):
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, context.DeadlineExceeded):
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, errs.ErrTransactionFailure):
			statusCode = http.StatusServiceUnavailable
		default:
			statusCode = http.StatusInternalServerError
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Respond with success message
	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
	}{
		Message: "User deleted successfully",
	}, http.StatusOK)
}

