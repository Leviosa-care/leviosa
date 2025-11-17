package couponHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) UpdateCoupon(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httpx.RespondWithError(w, errors.New("unsupported media type: please send 'application/json'"), http.StatusUnsupportedMediaType)
		return
	}

	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("Handler: Processing update_coupon", "coupon_id", "")

	couponID := strings.TrimPrefix(r.URL.Path, "/admin/coupons/")
	if couponID == "" || strings.Contains(couponID, "/") {
		httpx.RespondWithError(w, errors.New("invalid URL path"), http.StatusBadRequest)
		return
	}

	var updateRequest domain.UpdateCouponRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&updateRequest); err != nil {
		logger.Error("Handler: Error decoding JSON body", "error", err)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	err = h.svc.UpdateCoupon(ctx, couponID, &updateRequest)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "update coupon")
		return
	}

	logger.Info("Handler: Coupon update successful", "coupon_id", couponID)
	httpx.RespondWithJSON(
		w,
		struct {
			Message string `json:"message"`
		}{
			Message: "Coupon updated successfully!",
		},
		http.StatusOK,
	)
}
