package promotionCodeHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
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
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrAlreadyExists):
			httpx.RespondWithError(w, err, http.StatusConflict)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrDomainNotCreated):
			httpx.RespondWithError(w, errors.New("failed to create promotion code due to an unprocessable entity"), http.StatusUnprocessableEntity)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.Error("Handler: Internal server error during promotion code creation", "error", err)
			httpx.RespondWithError(w, errors.New("internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.Error("Handler: Unhandled error from service during promotion code creation", "error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
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
