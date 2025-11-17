package promotionCodeHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
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
		httpx.RespondWithServiceError(w, logger, ctx, err, "delete promotion code")
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
