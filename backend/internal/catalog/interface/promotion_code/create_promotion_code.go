package promotionCodeHandler

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

func (h *handler) CreatePromotionCode(w http.ResponseWriter, r *http.Request) {
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

	logger.Info("Handler: Processing create_promotion_code", "promotion_code_id", "")

	var promotionCode domain.CreatePromotionCodeRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&promotionCode); err != nil {
		logger.Error("Handler: Error decoding JSON body", "error", err)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	promotionCodeID, err := h.svc.CreatePromotionCode(ctx, &promotionCode)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "create promotion code")
		return
	}

	logger.Info("Handler: Promotion code creation successful", "promotion_code_id", promotionCodeID)
	httpx.RespondWithJSON(
		w,
		struct {
			ID      string `json:"id"`
			Message string `json:"message"`
		}{
			ID:      promotionCodeID,
			Message: "Promotion code created successfully!",
		},
		http.StatusCreated,
	)
}
