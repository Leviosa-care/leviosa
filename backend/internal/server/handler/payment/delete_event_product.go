package payment

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/hengadev/leviosa/internal/server/handler"
	mw "github.com/hengadev/leviosa/internal/server/middleware"
	"github.com/hengadev/leviosa/pkg/ctxutil"
	"github.com/hengadev/leviosa/pkg/jsonio"

	"github.com/stripe/stripe-go/v79"
)

// TODO: finish the implementation of that thing.
// func DeleteEventProduct(p *payment.Service, e *event.Service) http.Handler {
func (a *AppInstance) DeleteEventProduct() http.Handler {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get event id
		ctx := r.Context()

		logger, err := ctxutil.GetLoggerFromContext(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "logger not found in context", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		eventID := r.PathValue("id")
		_, err = a.Svcs.Stripe.RemoveProduct(ctx, eventID)
		if err != nil {
			logger.ErrorContext(ctx, "failed to remove product", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
			return
		}
		type Response struct {
			EventID string `json:"eventid"`
		}
		if err = jsonio.Encode(w, http.StatusCreated, Response{eventID}); err != nil {
			logger.ErrorContext(ctx, "failed to encode eventID for product registered", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
			return
		}
	})
	return mw.EnableHeaders(handler, "Authorization")
}
