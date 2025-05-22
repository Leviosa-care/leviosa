package productHandler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/ctxutil"
)

func (a *AppInstance) DeleteOffer(w http.ResponseWriter, r *http.Request) {
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

	// this product ID is an int
	offerIDStr := r.URL.Query().Get("id")
	offerID, err := strconv.Atoi(offerIDStr)
	if err != nil {
		// TODO: do some error thing so that the user can know that they an invalid productID
	}

	// TODO: maybe return product name so that I can return the right message
	err = a.Svcs.Product.RemoveOffer(ctx, offerID)
	if err != nil {
		// TODO: handle the remaining the update product
		switch {
		case errors.Is(err, domain.ErrNotFound):
		case errors.Is(err, rp.ErrContext):
		}
		return
	}

	http.Error(w, fmt.Sprintf("product with ID %d successfully deleted", offerID), http.StatusNoContent)
}
