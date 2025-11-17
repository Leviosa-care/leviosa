package couponHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetValidCoupons(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)

	coupons, err := h.svc.GetValidCoupons(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get valid coupons")
		return
	}

	logger.Info("Handler: Valid coupons retrieval successful")
	httpx.RespondWithJSON(w, coupons, http.StatusOK)
}
