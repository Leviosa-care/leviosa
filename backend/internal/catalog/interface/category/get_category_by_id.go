package categoryHandler

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categoryID := strings.TrimPrefix(r.URL.Path, "/categories/")
	if categoryID == "" || strings.Contains(categoryID, "/") {
		httpx.RespondWithError(w, errors.New("invalid URL path"), http.StatusBadRequest)
		return
	}

	categoryWithImage, err := h.aggr.GetCategoryByIDWithImage(ctx, categoryID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			log.Printf("Handler: Internal server error during image retrieval: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred while getting image"), http.StatusInternalServerError)
		case errors.Is(err, errs.ErrDomainNotFound):
			// The category ID from the URL was not found.
			httpx.RespondWithError(w, err, http.StatusNotFound)
		default:
			log.Printf("Handler: Unhandled error from image service: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred while getting image"), http.StatusInternalServerError)
		}
	}

	httpx.RespondWithJSON(w, categoryWithImage, http.StatusOK)
}

// NOTE: the old thing
// func (h *handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
//
// 	categoryID := strings.TrimPrefix(r.URL.Path, "/categories/")
// 	if categoryID == "" || strings.Contains(categoryID, "/") {
// 		httpx.RespondWithError(w, errors.New("invalid URL path"), http.StatusBadRequest)
// 		return
// 	}
//
// 	image, err := h.imageSvc.GetActiveImage(ctx, categoryID, string(domain.CategoryType))
// 	if err != nil {
// 		// A NotFoundErr for the image is expected if no image exists.
// 		// We handle this by setting the image to nil and continuing.
// 		if !errors.Is(err, errs.ErrDomainNotFound) {
// 			// All other errors (invalid ID, DB error, etc.) should be handled as a failure.
// 			switch {
// 			case errors.Is(err, errs.ErrInvalidValue):
// 				httpx.RespondWithError(w, err, http.StatusBadRequest)
// 			case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
// 				log.Printf("Handler: Internal server error during image retrieval: %v", err)
// 				httpx.RespondWithError(w, errors.New("an internal server error occurred while getting image"), http.StatusInternalServerError)
// 			default:
// 				log.Printf("Handler: Unhandled error from image service: %v", err)
// 				httpx.RespondWithError(w, errors.New("an unexpected error occurred while getting image"), http.StatusInternalServerError)
// 			}
// 			return
// 		}
// 		image = nil
// 	}
//
// 	category, err := h.svc.GetCategoryByID(ctx, categoryID)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, errs.ErrInvalidValue):
// 			// This covers an empty ID.
// 			httpx.RespondWithError(w, err, http.StatusBadRequest)
// 		case errors.Is(err, errs.ErrDomainNotFound):
// 			// The category ID from the URL was not found.
// 			httpx.RespondWithError(w, err, http.StatusNotFound)
// 		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
// 			// A general internal server error occurred (DB, corrupt data, etc.).
// 			log.Printf("Handler: Internal server error during category retrieval: %v", err)
// 			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
// 		default:
// 			// A catch-all for any other unhandled errors.
// 			log.Printf("Handler: Unhandled error from service during category retrieval: %v", err)
// 			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
// 		}
// 		return
// 	}
//
// 	type Res struct {
// 		Image    *domain.Image    `json:"image"`
// 		Category *domain.Category `json:"category"`
// 	}
//
// 	httpx.RespondWithJSON(w, Res{Image: image, Category: category}, http.StatusOK)
// }
