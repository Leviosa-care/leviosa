package imageHandler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/hengadev/errsx"
)

func (h *handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("Handler: Processing upload_image", "image_id", "")

	const maxFormMemory = 32 << 20            // 32 MB
	err = r.ParseMultipartForm(maxFormMemory) // 32 MB max memory buffer
	if err != nil {
		logger.Error("Handler: Error parsing multipart form", "error", err)
		httpx.RespondWithError(w, errors.New("failed to parse form data, request too large or invalid"), http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	file, header, err := r.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			httpx.RespondWithError(w, errors.New("required form file 'image' is missing"), http.StatusBadRequest)
		} else {
			logger.Error("Handler: Error retrieving form file", "error", err)
			httpx.RespondWithError(w, errors.New("failed to retrieve form file"), http.StatusBadRequest)
		}
		return
	}
	defer file.Close()

	// header.Size is the file size in bytes (int64).
	fileSize := header.Size
	contentType := header.Header.Get("Content-Type")

	// fields check
	var fieldErrs errsx.Map
	parentID := r.FormValue("parent_id")
	if parentID == "" {
		fieldErrs.Set("parent ID", "required form field 'parent_id' is missing")
	}

	parentType := r.FormValue("parent_type")
	if parentType == "" {
		fieldErrs.Set("parent type", "required form field 'parent_type' is missing")
	}

	title := r.FormValue("title")
	if title == "" {
		fieldErrs.Set("title", "required form field 'title' is missing")
	}

	isActiveStr := r.FormValue("is_active")
	isActive, err := strconv.ParseBool(isActiveStr)
	if err != nil {
		isActive = false
	}

	if fieldErrs != nil {
		httpx.RespondWithError(w, fieldErrs.AsError(), http.StatusBadRequest)
		return
	}

	req := &domain.CreateImageRequest{
		ParentID:   parentID,
		ParentType: domain.ParentType(parentType),
		Title:      title,
		IsActive:   &isActive, // & to make it a pointer
	}

	imageID, err := h.svc.AddImage(ctx, req, file, fileSize, contentType)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "upload image")
		return
	}

	logger.Info("Handler: Image created successfully", "image_id", imageID)
	httpx.RespondWithJSON(w, struct {
		ID string `json:"id"`
	}{ID: imageID}, http.StatusCreated)
}
