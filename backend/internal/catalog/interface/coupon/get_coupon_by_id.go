package couponHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetCouponByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("Handler: Processing get_coupon", "coupon_id", "")

	couponID := strings.TrimPrefix(r.URL.Path, "/admin/coupons/")
	if couponID == "" || strings.Contains(couponID, "/") {
		httpx.RespondWithError(w, errors.New("invalid URL path"), http.StatusBadRequest)
		return
	}

	coupon, err := h.svc.GetCouponByID(ctx, couponID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get coupon by ID")
		return
	}

	logger.Info("Handler: Coupon retrieval successful", "coupon_id", couponID)
	httpx.RespondWithJSON(w, coupon, http.StatusOK)
}
