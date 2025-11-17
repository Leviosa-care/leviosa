package couponHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAllCoupons(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)

	coupons, err := h.svc.GetAllCoupons(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get all coupons")
		return
	}

	logger.Info("Handler: All coupons retrieval successful")
	httpx.RespondWithJSON(w, coupons, http.StatusOK)
}
