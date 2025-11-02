package imageHandler

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

func (h *handler) SetActiveImage(w http.ResponseWriter, r *http.Request) {
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

	logger.Info("Handler: Processing set_active_image", "image_id", "")

	var request domain.ImageModifierRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevent clients from sending unexpected fields
	if err := decoder.Decode(&request); err != nil {
		logger.Error("Handler: Error decoding JSON body", "error", err)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	if err := h.svc.SetActiveImage(ctx, &request); err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrUnexpectedError):
			logger.Error("Handler: Internal server error while setting image to active", "error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		case errors.Is(err, errs.ErrQueryFailed):
			logger.Error("Handler: Internal database error while setting image to active", "error", err)
			httpx.RespondWithError(w, errors.New("an internal database error occurred"), http.StatusInternalServerError)
		default:
			logger.Error("Handler: Unhandled error from service while setting image to active", "error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.Info("Handler: Image set as active successfully", "image_id", request.ImageID)
	w.WriteHeader(http.StatusNoContent)
}
