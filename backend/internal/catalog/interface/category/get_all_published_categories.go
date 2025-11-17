package categoryHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAllPublishedCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get all published categories",
		"operation", "get_all_published_categories",
		"method", r.Method,
		"path", r.URL.Path)

	categoryWithImages, err := h.aggr.GetAllPublishedCategoriesWithImages(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get all published categories")
		return
	}

	count := len(categoryWithImages)
	logger.InfoContext(ctx, "Handler: get all published categories completed",
		"operation", "get_all_published_categories",
		"count", count,
		"status_code", http.StatusOK)

	// httpx.RespondWithJSON(w, categories, http.StatusOK)
	httpx.RespondWithJSON(w, categoryWithImages, http.StatusOK)
}
