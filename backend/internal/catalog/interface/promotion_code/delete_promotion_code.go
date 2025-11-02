package promotionCodeHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
)

func (h *handler) DeletePromotionCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("Handler: Processing delete_promotion_code", "promotion_code_id", "")

	promotionCodeID := strings.TrimPrefix(r.URL.Path, "/admin/promotion-codes/")
	if promotionCodeID == "" || strings.Contains(promotionCodeID, "/") {
		httpx.RespondWithError(w, errors.New("invalid URL path"), http.StatusBadRequest)
		return
	}

	err = h.svc.DeletePromotionCode(ctx, promotionCodeID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.Error("Handler: Internal server error during promotion code deletion", "error", err)
			httpx.RespondWithError(w, errors.New("internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.Error("Handler: Unhandled error from service during promotion code deletion", "error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.Info("Handler: Promotion code deletion successful", "promotion_code_id", promotionCodeID)
	httpx.RespondWithJSON(
		w,
		struct {
			Message string `json:"message"`
		}{
			Message: "Promotion code deleted successfully!",
		},
		http.StatusOK,
	)
}
