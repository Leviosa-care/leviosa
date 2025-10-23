package productHandler

import (
	"errors"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAllPublishedProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	products, err := h.aggr.GetAllPublishedProducts(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrQueryFailed):
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
		default:
			httpx.RespondWithError(w, errors.New("internal server occurred"), http.StatusInternalServerError)
		}
		return
	}

	httpx.RespondWithJSON(w, products, http.StatusOK)
}
