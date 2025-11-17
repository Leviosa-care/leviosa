package couponHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) DeactivateCoupon(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	logger.Info("Handler: Processing deactivate_coupon", "coupon_id", "")

	// Extract ID from URL path like /admin/coupons/{id}/deactivate
	urlPath := r.URL.Path
	parts := strings.Split(urlPath, "/")
	if len(parts) < 4 || parts[len(parts)-1] != "deactivate" {
		httpx.RespondWithError(w, errors.New("invalid URL path"), http.StatusBadRequest)
		return
	}

	couponID := parts[len(parts)-2]
	if couponID == "" {
		httpx.RespondWithError(w, errors.New("invalid coupon ID"), http.StatusBadRequest)
		return
	}

	err = h.svc.DeactivateCoupon(ctx, couponID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "deactivate coupon")
		return
	}

	logger.Info("Handler: Coupon deactivation successful", "coupon_id", couponID)
	httpx.RespondWithJSON(
		w,
		struct {
			Message string `json:"message"`
		}{
			Message: "Coupon deactivated successfully!",
		},
		http.StatusOK,
	)
}
