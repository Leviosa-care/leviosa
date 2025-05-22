package productHandler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/ctxutil"
)

func (a *AppInstance) DeleteProduct(w http.ResponseWriter, r *http.Request) {
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

	productID := r.URL.Query().Get("id")
	if productID == "" {
		// TODO: do some logger for when the ID provided is invalid
	}

	// TODO: maybe return product name so that I can return the right message
	err = a.Svcs.Product.RemoveProduct(ctx, productID)
	if err != nil {
		// TODO: handle the remaining the update product
		switch {
		case errors.Is(err, domain.ErrNotFound):
		case errors.Is(err, rp.ErrContext):
		}
		return
	}

	http.Error(w, fmt.Sprintf("product with ID %s successfully deleted", productID), http.StatusNoContent)
}
