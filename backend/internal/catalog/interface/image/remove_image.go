package imageHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) RemoveImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request domain.ImageModifierRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevent clients from sending unexpected fields
	if err := decoder.Decode(&request); err != nil {
		log.Printf("Handler: Error decoding JSON body: %v", err)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteImage(ctx, &request); err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrExternalService):
			httpx.RespondWithError(w, fmt.Errorf("failed to delete image due to an external service error: %w", err), http.StatusServiceUnavailable)
		case errors.Is(err, errs.ErrUnexpectedError):
			log.Printf("Handler: Internal server error during image upload: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			log.Printf("Handler: Unhandled error from service during image deletion: %v", err) // Corrected log message
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

	log.Printf("Handler: Image %s deleted successfully.", request.ImageID)
}
