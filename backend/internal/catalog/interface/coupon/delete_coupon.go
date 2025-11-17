package couponHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) DeleteCoupon(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("Handler: Processing delete_coupon", "coupon_id", "")

	couponID := strings.TrimPrefix(r.URL.Path, "/admin/coupons/")
	if couponID == "" || strings.Contains(couponID, "/") {
		httpx.RespondWithError(w, errors.New("invalid URL path"), http.StatusBadRequest)
		return
	}

	err = h.svc.DeleteCoupon(ctx, couponID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "delete coupon")
		return
	}

	logger.Info("Handler: Coupon deletion successful", "coupon_id", couponID)
	httpx.RespondWithJSON(
		w,
		struct {
			Message string `json:"message"`
		}{
			Message: "Coupon deleted successfully!",
		},
		http.StatusOK,
	)
}
