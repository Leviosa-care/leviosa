package couponHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetCouponByStripeID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)

	stripeID := strings.TrimPrefix(r.URL.Path, "/admin/coupons/stripe/")
	if stripeID == "" || strings.Contains(stripeID, "/") {
		httpx.RespondWithError(w, errors.New("invalid Stripe ID"), http.StatusBadRequest)
		return
	}

	coupon, err := h.svc.GetCouponByStripeID(ctx, stripeID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get coupon by Stripe ID")
		return
	}

	logger.Info("Handler: Coupon retrieval by Stripe ID successful", "stripe_id", stripeID)
	httpx.RespondWithJSON(w, coupon, http.StatusOK)
}
