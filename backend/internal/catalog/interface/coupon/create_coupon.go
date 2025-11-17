package couponHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) CreateCoupon(w http.ResponseWriter, r *http.Request) {
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

	logger.InfoContext(ctx, "Handler: Processing create coupon",
		"operation", "create_coupon",
		"method", r.Method,
		"path", r.URL.Path)

	var coupon domain.CreateCouponRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&coupon); err != nil {
		logger.ErrorContext(ctx, "Handler: create coupon failed",
			"operation", "create_coupon",
			"error_context", "invalid JSON request body",
			"status_code", http.StatusBadRequest,
			"error", err)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	couponID, err := h.svc.CreateCoupon(ctx, &coupon)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "create coupon")
		return
	}

	logger.InfoContext(ctx, "Handler: create coupon completed",
		"operation", "create_coupon",
		"coupon_id", couponID,
		"status_code", http.StatusCreated)

	httpx.RespondWithJSON(
		w,
		struct {
			ID      string `json:"id"`
			Message string `json:"message"`
		}{
			ID:      couponID,
			Message: "Coupon created successfully!",
		},
		http.StatusCreated,
	)
}
