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

func (h *handler) ValidatePromotionCode(w http.ResponseWriter, r *http.Request) {
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

	logger.Info("Handler: Processing validate_promotion_code", "promotion_code", "")

	var validateRequest domain.ValidatePromotionCodeRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&validateRequest); err != nil {
		logger.Error("Handler: Error decoding JSON body", "error", err)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	validationResponse, err := h.svc.ValidatePromotionCode(ctx, &validateRequest)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "validate promotion code")
		return
	}

	// Return validation response regardless of validity
	logger.Info("Handler: Promotion code validation successful", "promotion_code", validateRequest.Code)
	httpx.RespondWithJSON(w, validationResponse, http.StatusOK)
}
