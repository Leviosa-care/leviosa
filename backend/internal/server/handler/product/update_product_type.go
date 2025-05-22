package productHandler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/product"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/ctxutil"
	"github.com/hengadev/leviosa/pkg/jsonio"
)

func (a *AppInstance) UpdateProductType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "logger not found in context", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := ctxutil.ValidateRoleInContext(ctx, models.ADMINISTRATOR); err != nil {
		logger.ErrorContext(ctx, "validate role from context", "error", err)
		http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
		return
	}

	productType, err := jsonio.DecodeValid[productService.Offer](ctx, r.Body)
	if err != nil {
		switch {
		case errors.Is(err, jsonio.ErrDecodeJSON):
			logger.WarnContext(ctx, "decode product  type", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		default:
			logger.WarnContext(ctx, "invalid product type creation", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
		return
	}

	// TODO: now get the product from this thing
	err = a.Svcs.Product.UpdateProductType(ctx, &productType)
	if err != nil {
		// TODO: handle the remaining the update product
		switch {
		case errors.Is(err, domain.ErrNotFound):
		case errors.Is(err, rp.ErrContext):
		}
		return
	}

	http.Error(w, fmt.Sprintf("product with name %s udpated", productType.Name), http.StatusNoContent)
}
