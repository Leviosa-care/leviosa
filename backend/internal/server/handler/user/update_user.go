package userHandler

import (
	// "errors"
	"net/http"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/ctxutil"
	"github.com/hengadev/leviosa/pkg/jsonio"
)

func (a *AppInstance) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "logger not found in context", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// userID, ok := ctx.Value(ctxutil.UserIDKey).(string)
	// if !ok {
	// 	logger.ErrorContext(ctx, "user ID not found in context")
	// 	jsonio.WriteResponse(w, errors.New("failed to get user ID from context").Error(), http.StatusInternalServerError)
	// 	return
	// }

	// use a custom valid for the updtate thing
	user, err := jsonio.Decode[models.User](r.Body)
	if err != nil {
		logger.ErrorContext(ctx, "failed to decode user", "error", err)
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return
	}
	// modify user
	if err = a.Svcs.User.UpdateAccount(ctx, &user); err != nil {
		switch err {
		// TODO: handle the validation error to just send back the fields that are not updated because prohibited from updates
		default:
			logger.ErrorContext(ctx, "failed to modify the user", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
			return
		}
	}
}
