package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/middleware"
)

func (h *handler) GetCompanyInstagram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		middleware.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	response, err := h.svc.GetCompanyInstagram(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			middleware.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			middleware.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.ErrorContext(ctx, fmt.Sprintf("Handler: Internal server error during company instagram retrieval: %v", err))
			middleware.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.ErrorContext(ctx, fmt.Sprintf("Handler: Unhandled error from service during company instagram retrieval: %v", err))
			middleware.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	middleware.RespondWithJSON(w, response, http.StatusOK)
}

