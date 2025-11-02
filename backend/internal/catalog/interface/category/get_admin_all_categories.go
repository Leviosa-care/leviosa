package categoryHandler

import (
	"errors"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAdminAllCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get admin all categories",
		"operation", "get_admin_all_categories",
		"method", r.Method,
		"path", r.URL.Path)

	categoryWithImages, err := h.aggr.GetAdminAllCategoriesWithImages(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			logger.ErrorContext(ctx, "Handler: get admin all categories failed",
				"operation", "get_admin_all_categories",
				"error_context", "invalid value error",
				"status_code", http.StatusBadRequest,
				"error", err)
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrQueryFailed):
			logger.ErrorContext(ctx, "Handler: get admin all categories failed",
				"operation", "get_admin_all_categories",
				"error_context", "query failed",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
		default:
			logger.ErrorContext(ctx, "Handler: get admin all categories failed",
				"operation", "get_admin_all_categories",
				"error_context", "internal server error",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("internal server occurred"), http.StatusInternalServerError)
		}
		return
	}

	count := len(categoryWithImages)
	logger.InfoContext(ctx, "Handler: get admin all categories completed",
		"operation", "get_admin_all_categories",
		"count", count,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, categoryWithImages, http.StatusOK)
}

// NOTE: the old way of doing things, N+1 query
// func (h *handler) GetAdminAllCategories(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	categories, err := h.svc.GetAllCategories(ctx)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
// 			// A general internal server error occurred (DB, corrupt data, etc.).
// 			log.Printf("Handler: Internal server error during categories retrieval: %v", err)
// 			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
// 		default:
// 			// A catch-all for any other unhandled errors.
// 			log.Printf("Handler: Unhandled error from service during categories retrieval: %v", err)
// 			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
// 		}
// 		return
// 	}
//
// 	responseSlice := make([]domain.CategoryWithImage, 0, len(categories))
//
// 	// Use a wait group to concurrently fetch images, if you want to optimize for speed.
// 	// For this example, we'll keep it simple and sequential to avoid unnecessary complexity,
// 	// but concurrency is a good consideration for performance.
// 	for _, category := range categories {
// 		// Attempt to get the active image for the current category.
// 		// image, err := h.image.GetActiveImage(ctx, category.ID.String(), string(domain.CategoryType))
// 		image, err := h.aggr.GetActiveImage(ctx, category.ID.String(), string(domain.CategoryType))
//
// 		// If no image is found, treat it as a success but with a nil image.
// 		if err != nil && !errors.Is(err, errs.ErrDomainNotFound) {
// 			// If there's an error other than "not found," this is a server error.
// 			log.Printf("Handler: Internal server error during image retrieval for category %s: %v", category.ID, err)
// 			httpx.RespondWithError(w, errors.New("an internal server error occurred while getting images"), http.StatusInternalServerError)
// 			return
// 		}
//
// 		// If the image was not found, the image variable will be nil.
// 		responseSlice = append(responseSlice, domain.CategoryWithImage{
// 			Category: category,
// 			Image:    image,
// 		})
// 	}
//
// 	// httpx.RespondWithJSON(w, categories, http.StatusOK)
// 	httpx.RespondWithJSON(w, responseSlice, http.StatusOK)
// }
