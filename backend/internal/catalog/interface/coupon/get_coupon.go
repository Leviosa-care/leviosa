package couponHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
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
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.Error("Handler: Internal server error during coupon retrieval", "error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.Error("Handler: Unhandled error from service during coupon retrieval", "error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.Info("Handler: Coupon retrieval successful", "coupon_id", couponID)
	httpx.RespondWithJSON(w, coupon, http.StatusOK)
}

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
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.Error("Handler: Internal server error during coupon retrieval by Stripe ID", "error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.Error("Handler: Unhandled error from service during coupon retrieval by Stripe ID", "error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.Info("Handler: Coupon retrieval by Stripe ID successful", "stripe_id", stripeID)
	httpx.RespondWithJSON(w, coupon, http.StatusOK)
}

func (h *handler) GetAllCoupons(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)

	coupons, err := h.svc.GetAllCoupons(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.Error("Handler: Internal server error during coupons retrieval", "error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.Error("Handler: Unhandled error from service during coupons retrieval", "error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.Info("Handler: All coupons retrieval successful")
	httpx.RespondWithJSON(w, coupons, http.StatusOK)
}

func (h *handler) GetValidCoupons(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)

	coupons, err := h.svc.GetValidCoupons(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.Error("Handler: Internal server error during valid coupons retrieval", "error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.Error("Handler: Unhandled error from service during valid coupons retrieval", "error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.Info("Handler: Valid coupons retrieval successful")
	httpx.RespondWithJSON(w, coupons, http.StatusOK)
}
