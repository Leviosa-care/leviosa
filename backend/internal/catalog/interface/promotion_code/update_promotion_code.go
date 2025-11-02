package promotionCodeHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
)

func (h *handler) UpdatePromotionCode(w http.ResponseWriter, r *http.Request) {
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

	logger.Info("Handler: Processing update_promotion_code", "promotion_code_id", "")

	promotionCodeID := strings.TrimPrefix(r.URL.Path, "/admin/promotion-codes/")
	if promotionCodeID == "" || strings.Contains(promotionCodeID, "/") {
		httpx.RespondWithError(w, errors.New("invalid URL path"), http.StatusBadRequest)
		return
	}

	var updateRequest domain.UpdatePromotionCodeRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&updateRequest); err != nil {
		logger.Error("Handler: Error decoding JSON body", "error", err)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	err = h.svc.UpdatePromotionCode(ctx, promotionCodeID, &updateRequest)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.Error("Handler: Internal server error during promotion code update", "error", err)
			httpx.RespondWithError(w, errors.New("internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.Error("Handler: Unhandled error from service during promotion code update", "error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.Info("Handler: Promotion code update successful", "promotion_code_id", promotionCodeID)
	httpx.RespondWithJSON(
		w,
		struct {
			Message string `json:"message"`
		}{
			Message: "Promotion code updated successfully!",
		},
		http.StatusOK,
	)
}

func (h *handler) DeactivatePromotionCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	logger.Info("Handler: Processing deactivate_promotion_code", "promotion_code_id", "")

	// Extract ID from URL path like /admin/promotion-codes/{id}/deactivate
	urlPath := r.URL.Path
	parts := strings.Split(urlPath, "/")
	if len(parts) < 4 || parts[len(parts)-1] != "deactivate" {
		httpx.RespondWithError(w, errors.New("invalid URL path"), http.StatusBadRequest)
		return
	}

	promotionCodeID := parts[len(parts)-2]
	if promotionCodeID == "" {
		httpx.RespondWithError(w, errors.New("invalid promotion code ID"), http.StatusBadRequest)
		return
	}

	err = h.svc.DeactivatePromotionCode(ctx, promotionCodeID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.Error("Handler: Internal server error during promotion code deactivation", "error", err)
			httpx.RespondWithError(w, errors.New("internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.Error("Handler: Unhandled error from service during promotion code deactivation", "error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.Info("Handler: Promotion code deactivation successful", "promotion_code_id", promotionCodeID)
	httpx.RespondWithJSON(
		w,
		struct {
			Message string `json:"message"`
		}{
			Message: "Promotion code deactivated successfully!",
		},
		http.StatusOK,
	)
}
