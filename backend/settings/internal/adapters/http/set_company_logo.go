package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/middleware"
)

func (h *handler) SetCompanyLogo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		middleware.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	const maxFormMemory = 32 << 20            // 32 MB
	err = r.ParseMultipartForm(maxFormMemory) // 32 MB max memory buffer
	if err != nil {
		logger.DebugContext(ctx, fmt.Sprintf("Handler: Error parsing multipart form: %v", err))
		middleware.RespondWithError(w, errors.New("failed to parse form data, request too large or invalid"), http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	file, header, err := r.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			middleware.RespondWithError(w, errors.New("required form file 'image' is missing"), http.StatusBadRequest)
		} else {
			logger.ErrorContext(ctx, fmt.Sprintf("Handler: Error retrieving form file: %v", err))
			middleware.RespondWithError(w, errors.New("failed to retrieve form file"), http.StatusBadRequest)
		}
		return
	}
	defer file.Close()

	// header.Size is the file size in bytes (int64).
	fileSize := header.Size
	contentType := header.Header.Get("Content-Type")

	if _, err := h.svc.SetCompanyLogo(ctx, file, fileSize, contentType); err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			middleware.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrExternalService):
			middleware.RespondWithError(w, fmt.Errorf("failed to upload image due to an external service error: %w", err), http.StatusServiceUnavailable)
		case errors.Is(err, errs.ErrUnexpectedError), errors.Is(err, errs.ErrUnexpectedError):
			logger.ErrorContext(ctx, fmt.Sprintf("Handler: Internal server error during image upload: %v", err))
			middleware.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.ErrorContext(ctx, fmt.Sprintf("Handler: Unhandled error from service during image upload: %v", err))
			middleware.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	logger.InfoContext(ctx, fmt.Sprintf("Handler: Image created successfully."))
}
