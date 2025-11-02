package categoryHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) RemoveCategory(w http.ResponseWriter, r *http.Request) {
	// TODO: this an admin only request

	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing remove category",
		"operation", "remove_category",
		"method", r.Method,
		"path", r.URL.Path)
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

	if err := h.aggr.RemoveCategoryWithImages(ctx, categoryID); err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrExternalService):
			httpx.RespondWithError(w, errors.New("failed to delete category images from external storage"), http.StatusServiceUnavailable)
		case errors.Is(err, errs.ErrCategoryHasProducts): // Or errors.Is(err, errs.ErrConflict) if you only return ErrConflict from service
			httpx.RespondWithError(w, errors.New("cannot delete category: it still has associated products"), http.StatusConflict)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.ErrorContext(ctx, "Handler: remove category failed",
				"operation", "remove_category",
				"error_context", "internal server error during image deletion",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred during image cleanup"), http.StatusInternalServerError)
		default:
			logger.ErrorContext(ctx, "Handler: remove category failed",
				"operation", "remove_category",
				"error_context", "unexpected error from image service during deletion",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.InfoContext(ctx, "Handler: remove category completed",
		"operation", "remove_category",
		"category_id", categoryID,
		"status_code", http.StatusOK)

	w.WriteHeader(http.StatusNoContent)
}

// NOTE: the old way of doing
// func (h *handler) RemoveCategory(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
//
// 	// TODO: this an admin only request
// 	parts := strings.Split(r.URL.Path, "/")
// 	if len(parts) != 4 || parts[0] != "" || parts[1] != "admin" || parts[2] != "categories" {
// 		httpx.RespondWithError(w, errors.New("invalid URL path format. Expected /admin/categories/{id}"), http.StatusBadRequest)
// 		return
// 	}
// 	categoryID := parts[3] // The ID should be the last part
//
// 	if categoryID == "" {
// 		httpx.RespondWithError(w, errors.New("category ID is missing from the URL"), http.StatusBadRequest)
// 		return
// 	}
//
// 	if err := h.imageSvc.DeleteImages(ctx, categoryID, "category"); err != nil {
// 		switch {
// 		case errors.Is(err, errs.ErrInvalidValue):
// 			httpx.RespondWithError(w, err, http.StatusBadRequest)
// 		case errors.Is(err, errs.ErrExternalService):
// 			httpx.RespondWithError(w, errors.New("failed to delete category images from external storage"), http.StatusServiceUnavailable)
// 		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
// 			log.Printf("Handler: Internal server error during image deletion: %v", err)
// 			httpx.RespondWithError(w, errors.New("an internal server error occurred during image cleanup"), http.StatusInternalServerError)
// 		default:
// 			log.Printf("Handler: Unhandled error from image service during deletion: %v", err)
// 			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
// 		}
// 		return
// 	}
//
// 	if err := h.svc.RemoveCategory(ctx, categoryID); err != nil {
// 		switch {
// 		case errors.Is(err, errs.ErrInvalidValue):
// 			httpx.RespondWithError(w, err, http.StatusBadRequest)
// 		case errors.Is(err, errs.ErrDomainNotFound):
// 			httpx.RespondWithError(w, err, http.StatusNotFound)
// 		case errors.Is(err, errs.ErrExternalService):
// 			httpx.RespondWithError(w, errors.New("failed to delete category images due to external storage issue"), http.StatusServiceUnavailable)
// 		case errors.Is(err, errs.ErrCategoryHasProducts): // Or errors.Is(err, errs.ErrConflict) if you only return ErrConflict from service
// 			httpx.RespondWithError(w, errors.New("cannot delete category: it still has associated products"), http.StatusConflict)
// 		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
// 			log.Printf("Handler: Internal server error during category deletion: %v", err)
// 			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
// 		default:
// 			log.Printf("Handler: Unhandled error from service during category deletion: %v", err)
// 			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
// 		}
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusNoContent)
// 	log.Println("Handler: Category removed successfully.")
// }
