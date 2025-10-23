package couponHandler

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetCouponByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
			log.Printf("Handler: Internal server error during coupon retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			log.Printf("Handler: Unhandled error from service during coupon retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	httpx.RespondWithJSON(w, coupon, http.StatusOK)
}

func (h *handler) GetCouponByStripeID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
			log.Printf("Handler: Internal server error during coupon retrieval by Stripe ID: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			log.Printf("Handler: Unhandled error from service during coupon retrieval by Stripe ID: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	httpx.RespondWithJSON(w, coupon, http.StatusOK)
}

func (h *handler) GetAllCoupons(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	coupons, err := h.svc.GetAllCoupons(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			log.Printf("Handler: Internal server error during coupons retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			log.Printf("Handler: Unhandled error from service during coupons retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	httpx.RespondWithJSON(w, coupons, http.StatusOK)
}

func (h *handler) GetValidCoupons(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	coupons, err := h.svc.GetValidCoupons(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			log.Printf("Handler: Internal server error during valid coupons retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			log.Printf("Handler: Unhandled error from service during valid coupons retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	httpx.RespondWithJSON(w, coupons, http.StatusOK)
}
