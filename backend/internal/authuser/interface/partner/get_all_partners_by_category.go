package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAllPartnersByCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract category ID from URL path parameter
	categoryID := r.PathValue("id")
	if categoryID == "" {
		logger.WarnContext(ctx, "Handler: Missing category ID in request",
			"operation", "get_partners_by_category",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("category ID is required"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get partners by category request",
		"operation", "get_partners_by_category",
		"method", r.Method,
		"path", r.URL.Path,
		"category_id", categoryID,
		"user_agent", r.Header.Get("User-Agent"))

	partners, err := h.svc.GetAllPartnersByCategory(ctx, categoryID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get all partners by category")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get partners by category completed",
		"operation", "get_partners_by_category",
		"method", r.Method,
		"path", r.URL.Path,
		"category_id", categoryID,
		"status_code", http.StatusOK,
		"partner_count", len(partners))

	httpx.RespondWithJSON(w, partners, http.StatusOK)
}
