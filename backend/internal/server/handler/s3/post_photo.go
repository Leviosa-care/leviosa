package photohandler

import (
	"log/slog"
	"net/http"

	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/ctxutil"
	"github.com/hengadev/leviosa/pkg/serverutil"
)

// func PostPhoto(ph *photo.Service) http.Handler {
func (a *AppInstance) PostPhoto() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger, err := ctxutil.GetLoggerFromContext(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "logger not found in context", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: this is not something general I need to remove that from the posting thing
		eventID := r.PathValue("id")
		if eventID == "" {
			// return  some error because I need that eventID
		}

		// Limit upload size to 10MB
		if err = r.ParseMultipartForm(10 << 20); err != nil {
			// add some logger to if the file is too big
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		file, fileheader, err := r.FormFile("photo")
		if err != nil {
			logger.ErrorContext(ctx, "failed to get the photo information from form", "error", err)
			http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
			return
		}

		// post file to the right bucket
		url, err := a.Svcs.Media.PostFile(ctx, file, fileheader.Filename, eventID)
		if err != nil {
			logger.ErrorContext(ctx, "failed to post file", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
			return
		}
		// TODO: see if I can send the url without making an object.
		// send the url back to the user
		// type Response struct {
		// 	URL string `json:"url"`
		// }
		// if err := serverutil.Encode(w, http.StatusSeeOther, Response{URL: url}); err != nil {
		if err := serverutil.Encode(w, http.StatusSeeOther, url); err != nil {
			logger.ErrorContext(ctx, "failed to send url back to client", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
			return
		}
	})
}
