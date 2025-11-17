package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAllPartnersByProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract product ID from URL path parameter
	productID := r.PathValue("id")
	if productID == "" {
		logger.WarnContext(ctx, "Handler: Missing product ID in request",
			"operation", "get_partners_by_product",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("product ID is required"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get partners by product request",
		"operation", "get_partners_by_product",
		"method", r.Method,
		"path", r.URL.Path,
		"product_id", productID,
		"user_agent", r.Header.Get("User-Agent"))

	partners, err := h.svc.GetAllPartnersByProduct(ctx, productID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get partners by product")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get partners by product completed",
		"operation", "get_partners_by_product",
		"method", r.Method,
		"path", r.URL.Path,
		"product_id", productID,
		"status_code", http.StatusOK,
		"partner_count", len(partners))

	httpx.RespondWithJSON(w, partners, http.StatusOK)
}
