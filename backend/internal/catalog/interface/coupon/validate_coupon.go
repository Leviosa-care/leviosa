package couponHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) ValidateCoupon(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httpx.RespondWithError(w, errors.New("unsupported media type: please send 'application/json'"), http.StatusUnsupportedMediaType)
		return
	}

	ctx := r.Context()

	var validateRequest struct {
		StripeCouponID string `json:"stripeCouponId"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&validateRequest); err != nil {
		log.Printf("Handler: Error decoding JSON body: %v", err)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	coupon, err := h.svc.ValidateCoupon(ctx, validateRequest.StripeCouponID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			// Return validation failure response instead of error for business validation
			httpx.RespondWithJSON(w, struct {
				Valid  bool   `json:"valid"`
				Reason string `json:"reason"`
			}{
				Valid:  false,
				Reason: err.Error(),
			}, http.StatusOK)
		case errors.Is(err, errs.ErrDomainNotFound):
			// Return validation failure response for not found
			httpx.RespondWithJSON(w, struct {
				Valid  bool   `json:"valid"`
				Reason string `json:"reason"`
			}{
				Valid:  false,
				Reason: "coupon not found",
			}, http.StatusOK)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			log.Printf("Handler: Internal server error during coupon validation: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			log.Printf("Handler: Unhandled error from service during coupon validation: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	// Return successful validation with coupon details
	httpx.RespondWithJSON(w, struct {
		Valid  bool        `json:"valid"`
		Coupon interface{} `json:"coupon"`
	}{
		Valid:  true,
		Coupon: coupon,
	}, http.StatusOK)
}

