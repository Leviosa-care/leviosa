package promotionCodeHandler

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetPromotionCodeByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	promotionCodeID := strings.TrimPrefix(r.URL.Path, "/admin/promotion-codes/")
	if promotionCodeID == "" || strings.Contains(promotionCodeID, "/") {
		httpx.RespondWithError(w, errors.New("invalid URL path"), http.StatusBadRequest)
		return
	}

	promotionCode, err := h.svc.GetPromotionCodeByID(ctx, promotionCodeID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			log.Printf("Handler: Internal server error during promotion code retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			log.Printf("Handler: Unhandled error from service during promotion code retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	httpx.RespondWithJSON(w, promotionCode, http.StatusOK)
}

func (h *handler) GetPromotionCodeByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract code from URL path
	urlPath := r.URL.Path
	var code string
	if strings.HasPrefix(urlPath, "/admin/promotion-codes/code/") {
		code = strings.TrimPrefix(urlPath, "/admin/promotion-codes/code/")
	} else if strings.HasPrefix(urlPath, "/promotion-codes/code/") {
		code = strings.TrimPrefix(urlPath, "/promotion-codes/code/")
	}

	if code == "" || strings.Contains(code, "/") {
		httpx.RespondWithError(w, errors.New("invalid promotion code"), http.StatusBadRequest)
		return
	}

	promotionCode, err := h.svc.GetPromotionCodeByCode(ctx, code)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			log.Printf("Handler: Internal server error during promotion code retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			log.Printf("Handler: Unhandled error from service during promotion code retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	httpx.RespondWithJSON(w, promotionCode, http.StatusOK)
}

func (h *handler) GetAllPromotionCodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	promotionCodes, err := h.svc.GetAllPromotionCodes(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			log.Printf("Handler: Internal server error during promotion codes retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			log.Printf("Handler: Unhandled error from service during promotion codes retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	httpx.RespondWithJSON(w, promotionCodes, http.StatusOK)
}

func (h *handler) GetActivePromotionCodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	promotionCodes, err := h.svc.GetActivePromotionCodes(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			log.Printf("Handler: Internal server error during active promotion codes retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			log.Printf("Handler: Unhandled error from service during active promotion codes retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	httpx.RespondWithJSON(w, promotionCodes, http.StatusOK)
}

func (h *handler) GetPromotionCodeWithCoupon(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := strings.TrimPrefix(r.URL.Path, "/promotion-codes/code/")
	if code == "" || strings.Contains(code, "/") {
		httpx.RespondWithError(w, errors.New("invalid promotion code"), http.StatusBadRequest)
		return
	}

	promotionCodeWithCoupon, err := h.svc.GetPromotionCodeWithCoupon(ctx, code)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			log.Printf("Handler: Internal server error during promotion code with coupon retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			log.Printf("Handler: Unhandled error from service during promotion code with coupon retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	httpx.RespondWithJSON(w, promotionCodeWithCoupon, http.StatusOK)
}
