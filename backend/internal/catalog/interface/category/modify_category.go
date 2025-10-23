package categoryHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) ModifyCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO: this an admin only request

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 || parts[0] != "" || parts[1] != "admin" || parts[2] != "categories" {
		httpx.RespondWithError(w, errors.New("invalid URL path format. Expected /admin/categories/{id}"), http.StatusBadRequest)
		return
	}
	categoryID := parts[3] // The ID should be the last part
	if categoryID == "" {
		httpx.RespondWithError(w, errors.New("category ID is missing from the URL"), http.StatusBadRequest)
		return
	}

	var req domain.UpdateCategoryRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevent clients from sending unexpected fields
	if err := decoder.Decode(&req); err != nil {
		log.Printf("Handler: Error decoding JSON body: %v", err)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	if req.Name == nil && req.Description == nil && req.Status == nil && req.Metadata == nil {
		httpx.RespondWithError(w, errors.New("no updatable fields provided in request body"), http.StatusBadRequest)
		return
	}

	req.ID = categoryID

	if err := h.svc.UpdateCategory(ctx, &req); err != nil {
		// TODO: better error handling to return the proper status
		// 204 (no content) : since there is no body in response, if there was return 200
		// 400 (bad request) : mal formed input, invalid format
		// 404 (not found) : no category with given ID
		// 409 (conflit) : unique constraint violated
		// 500 (Internal Server error) : server error, something broke
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)

		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)

		case errors.Is(err, errs.ErrConflict):
			// This covers unique constraint violations (e.g., duplicate name) or other errs-level conflicts.
			httpx.RespondWithError(w, err, http.StatusConflict)

		case errors.Is(err, errs.ErrExternalService):
			// If you had an external service call in update (e.g., image re-upload)
			httpx.RespondWithError(w, errors.New("external service error during category update"), http.StatusServiceUnavailable)

		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			log.Printf("Handler: Internal server error during category update: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)

		default:
			log.Printf("Handler: Unhandled error from service during category update: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
