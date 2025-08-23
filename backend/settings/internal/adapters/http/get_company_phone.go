package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
)

func (h *handler) GetCompanyPhone(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	response, err := h.svc.GetCompanyTelephone(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			log.Printf("Handler: Internal server error during company phone retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			log.Printf("Handler: Unhandled error from service during company phone retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}