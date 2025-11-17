package couponHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) ValidateCoupon(w http.ResponseWriter, r *http.Request) {
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

	logger.Info("Handler: Processing validate_coupon", "coupon_code", "")

	var validateRequest struct {
		StripeCouponID string `json:"stripeCouponId"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&validateRequest); err != nil {
		logger.Error("Handler: Error decoding JSON body", "error", err)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	coupon, err := h.svc.ValidateCoupon(ctx, validateRequest.StripeCouponID)
	if err != nil {
		// Special handling for validation errors - return JSON response instead of error
		if errors.Is(err, errs.ErrInvalidValue) {
			// Return validation failure response instead of error for business validation
			httpx.RespondWithJSON(w, struct {
				Valid  bool   `json:"valid"`
				Reason string `json:"reason"`
			}{
				Valid:  false,
				Reason: err.Error(),
			}, http.StatusOK)
			return
		}
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			// Return validation failure response for not found
			httpx.RespondWithJSON(w, struct {
				Valid  bool   `json:"valid"`
				Reason string `json:"reason"`
			}{
				Valid:  false,
				Reason: "coupon not found",
			}, http.StatusOK)
			return
		}
		// For all other errors, use the centralized error handler
		httpx.RespondWithServiceError(w, logger, ctx, err, "validate coupon")
		return
	}

	// Return successful validation with coupon details
	logger.Info("Handler: Coupon validation successful", "coupon_code", validateRequest.StripeCouponID)
	httpx.RespondWithJSON(w, struct {
		Valid  bool        `json:"valid"`
		Coupon interface{} `json:"coupon"`
	}{
		Valid:  true,
		Coupon: coupon,
	}, http.StatusOK)
}
