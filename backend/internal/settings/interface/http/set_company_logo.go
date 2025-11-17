package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) SetCompanyLogo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing set company logo request",
		"operation", "set_company_logo",
		"method", r.Method,
		"path", r.URL.Path)

	const maxFormMemory = 32 << 20            // 32 MB
	err = r.ParseMultipartForm(maxFormMemory) // 32 MB max memory buffer
	if err != nil {
		logger.DebugContext(ctx, fmt.Sprintf("Handler: Error parsing multipart form: %v", err))
		httpx.RespondWithError(w, errors.New("failed to parse form data, request too large or invalid"), http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	file, header, err := r.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			httpx.RespondWithError(w, errors.New("required form file 'image' is missing"), http.StatusBadRequest)
		} else {
			logger.ErrorContext(ctx, fmt.Sprintf("Handler: Error retrieving form file: %v", err))
			httpx.RespondWithError(w, errors.New("failed to retrieve form file"), http.StatusBadRequest)
		}
		return
	}
	defer file.Close()

	// header.Size is the file size in bytes (int64).
	fileSize := header.Size
	contentType := header.Header.Get("Content-Type")

	if _, err := h.svc.SetCompanyLogo(ctx, file, fileSize, contentType); err != nil {
		if errors.Is(err, errs.ErrExternalService) {
			httpx.RespondWithError(w, fmt.Errorf("failed to upload image due to an external service error: %w", err), http.StatusServiceUnavailable)
		} else {
			httpx.RespondWithServiceError(w, logger, ctx, err, "upload company logo")
		}
		return
	}

	logger.InfoContext(ctx, "Handler: Set company logo completed",
		"operation", "set_company_logo",
		"status_code", http.StatusCreated)

	w.WriteHeader(http.StatusCreated)
}
