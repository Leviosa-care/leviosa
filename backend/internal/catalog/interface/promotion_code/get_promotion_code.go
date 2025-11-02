package promotionCodeHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
)

func (h *handler) GetPromotionCodeByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("Handler: Processing get_promotion_code", "promotion_code_id", "")

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
			logger.Error("Handler: Internal server error during promotion code retrieval", "error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.Error("Handler: Unhandled error from service during promotion code retrieval", "error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.Info("Handler: Promotion code retrieval successful", "promotion_code_id", promotionCodeID)
	httpx.RespondWithJSON(w, promotionCode, http.StatusOK)
}

func (h *handler) GetPromotionCodeByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)

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
			logger.Error("Handler: Internal server error during promotion code retrieval", "error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.Error("Handler: Unhandled error from service during promotion code retrieval", "error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.Info("Handler: Promotion code retrieval by code successful", "promotion_code", code)
	httpx.RespondWithJSON(w, promotionCode, http.StatusOK)
}

func (h *handler) GetAllPromotionCodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)

	promotionCodes, err := h.svc.GetAllPromotionCodes(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.Error("Handler: Internal server error during promotion codes retrieval", "error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.Error("Handler: Unhandled error from service during promotion codes retrieval", "error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.Info("Handler: All promotion codes retrieval successful")
	httpx.RespondWithJSON(w, promotionCodes, http.StatusOK)
}

func (h *handler) GetActivePromotionCodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)

	promotionCodes, err := h.svc.GetActivePromotionCodes(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.Error("Handler: Internal server error during active promotion codes retrieval", "error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.Error("Handler: Unhandled error from service during active promotion codes retrieval", "error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.Info("Handler: Active promotion codes retrieval successful")
	httpx.RespondWithJSON(w, promotionCodes, http.StatusOK)
}

func (h *handler) GetPromotionCodeWithCoupon(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)

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
			logger.Error("Handler: Internal server error during promotion code with coupon retrieval", "error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.Error("Handler: Unhandled error from service during promotion code with coupon retrieval", "error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.Info("Handler: Promotion code with coupon retrieval successful", "promotion_code", code)
	httpx.RespondWithJSON(w, promotionCodeWithCoupon, http.StatusOK)
}
